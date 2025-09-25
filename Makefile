.PHONY: help build run test clean docker-build docker-up docker-down migrate/up migrate/down

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variables
BINARY_NAME=authorization
DOCKER_IMAGE=authorization-app
GO_VERSION=1.25

# Database connection variables from .env
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_NAME ?= go_login
POSTGRES_EXTERNAL_PORT ?= 8888
APP_EXTERNAL_PORT ?= 8080

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_/-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/server

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

run-dev: ## Run the application in development mode
	@echo "Running in development mode..."
	@go run ./cmd/server

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -f docker/Dockerfile -t $(DOCKER_IMAGE) .

docker-up: env-validate ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-down: ## Stop services with Docker Compose
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

docker-down-volumes: ## Stop services and remove volumes
	@echo "Stopping services and removing volumes..."
	@docker-compose down -v

docker-logs: ## Show logs from Docker Compose services
	@docker-compose logs -f

docker-logs-app: ## Show logs from app service only
	@docker-compose logs -f app

docker-logs-db: ## Show logs from postgres service only  
	@docker-compose logs -f postgres

docker-rebuild: docker-down docker-build docker-up ## Rebuild and restart services

docker-reset: docker-down-volumes docker-up ## Reset everything (remove volumes and restart)

# Database migrations (manual - since we're using GORM AutoMigrate)
migrate/up: ## Apply database migrations (manual SQL)
	@echo "Applying database migrations..."
	@echo "Note: This project uses GORM AutoMigrate. Manual SQL migrations are in migrations/ folder."
	@echo "To run manually: psql -h localhost -p $(POSTGRES_EXTERNAL_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/0001_create_users.up.sql"

migrate/down: ## Rollback database migrations (manual SQL)
	@echo "Rolling back database migrations..."
	@echo "Note: This project uses GORM AutoMigrate. Manual SQL migrations are in migrations/ folder."
	@echo "To run manually: psql -h localhost -p $(POSTGRES_EXTERNAL_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/0002_create_refresh_tokens.down.sql"

# Dependencies
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Code quality
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linters (requires golangci-lint)
	@echo "Running linters..."
	@golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Production
build-prod: ## Build production binary
	@echo "Building production binary..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/$(BINARY_NAME) ./cmd/server

# Quick development workflow
dev: clean fmt vet test build ## Run full development workflow (clean, format, vet, test, build)

# Database helpers (requires running PostgreSQL)
db-connect: ## Connect to PostgreSQL database
	@echo "Connecting to database..."
	@psql -h localhost -p $(POSTGRES_EXTERNAL_PORT) -U $(DB_USER) -d $(DB_NAME)

db-reset: ## Reset database (drop and recreate)
	@echo "Resetting database..."
	@docker-compose exec postgres psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@docker-compose exec postgres psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"

# Health checks
health: ## Check service health
	@echo "Checking service health..."
	@curl -f http://localhost:$(APP_EXTERNAL_PORT)/health || echo "Service is not running"

# Environment
env-copy: ## Copy .env.example to .env
	@echo "Copying .env.example to .env..."
	@cp .env.example .env
	@echo "Don't forget to update the values in .env!"

env-show: ## Show current environment variables
	@echo "Current environment variables from .env:"
	@echo "DB_HOST=$(DB_HOST)"
	@echo "DB_PORT=$(DB_PORT)" 
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_NAME=$(DB_NAME)"
	@echo "POSTGRES_EXTERNAL_PORT=$(POSTGRES_EXTERNAL_PORT)"
	@echo "APP_EXTERNAL_PORT=$(APP_EXTERNAL_PORT)"
	@echo "APP_ENV=$(APP_ENV)"

env-validate: ## Validate .env file exists and has required variables
	@echo "Validating .env file..."
	@test -f .env || (echo "Error: .env file not found. Run 'make env-copy' first." && exit 1)
	@test -n "$(DB_PASS)" || (echo "Error: DB_PASS not set in .env" && exit 1)
	@test -n "$(JWT_SECRET)" || (echo "Error: JWT_SECRET not set in .env" && exit 1)
	@echo "✓ .env file is valid"
