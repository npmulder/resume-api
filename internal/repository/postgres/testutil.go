package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/npmulder/resume-api/internal/config"
	"github.com/npmulder/resume-api/internal/database"
)

// TestDB represents a test database connection
type TestDB struct {
	*database.DB
	cleanup func()
}

// setupTestDB creates a test database connection with proper cleanup
func setupTestDB(t *testing.T) *TestDB {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	cfg := &config.DatabaseConfig{
		Host:               getTestEnv("TEST_DB_HOST", "localhost"),
		Port:               5432,
		Name:               getTestEnv("TEST_DB_NAME", "resume_api_test"),
		User:               getTestEnv("TEST_DB_USER", "dev"),
		Password:           getTestEnv("TEST_DB_PASSWORD", "devpass"),
		SSLMode:            "disable",
		MaxConnections:     5,
		MaxIdleConnections: 2,
		ConnMaxLifetime:    30 * time.Minute,
		ConnMaxIdleTime:    5 * time.Minute,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg, nil)
	require.NoError(t, err, "Failed to connect to test database")

	// Test the connection
	err = db.Ping(ctx)
	require.NoError(t, err, "Failed to ping test database")

	return &TestDB{
		DB: db,
		cleanup: func() {
			db.Close()
		},
	}
}

// Close cleans up the test database connection
func (tdb *TestDB) Close() {
	tdb.cleanup()
}

// CleanupTables removes all data from tables for clean test state
func (tdb *TestDB) CleanupTables(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clean tables in correct order due to potential foreign keys
	tables := []string{
		"projects",
		"education", 
		"achievements",
		"skills",
		"experiences",
		"profiles",
	}

	for _, table := range tables {
		_, err := tdb.Pool().Exec(ctx, "DELETE FROM "+table)
		require.NoError(t, err, "Failed to clean table: %s", table)
	}
}

// getTestEnv gets environment variable for tests with fallback
func getTestEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Helper functions for creating test data

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int
func intPtr(i int) *int {
	return &i
}

// boolPtr returns a pointer to the given bool
func boolPtr(b bool) *bool {
	return &b
}

// timePtr returns a pointer to the given time
func timePtr(t time.Time) *time.Time {
	return &t
}