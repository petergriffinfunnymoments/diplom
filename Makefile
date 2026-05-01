.PHONY: build run test clean migrate-up migrate-down docker-up docker-down run-all

# Variables
APP_NAME=payment-gateway
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Build all services
build:
	@echo "Building all services..."
	go build $(LDFLAGS) -o bin/api-gateway ./cmd/api-gateway
	go build $(LDFLAGS) -o bin/orchestrator ./cmd/orchestrator
	go build $(LDFLAGS) -o bin/validator ./cmd/validator
	go build $(LDFLAGS) -o bin/antifraud ./cmd/antifraud
	go build $(LDFLAGS) -o bin/tokenizer ./cmd/tokenizer
	go build $(LDFLAGS) -o bin/adapter ./cmd/adapter
	go build $(LDFLAGS) -o bin/notification ./cmd/notification
	go build $(LDFLAGS) -o bin/logger ./cmd/logger
	@echo "Build completed successfully!"

# Build specific service
build-service:
ifndef SERVICE
	$(error SERVICE is required. Usage: make build-service SERVICE=api-gateway)
endif
	@echo "Building $(SERVICE)..."
	go build $(LDFLAGS) -o bin/$(SERVICE) ./cmd/$(SERVICE)
	@echo "Build completed!"

# Run all services (for development)
run-all:
	@echo "Starting all services..."
	go run ./cmd/api-gateway &
	go run ./cmd/orchestrator &
	go run ./cmd/validator &
	go run ./cmd/antifraud &
	go run ./cmd/tokenizer &
	go run ./cmd/adapter &
	go run ./cmd/notification &
	go run ./cmd/logger &
	@echo "All services started. Press Ctrl+C to stop."
	wait

# Run specific service
run:
ifndef SERVICE
	$(error SERVICE is required. Usage: make run SERVICE=api-gateway)
endif
	@echo "Running $(SERVICE)..."
	go run ./cmd/$(SERVICE)

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean completed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed!"

# Database migrations
migrate-up:
ifndef DATABASE_URL
	$(error DATABASE_URL is required. Usage: make migrate-up DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable")
endif
	@echo "Running migrations up..."
	migrate -path migrations -database "$(DATABASE_URL)" up
	@echo "Migrations completed!"

migrate-down:
ifndef DATABASE_URL
	$(error DATABASE_URL is required. Usage: make migrate-down DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable")
endif
	@echo "Running migrations down..."
	migrate -path migrations -database "$(DATABASE_URL)" down
	@echo "Migrations rollback completed!"

migrate-create:
ifndef NAME
	$(error NAME is required. Usage: make migrate-create NAME=create_users_table)
endif
	@echo "Creating migration: $(NAME)..."
	touch migrations/$(NAME).up.sql
	touch migrations/$(NAME).down.sql
	@echo "Migration files created!"

# Docker commands
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Docker containers started!"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Docker containers stopped!"

docker-logs:
	docker-compose logs -f

docker-restart:
	docker-compose restart

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "Linting completed!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Formatting completed!"

# Generate mocks (if using gomock)
mocks:
	@echo "Generating mocks..."
	go generate ./...
	@echo "Mocks generated!"

# Help
help:
	@echo "Payment Gateway Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build              - Build all services"
	@echo "  make build-service SERVICE=<name> - Build specific service"
	@echo "  make run SERVICE=<name> - Run specific service"
	@echo "  make run-all            - Run all services (development)"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make deps               - Install dependencies"
	@echo "  make migrate-up DATABASE_URL=<url> - Run migrations up"
	@echo "  make migrate-down DATABASE_URL=<url> - Run migrations down"
	@echo "  make migrate-create NAME=<name> - Create new migration"
	@echo "  make docker-up          - Start Docker containers"
	@echo "  make docker-down        - Stop Docker containers"
	@echo "  make docker-logs        - View Docker logs"
	@echo "  make lint               - Run linter"
	@echo "  make fmt                - Format code"
	@echo "  make mocks              - Generate mocks"
	@echo "  make help               - Show this help"
