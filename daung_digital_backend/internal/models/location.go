package models

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a location entity
type Location struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Address   string     `json:"address" db:"address"`
	Phone     *string    `json:"phone,omitempty" db:"phone"`
	Latitude  float64    `json:"latitude" db:"latitude"`
	Longitude float64    `json:"longitude" db:"longitude"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// LocationCreate represents request payload for creating a location
type LocationCreate struct {
	Name      string  `json:"name" binding:"required,min=1,max=255"`
	Address   string  `json:"address" binding:"required,min=1,max=1000"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// LocationUpdate represents request payload for updating a location
type LocationUpdate struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Address   *string `json:"address,omitempty" binding:"omitempty,min=1,max=1000"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Latitude  *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`
}

// NearbyQuery represents query parameters for nearby locations
type NearbyQuery struct {
	Latitude  float64 `form:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `form:"longitude" binding:"required,min=-180,max=180"`
	Radius    float64 `form:"radius" binding:"required,min=0.1,max=100"`
	Limit     int     `form:"limit" binding:"omitempty,min=1,max=100"`
}

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Cursor *string `form:"cursor"`
	Limit  int     `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort   string  `form:"sort" binding:"omitempty,oneof=created_at updated_at name"`
	Order  string  `form:"order" binding:"omitempty,oneof=asc desc"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       []Location `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}

// NearbyResponse represents response for nearby locations
type NearbyResponse struct {
	Data     []Location `json:"data"`
	Metadata NearbyMetadata `json:"metadata"`
}

// NearbyMetadata represents metadata for nearby query
type NearbyMetadata struct {
	Center  Center `json:"center"`
	Radius  float64 `json:"radius"`
	Count   int    `json:"count"`
}

// Center represents center coordinates
type Center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
