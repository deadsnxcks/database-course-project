package postgresql

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	"errors"
	"fmt"
	"time"

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

	if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("%s: ping failed: %w", op, err)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeInUse)
		}
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
	err := s.pool.QueryRow(ctx, `
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
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoTypeNotFound)
	}

	return nil
}

func (s *Storage) Operations(
	ctx context.Context,
) ([]models.Operation, error) {
	const op = "storage.postgresql.Operation"

	var operations []models.Operation

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, created_at
		FROM operation
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Operation
		if err := rows.Scan(&o.ID, &o.Title, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		operations = append(operations, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return operations, nil
}

func (s *Storage) SaveOperation(
	ctx context.Context,
	operation models.Operation,
) (int64, error) {
	const op = "storage.postgresql.SaveOperation"

	var id int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO operation (title)
		VALUES ($1)
		RETURNING id
	`, operation.Title).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrOperationExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteOperation(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.postgresql.DeleteOperation"

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM operation
		WHERE id = $1
	`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrOperationInUse)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrOperationNotFound)
	}

	return nil
}

func (s *Storage) Operation(
	ctx context.Context,
	id int64,
) (models.Operation, error) {
	const op = "storage.postgresql.Operation"

	var o models.Operation
	err := s.pool.QueryRow(ctx, `
		SELECT id, title, created_at
		FROM operation
		WHERE id = $1
	`, id).Scan(&o.ID, &o.Title, &o.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Operation{}, fmt.Errorf("%s: %w", op, storage.ErrOperationNotFound)
		}

		return models.Operation{}, fmt.Errorf("%s: %w", op, err)
	}

	return o, nil
}

func (s *Storage) UpdateOperation(
	ctx context.Context,
	id int64,
	title *string,
	createdAt *time.Time,
) error {
	const op = "storage.postgresql.UpdateOperation"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE operation
		SET title = COALESCE($1, title),
			created_at = COALESCE($2, created_at)
		WHERE id = $3
	`, title, createdAt, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrOperationNotFound)
	}

	return nil
}

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

func (s *Storage) CargoByID(
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

func (s *Storage) CargosByVesselID(
	ctx context.Context,
	vesselID int64,
) ([]models.Cargo, error) {
	const op = "storage.postgresql.Cargo"

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, type_id, weight, volume, vessel_id
		FROM cargo
		WHERE vessel_id = $1
	`, vesselID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cargosByVeselID []models.Cargo
	for rows.Next() {
		var c models.Cargo
		if err := rows.Scan(&c.ID, 
			&c.Title, 
			&c.TypeID, 
			&c.Weight, 
			&c.Volume, 
			&c.VesselID,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cargosByVeselID = append(cargosByVeselID, c)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargosByVeselID, nil
}

func (s *Storage) CargosByTypeID(
	ctx context.Context,
	typeID int64,
) ([]models.Cargo, error) {
	const op = "storage.postgresql.Cargo"

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, type_id, weight, volume, vessel_id
		FROM cargo
		WHERE type_id = $1
	`, typeID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w",op, err)
	}
	defer rows.Close()

	var cargosByTypeID []models.Cargo
	for rows.Next() {
		var c models.Cargo
		if err := rows.Scan(&c.ID, 
			&c.Title, 
			&c.TypeID, 
			&c.Weight, 
			&c.Volume, 
			&c.VesselID,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cargosByTypeID = append(cargosByTypeID, c)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargosByTypeID, nil
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
	`, title, typeID, weight, volume, vesselID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
	}

	return nil
}

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
	sl models.StorageLocation,
) (int64, error) {
	const op = "storage.postgresql.SaveStorageLoc"

	var id int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO storage_loc (cargo_type_id, max_weight, max_volume, cargo_id, date_of_placement)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`, sl.CargoTypeID, sl.MaxWeight, sl.MaxVolume, sl.CargoID, sl.DateOfPlacement).Scan(&id)
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

	cmdTag, err := s.pool.Exec(ctx, `
		DELETE FROM storage_loc
		WHERE id = $1
	`)
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

	_, err := s.pool.Exec(ctx, `
		UPDATE storage_loc sl
		SET cargo_id = c.id
			date_of_placement = $3
		FROM cargo c
		WHERE 
			sl.id = $1
			AND c.id = $2
			AND sl.cargo_id IS NULL
			AND sl.cargo_type_id = c.type_id
			AND sl.max_weight <= c.weight
			AND sl.max_volume <= c.volume
	`, storageLocID, cargoID, date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var exists bool
			checkErr := s.pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM storage_loc WHERE id = $1
				)
			`, storageLocID).Scan(&exists)

			if checkErr != nil {
				return fmt.Errorf("%s: %w", op, checkErr)
			}

			if !exists {
				return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
			}

			checkErr = s.pool.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM cargo WHERE id = $1
				)
			`, cargoID).Scan(&exists)

			if checkErr != nil {
				return fmt.Errorf("%s: %w", op, checkErr)
			}

			if !exists {
				return fmt.Errorf("%s: %w", op, storage.ErrCargoNotFound)
			}

			return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotSuitable)
		}

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
		fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() > 0 {
		return nil
	}

	var exists bool
	err = s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM storage_loc WHERE id = $1
		)
	`, id).Scan(&exists)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return fmt.Errorf("%s: %w", op, storage.ErrStorageLocNotFound)
	}

	return fmt.Errorf("%s: %w", op, storage.ErrStorageLocAlreadyEmpty)
}

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