.PHONY: build test lint clean run docker-build docker-up docker-down

# Variables
BINARY_NAME=server
BUILD_DIR=bin
GO=go
GOFLAGS=-v

# Build the application
build:
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run:
	$(GO) run ./cmd/server

# Run tests
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	$(GO) tool cover -html=coverage.out

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GO) fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# Install dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Docker build
docker-build:
	docker build -t location-api .

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down

# Docker compose logs
docker-logs:
	docker-compose logs -f

# Run with race detector
run-race:
	$(GO) run -race ./cmd/server

# Generate go.sum
go-sum:
	$(GO) mod tidy
