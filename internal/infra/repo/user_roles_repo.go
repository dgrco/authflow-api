package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrco/quikslate/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *PgRepository) AssignRole(ctx context.Context, userID, businessID string, locationID *string, role domain.Role) error {
	query := `
		INSERT INTO user_roles (user_id, business_id, location_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, business_id) DO UPDATE SET role = EXCLUDED.role, updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query, userID, businessID, locationID, role)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

func (r *PgRepository) GetUserRole(ctx context.Context, userID, businessID string) (domain.UserRole, error) {
	query := `
		SELECT id, user_id, business_id, location_id, role, created_at, updated_at
		FROM user_roles
		WHERE user_id = $1 AND business_id = $2
	`

	var urole domain.UserRole
	err := r.pool.QueryRow(ctx, query, userID, businessID).Scan(
		&urole.ID,
		&urole.UserID,
		&urole.BusinessID,
		&urole.LocationID,
		&urole.Role,
		&urole.CreatedAt,
		&urole.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.UserRole{}, domain.ErrNotFound
		default:
			return domain.UserRole{}, fmt.Errorf("failed to get user role: %w", err)
		}
	}

	return urole, nil
}

func (r *PgRepository) RemoveRole(ctx context.Context, userID, businessID string) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND business_id = $2
	`

	cmdTag, err := r.pool.Exec(ctx, query, userID, businessID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

var _ domain.UserRoleRepository = (*PgRepository)(nil)
