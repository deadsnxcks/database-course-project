package cargotypeservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

const (
	opStart = "services.cargo-type"
)

type CargoTypeService struct {
	log *slog.Logger
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
		log: log,
		ctProvider: ctProvider,
	}
}

func (c *CargoTypeService) List(ctx context.Context) ([]models.CargoType, error) {
	const op = opStart + ".List"

	log := c.log.With(slog.String("op", op))
	log.Info("Listing cargo types")

	types, err := c.ctProvider.CargoTypes(ctx)
	if err != nil {
		log.Error("failed to list cargo types", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return types, nil
}

func (c *CargoTypeService) Get(ctx context.Context, id int64) (models.CargoType, error) {
	const op = opStart + ".Get"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return models.CargoType{}, fmt.Errorf("%s: invalid id", op)
	}

	ct, err := c.ctProvider.CargoType(ctx, id)
	if err != nil {
		log.Error("failed to get cargo type", sl.Err(err))
		return models.CargoType{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo type retrieved", slog.Int64("id", id))
	return ct, nil
}

func (c *CargoTypeService) Create(ctx context.Context, cargoType models.CargoType) (int64, error) {
	const op = opStart + ".Create"

	log := c.log.With(slog.String("op", op), slog.String("title", cargoType.Title))

	if cargoType.Title == "" {
		return 0, fmt.Errorf("%s: title is required", op)
	}
	if cargoType.ProcessCost <= 0 {
		return 0, fmt.Errorf("%s: processCost must be positive", op)
	}

	id, err := c.ctProvider.SaveCargoType(ctx, cargoType)
	if err != nil {
		log.Error("failed to create cargo type", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo type created", slog.Int64("id", id))
	return id, nil
}

func (c *CargoTypeService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	err := c.ctProvider.DeleteCargoType(ctx, id)
	if err != nil {
		log.Error("failed to delete cargo type", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo type deleted")
	return nil
}

func (c *CargoTypeService) Update(
	ctx context.Context,
	id int64,
	title *string,
	processCost *float64,
) error {
	const op = opStart + ".Update"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	err := c.ctProvider.UpdateCargoType(ctx, id, title, processCost)
	if err != nil {
		log.Error("failed to update cargo type", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo type updated")
	return nil
}