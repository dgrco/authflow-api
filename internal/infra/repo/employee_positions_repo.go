package repo

import (
	"context"
	"fmt"

	"github.com/dgrco/quikslate/internal/domain"
)

func (r *PgRepository) AddPosition(ctx context.Context, userID, positionID string) error {
	query := `
		INSERT INTO employee_positions (user_id, position_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, position_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID, positionID)
	if err != nil {
		return fmt.Errorf("failed to add employee position: %w", err)
	}
	return nil
}

func (r *PgRepository) RemovePosition(ctx context.Context, userID, positionID string) error {
	query := `
		DELETE FROM employee_positions
		WHERE user_id = $1 AND position_id = $2
	`

	cmdTag, err := r.pool.Exec(ctx, query, userID, positionID)
	if err != nil {
		return fmt.Errorf("failed to remove employee position: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PgRepository) GetPositionsByUserID(ctx context.Context, userID string) ([]domain.EmployeePosition, error) {
	query := `
		SELECT user_id, position_id
		FROM employee_positions
		WHERE user_id = $1
	`

	var positions []domain.EmployeePosition
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee positions by user ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.EmployeePosition
		err := rows.Scan(&p.UserID, &p.PositionID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee position: %w", err)
		}
		positions = append(positions, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate employee positions: %w", err)
	}

	return positions, nil
}

var _ domain.EmployeePositionRepository = (*PgRepository)(nil)
