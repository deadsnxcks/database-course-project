package operationservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

const opStart = "services.operation"

type OperationService struct {
	log *slog.Logger
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
		log: log,
		oProvider: oProvider,
	}
}

func (o *OperationService) List(ctx context.Context) ([]models.Operation, error) {
	const op = opStart + ".List"

	log := o.log.With(slog.String("op", op))
	log.Info("Listing operations")

	ops, err := o.oProvider.Operations(ctx)
	if err != nil {
		log.Error("failed to list operations", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ops, nil
}

func (o *OperationService) Get(ctx context.Context, id int64) (models.Operation, error) {
	const op = opStart + ".Get"

	log := o.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return models.Operation{}, fmt.Errorf("%s: invalid id", op)
	}

	opModel, err := o.oProvider.Operation(ctx, id)
	if err != nil {
		log.Error("failed to get operation", sl.Err(err))
		return models.Operation{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Operation retrieved", slog.Int64("id", id))
	return opModel, nil
}

func (o *OperationService) Create(ctx context.Context, title string) (int64, error) {
	const op = opStart + ".Create"

	log := o.log.With(slog.String("op", op), slog.String("title", title))

	if title == "" {
		return 0, fmt.Errorf("%s: title is required", op)
	}

	id, err := o.oProvider.SaveOperation(ctx, models.Operation{Title: title})
	if err != nil {
		log.Error("failed to create operation", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Operation created", slog.Int64("id", id))
	return id, nil
}

func (o *OperationService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	log := o.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := o.oProvider.DeleteOperation(ctx, id); err != nil {
		log.Error("failed to delete operation", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Operation deleted", slog.Int64("id", id))
	return nil
}

func (o *OperationService) Update(ctx context.Context, id int64, title *string) error {
	const op = opStart + ".Update"

	log := o.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := o.oProvider.UpdateOperation(ctx, id, title); err != nil {
		log.Error("failed to update operation", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Operation updated", slog.Int64("id", id))
	return nil
}