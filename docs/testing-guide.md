# Testing Guide

## Overview

The Resume API project includes comprehensive testing with Docker integration for easy setup and consistent results across different environments.

## Quick Start

### Prerequisites
- **Docker**: Required for integration tests
- **Go 1.21+**: For running tests
- **Make**: Optional but recommended for convenience

### Run Tests

```bash
# Quick compilation check (no database required)
make test-short

# Full integration tests with Docker (recommended)
make test

# Alternative: Docker Compose approach
make test-compose
```

## Test Types

### 1. Unit Tests (Short Mode)
- **Purpose**: Compilation verification and basic logic testing
- **Duration**: Fast (~1-2 seconds)
- **Requirements**: Go compiler only
- **Command**: `make test-short`

### 2. Integration Tests (Docker)
- **Purpose**: Full repository testing with real PostgreSQL
- **Duration**: Medium (~30-60 seconds)
- **Requirements**: Docker daemon running
- **Command**: `make test`

### 3. Development Testing
- **Purpose**: Continuous testing during development
- **Duration**: Variable
- **Requirements**: Docker Compose
- **Setup**: `make dev-db` then run tests manually

## Available Commands

### Testing Commands
```bash
make test-short      # Quick compilation check
make test           # Full Docker integration tests  
make test-docker    # Same as 'make test'
make test-compose   # Docker Compose approach
```

### Development Commands
```bash
make dev-db         # Start development database
make dev-db-admin   # Start database + pgAdmin UI
make dev-db-stop    # Stop development database
make migrate-up     # Apply database migrations
make migrate-down   # Rollback migrations
```

### Utility Commands
```bash
make build          # Build all packages
make clean          # Clean up Docker resources
make deps           # Update Go dependencies
make lint           # Run code linter
```

## Test Configuration

### Environment Variables
```bash
# Database connection (automatically set by test scripts)
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=resume_api_test
TEST_DB_USER=dev
TEST_DB_PASSWORD=devpass

# Docker configuration
POSTGRES_VERSION=15-alpine  # PostgreSQL Docker image
```

### Docker Configuration
- **Image**: `postgres:15-alpine`
- **Container**: `resume-api-test-db`
- **Port**: `5433` (avoids conflicts with local PostgreSQL)
- **Auto-cleanup**: Yes (containers removed after tests)

## Repository Test Coverage

| Repository | Test Cases | Coverage Areas |
|------------|------------|----------------|
| **Profile** | 7 tests | CRUD, validation, constraints |
| **Experience** | 13 tests | Filtering, pagination, dates |
| **Skill** | 12 tests | Categories, levels, features |
| **Achievement** | 11 tests | Years, categories, metrics |
| **Education** | 13 tests | Types, credentials, expiry |
| **Project** | 14 tests | JSONB, technologies, status |

**Total: 70 comprehensive test cases**

## Test Structure

### Typical Test Flow
1. **Setup**: Clean database state
2. **Create**: Insert test data
3. **Test**: Execute repository operations
4. **Verify**: Assert expected results
5. **Cleanup**: Automatic via test utilities

### Example Test
```go
func TestProfileRepository(t *testing.T) {
    testDB := setupTestDB(t)
    defer testDB.Close()
    
    repo := NewProfileRepository(testDB.Pool())
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    t.Run("CreateProfile", func(t *testing.T) {
        testDB.CleanupTables(t)  // Clean state
        
        profile := &models.Profile{
            Name:  "John Doe",
            Email: "john@example.com",
        }
        
        err := repo.CreateProfile(ctx, profile)
        require.NoError(t, err)
        assert.NotZero(t, profile.ID)
    })
}
```

## Debugging Tests

### View Database During Tests
```bash
# Start development database
make dev-db-admin

# Run migrations
make migrate-up

# Access pgAdmin at http://localhost:5050
# Login: admin@test.com / admin
```

### Manual Database Connection
```bash
# Connect to test database
docker exec -it resume-api-test-db psql -U dev -d resume_api_test

# Or using docker-compose
docker-compose -f docker-compose.test.yml exec test-db psql -U dev -d resume_api_test
```

### View Container Logs
```bash
# Docker container logs
docker logs resume-api-test-db

# Docker Compose logs
docker-compose -f docker-compose.test.yml logs test-db
```

## Troubleshooting

### Common Issues

#### Docker Not Running
```
❌ Docker daemon is not running
```
**Solution**: Start Docker Desktop or Docker daemon

#### Port Conflicts
```
❌ Port 5433 already in use
```
**Solution**: Change `TEST_DB_PORT` environment variable or stop conflicting services

#### Migration Failures
```
❌ Migration failed
```
**Solution**: Check migration files in `./migrations/` directory

#### Test Timeouts
```
❌ Tests timing out
```
**Solution**: Increase timeout in test scripts or check Docker performance

### Reset Everything
```bash
# Nuclear option: clean everything and restart
make clean
docker system prune -f
make test
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Run tests
        run: make test
```

### Local Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit
make test-short
```

## Performance Considerations

### Test Speed Optimization
- **Short tests**: Use for quick feedback during development
- **Parallel tests**: Consider `go test -parallel` for larger test suites
- **Test focus**: Use `go test -run TestSpecific` for targeted testing
- **Docker resources**: Allocate sufficient memory/CPU to Docker

### Database Performance
- **Connection pooling**: Tests use minimal connections (5 max)
- **Table cleanup**: Only clears data, preserves schema
- **Transactions**: Tests use rollback patterns where possible

---

*For more detailed information, see [Repository Testing Summary](./repository-testing-summary.md)*