# Resume API Development Guidelines

This document provides essential information for developers working on the Resume API project. It includes build/configuration instructions, testing information, and additional development details specific to this project.

## Build/Configuration Instructions

### Prerequisites
- Go 1.21+ installed
- PostgreSQL 15+ (local or containerized)
- Docker (optional, for containerized database)

### Environment Configuration
The application uses environment variables for configuration, which can be set in a `.env` file:

```
# Copy the example environment file
cp .env.example .env
# Edit the .env file with your specific settings
```

Key configuration parameters:
- `SERVER_PORT`: HTTP server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`: Database connection details
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

### Database Setup

#### Option 1: Local PostgreSQL
```
createdb resume_api_dev
```

#### Option 2: Docker Container (Recommended)
```
# Start a development database
make dev-db

# Or with pgAdmin for database inspection
make dev-db-admin
```

This will start a PostgreSQL container with the following connection details:
- Host: localhost
- Port: 5433 (to avoid conflicts with local PostgreSQL)
- Database: resume_api_test
- User: dev
- Password: devpass

### Running Migrations
```
# Apply all migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### Building and Running the API
```
# Build the application
make build

# Run the API directly with Go
go run cmd/api/main.go
```

## Testing Information

### Test Configuration
Tests that require a database connection use environment variables for configuration:
- `TEST_DB_HOST`: Database host (default: localhost)
- `TEST_DB_PORT`: Database port (default: 5433)
- `TEST_DB_NAME`: Database name (default: resume_api_test)
- `TEST_DB_USER`: Database user (default: dev)
- `TEST_DB_PASSWORD`: Database password (default: devpass)

### Running Tests

#### Quick Tests (No Database)
```
# Run tests without database dependencies
make test-short
```

#### Full Tests with Database
```
# Run tests with Docker PostgreSQL (recommended)
make test

# Alternative: Run tests with Docker Compose
make test-compose
```

The test scripts handle:
1. Starting a PostgreSQL container
2. Running migrations
3. Executing tests
4. Cleaning up the container

### Writing Tests

#### Unit Tests Example
Here's an example of a table-driven test for a utility function:

```
// Example unit test for a utility function
func TestTruncate(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        maxLength int
        expected  string
    }{
        {
            name:      "short string not truncated",
            input:     "Hello",
            maxLength: 10,
            expected:  "Hello",
        },
        {
            name:      "long string truncated",
            input:     "Hello, world! This is a test.",
            maxLength: 10,
            expected:  "Hello, wor...",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Truncate(tt.input, tt.maxLength)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### Repository Tests
Repository tests require a database connection. The project includes helper functions for setting up test databases:

```
// Example repository test
func TestProfileRepository_GetProfile(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    // Test implementation using database
}
```

#### Service Tests with Mocks
Service tests use mocks to isolate the service layer from external dependencies:

```
// Example service test with mocks
func TestResumeService_GetProfile(t *testing.T) {
    mockProfileRepo := new(MockProfileRepository)
    mockRepos := repository.Repositories{Profile: mockProfileRepo}
    service := NewResumeService(mockRepos)

    expectedProfile := &models.Profile{ID: 1, Name: "Test User"}
    mockProfileRepo.On("GetProfile", mock.Anything).Return(expectedProfile, nil)

    profile, err := service.GetProfile(context.Background())

    assert.NoError(t, err)
    assert.Equal(t, expectedProfile, profile)
    mockProfileRepo.AssertExpectations(t)
}
```

## Additional Development Information

### Project Structure
The project follows a clean architecture approach with clear separation of concerns:

```
resume-api/
├── cmd/                # Application entry points
│   ├── api/            # Main API server
│   └── migrate/        # Database migration tool
├── internal/           # Private application code
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   ├── services/       # Business logic layer
│   └── utils/          # Utility functions
├── migrations/         # SQL migration files
└── scripts/            # Utility scripts
```

### Code Style Guidelines

#### Naming Conventions
- Interfaces: Use noun + "er" suffix (e.g., `ProfileRepository`)
- Structs: Use CamelCase, descriptive names (e.g., `ResumeService`)
- Functions: Use CamelCase, verb-based names (e.g., `GetProfile`)

#### Error Handling
- Always wrap errors with context using `fmt.Errorf("context: %w", err)`
- Define custom errors for expected error conditions
- Check specific errors using `errors.Is()` or `errors.As()`

#### Context Usage
- Always accept context as the first parameter in functions that perform I/O
- Pass context through all layers (handlers → services → repositories)

#### Interface Design
- Keep interfaces small and focused
- Define interfaces where they're used, not where they're implemented

### Database Patterns
The project uses the repository pattern for data access:
- Each entity has its own repository interface and implementation
- SQL queries are defined in the repository implementations
- Repositories use prepared statements and parameterized queries
- Migrations are versioned and can be applied/rolled back

### Makefile Commands
The project includes a Makefile with convenient commands:

```
# Show available commands
make help

# Development database
make dev-db          # Start development database
make dev-db-admin    # Start database with pgAdmin
make dev-db-stop     # Stop development database

# Testing
make test            # Run tests with Docker PostgreSQL
make test-short      # Run tests without database

# Build & Quality
make build           # Build all packages
make lint            # Run linter

# Utilities
make deps            # Download and tidy dependencies
make tools           # Install development tools
```

### API Testing
You can test the API endpoints using curl or httpie:

```
# Get profile
curl http://localhost:8080/api/v1/profile

# Get experiences with filtering
curl "http://localhost:8080/api/v1/experiences?company=Derivco&limit=5"

# Get skills by category
curl "http://localhost:8080/api/v1/skills?category=Languages"
```

### Troubleshooting

#### Database Connection Issues
- Check if PostgreSQL is running: `pg_isready`
- Verify connection string in .env
- Check network connectivity and firewall settings

#### Migration Errors
- Check migration file syntax
- Ensure database user has proper permissions
- Verify migration hasn't been partially applied

#### Port Already in Use
```
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```