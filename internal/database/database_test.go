package database

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npmulder/resume-api/internal/config"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	cfg := getTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Reduce noise in tests
	}))

	t.Run("successful connection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db, err := New(ctx, cfg, logger)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()

		// Test basic functionality
		assert.NotNil(t, db.Pool())
		
		// Test ping
		err = db.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("connection with invalid config", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		invalidCfg := *cfg
		invalidCfg.Port = 99999 // Invalid port

		db, err := New(ctx, &invalidCfg, logger)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

func TestDatabaseHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	cfg := getTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := New(ctx, cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	t.Run("health check passes", func(t *testing.T) {
		status, err := db.Health(ctx)
		require.NoError(t, err)
		require.NotNil(t, status)

		assert.Equal(t, "healthy", status.Status)
		assert.NotZero(t, status.Timestamp)
		assert.Greater(t, status.ResponseTime, time.Duration(0))
		assert.NotEmpty(t, status.Version)
		assert.Greater(t, status.Connections.Maximum, 0)
	})

	t.Run("connection stats", func(t *testing.T) {
		stats := db.Stats()
		require.NotNil(t, stats)

		assert.Greater(t, int(stats.MaxConns()), 0)
		assert.GreaterOrEqual(t, int(stats.TotalConns()), 0)
	})
}

func TestDatabaseTransactions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	cfg := getTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := New(ctx, cfg, logger)
	require.NoError(t, err)
	defer db.Close()

	t.Run("successful transaction", func(t *testing.T) {
		err := db.WithTx(ctx, func(tx pgx.Tx) error {
			// Simple test query within transaction
			var result int
			return tx.QueryRow(ctx, "SELECT 1").Scan(&result)
		})
		assert.NoError(t, err)
	})

	t.Run("failed transaction rollback", func(t *testing.T) {
		err := db.WithTx(ctx, func(tx pgx.Tx) error {
			// Force an error to test rollback
			return assert.AnError
		})
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestMustNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database tests in short mode")
	}

	cfg := getTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	t.Run("successful MustNew", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var db *DB
		assert.NotPanics(t, func() {
			db = MustNew(ctx, cfg, logger)
		})
		require.NotNil(t, db)
		defer db.Close()
	})

	t.Run("MustNew panics on failure", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		invalidCfg := *cfg
		invalidCfg.Port = 99999

		assert.Panics(t, func() {
			MustNew(ctx, &invalidCfg, logger)
		})
	})
}

func BenchmarkDatabaseConnection(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping database benchmarks in short mode")
	}

	cfg := getTestConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	
	ctx := context.Background()

	b.Run("ping", func(b *testing.B) {
		db, err := New(ctx, cfg, logger)
		require.NoError(b, err)
		defer db.Close()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := db.Ping(ctx)
			require.NoError(b, err)
		}
	})

	b.Run("simple query", func(b *testing.B) {
		db, err := New(ctx, cfg, logger)
		require.NoError(b, err)
		defer db.Close()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result int
			err := db.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
			require.NoError(b, err)
		}
	})
}

// getTestConfig returns a test database configuration
func getTestConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		Host:               getEnv("TEST_DB_HOST", "localhost"),
		Port:               5432,
		Name:               getEnv("TEST_DB_NAME", "resume_api_test"),
		User:               getEnv("TEST_DB_USER", "dev"),
		Password:           getEnv("TEST_DB_PASSWORD", "devpass"),
		SSLMode:            "disable",
		MaxConnections:     5,  // Lower for tests
		MaxIdleConnections: 2,
		ConnMaxLifetime:    30 * time.Minute,
		ConnMaxIdleTime:    5 * time.Minute,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}