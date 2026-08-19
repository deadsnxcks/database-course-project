package cargoservice

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
	opStart = "services.cargo"
)

type CargoService struct {
	log       *slog.Logger
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
		log:       log,
		cProvider: cProvider,
	}
}

func (c *CargoService) List(
	ctx context.Context,
) ([]models.Cargo, error) {
	const op = opStart + ".List"

	cargos, err := c.cProvider.Cargos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargos, nil
}

func (c *CargoService) Get(
	ctx context.Context,
	id int64,
) (models.Cargo, error) {
	const op = opStart + ".Get"

	cargo, err := c.cProvider.Cargo(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrCargoNotFound) {
			return models.Cargo{}, fmt.Errorf("%s: %w", op, domain.ErrCargoNotFound)
		}
		return models.Cargo{}, fmt.Errorf("%s: %w", op, err)
	}

	return cargo, nil
}

func (c *CargoService) Create(
	ctx context.Context,
	cargo models.Cargo,
) (int64, error) {
	const op = opStart + ".Create"

	id, err := c.cProvider.SaveCargo(ctx, cargo)
	if err != nil {
		if errors.Is(err, storage.ErrRelatedEntityNotFound) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrRelatedEntityNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	c.log.Info("cargo created", slog.Int64("id", id))

	return id, nil
}

func (c *CargoService) Delete(
	ctx context.Context,
	id int64,
) error {
	const op = opStart + ".Delete"

	if err := c.cProvider.DeleteCargo(ctx, id); err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoNotFound)
		case errors.Is(err, storage.ErrCargoInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoInUse)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

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

	if err := c.cProvider.UpdateCargo(
		ctx,
		id,
		title,
		typeID,
		weight,
		volume,
		vesselID,
	); err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoNotFound)
		case errors.Is(err, storage.ErrCargoExists):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoExists)
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrRelatedEntityNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
