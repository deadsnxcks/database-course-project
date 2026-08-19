package operationservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/deadsnxcks/dbcp/server/internal/domain"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
)

const opStart = "services.operation"

type OperationService struct {
	log       *slog.Logger
	oProvider OperationProvider
}

type OperationProvider interface {
	Operations(ctx context.Context) ([]models.Operation, error)
	SaveOperation(ctx context.Context, operation models.Operation) (int64, error)
	DeleteOperation(ctx context.Context, id int64) error
	Operation(ctx context.Context, id int64) (models.Operation, error)
	UpdateOperation(
		ctx context.Context,
		id int64,
		title *string,
	) error
}

func New(
	log *slog.Logger,
	oProvider OperationProvider,
) *OperationService {
	return &OperationService{
		log:       log,
		oProvider: oProvider,
	}
}

func (o *OperationService) List(ctx context.Context) ([]models.Operation, error) {
	const op = opStart + ".List"

	ops, err := o.oProvider.Operations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ops, nil
}

func (o *OperationService) Get(ctx context.Context, id int64) (models.Operation, error) {
	const op = opStart + ".Get"

	opModel, err := o.oProvider.Operation(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrOperationNotFound) {
			return models.Operation{}, fmt.Errorf("%s: %w", op, domain.ErrOperationNotFound)
		}

		return models.Operation{}, fmt.Errorf("%s: %w", op, err)
	}

	return opModel, nil
}

func (o *OperationService) Create(ctx context.Context, title string) (int64, error) {
	const op = opStart + ".Create"

	id, err := o.oProvider.SaveOperation(ctx, models.Operation{Title: title})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	o.log.Info("operation created", slog.Int64("id", id))

	return id, nil
}

func (o *OperationService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	if err := o.oProvider.DeleteOperation(ctx, id); err != nil {
		switch {
		case errors.Is(err, storage.ErrOperationInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrOperationInUse)
		case errors.Is(err, storage.ErrOperationNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrOperationNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (o *OperationService) Update(ctx context.Context, id int64, title *string) error {
	const op = opStart + ".Update"

	if err := o.oProvider.UpdateOperation(ctx, id, title); err != nil {
		if errors.Is(err, storage.ErrOperationNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrOperationNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
