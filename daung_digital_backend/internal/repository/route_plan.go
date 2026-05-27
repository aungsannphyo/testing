package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/daung-digital/location-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoutePlanRepository defines the interface for route plan data access
type RoutePlanRepository interface {
	Create(ctx context.Context, routePlan *models.RoutePlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.RoutePlan, error)
	List(ctx context.Context, limit int, cursor *string, sort string, order string) ([]models.RoutePlan, *string, bool, error)
	Update(ctx context.Context, routePlan *models.RoutePlan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresRoutePlanRepository implements RoutePlanRepository for PostgreSQL
type PostgresRoutePlanRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRoutePlanRepository creates a new PostgreSQL route plan repository
func NewPostgresRoutePlanRepository(pool *pgxpool.Pool) RoutePlanRepository {
	return &PostgresRoutePlanRepository{pool: pool}
}

// Create inserts a new route plan into the database
func (r *PostgresRoutePlanRepository) Create(ctx context.Context, routePlan *models.RoutePlan) error {
	query := `
		INSERT INTO route_plans (id, name, location_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	now := time.Now()
	routePlan.CreatedAt = now
	routePlan.UpdatedAt = now

	err := r.pool.QueryRow(ctx, query,
		routePlan.ID,
		routePlan.Name,
		routePlan.LocationIDs,
		now,
		now,
	).Scan(&routePlan.CreatedAt, &routePlan.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create route plan: %w", err)
	}

	return nil
}

// GetByID retrieves a route plan by its ID
func (r *PostgresRoutePlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RoutePlan, error) {
	query := `
		SELECT id, name, location_ids, created_at, updated_at
		FROM route_plans
		WHERE id = $1
	`

	var routePlan models.RoutePlan
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&routePlan.ID,
		&routePlan.Name,
		&routePlan.LocationIDs,
		&routePlan.CreatedAt,
		&routePlan.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("route plan not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get route plan: %w", err)
	}

	return &routePlan, nil
}

// List retrieves route plans with pagination
func (r *PostgresRoutePlanRepository) List(ctx context.Context, limit int, cursor *string, sort string, order string) ([]models.RoutePlan, *string, bool, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	query := fmt.Sprintf(`
		SELECT id, name, location_ids, created_at, updated_at
		FROM route_plans
		ORDER BY %s %s
		LIMIT $1 + 1
	`, sort, order)

	rows, err := r.pool.Query(ctx, query, limit+1)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to list route plans: %w", err)
	}
	defer rows.Close()

	var routePlans []models.RoutePlan
	for rows.Next() {
		var routePlan models.RoutePlan
		if err := rows.Scan(
			&routePlan.ID,
			&routePlan.Name,
			&routePlan.LocationIDs,
			&routePlan.CreatedAt,
			&routePlan.UpdatedAt,
		); err != nil {
			return nil, nil, false, fmt.Errorf("failed to scan route plan: %w", err)
		}
		routePlans = append(routePlans, routePlan)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, false, fmt.Errorf("error iterating route plans: %w", err)
	}

	hasMore := len(routePlans) > limit
	if hasMore {
		routePlans = routePlans[:limit]
	}

	var nextCursor *string
	if hasMore && len(routePlans) > 0 {
		lastID := routePlans[len(routePlans)-1].ID.String()
		nextCursor = &lastID
	}

	return routePlans, nextCursor, hasMore, nil
}

// Update updates an existing route plan
func (r *PostgresRoutePlanRepository) Update(ctx context.Context, routePlan *models.RoutePlan) error {
	query := `
		UPDATE route_plans
		SET name = $2, location_ids = $3, updated_at = $4
		WHERE id = $1
		RETURNING updated_at
	`

	routePlan.UpdatedAt = time.Now()

	err := r.pool.QueryRow(ctx, query,
		routePlan.ID,
		routePlan.Name,
		routePlan.LocationIDs,
		routePlan.UpdatedAt,
	).Scan(&routePlan.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("route plan not found: %w", err)
		}
		return fmt.Errorf("failed to update route plan: %w", err)
	}

	return nil
}

// Delete removes a route plan from the database
func (r *PostgresRoutePlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM route_plans WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete route plan: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("route plan not found")
	}

	return nil
}
