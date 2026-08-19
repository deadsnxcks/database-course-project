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

func (s *Storage) Cargos(
	ctx context.Context,
) ([]models.Cargo, error) {
	const op = "storage.postgresql.Cargos"

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, type_id, weight, volume, vessel_id
		FROM cargo
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cargos []models.Cargo
	for rows.Next() {
		var c models.Cargo
		if err := rows.Scan(&c.ID,
			&c.Title,
			&c.TypeID,
			&c.Weight,
			&c.Volume,
			&c.VesselID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cargos = append(cargos, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargos, nil
}

func (s *Storage) SaveCargo(
	ctx context.Context,
	cargo models.Cargo,
) (int64, error) {
	const op = "storage.postgresql.SaveCargo"

	var id int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO cargo (title, type_id, weight, volume, vessel_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, cargo.Title, cargo.TypeID, cargo.Weight, cargo.Volume, cargo.VesselID).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteCargo(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.DeleteCargo"

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM cargo	
		WHERE id = $1
	`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrCargoInUse)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
	}

	return nil
}

func (s *Storage) Cargo(
	ctx context.Context,
	id int64,
) (models.Cargo, error) {
	const op = "storage.postgresql.CargoByID"

	var c models.Cargo

	err := s.pool.QueryRow(ctx, `
		SELECT id, title, type_id, weight, volume, vessel_id
		FROM cargo
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Title, &c.TypeID, &c.Weight, &c.Volume, &c.VesselID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Cargo{}, fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
		}
		return models.Cargo{}, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (s *Storage) UpdateCargo(
	ctx context.Context,
	id int64,
	title *string,
	typeID *int64,
	weight *float64,
	volume *float64,
	vesselID *int64,
) error {
	const op = "storage.postgresql.UpdateCargo"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE cargo
		SET title = COALESCE($1, title),
			type_id = COALESCE($2, type_id),
			weight = COALESCE($3, weight),
			volume = COALESCE($4, volume),
			vessel_id = COALESCE($5, vessel_id)
		WHERE id = $6
	`, title, typeID, weight, volume, vesselID, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrCargoExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
	}

	return nil
}
