#!/bin/bash

# Repository Test Runner Script with Docker Compose
# This script runs repository tests using docker-compose for PostgreSQL

set -e

echo "🧪 Resume API Repository Tests (Docker Compose)"
echo "=============================================="

# Check if Docker and Docker Compose are available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed or not in PATH"
    echo "   Please install Docker to run integration tests"
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null 2>&1; then
    echo "❌ Docker Compose is not installed or not in PATH"
    echo "   Please install Docker Compose to run integration tests"
    echo "   Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

# Determine which compose command to use
COMPOSE_CMD="docker-compose"
if docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
fi

echo "📋 Using compose command: $COMPOSE_CMD"

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running"
    echo "   Please start Docker and try again"
    exit 1
fi

# Configuration
COMPOSE_FILE="docker-compose.test.yml"
SERVICE_NAME="test-db"
TEST_DB_HOST="localhost"
TEST_DB_PORT="5433"
TEST_DB_NAME="resume_api_test"
TEST_DB_USER="dev"
TEST_DB_PASSWORD="devpass"

echo "📋 Docker Compose Test Configuration:"
echo "   Compose File: $COMPOSE_FILE"
echo "   Service: $SERVICE_NAME"
echo "   Host: $TEST_DB_HOST"
echo "   Port: $TEST_DB_PORT"
echo "   Database: $TEST_DB_NAME"
echo "   User: $TEST_DB_USER"

# Function to cleanup
cleanup() {
    echo "🧹 Cleaning up Docker Compose services..."
    $COMPOSE_CMD -f $COMPOSE_FILE down -v &> /dev/null || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Check if compose file exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo "❌ Docker Compose file not found: $COMPOSE_FILE"
    echo "   Please ensure $COMPOSE_FILE exists in the project root"
    exit 1
fi

# Stop any existing services
echo "🔄 Stopping any existing test services..."
$COMPOSE_CMD -f $COMPOSE_FILE down -v &> /dev/null || true

# Start PostgreSQL service
echo "🐳 Starting PostgreSQL service with Docker Compose..."
$COMPOSE_CMD -f $COMPOSE_FILE up -d $SERVICE_NAME

# Wait for PostgreSQL to be healthy
echo "⏳ Waiting for PostgreSQL to be healthy..."
MAX_ATTEMPTS=60
ATTEMPT=1

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
    if $COMPOSE_CMD -f $COMPOSE_FILE ps $SERVICE_NAME | grep -q "healthy"; then
        echo "   ✅ PostgreSQL is healthy (attempt $ATTEMPT/$MAX_ATTEMPTS)"
        break
    fi
    
    if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
        echo "   ❌ PostgreSQL failed to become healthy after $MAX_ATTEMPTS attempts"
        echo "   Service logs:"
        $COMPOSE_CMD -f $COMPOSE_FILE logs $SERVICE_NAME
        exit 1
    fi
    
    echo "   ⏳ Attempt $ATTEMPT/$MAX_ATTEMPTS - waiting for PostgreSQL to be healthy..."
    sleep 2
    ATTEMPT=$((ATTEMPT + 1))
done

# Verify database connection
echo "🔗 Testing database connection..."
if $COMPOSE_CMD -f $COMPOSE_FILE exec -T $SERVICE_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME -c "SELECT version();" &> /dev/null; then
    echo "   ✅ Database connection successful"
else
    echo "   ❌ Cannot connect to database"
    echo "   Service logs:"
    $COMPOSE_CMD -f $COMPOSE_FILE logs $SERVICE_NAME
    exit 1
fi

# Show service status
echo "📊 Service Status:"
$COMPOSE_CMD -f $COMPOSE_FILE ps

# Run migrations
echo "🔄 Running database migrations..."
export TEST_DB_HOST=$TEST_DB_HOST
export TEST_DB_PORT=$TEST_DB_PORT
export TEST_DB_NAME=$TEST_DB_NAME
export TEST_DB_USER=$TEST_DB_USER
export TEST_DB_PASSWORD=$TEST_DB_PASSWORD

# Update the DATABASE_URL to use the test container
export DATABASE_URL="postgres://$TEST_DB_USER:$TEST_DB_PASSWORD@$TEST_DB_HOST:$TEST_DB_PORT/$TEST_DB_NAME?sslmode=disable"

if [ -f "./cmd/migrate/main.go" ]; then
    echo "   Running migrations with: $DATABASE_URL"
    go run ./cmd/migrate/main.go up
    echo "   ✅ Migrations completed successfully"
else
    echo "   ⚠️  Migration tool not found at ./cmd/migrate/main.go"
    echo "   Creating tables manually..."
    
    # Run migrations directly if migrate tool not found
    if [ -d "./migrations" ]; then
        for migration in ./migrations/*.up.sql; do
            if [ -f "$migration" ]; then
                echo "   Applying migration: $(basename $migration)"
                $COMPOSE_CMD -f $COMPOSE_FILE exec -T $SERVICE_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME < "$migration"
            fi
        done
        echo "   ✅ Manual migrations completed"
    else
        echo "   ⚠️  No migrations directory found"
    fi
fi

# Run tests
echo ""
echo "🧪 Running repository tests..."
echo "=========================================="

# Set environment variables for tests
export TEST_DB_HOST=$TEST_DB_HOST
export TEST_DB_PORT=$TEST_DB_PORT
export TEST_DB_NAME=$TEST_DB_NAME
export TEST_DB_USER=$TEST_DB_USER
export TEST_DB_PASSWORD=$TEST_DB_PASSWORD

# Run tests with verbose output and proper timeout
echo "Environment variables set:"
echo "   TEST_DB_HOST=$TEST_DB_HOST"
echo "   TEST_DB_PORT=$TEST_DB_PORT"
echo "   TEST_DB_NAME=$TEST_DB_NAME"
echo "   TEST_DB_USER=$TEST_DB_USER"
echo ""

if go test ./internal/repository/postgres/ -v -timeout=120s; then
    echo ""
    echo "✅ All repository tests completed successfully!"
else
    echo ""
    echo "❌ Some tests failed!"
    exit 1
fi

echo ""
echo "📊 Test Coverage Summary:"
echo "   • ProfileRepository: CRUD operations, validation, error handling"
echo "   • ExperienceRepository: Filtering, pagination, date ranges"
echo "   • SkillRepository: Categories, levels, featured skills"
echo "   • AchievementRepository: Years, categories, impact metrics"
echo "   • EducationRepository: Types, statuses, credentials"
echo "   • ProjectRepository: JSONB technologies, status filtering"

echo ""
echo "🐳 Docker Compose Information:"
echo "   Services will be automatically stopped and removed when script exits"
echo "   To manually connect to the test database:"
echo "   $COMPOSE_CMD -f $COMPOSE_FILE exec $SERVICE_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME"
echo ""
echo "   To start pgAdmin for database inspection:"
echo "   $COMPOSE_CMD -f $COMPOSE_FILE --profile admin up -d pgadmin"
echo "   Then visit: http://localhost:5050 (admin@test.com / admin)"

echo ""
echo "🎉 Repository layer testing complete!"

# Services cleanup happens automatically via trap