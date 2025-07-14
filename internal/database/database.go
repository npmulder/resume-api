// Package database provides PostgreSQL database connection and management utilities
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/npmulder/resume-api/internal/config"
)

// DB wraps a pgx connection pool with additional functionality
type DB struct {
	pool   *TracedPool
	config *config.DatabaseConfig
	logger *slog.Logger
}

// New creates a new database connection with the given configuration
func New(ctx context.Context, cfg *config.DatabaseConfig, logger *slog.Logger) (*DB, error) {
	if logger == nil {
		logger = slog.Default()
	}

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = 1 // Always keep at least one connection
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Configure connection settings
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"application_name": "resume-api",
	}

	// Set up logging for database connections
	poolConfig.ConnConfig.Tracer = &queryTracer{logger: logger}

	logger.Info("Connecting to database",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.Name),
		slog.String("user", cfg.User),
		slog.Int("max_connections", cfg.MaxConnections),
		slog.Duration("max_lifetime", cfg.ConnMaxLifetime),
	)

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Wrap the pool with tracing
	tracedPool := NewTracedPool(pool)

	db := &DB{
		pool:   tracedPool,
		config: cfg,
		logger: logger,
	}

	// Test the connection
	if err := db.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return db, nil
}

// Pool returns the underlying pgx connection pool
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool.Pool()
}

// Ping tests the database connection
func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.pool.Ping(ctx)
}

// Close closes all connections in the pool
func (db *DB) Close() {
	db.logger.Info("Closing database connections")
	db.pool.Close()
}

// Stats returns connection pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// Health performs a comprehensive health check
func (db *DB) Health(ctx context.Context) (*HealthStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	status := &HealthStatus{
		Timestamp: time.Now(),
	}

	// Test basic connectivity
	start := time.Now()
	if err := db.Ping(ctx); err != nil {
		status.Status = "unhealthy"
		status.Error = err.Error()
		status.ResponseTime = time.Since(start)
		return status, err
	}
	status.ResponseTime = time.Since(start)

	// Get pool statistics
	stats := db.Stats()
	status.Connections = ConnectionStats{
		Total:     int(stats.TotalConns()),
		Idle:      int(stats.IdleConns()),
		Used:      int(stats.AcquiredConns()),
		Maximum:   int(stats.MaxConns()),
		Acquiring: 0, // This field is not available in pgx v5
	}

	// Test a simple query
	var result int
	err := db.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		status.Status = "degraded"
		status.Error = fmt.Sprintf("query test failed: %v", err)
		return status, err
	}

	if result != 1 {
		status.Status = "degraded"
		status.Error = "unexpected query result"
		return status, fmt.Errorf("unexpected query result: %d", result)
	}

	// Check database version
	var version string
	err = db.pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		status.Status = "degraded"
		status.Error = fmt.Sprintf("version check failed: %v", err)
		db.logger.Warn("Failed to get database version", "error", err)
	} else {
		status.Version = version
	}

	status.Status = "healthy"
	return status, nil
}

// HealthStatus represents the health status of the database
type HealthStatus struct {
	Status        string            `json:"status"`
	Timestamp     time.Time         `json:"timestamp"`
	ResponseTime  time.Duration     `json:"response_time"`
	Version       string            `json:"version,omitempty"`
	Connections   ConnectionStats   `json:"connections"`
	Error         string            `json:"error,omitempty"`
}

// ConnectionStats represents connection pool statistics
type ConnectionStats struct {
	Total     int `json:"total"`
	Idle      int `json:"idle"`
	Used      int `json:"used"`
	Maximum   int `json:"maximum"`
	Acquiring int `json:"acquiring"`
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

// WithTx executes a function within a database transaction
func (db *DB) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			db.logger.Error("Failed to rollback transaction",
				"original_error", err,
				"rollback_error", rbErr,
			)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// queryTracer implements pgx.QueryTracer for logging database queries
type queryTracer struct {
	logger *slog.Logger
}

func (t *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// Store start time in context for duration calculation
	return context.WithValue(ctx, "query_start", time.Now())
}

func (t *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// Only log slow queries or errors in production
	startTime, ok := ctx.Value("query_start").(time.Time)
	if !ok {
		return // No start time available
	}
	duration := time.Since(startTime)

	if data.Err != nil {
		t.logger.Error("Database query failed",
			slog.Duration("duration", duration),
			slog.String("error", data.Err.Error()),
		)
	} else if duration > 100*time.Millisecond {
		t.logger.Warn("Slow database query",
			slog.Duration("duration", duration),
		)
	} else {
		t.logger.Debug("Database query executed",
			slog.Duration("duration", duration),
		)
	}
}

// MustNew creates a new database connection and panics if it fails
// Use this in main.go where database failure should stop the application
func MustNew(ctx context.Context, cfg *config.DatabaseConfig, logger *slog.Logger) *DB {
	db, err := New(ctx, cfg, logger)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}
