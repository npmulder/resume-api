#!/bin/bash

# Repository Test Runner Script with Docker
# This script runs repository tests using a PostgreSQL Docker container

set -e

echo "🧪 Resume API Repository Tests (Docker)"
echo "======================================"

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed or not in PATH"
    echo "   Please install Docker to run integration tests"
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running"
    echo "   Please start Docker and try again"
    exit 1
fi

# Configuration
POSTGRES_VERSION=${POSTGRES_VERSION:-15-alpine}
CONTAINER_NAME="resume-api-test-db"
TEST_DB_HOST="localhost"
TEST_DB_PORT=${TEST_DB_PORT:-5433}  # Use different port to avoid conflicts
TEST_DB_NAME=${TEST_DB_NAME:-resume_api_test}
TEST_DB_USER=${TEST_DB_USER:-dev}
TEST_DB_PASSWORD=${TEST_DB_PASSWORD:-devpass}

echo "📋 Docker PostgreSQL Test Configuration:"
echo "   PostgreSQL Version: $POSTGRES_VERSION"
echo "   Container Name: $CONTAINER_NAME"
echo "   Host: $TEST_DB_HOST"
echo "   Port: $TEST_DB_PORT"
echo "   Database: $TEST_DB_NAME"
echo "   User: $TEST_DB_USER"

# Function to cleanup container
cleanup() {
    echo "🧹 Cleaning up Docker container..."
    docker stop $CONTAINER_NAME &> /dev/null || true
    docker rm $CONTAINER_NAME &> /dev/null || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Stop and remove existing container if it exists
echo "🔄 Checking for existing test container..."
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "   Removing existing container: $CONTAINER_NAME"
    docker stop $CONTAINER_NAME &> /dev/null || true
    docker rm $CONTAINER_NAME &> /dev/null || true
fi

# Start PostgreSQL container
echo "🐳 Starting PostgreSQL Docker container..."
docker run -d \
    --name $CONTAINER_NAME \
    -e POSTGRES_DB=$TEST_DB_NAME \
    -e POSTGRES_USER=$TEST_DB_USER \
    -e POSTGRES_PASSWORD=$TEST_DB_PASSWORD \
    -p $TEST_DB_PORT:5432 \
    postgres:$POSTGRES_VERSION

echo "   ✅ Container started: $CONTAINER_NAME"

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
MAX_ATTEMPTS=30
ATTEMPT=1

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
    if docker exec $CONTAINER_NAME pg_isready -U $TEST_DB_USER -d $TEST_DB_NAME &> /dev/null; then
        echo "   ✅ PostgreSQL is ready (attempt $ATTEMPT/$MAX_ATTEMPTS)"
        break
    fi
    
    if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
        echo "   ❌ PostgreSQL failed to start after $MAX_ATTEMPTS attempts"
        echo "   Container logs:"
        docker logs $CONTAINER_NAME
        exit 1
    fi
    
    echo "   ⏳ Attempt $ATTEMPT/$MAX_ATTEMPTS - waiting for PostgreSQL..."
    sleep 2
    ATTEMPT=$((ATTEMPT + 1))
done

# Verify database connection
echo "🔗 Testing database connection..."
if docker exec $CONTAINER_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME -c "SELECT version();" &> /dev/null; then
    echo "   ✅ Database connection successful"
else
    echo "   ❌ Cannot connect to database"
    echo "   Container logs:"
    docker logs $CONTAINER_NAME
    exit 1
fi

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
                docker exec -i $CONTAINER_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME < "$migration"
            fi
        done
        echo "   ✅ Manual migrations completed"
    else
        echo "   ⚠️  No migrations directory found"
    fi
fi

# Display container info
echo "📊 Container Information:"
echo "   Container ID: $(docker ps --format '{{.ID}}' --filter name=$CONTAINER_NAME)"
echo "   Image: postgres:$POSTGRES_VERSION"
echo "   Status: $(docker ps --format '{{.Status}}' --filter name=$CONTAINER_NAME)"

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
echo "🐳 Docker Container Information:"
echo "   Container will be automatically removed when script exits"
echo "   To manually connect to the test database:"
echo "   docker exec -it $CONTAINER_NAME psql -U $TEST_DB_USER -d $TEST_DB_NAME"

echo ""
echo "🎉 Repository layer testing complete!"

# Container cleanup happens automatically via trap