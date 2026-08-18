package opercargoservice

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
	opStart = "services.operationcargo"
)

type OperationCargoService struct {
	log        *slog.Logger
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
		log:        log,
		ocProvider: ocProvider,
	}
}

func (s *OperationCargoService) List(
	ctx context.Context,
) ([]models.OperationCargo, error) {
	const op = opStart + ".List"

	opsCargos, err := s.ocProvider.OperationsCargos(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return opsCargos, nil
}

func (s *OperationCargoService) Create(
	ctx context.Context,
	operID, cargoID int64,
) error {
	const op = opStart + ".Create"

	operCargo := models.OperationCargo{
		OperationID: operID,
		CargoID:     cargoID,
	}

	if err := s.ocProvider.SaveOperationCargo(ctx, operCargo); err != nil {
		switch {
		case errors.Is(err, storage.ErrOperCargoAlreadyExist):
			return fmt.Errorf("%s: %w", op, domain.ErrOperCargoAlreadyExist)
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrRelatedEntityNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	s.log.Info("operation_cargo created", slog.Int64("operation_id", operID), slog.Int64("cargo_id", cargoID))

	return nil
}

func (s *OperationCargoService) Delete(
	ctx context.Context,
	operID,
	cargoID int64,
) error {
	const op = opStart + ".Delete"

	operCargo := models.OperationCargo{
		OperationID: operID,
		CargoID:     cargoID,
	}

	if err := s.ocProvider.DeleteOperationCargo(ctx, operCargo); err != nil {
		if errors.Is(err, storage.ErrOperCargoNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrOperCargoNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
