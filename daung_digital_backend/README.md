# Location API Backend

Go backend for the Location Management API using Gin framework and PostgreSQL.

## Features

- RESTful API for location management
- JWT authentication
- PostgreSQL with pgx driver
- Geospatial queries (nearby locations)
- Cursor-based pagination
- Structured logging with Zap
- Graceful shutdown

## Project Structure

```
daung_digital_backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Configuration management
├── internal/
│   ├── database/
│   │   └── postgres.go      # Database connection and migrations
│   ├── handler/
│   │   └── location.go      # HTTP handlers
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication
│   │   ├── cors.go          # CORS middleware
│   │   └── logger.go        # Request logging
│   ├── models/
│   │   └── location.go      # Data models
│   ├── repository/
│   │   └── location.go      # Data access layer
│   └── service/
│       └── location.go      # Business logic
├── go.mod
├── go.sum
└── .env.example
```

## Prerequisites

- Go 1.21+
- PostgreSQL 12+ with earthdistance extension

## Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Edit `.env` with your configuration:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=location_db
JWT_SECRET=your_secret_key
```

3. Install dependencies:
```bash
go mod download
```

4. Run the server:
```bash
go run cmd/server/main.go
```

## API Endpoints

All endpoints require JWT authentication in the `Authorization` header.

### Locations

- `POST /v1/locations` - Create a new location
- `GET /v1/locations` - List all locations (paginated)
- `GET /v1/locations/:id` - Get a single location
- `PATCH /v1/locations/:id` - Update a location
- `DELETE /v1/locations/:id` - Delete a location
- `GET /v1/locations/nearby` - Find locations near coordinates

## Example Requests

### Create Location
```bash
curl -X POST http://localhost:8080/v1/locations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Central Park",
    "address": "59th to 110th St, New York, NY 10026",
    "phone": "+1-212-310-6600",
    "latitude": 40.7829,
    "longitude": -73.9654
  }'
```

### List Locations
```bash
curl -X GET http://localhost:8080/v1/locations?limit=20&sort=created_at&order=desc \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Find Nearby
```bash
curl -X GET "http://localhost:8080/v1/locations/nearby?latitude=40.7829&longitude=-73.9654&radius=10&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Database

The application uses PostgreSQL with the following extensions:
- `uuid-ossp` - For UUID generation
- `earthdistance` - For geospatial queries

### Manual Setup (Optional)

If you need to set up the database manually:

```sql
CREATE DATABASE location_db;

\c location_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    phone VARCHAR(20),
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_locations_created_at ON locations(created_at DESC);
CREATE INDEX idx_locations_updated_at ON locations(updated_at DESC);
CREATE INDEX idx_locations_name ON locations(name);
```

## Development

### Run tests
```bash
go test ./... -v -race
```

### Run linter
```bash
golangci-lint run
```

### Format code
```bash
go fmt ./...
```

## Building

```bash
go build -o bin/server cmd/server/main.go
```

## Docker

Build and run with Docker:

```bash
docker build -t location-api .
docker run -p 8080:8080 --env-file .env location-api
```

## License

MIT
