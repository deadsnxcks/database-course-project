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
) error {
	const op = "storage.postgresql.UpdateOperation"

	cmdTag, err := s.pool.Exec(ctx, `
		UPDATE operation
		SET title = COALESCE($1, title)
		WHERE id = $2
	`, title, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrOperationNotFound)
	}

	return nil
}
