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

// LocationRepository defines the interface for location data access
type LocationRepository interface {
	Create(ctx context.Context, location *models.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Location, error)
	List(ctx context.Context, limit int, cursor *string, sort string, order string) ([]models.Location, *string, bool, error)
	Update(ctx context.Context, location *models.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindNearby(ctx context.Context, lat, long, radius float64, limit int) ([]models.Location, error)
}

// PostgresLocationRepository implements LocationRepository for PostgreSQL
type PostgresLocationRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresLocationRepository creates a new PostgreSQL location repository
func NewPostgresLocationRepository(pool *pgxpool.Pool) LocationRepository {
	return &PostgresLocationRepository{pool: pool}
}

// Create inserts a new location into the database
func (r *PostgresLocationRepository) Create(ctx context.Context, location *models.Location) error {
	query := `
		INSERT INTO locations (id, name, address, phone, latitude, longitude, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now
	
	err := r.pool.QueryRow(ctx, query,
		location.ID,
		location.Name,
		location.Address,
		location.Phone,
		location.Latitude,
		location.Longitude,
		now,
		now,
	).Scan(&location.CreatedAt, &location.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}
	
	return nil
}

// GetByID retrieves a location by its ID
func (r *PostgresLocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	query := `
		SELECT id, name, address, phone, latitude, longitude, created_at, updated_at
		FROM locations
		WHERE id = $1
	`
	
	var location models.Location
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&location.ID,
		&location.Name,
		&location.Address,
		&location.Phone,
		&location.Latitude,
		&location.Longitude,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("location not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	
	return &location, nil
}

// List retrieves locations with pagination
func (r *PostgresLocationRepository) List(ctx context.Context, limit int, cursor *string, sort string, order string) ([]models.Location, *string, bool, error) {
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
		SELECT id, name, address, phone, latitude, longitude, created_at, updated_at
		FROM locations
		ORDER BY %s %s
		LIMIT $1 + 1
	`, sort, order)
	
	rows, err := r.pool.Query(ctx, query, limit+1)
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to list locations: %w", err)
	}
	defer rows.Close()
	
	var locations []models.Location
	for rows.Next() {
		var location models.Location
		if err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.Phone,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt,
			&location.UpdatedAt,
		); err != nil {
			return nil, nil, false, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}
	
	if err := rows.Err(); err != nil {
		return nil, nil, false, fmt.Errorf("error iterating locations: %w", err)
	}
	
	hasMore := len(locations) > limit
	if hasMore {
		locations = locations[:limit]
	}
	
	var nextCursor *string
	if hasMore && len(locations) > 0 {
		lastID := locations[len(locations)-1].ID.String()
		nextCursor = &lastID
	}
	
	return locations, nextCursor, hasMore, nil
}

// Update updates an existing location
func (r *PostgresLocationRepository) Update(ctx context.Context, location *models.Location) error {
	query := `
		UPDATE locations
		SET name = $2, address = $3, phone = $4, latitude = $5, longitude = $6, updated_at = $7
		WHERE id = $1
		RETURNING updated_at
	`
	
	location.UpdatedAt = time.Now()
	
	err := r.pool.QueryRow(ctx, query,
		location.ID,
		location.Name,
		location.Address,
		location.Phone,
		location.Latitude,
		location.Longitude,
		location.UpdatedAt,
	).Scan(&location.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("location not found: %w", err)
		}
		return fmt.Errorf("failed to update location: %w", err)
	}
	
	return nil
}

// Delete removes a location from the database
func (r *PostgresLocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM locations WHERE id = $1`
	
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("location not found")
	}
	
	return nil
}

// FindNearby finds locations within a given radius of coordinates
func (r *PostgresLocationRepository) FindNearby(ctx context.Context, lat, long, radius float64, limit int) ([]models.Location, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	// Use Haversine formula for distance calculation
	query := `
		SELECT id, name, address, phone, latitude, longitude, created_at, updated_at
		FROM locations
		WHERE earthdistance(ll_to_earth($1, $2) <@ earth_box(ll_to_earth($1, $2), $3))
		ORDER BY earth_distance(ll_to_earth($1, $2), ll_to_earth(latitude, longitude))
		LIMIT $4
	`
	
	rows, err := r.pool.Query(ctx, query, lat, long, radius*1000, limit) // radius in meters
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}
	defer rows.Close()
	
	var locations []models.Location
	for rows.Next() {
		var location models.Location
		if err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
			&location.Phone,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt,
			&location.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating locations: %w", err)
	}
	
	return locations, nil
}
