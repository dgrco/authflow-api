package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrco/quikslate/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *PgRepository) CreateShift(
	ctx context.Context,
	userID *string,
	locationID, positionID string,
	status domain.ShiftStatus,
	startTime, endTime time.Time,
) (domain.Shift, error) {
	query := `
		INSERT INTO shifts (user_id, location_id, position_id, status, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, location_id, position_id, status, start_time, end_time, created_at, updated_at
	`

	var s domain.Shift
	err := r.pool.QueryRow(ctx, query, userID, locationID, positionID, status, startTime, endTime).Scan(
		&s.ID,
		&s.UserID,
		&s.LocationID,
		&s.PositionID,
		&s.Status,
		&s.StartTime,
		&s.EndTime,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return domain.Shift{}, fmt.Errorf("failed to create shift: %w", err)
	}

	return s, nil
}

func (r *PgRepository) GetShiftByID(ctx context.Context, id string) (domain.Shift, error) {
	query := `
		SELECT id, user_id, location_id, position_id, status, start_time, end_time, created_at, updated_at
		FROM shifts
		WHERE id = $1
	`

	var s domain.Shift
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.LocationID,
		&s.PositionID,
		&s.Status,
		&s.StartTime,
		&s.EndTime,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Shift{}, domain.ErrNotFound
		default:
			return domain.Shift{}, fmt.Errorf("failed to get shift by ID: %w", err)
		}
	}

	return s, nil
}

func (r *PgRepository) GetShiftsByLocationID(ctx context.Context, locationID string) ([]domain.Shift, error) {
	query := `
		SELECT id, user_id, location_id, position_id, status, start_time, end_time, created_at, updated_at
		FROM shifts
		WHERE location_id = $1
	`

	var shifts []domain.Shift
	rows, err := r.pool.Query(ctx, query, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts by location ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Shift
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.LocationID,
			&s.PositionID,
			&s.Status,
			&s.StartTime,
			&s.EndTime,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate shifts: %w", err)
	}

	return shifts, nil
}

func (r *PgRepository) UpdateShiftByID(ctx context.Context, id string, update domain.ShiftUpdate) error {
	builder := newUpdateBuilder()
	if update.Status != nil {
		builder.Add("status", *update.Status)
	}
	if update.StartTime != nil {
		builder.Add("start_time", *update.StartTime)
	}
	if update.EndTime != nil {
		builder.Add("end_time", *update.EndTime)
	}
	if builder.IsEmpty() {
		return nil
	}

	query, args := builder.Build("shifts", "id", id)
	cmdTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update shift by ID: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PgRepository) UnassignShift(ctx context.Context, id string) error {
	query := `
		UPDATE shifts
		SET user_id = NULL, status = 'uncovered', updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to unassign shift: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Soft-delete a shift (use this over DeleteShift most of the time)
func (r *PgRepository) CancelShift(ctx context.Context, id string) error {
	query := `
		UPDATE shifts
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel shift: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Hard-delete a shift (should only be used for admin purposes)
func (r *PgRepository) DeleteShift(ctx context.Context, id string) error {
	query := `
		DELETE FROM shifts
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete shift: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

var _ domain.ShiftRepository = (*PgRepository)(nil)
