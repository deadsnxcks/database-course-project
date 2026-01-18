package opercargoservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"dbcp/internal/storage"
	"fmt"
	"log/slog"
)

const (
	opStart = "services.operationcargo"
)

type OperationCargoService struct {
	log *slog.Logger
	ocProvider OperationCargoProvider
}

type OperationCargoProvider interface {
	OperationsCargos(ctx context.Context) ([]models.OperationCargo, error)
	SaveOperationCargo(ctx context.Context, operCargo models.OperationCargo) error
	DeleteOperationCargo(ctx context.Context, operCargo models.OperationCargo) error
}

func New(
	log *slog.Logger,
	ocProvider OperationCargoProvider,
) *OperationCargoService {
	return &OperationCargoService{
		log: log,
		ocProvider: ocProvider,
	}
}

func (s *OperationCargoService) List(
	ctx context.Context,
) ([]models.OperationCargo, error) {
	const op = opStart + ".List"

	log := s.log.With(slog.String("op", op))
	log.Info("Listing operation cargos")

	opsCargos, err := s.ocProvider.OperationsCargos(ctx)
	if err != nil {
		log.Error("failed to list operation cargos", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return opsCargos, nil
}

func (s *OperationCargoService) Create(
	ctx context.Context, 
	operID, cargoID int64,
) error {
	const op = opStart + ".Create"

	log := s.log.With(
		slog.String("op", op), 
		slog.Int64("operation_id", operID), 
		slog.Int64("cargo_id", cargoID),
	)

	if operID <= 0 || cargoID <= 0 {
		return fmt.Errorf("%s: invalid operation_id or cargo_id", op)
	}

	operCargo := models.OperationCargo{
		OperationID: operID,
		CargoID:     cargoID,
	}

	if err := s.ocProvider.SaveOperationCargo(ctx, operCargo); err != nil {
		log.Error("failed to create operation cargo", sl.Err(err))
		if err == storage.ErrOperCargoAlreadyExist {
			return fmt.Errorf("%s: %w", op, storage.ErrOperCargoAlreadyExist)
		}
		if err == storage.ErrRelatedEntityNotFound {
			return fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("OperationCargo created")
	return nil
}

func (s *OperationCargoService) Delete(
	ctx context.Context, 
	operID, 
	cargoID int64,
) error {
	const op = opStart + ".Delete"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("operation_id", operID), 
		slog.Int64("cargo_id", cargoID),
	)

	if operID <= 0 || cargoID <= 0 {
		return fmt.Errorf("%s: invalid operation_id or cargo_id", op)
	}

	operCargo := models.OperationCargo{
		OperationID: operID,
		CargoID:     cargoID,
	}

	if err := s.ocProvider.DeleteOperationCargo(ctx, operCargo); err != nil {
		log.Error("failed to delete operation cargo", sl.Err(err))
		if err == storage.ErrOperCargoNotFound {
			return fmt.Errorf("%s: %w", op, storage.ErrOperCargoNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("OperationCargo deleted")
	return nil
}