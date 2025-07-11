# Resume API Makefile
# Provides convenient commands for development and testing

.PHONY: help test test-short test-docker test-compose dev-db clean build lint

# Default target
help:
	@echo "Resume API Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Testing:"
	@echo "  test-short       Run tests without database (compilation check only)"
	@echo "  test-docker      Run full tests with Docker PostgreSQL container"
	@echo "  test-compose     Run full tests with Docker Compose"
	@echo "  test             Alias for test-docker (recommended)"
	@echo ""
	@echo "Development:"
	@echo "  dev-db          Start development database with Docker Compose"
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

# Testing commands
test-short:
	@echo "🧪 Running short tests (compilation check only)..."
	go test -short ./... -v

test-docker:
	@echo "🐳 Running tests with Docker PostgreSQL..."
	./scripts/test-repositories.sh

test-compose:
	@echo "🐳 Running tests with Docker Compose..."
	./scripts/test-repositories-compose.sh

test: test-docker

# Development database
dev-db:
	@echo "🐳 Starting development database..."
	docker-compose -f docker-compose.test.yml up -d test-db
	@echo "✅ Development database started on port 5433"
	@echo "   Connection: postgres://dev:devpass@localhost:5433/resume_api_test"

dev-db-admin:
	@echo "🐳 Starting development database with pgAdmin..."
	docker-compose -f docker-compose.test.yml --profile admin up -d
	@echo "✅ Development database and pgAdmin started"
	@echo "   Database: postgres://dev:devpass@localhost:5433/resume_api_test"
	@echo "   pgAdmin: http://localhost:5050 (admin@test.com / admin)"

dev-db-stop:
	@echo "🛑 Stopping development database..."
	docker-compose -f docker-compose.test.yml down -v

# Migration commands
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

# Build and quality
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

# Cleanup
clean:
	@echo "🧹 Cleaning up Docker containers and volumes..."
	docker-compose -f docker-compose.test.yml down -v --remove-orphans || true
	docker container rm resume-api-test-db || true
	docker volume prune -f || true
	@echo "✅ Cleanup completed"

# Dependencies and tools
deps:
	@echo "📦 Downloading and tidying dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

tools:
	@echo "🛠️  Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Convenience targets
.PHONY: up down logs
up: dev-db
down: dev-db-stop
logs:
	docker-compose -f docker-compose.test.yml logs -f test-db