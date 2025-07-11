# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Resume API is a Go-based REST API that serves resume/CV data from PostgreSQL. This is a learning project designed to demonstrate Go best practices, clean architecture, and modern deployment patterns. The API will be deployed on Kubernetes in a homelab environment and consumed by a frontend built with Lovable.

## Development Commands

### Core Development
```bash
# Run the API server
go run cmd/api/main.go

# Run with specific environment
ENV=development go run cmd/api/main.go

# Build for production
go build -o bin/api cmd/api/main.go

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/services/

# Format and lint code
go fmt ./...
golangci-lint run

# Update dependencies
go mod tidy
```

### Database Operations
```bash
# Run database migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down 1

# Seed database with resume data
go run scripts/seed.go

# Create new migration
migrate create -ext sql -dir migrations -seq migration_name
```

### Testing and Quality
```bash
# Run all tests with race detection
go test -race ./...

# Generate test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests
go test -tags=integration ./tests/

# Generate mocks (if using mockery)
go generate ./...
```

## Project Architecture

The project follows Clean Architecture with these layers:

### Directory Structure
```
cmd/api/           # Application entry point
internal/
├── handlers/      # HTTP request handlers (Gin/Echo)
├── services/      # Business logic layer
├── repository/    # Data access layer (PostgreSQL)
├── models/        # Domain models and DTOs
├── database/      # DB connection and migrations
├── middleware/    # HTTP middleware (CORS, logging, etc.)
├── config/        # Configuration management
└── validator/     # Custom validation logic

docs/              # Design docs, API specs, task tracking
deployments/       # Kubernetes manifests
scripts/           # Database utilities
tests/             # Integration tests
migrations/        # Database migration files
```

### Key Design Patterns
- **Repository Pattern**: Data access abstraction with interfaces
- **Dependency Injection**: Services injected into handlers via interfaces
- **Context Propagation**: All operations accept `context.Context`
- **Error Wrapping**: Errors wrapped with context using `fmt.Errorf`
- **Clean Architecture**: Dependencies flow inward through interfaces

## API Endpoints

The API serves resume data through these endpoints:

```
GET /api/v1/profile      # Personal information
GET /api/v1/experiences  # Work history (filterable)
GET /api/v1/skills       # Skills by category
GET /api/v1/achievements # Key achievements  
GET /api/v1/education    # Education & certifications
GET /health              # Health check for Kubernetes
```

## Database Schema

PostgreSQL database with these main tables:
- `profiles` - Personal information and summary
- `experiences` - Work history with date ranges
- `skills` - Categorized skills with levels
- `achievements` - Key accomplishments
- `education` - Education and certifications

## Configuration

Environment-based configuration using Viper:

```bash
# Required environment variables
DB_HOST=localhost
DB_PORT=5432
DB_NAME=resume_api_dev
DB_USER=dev
DB_PASSWORD=devpass
SERVER_PORT=8080
LOG_LEVEL=info
```

## Testing Strategy

- **Unit Tests**: Repository, service, and handler layers with mocking
- **Integration Tests**: End-to-end API testing with test database
- **Table-Driven Tests**: Go testing best practices
- **Test Coverage**: Target 80%+ coverage

## Deployment

### Local Development
```bash
# Start PostgreSQL (Docker)
docker run --name resume-postgres \
  -e POSTGRES_DB=resume_api_dev \
  -e POSTGRES_USER=dev \
  -e POSTGRES_PASSWORD=devpass \
  -p 5432:5432 -d postgres:15

# Run migrations and seed data
go run cmd/migrate/main.go up
go run scripts/seed.go

# Start API server
go run cmd/api/main.go
```

### Production (Kubernetes)
- Containerized deployment with multi-stage Dockerfile
- Kubernetes manifests in `/deployments/`
- ConfigMap for configuration
- Secrets for sensitive data
- Health probes configured

## Learning Objectives

This project demonstrates:
- Go web development with popular frameworks
- Clean Architecture and dependency injection
- PostgreSQL integration and migrations
- RESTful API design and error handling
- Testing patterns and mocking
- Containerization and Kubernetes deployment
- Modern Go project structure and tooling

## Documentation

- **System Design**: `/docs/design/system-design.md`
- **Development Guide**: `/docs/development.md`
- **Task Tracking**: `/docs/tasks.md`
- **API Documentation**: Generated from OpenAPI specs

## Important Notes

- The `/resume/` folder contains personal CV files and is gitignored
- Database migrations must be run before starting the API
- All code follows Go best practices and clean architecture principles
- The project serves as both a functional API and a Go learning exercise