package postgresql

import (
	"context"
	"fmt"

	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
)

func (s *Storage) CargoDetailReport(
	ctx context.Context,
) ([]models.CargoDetailItem, error) {
	const op = "storage.postgresql.CargoDetailReport"

	rows, err := s.pool.Query(ctx, `
        SELECT 
            c.title AS cargo_name,
            c.weight AS weight,
            ct.title AS cargo_type,
            v.title AS vessel_type,
            o.created_at AS unloading_date
        FROM cargo c
        JOIN cargo_type ct ON c.type_id = ct.id
        JOIN vessel v ON c.vessel_id = v.id
        JOIN operation_cargo oc ON c.id = oc.cargo_id
        JOIN operation o ON oc.operation_id = o.id
        WHERE o.title = 'Выгрузка'
        ORDER BY o.created_at DESC
    `)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []models.CargoDetailItem
	for rows.Next() {
		var item models.CargoDetailItem

		err := rows.Scan(
			&item.CargoName,
			&item.Weight,
			&item.CargoType,
			&item.VesselName,
			&item.UnloadingDate,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}

func (s *Storage) CargoTypeReport(
	ctx context.Context,
) ([]models.CargoTypeItem, error) {
	const op = "storage.postgresql.CargoTypeReport"

	rows, err := s.pool.Query(ctx, `
		SELECT 
			ct.title AS cargo_type_name,
			COUNT(c.id) AS cargo_count,
			SUM(c.weight) AS total_weight_tons,
			SUM(c.weight) * ct.process_cost AS total_process_cost
		FROM cargo c
		JOIN cargo_type ct ON c.type_id = ct.id
		GROUP BY ct.id, ct.title, ct.process_cost
		ORDER BY total_weight_tons DESC;
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []models.CargoTypeItem
	for rows.Next() {
		var item models.CargoTypeItem

		err := rows.Scan(
			&item.CargoTypeName,
			&item.CargoCount,
			&item.TotalWeight,
			&item.TotalProcessCost,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}
