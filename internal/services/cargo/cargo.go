package cargoservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

const (
	opStart = "services.cargo"
)

type CargoService struct {
	log *slog.Logger
	cProvider CargoProvider
}

type CargoProvider interface {
	Cargos(ctx context.Context) ([]models.Cargo, error)
	SaveCargo(ctx context.Context, cargo models.Cargo) (int64, error)
	DeleteCargo(ctx context.Context, id int64) error
	Cargo(ctx context.Context, id int64) (models.Cargo, error)
	UpdateCargo(
		ctx context.Context,
		id int64,
		title *string,
		typeID *int64,
		weight *float64,
		volume *float64,
		vesselID *int64,
	) error
}

func New(
	log *slog.Logger,
	cProvider CargoProvider,
) *CargoService {
	return &CargoService{
		log: log,
		cProvider: cProvider,
	}
}

func (c *CargoService) List(
	ctx context.Context,
) ([]models.Cargo, error) {
	const op = opStart + ".List"

	log := c.log.With(slog.String("op", op))
	log.Info("Listing cargos")

	cargos, err := c.cProvider.Cargos(ctx)
	if err != nil {
		log.Error("failed to list cargos", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargos, nil
}

func (c *CargoService) Get(
	ctx context.Context, 
	id int64,
) (models.Cargo, error) {
	const op = opStart + ".Get"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return models.Cargo{}, fmt.Errorf("%s: invalid id", op)
	}

	cargo, err := c.cProvider.Cargo(ctx, id)
	if err != nil {
		log.Error("failed to get cargo", sl.Err(err))
		return models.Cargo{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo received", slog.Int64("id", id))
	return cargo, nil
}

func (c *CargoService) Create(
	ctx context.Context, 
	cargo models.Cargo,
) (int64, error) {
	const op = opStart + ".Create"

	log := c.log.With(
		slog.String("op", op),
		slog.String("title", cargo.Title),
	)

	if cargo.Title == "" {
		return 0, fmt.Errorf("%s: title is required", op)
	}
	if cargo.TypeID <= 0 {
		return 0, fmt.Errorf("%s: cargoTypeID is required", op)
	}
	if cargo.Weight <= 0 {
		return 0, fmt.Errorf("%s: weight must be positive", op)
	}

	id, err := c.cProvider.SaveCargo(ctx, cargo)
	if err != nil {
		log.Error("failed to create cargo", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo created", slog.Int64("id", id))
	return id, nil
}

func (c *CargoService) Delete(
	ctx context.Context, 
	id int64,
) error {
	const op = opStart + ".Delete"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := c.cProvider.DeleteCargo(ctx, id); err != nil {
		log.Error("failed to delete cargo", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo deleted")
	return nil
}

func (c *CargoService) Update(
	ctx context.Context,
	id int64,
	title *string,
	typeID *int64,
	weight *float64,
	volume *float64,
	vesselID *int64,
) error {
	const op = opStart + ".Update"

	log := c.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := c.cProvider.UpdateCargo(
		ctx,
		id,
		title,
		typeID,
		weight,
		volume,
		vesselID,
	); err != nil {
		log.Error("failed to update cargo", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Cargo updated")
	return nil
}
