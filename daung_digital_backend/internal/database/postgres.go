package database

import (
	"context"
	"fmt"
	"time"

	"github.com/daung-digital/location-api/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// NewPool creates a new PostgreSQL connection pool
func NewPool(cfg *config.DatabaseConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.MaxLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxLifetime
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return pool, nil
}

// RunMigrations runs database migrations
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE EXTENSION IF NOT EXISTS cube;`,
		`CREATE EXTENSION IF NOT EXISTS earthdistance;`,
		`CREATE TABLE IF NOT EXISTS locations (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			phone VARCHAR(20),
			latitude DECIMAL(10, 8) NOT NULL,
			longitude DECIMAL(11, 8) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_locations_created_at ON locations(created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_locations_updated_at ON locations(updated_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_locations_name ON locations(name);`,
		`CREATE TABLE IF NOT EXISTS route_plans (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			location_ids UUID[] NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_route_plans_created_at ON route_plans(created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_route_plans_name ON route_plans(name);`,
	}

	for _, migration := range migrations {
		if _, err := pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}
