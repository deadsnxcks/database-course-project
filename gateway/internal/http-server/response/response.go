package response

import (
	"log/slog"
	"net/http"

	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"

	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var httpByCode = map[codes.Code]int{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.FailedPrecondition: http.StatusConflict,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.Unavailable:        http.StatusServiceUnavailable,
}

func GRPCError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error("non-grpc error from client", sl.Err(err))
		writeError(w, r, http.StatusInternalServerError, "Internal", "internal error")

		return
	}

	httpStatus, known := httpByCode[st.Code()]
	if !known {
		httpStatus = http.StatusInternalServerError
	}

	if httpStatus == http.StatusInternalServerError {
		log.Error("grpc call failed", sl.Err(err))
		writeError(w, r, httpStatus, st.Code().String(), "internal error")

		return
	}

	log.Debug("grpc call rejected",
		slog.String("code", st.Code().String()),
		slog.String("message", st.Message()),
	)
	writeError(w, r, httpStatus, st.Code().String(), st.Message())
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	writeError(w, r, http.StatusBadRequest, codes.InvalidArgument.String(), msg)
}

func writeError(w http.ResponseWriter, r *http.Request, httpStatus int, code, msg string) {
	render.Status(r, httpStatus)
	render.JSON(w, r, ErrorBody{Code: code, Message: msg})
}