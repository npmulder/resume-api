# Development Guide - Resume API

## Quick Start

### Prerequisites
- Go 1.21+ installed
- PostgreSQL 15+ running locally or in container
- Docker (optional, for containerized database)
- Git

### Local Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd resume-api
   go mod tidy
   ```

2. **Database Setup**:
   ```bash
   # Option 1: Local PostgreSQL
   createdb resume_api_dev
   
   # Option 2: Docker container
   docker run --name resume-postgres \
     -e POSTGRES_DB=resume_api_dev \
     -e POSTGRES_USER=dev \
     -e POSTGRES_PASSWORD=devpass \
     -p 5432:5432 -d postgres:15
   ```

3. **Environment Configuration**:
   ```bash
   cp .env.example .env
   # Edit .env with your database settings
   ```

4. **Run Database Migrations**:
   ```bash
   go run cmd/migrate/main.go up
   ```

5. **Seed Database**:
   ```bash
   go run scripts/seed.go
   ```

6. **Run the API**:
   ```bash
   go run cmd/api/main.go
   ```

## Project Structure

```
resume-api/
├── cmd/
│   ├── api/            # Main application entry point
│   └── migrate/        # Database migration runner
├── internal/           # Private application code
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models and DTOs
│   ├── repository/     # Data access layer
│   ├── services/       # Business logic layer
│   └── validator/      # Custom validation logic
├── pkg/               # Public library code (if any)
├── docs/              # Documentation
├── deployments/       # Kubernetes and Docker configs
├── scripts/           # Utility scripts
├── tests/             # Test files
├── migrations/        # Database migration files
└── resume/            # Personal resume files (gitignored)
```

## Development Standards

### Go Code Style

**1. Naming Conventions**:
```go
// Interfaces: noun + "er" suffix
type ProfileRepository interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
}

// Structs: CamelCase, descriptive
type ResumeService struct {
    profileRepo ProfileRepository
    logger      *slog.Logger
}

// Functions: CamelCase, verb-based
func (s *ResumeService) GetProfile(ctx context.Context) (*models.Profile, error) {
    // implementation
}
```

**2. Error Handling**:
```go
// Always wrap errors with context
func (r *profileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
    var profile models.Profile
    err := r.db.GetContext(ctx, &profile, query)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrProfileNotFound
        }
        return nil, fmt.Errorf("failed to get profile from database: %w", err)
    }
    return &profile, nil
}

// Define custom errors
var (
    ErrProfileNotFound = errors.New("profile not found")
    ErrInvalidInput    = errors.New("invalid input")
)
```

**3. Context Usage**:
```go
// Always accept context as first parameter
func (s *ResumeService) GetExperiences(ctx context.Context, filters ExperienceFilters) ([]*models.Experience, error) {
    // Pass context to repository layer
    return s.experienceRepo.GetExperiences(ctx, filters)
}
```

**4. Interface Design**:
```go
// Keep interfaces small and focused
type ProfileRepository interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
    UpdateProfile(ctx context.Context, profile *models.Profile) error
}

// Define interfaces where they're used (handlers), not where they're implemented
```

### Database Patterns

**1. Repository Pattern**:
```go
type profileRepository struct {
    db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
    return &profileRepository{db: db}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
    // Use named queries for better maintainability
    query := `
        SELECT id, name, title, email, phone, location, linkedin, summary, updated_at
        FROM profiles 
        LIMIT 1`
    
    var profile models.Profile
    err := r.db.GetContext(ctx, &profile, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get profile: %w", err)
    }
    return &profile, nil
}
```

**2. Migration Files**:
```sql
-- migrations/001_create_profiles.up.sql
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    location VARCHAR(255),
    linkedin VARCHAR(255),
    summary TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/001_create_profiles.down.sql
DROP TABLE IF EXISTS profiles;
```

### Testing Patterns

**1. Table-Driven Tests**:
```go
func TestProfileService_GetProfile(t *testing.T) {
    tests := []struct {
        name          string
        mockReturn    *models.Profile
        mockError     error
        expectedError error
    }{
        {
            name: "successful retrieval",
            mockReturn: &models.Profile{
                ID:   1,
                Name: "John Doe",
            },
            mockError:     nil,
            expectedError: nil,
        },
        {
            name:          "repository error",
            mockReturn:    nil,
            mockError:     errors.New("db error"),
            expectedError: errors.New("db error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**2. Mock Interfaces**:
```go
//go:generate mockery --name=ProfileRepository --output=mocks
type ProfileRepository interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
}
```

### Configuration Management

**Environment Variables**:
```bash
# .env.example
# Server Configuration
SERVER_PORT=8080
SERVER_TIMEOUT=30s
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=resume_api_dev
DB_USER=dev
DB_PASSWORD=devpass
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

**Configuration Struct**:
```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Port    int           `mapstructure:"port"`
    Host    string        `mapstructure:"host"`
    Timeout time.Duration `mapstructure:"timeout"`
}
```

## Common Development Tasks

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestProfileService_GetProfile ./internal/services/
```

### Database Operations
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_experiences_table

# Run migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down 1

# Check migration status
go run cmd/migrate/main.go version
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for race conditions
go test -race ./...

# Check for unused dependencies
go mod tidy
```

### Building and Running
```bash
# Build for current platform
go build -o bin/api cmd/api/main.go

# Build for Linux (for containers)
GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go

# Run with specific environment
ENV=development go run cmd/api/main.go

# Run with custom config file
CONFIG_FILE=config.local.yaml go run cmd/api/main.go
```

## API Testing

### Using curl
```bash
# Get profile
curl http://localhost:8080/api/v1/profile

# Get experiences with filtering
curl "http://localhost:8080/api/v1/experiences?company=Derivco&limit=5"

# Get skills by category
curl "http://localhost:8080/api/v1/skills?category=Languages"

# Health check
curl http://localhost:8080/health
```

### Using httpie
```bash
# Install httpie: pip install httpie

# Get profile
http GET localhost:8080/api/v1/profile

# Get experiences
http GET localhost:8080/api/v1/experiences company==Derivco limit==5

# Check health
http GET localhost:8080/health
```

## Debugging

### Logging
```go
// Use structured logging throughout
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// In handlers
logger.Info("handling request",
    slog.String("method", r.Method),
    slog.String("path", r.URL.Path),
    slog.String("user_agent", r.UserAgent()),
)

// In services
logger.Error("failed to get profile",
    slog.String("error", err.Error()),
    slog.String("operation", "GetProfile"),
)
```

### Using Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main application
dlv debug cmd/api/main.go

# Debug tests
dlv test ./internal/services/
```

## Performance Considerations

### Database Connection Pool
```go
// Configure connection pool
config := pgxpool.Config{
    MaxConns:        25,
    MinConns:        5,
    MaxConnLifetime: time.Hour,
    MaxConnIdleTime: time.Minute * 30,
}
```

### HTTP Timeouts
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

## Troubleshooting

### Common Issues

**1. Database Connection Failed**:
- Check PostgreSQL is running: `pg_isready`
- Verify connection string in .env
- Check network connectivity and firewall

**2. Migration Errors**:
- Check migration file syntax
- Ensure database user has proper permissions
- Verify migration hasn't been partially applied

**3. Port Already in Use**:
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

**4. Module Issues**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify module integrity
go mod verify
```

### Useful Commands

```bash
# Check Go version and environment
go version
go env

# List all dependencies
go list -m all

# Check for security vulnerabilities
go list -json -deps ./... | nancy sleuth

# Generate mocks (if using mockery)
go generate ./...

# Cross-compile for different platforms
GOOS=windows GOARCH=amd64 go build -o bin/api.exe cmd/api/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/api-mac cmd/api/main.go
```

This guide provides the foundation for productive development on the Resume API project while learning Go best practices.