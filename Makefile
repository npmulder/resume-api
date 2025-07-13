# Resume API Makefile
# Provides convenient commands for development and testing

# Configuration variables
POSTGRES_VERSION := 15-alpine
TEST_DB_HOST := localhost
TEST_DB_PORT := 5433
TEST_DB_NAME := resume_api_test
TEST_DB_USER := dev
TEST_DB_PASSWORD := devpass
DOCKER_COMPOSE_FILE := docker-compose.test.yml

# Define all phony targets (targets that don't create files)
.PHONY: help test test-short test-docker test-compose test-integration \
        dev-db dev-db-admin dev-db-stop \
        migrate-up migrate-down \
        build lint clean deps tools \
        up down logs

# Default target
help:
	@echo "Resume API Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Testing:"
	@echo "  test-short       Run tests without database (compilation check only)"
	@echo "  test-docker      Run full tests with Docker PostgreSQL container"
	@echo "  test-compose     Run full tests with Docker Compose"
	@echo "  test-integration Run integration tests with Docker PostgreSQL container"
	@echo "  test             Alias for test-docker (recommended)"
	@echo ""
	@echo "Development:"
	@echo "  dev-db          Start development database with Docker Compose"
	@echo "  dev-db-admin    Start development database with pgAdmin"
	@echo "  dev-db-stop     Stop development database"
	@echo "  migrate-up      Run database migrations"
	@echo "  migrate-down    Rollback database migrations"
	@echo ""
	@echo "Build & Quality:"
	@echo "  build           Build all packages"
	@echo "  lint            Run linter (if available)"
	@echo "  clean           Clean up Docker containers and volumes"
	@echo ""
	@echo "Utilities:"
	@echo "  deps            Download and tidy dependencies"
	@echo "  tools           Install development tools"
	@echo "  up              Alias for dev-db"
	@echo "  down            Alias for dev-db-stop"
	@echo "  logs            View database logs"

#
# Testing commands
#
test-short:
	@echo "🧪 Running short tests (compilation check only)..."
	go test -short ./... -v

test-docker:
	@echo "🐳 Running tests with Docker PostgreSQL..."
	./scripts/test-repositories.sh

test-compose:
	@echo "🐳 Running tests with Docker Compose..."
	./scripts/test-repositories-compose.sh

test-integration:
	@echo "🧪 Running integration tests with Docker PostgreSQL..."
	./scripts/test-integration.sh

# Main test alias
test: test-docker

#
# Development database commands
#
dev-db:
	@echo "🐳 Starting development database..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d test-db
	@echo "✅ Development database started on port $(TEST_DB_PORT)"
	@echo "   Connection: postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@$(TEST_DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)"

dev-db-admin:
	@echo "🐳 Starting development database with pgAdmin..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) --profile admin up -d
	@echo "✅ Development database and pgAdmin started"
	@echo "   Database: postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@$(TEST_DB_HOST):$(TEST_DB_PORT)/$(TEST_DB_NAME)"
	@echo "   pgAdmin: http://localhost:5050 (admin@test.com / admin)"

dev-db-stop:
	@echo "🛑 Stopping development database..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v

#
# Migration commands
#
migrate-up:
	@echo "⬆️  Running database migrations..."
	@if [ -f "./cmd/migrate/main.go" ]; then \
		go run ./cmd/migrate/main.go up; \
	else \
		echo "❌ Migration tool not found at ./cmd/migrate/main.go"; \
	fi

migrate-down:
	@echo "⬇️  Rolling back database migrations..."
	@if [ -f "./cmd/migrate/main.go" ]; then \
		go run ./cmd/migrate/main.go down; \
	else \
		echo "❌ Migration tool not found at ./cmd/migrate/main.go"; \
	fi

#
# Build and quality commands
#
build:
	@echo "🔨 Building all packages..."
	go build ./...
	@echo "✅ Build completed successfully"

lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Running go vet instead..."; \
		go vet ./...; \
	fi

clean:
	@echo "🧹 Cleaning up Docker containers and volumes..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans || true
	docker container rm resume-api-test-db || true
	docker volume prune -f || true
	@echo "✅ Cleanup completed"

#
# Dependency and tool commands
#
deps:
	@echo "📦 Downloading and tidying dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

tools:
	@echo "🛠️  Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

#
# Convenience aliases
#
up: dev-db
down: dev-db-stop
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f test-db
