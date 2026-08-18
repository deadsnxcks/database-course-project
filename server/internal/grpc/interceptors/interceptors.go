package interceptors

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/deadsnxcks/dbcp/server/internal/domain"
	"github.com/deadsnxcks/dbcp/server/internal/lib/logger/sl"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// codeOf — единственное место, где домен встречается с транспортом.
// Появится REST-фасад — рядом ляжет httpStatusOf с теми же категориями.
func codeOf(k domain.Kind) codes.Code {
	switch k {
	case domain.KindNotFound:
		return codes.NotFound
	case domain.KindConflict:
		return codes.AlreadyExists
	case domain.KindInUse, domain.KindUnprocessable:
		return codes.FailedPrecondition
	default:
		return codes.Internal
	}
}

// ErrorMapping переводит доменную ошибку в gRPC-статус.
// ВНЕШНИЙ в цепочке: подменяет ошибку последним, уже после логирования.
func ErrorMapping() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Хендлер мог сам вернуть status (InvalidArgument при проверке
		// формы запроса) — такой ответ не трогаем.
		if _, ok := status.FromError(err); ok {
			return resp, err
		}

		var de *domain.Error
		if errors.As(err, &de) {
			return resp, status.Error(codeOf(de.Kind), de.Msg)
		}

		// Неопознанная ошибка — наружу обезличенно, подробности в логе.
		return resp, status.Error(codes.Internal, "internal error")
	}
}

// Logging пишет одну строку на вызов. ВНУТРЕННИЙ в цепочке:
// видит исходную ошибку со всей цепочкой op до подмены.
func Logging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("took", time.Since(start)),
		}

		switch {
		case err == nil:
			log.Debug("rpc handled", attrs...)

		case domain.KindOf(err) != domain.KindUnknown:
			// Ожидаемый отказ: не найдено, конфликт, правило не выполнено.
			// Это не сбой сервера — уровень Debug, алерты не поднимаются.
			log.Debug("rpc rejected",
				append(attrs,
					slog.String("kind", domain.KindOf(err).String()),
					sl.Err(err),
				)...)

		default:
			// Всё остальное — настоящая ошибка. sl.Err(err) печатает
			// полную цепочку op: services.vessel.Get: storage.postgresql...
			log.Error("rpc failed", append(attrs, sl.Err(err))...)
		}

		return resp, err
	}
}
