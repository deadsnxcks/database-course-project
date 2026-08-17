package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CargoTypes(
	ctx context.Context,
) ([]models.CargoType, error) {
	const op = "storage.postgresql.CargoTypes"

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, process_cost
		FROM cargo_type
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s query: %w", op, err)
	}
	defer rows.Close()

	var cargoTypes []models.CargoType
	for rows.Next() {
		var ct models.CargoType
		if err := rows.Scan(&ct.ID, &ct.Title, &ct.ProcessCost); err != nil {
			return nil, fmt.Errorf("%s rows: %w", op, err)
		}
		cargoTypes = append(cargoTypes, ct)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s rows: %w", op, err)
	}

	return cargoTypes, nil
}

func (s *Storage) SaveCargoType(
	ctx context.Context,
	cargoType models.CargoType,
) (int64, error) {
	const op = "storage.postgresql.SaveCargoType"

	var id int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO cargo_type (title, process_cost)
		VALUES($1, $2)
		RETURNING id
	`, cargoType.Title, cargoType.ProcessCost).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrCargoTypeExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteCargoType(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.DeleteCargoType"

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM cargo_type
		WHERE id = $1
	`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeInUse)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeNotFound)
	}

	return nil
}

func (s *Storage) CargoType(
	ctx context.Context,
	id int64,
) (models.CargoType, error) {
	const op = "storage.postgresql.CargoType"

	var ct models.CargoType
	err := s.pool.QueryRow(ctx, `
		SELECT id, title, process_cost
		FROM cargo_type
		WHERE id = $1
	`, id).Scan(&ct.ID, &ct.Title, &ct.ProcessCost)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.CargoType{}, fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
		}

		return models.CargoType{}, fmt.Errorf("%s: %w", op, err)
	}

	return ct, nil
}

func (s *Storage) UpdateCargoType(
	ctx context.Context,
	id int64,
	title *string,
	processCost *float64,
) error {
	const op = "storage.postgresql.UpdateCargoType"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE cargo_type
		SET
			title = COALESCE($1, title),
			process_cost = COALESCE($2, process_cost)
		WHERE id = $3
	`, title, processCost, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeNotFound)
	}

	return nil
}