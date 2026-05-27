package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string        `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port         int           `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host         string        `envconfig:"DB_HOST" default:"localhost"`
	Port         int           `envconfig:"DB_PORT" default:"5432"`
	User         string        `envconfig:"DB_USER" required:"true"`
	Password     string        `envconfig:"DB_PASSWORD" required:"true"`
	Database     string        `envconfig:"DB_NAME" required:"true"`
	SSLMode      string        `envconfig:"DB_SSLMODE" default:"disable"`
	MaxOpenConns int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	MaxLifetime  time.Duration `envconfig:"DB_MAX_LIFETIME" default:"5m"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
