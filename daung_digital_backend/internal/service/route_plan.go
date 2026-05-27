package service

import (
	"context"
	"fmt"

	"github.com/daung-digital/location-api/internal/models"
	"github.com/daung-digital/location-api/internal/repository"
	"github.com/google/uuid"
)

// RoutePlanService handles business logic for route plans
type RoutePlanService struct {
	repo repository.RoutePlanRepository
}

// NewRoutePlanService creates a new route plan service
func NewRoutePlanService(repo repository.RoutePlanRepository) *RoutePlanService {
	return &RoutePlanService{repo: repo}
}

// CreateRoutePlan creates a new route plan
func (s *RoutePlanService) CreateRoutePlan(ctx context.Context, req *models.RoutePlanCreate) (*models.RoutePlan, error) {
	routePlan := &models.RoutePlan{
		ID:          uuid.New(),
		Name:        req.Name,
		LocationIDs: req.LocationIDs,
	}

	if err := s.repo.Create(ctx, routePlan); err != nil {
		return nil, fmt.Errorf("failed to create route plan: %w", err)
	}

	return routePlan, nil
}

// GetRoutePlan retrieves a route plan by ID
func (s *RoutePlanService) GetRoutePlan(ctx context.Context, id uuid.UUID) (*models.RoutePlan, error) {
	routePlan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get route plan: %w", err)
	}
	return routePlan, nil
}

// ListRoutePlans retrieves route plans with pagination
func (s *RoutePlanService) ListRoutePlans(ctx context.Context, limit int, cursor *string, sort string, order string) ([]models.RoutePlan, *string, bool, error) {
	routePlans, nextCursor, hasMore, err := s.repo.List(ctx, limit, cursor, sort, order)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to list route plans: %w", err)
	}

	return routePlans, nextCursor, hasMore, nil
}

// UpdateRoutePlan updates an existing route plan
func (s *RoutePlanService) UpdateRoutePlan(ctx context.Context, id uuid.UUID, req *models.RoutePlanUpdate) (*models.RoutePlan, error) {
	routePlan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get route plan: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		routePlan.Name = *req.Name
	}
	if req.LocationIDs != nil {
		routePlan.LocationIDs = *req.LocationIDs
	}

	if err := s.repo.Update(ctx, routePlan); err != nil {
		return nil, fmt.Errorf("failed to update route plan: %w", err)
	}

	return routePlan, nil
}

// DeleteRoutePlan deletes a route plan
func (s *RoutePlanService) DeleteRoutePlan(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete route plan: %w", err)
	}
	return nil
}
