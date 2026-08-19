package cargotypeservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/deadsnxcks/dbcp/server/internal/domain"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
)

const (
	opStart = "services.cargo-type"
)

type CargoTypeService struct {
	log        *slog.Logger
	ctProvider CargoTypeProvider
}

type CargoTypeProvider interface {
	CargoTypes(ctx context.Context) ([]models.CargoType, error)
	SaveCargoType(ctx context.Context, cargoType models.CargoType) (int64, error)
	DeleteCargoType(ctx context.Context, id int64) error
	CargoType(ctx context.Context, id int64) (models.CargoType, error)
	UpdateCargoType(
		ctx context.Context,
		id int64,
		title *string,
		processCost *float64,
	) error
}

func New(
	log *slog.Logger,
	ctProvider CargoTypeProvider,
) *CargoTypeService {
	return &CargoTypeService{
		log:        log,
		ctProvider: ctProvider,
	}
}

func (c *CargoTypeService) List(ctx context.Context) ([]models.CargoType, error) {
	const op = opStart + ".List"

	types, err := c.ctProvider.CargoTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return types, nil
}

func (c *CargoTypeService) Get(ctx context.Context, id int64) (models.CargoType, error) {
	const op = opStart + ".Get"

	ct, err := c.ctProvider.CargoType(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrCargoTypeNotFound) {
			return models.CargoType{}, fmt.Errorf("%s: %w", op, domain.ErrCargoTypeNotFound)
		}
		return models.CargoType{}, fmt.Errorf("%s: %w", op, err)
	}

	return ct, nil
}

func (c *CargoTypeService) Create(ctx context.Context, cargoType models.CargoType) (int64, error) {
	const op = opStart + ".Create"

	id, err := c.ctProvider.SaveCargoType(ctx, cargoType)
	if err != nil {
		if errors.Is(err, storage.ErrCargoTypeExists) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrCargoTypeExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("Cargo type created", slog.Int64("id", id))

	return id, nil
}

func (c *CargoTypeService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	err := c.ctProvider.DeleteCargoType(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoTypeInUse)
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoTypeNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (c *CargoTypeService) Update(
	ctx context.Context,
	id int64,
	title *string,
	processCost *float64,
) error {
	const op = opStart + ".Update"

	err := c.ctProvider.UpdateCargoType(ctx, id, title, processCost)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeExists):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoTypeExists)
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoTypeNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
