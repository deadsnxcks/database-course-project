package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) StorageLocations(
	ctx context.Context,
) ([]models.StorageLocation, error) {
	const op = "storage.postgresql.StorageLocations"

	rows, err := s.pool.Query(ctx, `
		SELECT id, cargo_type_id, max_weight,
			max_volume, cargo_id, date_of_placement
		FROM storage_loc
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var storageLocations []models.StorageLocation
	for rows.Next() {
		var sl models.StorageLocation
		if err := rows.Scan(&sl.ID,
			&sl.CargoTypeID,
			&sl.MaxWeight,
			&sl.MaxVolume,
			&sl.CargoID,
			&sl.DateOfPlacement,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		storageLocations = append(storageLocations, sl)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storageLocations, nil
}

func (s *Storage) SaveStorageLoc(
	ctx context.Context,
	cargoTypeID int64,
	maxWeight float64,
	maxVolume float64,
) (int64, error) {
	const op = "storage.postgresql.SaveStorageLoc"

	var id int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO storage_loc (cargo_type_id, max_weight, max_volume)
		VALUES($1, $2, $3)
		RETURNING id
	`, cargoTypeID, maxWeight, maxVolume).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteStorageLoc(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.DeleteStorageLoc"

	var inUse bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM storage_loc WHERE id = $1 AND cargo_id IS NOT NULL
		)
	`, id).Scan(&inUse)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if inUse {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocInUse)
	}

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM storage_loc
		WHERE id = $1 AND cargo_id IS NULL AND date_of_placement IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
	}

	return nil
}

func (s *Storage) StorageLocation(
	ctx context.Context,
	id int64,
) (models.StorageLocation, error) {
	const op = "storage.postgresql.StorageLocation"

	var sl models.StorageLocation
	err := s.pool.QueryRow(ctx, `
		SELECT id, cargo_type_id, max_weight, 
			max_volume, cargo_id, date_of_placement
		FROM storage_loc
		WHERE id = $1
	`, id).Scan(&sl.ID,
		&sl.CargoTypeID,
		&sl.MaxWeight,
		&sl.CargoID,
		&sl.DateOfPlacement,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.StorageLocation{}, fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
		}

		return models.StorageLocation{}, fmt.Errorf("%s: %w", op, err)
	}

	return sl, nil
}

func (s *Storage) UpdateStorageLoc(
	ctx context.Context,
	id int64,
	cargoTypeID *int64,
	maxWeight *float64,
	maxVolume *float64,
) error {
	const op = "storage.postgresql.UpdateStorageLoc"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE storage_loc
		SET cargo_type_id = COALESCE($1, cargo_type_id),
			max_weight = COALESCE($2, max_weight),
			max_volume = COALESCE($3, max_volume)
		WHERE id = $4
	`, cargoTypeID, maxWeight, maxVolume, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrRelatedEntityNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
	}

	return nil
}

func (s *Storage) UseStorageLoc(
	ctx context.Context,
	storageLocID int64,
	cargoID int64,
	date time.Time,
) error {
	const op = "storage.postgresql.UseStorageLoc"

	var isExist bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM storage_loc WHERE id = $1
		)
	`, storageLocID).Scan(&isExist)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !isExist {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
	}

	var isOccupied bool
	err = s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM storage_loc WHERE id = $1 AND cargo_id IS NOT NULL
		)
	`, storageLocID).Scan(&isOccupied)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if isOccupied {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocInUse)
	}

	var cargoExists bool
	err = s.pool.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM cargo WHERE id = $1
        )
    `, cargoID).Scan(&cargoExists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !cargoExists {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
	}

	var isCargoPlaced bool
	err = s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM storage_loc WHERE cargo_id = $1
		)
	`, cargoID).Scan(&isCargoPlaced)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if isCargoPlaced {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoAlreadyPlaced)
	}

	var isCargoType bool
	err = s.pool.QueryRow(ctx, `
		SELECT sl.cargo_type_id = c.type_id
		FROM storage_loc sl
		JOIN cargo c ON c.id = $2
		WHERE sl.id = $1
	`, storageLocID, cargoID).Scan(&isCargoType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !isCargoType {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocTypeNotSuitable)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx, `
		UPDATE storage_loc sl
		SET cargo_id = c.id,
			date_of_placement = $3
		FROM cargo c
		WHERE 
			sl.id = $1
			AND c.id = $2
			AND sl.cargo_id IS NULL
			AND sl.cargo_type_id = c.type_id
			AND sl.max_weight >= c.weight
			AND sl.max_volume >= c.volume
	`, storageLocID, cargoID, date)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotSuitable)
	}

	var operationID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO operation (title, created_at)
		VALUES ('Размещение на складе', $1)
		RETURNING id
	`, date).Scan(&operationID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO operation_cargo (operation_id, cargo_id)
		VALUES ($1, $2)
	`, operationID, cargoID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ResetStorageLoc(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.ResetStorageLoc"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE storage_loc
		SET cargo_id = NULL,
			date_of_placement = NULL
		WHERE id = $1 AND cargo_id IS NOT NULL
	`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		var exists bool
		err := s.pool.QueryRow(ctx, `
            SELECT EXISTS(SELECT 1 FROM storage_loc WHERE id = $1)
        `, id).Scan(&exists)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if !exists {
			return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
		}

		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocAlreadyEmpty)
	}

	return nil
}
