package models

import (
	"time"

	"github.com/google/uuid"
)

// RoutePlan represents a route plan entity
type RoutePlan struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	LocationIDs []uuid.UUID `json:"location_ids" db:"location_ids"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// RoutePlanCreate represents request payload for creating a route plan
type RoutePlanCreate struct {
	Name        string     `json:"name" binding:"required,min=1,max=255"`
	LocationIDs []uuid.UUID `json:"location_ids" binding:"required,min=2"`
}

// RoutePlanUpdate represents request payload for updating a route plan
type RoutePlanUpdate struct {
	Name        *string     `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	LocationIDs *[]uuid.UUID `json:"location_ids,omitempty" binding:"omitempty,min=2"`
}
