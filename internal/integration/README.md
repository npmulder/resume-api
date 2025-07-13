# Integration Tests

This directory contains integration tests for the Resume API. These tests verify the end-to-end functionality of the API, from HTTP handlers through services to repositories and back, using a real database.

## Test Coverage

The integration tests cover:

- **Profile**: End-to-end API testing with database
- **Experiences**: Filtering, pagination, date ranges
- **Skills**: Categories, levels, featured skills
- **Achievements**: Years, categories, impact metrics
- **Education**: Types, statuses, credentials
- **Projects**: Technologies, status filtering

## Running the Tests

### Using the Makefile (Recommended)

The easiest way to run the integration tests is to use the Makefile target:

```bash
make test-integration
```

This target:
1. Sets up a PostgreSQL Docker container for testing
2. Runs database migrations
3. Sets environment variables for tests
4. Runs the integration tests
5. Cleans up the container when done

### Using the Test Script (Alternative)

Alternatively, you can use the provided shell script:

```bash
./scripts/test-integration.sh
```

This script provides the same functionality as the Makefile target.

### Manual Testing

If you prefer to run the tests manually:

1. Ensure you have a PostgreSQL database available for testing
2. Set the following environment variables:
   ```bash
   export TEST_DB_HOST=localhost
   export TEST_DB_PORT=5433
   export TEST_DB_NAME=resume_api_test
   export TEST_DB_USER=dev
   export TEST_DB_PASSWORD=devpass
   ```
3. Run the migrations:
   ```bash
   go run ./cmd/migrate/main.go up
   ```
4. Run the tests:
   ```bash
   go test ./internal/integration/ -v
   ```

## Test Structure

The integration tests use a real database and test the entire flow from HTTP handlers through services to repositories and back. The tests:

1. Set up a test database connection
2. Create a test application with real repositories, services, and handlers
3. Create test data in the database
4. Make HTTP requests to the API endpoints
5. Verify the responses

This approach ensures that all components of the system work together correctly.

## Adding New Tests

To add new integration tests:

1. Add a new test function to `integration_test.go`
2. Follow the pattern of existing tests:
   - Set up a test database connection
   - Create test data
   - Make HTTP requests
   - Verify responses
3. Run the tests to ensure they pass
