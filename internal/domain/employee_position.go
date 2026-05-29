package domain

import "context"

type EmployeePosition struct {
	UserID     string `json:"user_id"`
	PositionID string `json:"position_id"`
}

type EmployeePositionRepository interface { 
	AddPosition(ctx context.Context, userID, positionID string) error
	RemovePosition(ctx context.Context, userID, positionID string) error
	GetPositionsByUserID(ctx context.Context, userID string) ([]EmployeePosition, error)
}
