package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// TableExists checks if a table exists in the database
func (db *DB) TableExists(ctx context.Context, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`
	
	err := db.pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if table exists: %w", err)
	}
	
	return exists, nil
}

// GetTableNames returns a list of all tables in the public schema
func (db *DB) GetTableNames(ctx context.Context) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name`
	
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table names: %w", err)
	}
	defer rows.Close()
	
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	
	return tables, nil
}

// CountRows returns the number of rows in a table
func (db *DB) CountRows(ctx context.Context, tableName string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", pgx.Identifier{tableName}.Sanitize())
	
	err := db.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows in table %s: %w", tableName, err)
	}
	
	return count, nil
}

// TruncateTable truncates a table (removes all rows)
func (db *DB) TruncateTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", pgx.Identifier{tableName}.Sanitize())
	
	_, err := db.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
	}
	
	return nil
}

// GetDatabaseSize returns the size of the database in bytes
func (db *DB) GetDatabaseSize(ctx context.Context) (int64, error) {
	var size int64
	query := "SELECT pg_database_size(current_database())"
	
	err := db.pool.QueryRow(ctx, query).Scan(&size)
	if err != nil {
		return 0, fmt.Errorf("failed to get database size: %w", err)
	}
	
	return size, nil
}

// LogDatabaseInfo logs useful database information
func (db *DB) LogDatabaseInfo(ctx context.Context) {
	// Get database version
	var version string
	if err := db.pool.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		db.logger.Warn("Failed to get database version", "error", err)
	} else {
		db.logger.Info("Database version", "version", version)
	}
	
	// Get database size
	if size, err := db.GetDatabaseSize(ctx); err != nil {
		db.logger.Warn("Failed to get database size", "error", err)
	} else {
		db.logger.Info("Database size", "bytes", size, "mb", size/1024/1024)
	}
	
	// Get table count
	if tables, err := db.GetTableNames(ctx); err != nil {
		db.logger.Warn("Failed to get table names", "error", err)
	} else {
		db.logger.Info("Database tables", "count", len(tables), "tables", tables)
		
		// Log row counts for each table
		for _, table := range tables {
			if count, err := db.CountRows(ctx, table); err != nil {
				db.logger.Warn("Failed to count rows", "table", table, "error", err)
			} else {
				db.logger.Debug("Table row count", "table", table, "rows", count)
			}
		}
	}
	
	// Log connection pool stats
	stats := db.Stats()
	db.logger.Info("Connection pool stats",
		slog.Int("total_conns", int(stats.TotalConns())),
		slog.Int("idle_conns", int(stats.IdleConns())),
		slog.Int("acquired_conns", int(stats.AcquiredConns())),
		slog.Int("max_conns", int(stats.MaxConns())),
	)
}

// IsHealthy performs a quick health check
func (db *DB) IsHealthy(ctx context.Context) bool {
	return db.Ping(ctx) == nil
}