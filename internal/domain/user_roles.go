package domain

import (
	"context"
	"time"
)

type Role string

const (
	Admin    Role = "admin"
	Manager  Role = "manager"
	Employee Role = "employee"
)

type UserRole struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BusinessID string    `json:"business_id"`
	LocationID *string   `json:"location_id"`
	Role       Role      `json:"user_role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, businessID string, locationID *string, role Role) error
	GetUserRole(ctx context.Context, userID, businessID string) (UserRole, error)
	RemoveRole(ctx context.Context, userID, businessID string) error
}
