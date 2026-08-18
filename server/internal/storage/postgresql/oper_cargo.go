package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) OperationsCargos(
	ctx context.Context,
) ([]models.OperationCargo, error) {
	const op = "storage.postgresql.OperationsCargos"

	rows, err := s.pool.Query(ctx, `
		SELECT operation_id, cargo_id
		FROM operation_cargo
		ORDER BY operation_id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var operCargos []models.OperationCargo
	for rows.Next() {
		var oc models.OperationCargo
		if err := rows.Scan(&oc.OperationID, &oc.CargoID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		operCargos = append(operCargos, oc)
	}

	return operCargos, nil
}

func (s *Storage) SaveOperationCargo(
	ctx context.Context,
	operCargo models.OperationCargo,
) error {
	const op = "storage.postgreql.SaveOperationCargo"

	_, err := s.pool.Exec(ctx, `
		INSERT INTO operation_cargo 
			(operation_id, cargo_id)
		VALUES ($1, $2)
	`, operCargo.OperationID, operCargo.CargoID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return fmt.Errorf("%s: %w", op, storage.ErrOperCargoAlreadyExist)
			case "23503":
				return fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteOperationCargo(
	ctx context.Context,
	operCargo models.OperationCargo,
) error {
	const op = "storage.postgresql.DeleteOperCargo"

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM operation_cargo oc
		WHERE oc.operation_id = $1 AND oc.cargo_id = $2
	`, operCargo.OperationID, operCargo.CargoID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrOperCargoNotFound)
	}

	return nil
}
