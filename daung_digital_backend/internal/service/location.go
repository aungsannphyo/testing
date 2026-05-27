package service

import (
	"context"
	"fmt"

	"github.com/daung-digital/location-api/internal/models"
	"github.com/daung-digital/location-api/internal/repository"
	"github.com/google/uuid"
)

// LocationService handles business logic for locations
type LocationService struct {
	repo repository.LocationRepository
}

// NewLocationService creates a new location service
func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

// CreateLocation creates a new location
func (s *LocationService) CreateLocation(ctx context.Context, req *models.LocationCreate) (*models.Location, error) {
	location := &models.Location{
		ID:        uuid.New(),
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	
	if err := s.repo.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}
	
	return location, nil
}

// GetLocation retrieves a location by ID
func (s *LocationService) GetLocation(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	return location, nil
}

// ListLocations retrieves locations with pagination
func (s *LocationService) ListLocations(ctx context.Context, limit int, cursor *string, sort string, order string) (*models.PaginatedResponse, error) {
	locations, nextCursor, hasMore, err := s.repo.List(ctx, limit, cursor, sort, order)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}
	
	return &models.PaginatedResponse{
		Data: locations,
		Pagination: models.Pagination{
			NextCursor: nextCursor,
			HasMore:    hasMore,
		},
	}, nil
}

// UpdateLocation updates an existing location
func (s *LocationService) UpdateLocation(ctx context.Context, id uuid.UUID, req *models.LocationUpdate) (*models.Location, error) {
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	
	// Apply updates
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Address != nil {
		location.Address = *req.Address
	}
	if req.Phone != nil {
		location.Phone = req.Phone
	}
	if req.Latitude != nil {
		location.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		location.Longitude = *req.Longitude
	}
	
	if err := s.repo.Update(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}
	
	return location, nil
}

// DeleteLocation deletes a location
func (s *LocationService) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}
	return nil
}

// FindNearby finds locations near given coordinates
func (s *LocationService) FindNearby(ctx context.Context, lat, long, radius float64, limit int) (*models.NearbyResponse, error) {
	locations, err := s.repo.FindNearby(ctx, lat, long, radius, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}
	
	return &models.NearbyResponse{
		Data: locations,
		Metadata: models.NearbyMetadata{
			Center: models.Center{
				Latitude:  lat,
				Longitude: long,
			},
			Radius: radius,
			Count:  len(locations),
		},
	}, nil
}
