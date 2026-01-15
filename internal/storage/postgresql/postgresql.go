package postgresql

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Storage, error) {
	const op = "storage.postgresql.New"

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) Vessels(
	ctx context.Context,
) ([]models.Vessel, error) {
	const op = "storage.postgresql.Vessels"

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, vessel_type, max_load
		FROM vessel
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s query: %w", op, err)
	}
	defer rows.Close()

	var vessels []models.Vessel
	for rows.Next() {
		var v models.Vessel
		if err := rows.Scan(&v.ID, &v.Title, &v.VesselType, &v.MaxLoad); err != nil {
			return nil, fmt.Errorf("%s rows: %w", op, err)
		}
		vessels = append(vessels, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s rows: %w", op, err)
	}

	return vessels, nil
}

func (s *Storage) SaveVessel(
	ctx context.Context, 
	vessel models.Vessel,
) (int64, error) {
	const op = "storage.postgresql.CreateVessel"

	var id int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO vessel (title, vessel_type, max_load)
		VALUES ($1, $2, $3)
		RETURNING id
	`, vessel.Title, vessel.VesselType, vessel.MaxLoad).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrVesselExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteVessel(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.DeleteVessel"

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM vessel 
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return storage.ErrVesselNotFound
	}

	return nil
}

func (s *Storage) Vessel(
	ctx context.Context,
	id int64,
) (models.Vessel, error) {
	const op = "storage.postgresql.Vessel"

	var v models.Vessel
	err := s.pool.QueryRow(ctx,`
		SELECT id, title, vessel_type, max_load
		FROM vessel
		WHERE id = $1
	`, id).Scan(&v.ID, &v.Title, &v.VesselType, &v.MaxLoad)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Vessel{}, fmt.Errorf("%s: %w", op, storage.ErrVesselNotFound)
		}

		return models.Vessel{}, fmt.Errorf("%s: %w", op, err)
	}

	return v, nil
}

func (s *Storage) UpdateVessel(
	ctx context.Context,
	id int64,
	title *string,
	vesselType *string,
	maxLoad *float64,
) error {
	const op = "storage.postgresql.UpdateVessel"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE vessel
		SET
			title = COALESCE($1, title),
			vessel_type = COALESCE($2, vessel_type),
			max_load = COALESCE($3, max_load)
		WHERE id = $4
	`, title, vesselType, maxLoad, id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrVesselExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrVesselNotFound)
	}

	return nil
}
