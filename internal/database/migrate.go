package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/npmulder/resume-api/internal/config"
)

// MigrateUp runs all up migrations
func MigrateUp(cfg *config.DatabaseConfig, logger *slog.Logger) error {
	m, err := createMigrator(cfg)
	if err != nil {
		return err
	}
	defer m.Close()

	logger.Info("Running database migrations up")
	
	err = m.Up()
	if err == migrate.ErrNoChange {
		logger.Info("No migrations to apply")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		logger.Warn("Failed to get migration version after up", "error", err)
	} else {
		logger.Info("Migrations completed",
			slog.Uint64("version", uint64(version)),
			slog.Bool("dirty", dirty),
		)
	}

	return nil
}

// MigrateDown runs all down migrations
func MigrateDown(cfg *config.DatabaseConfig, logger *slog.Logger) error {
	m, err := createMigrator(cfg)
	if err != nil {
		return err
	}
	defer m.Close()

	logger.Warn("Running database migrations down - this will destroy data!")
	
	err = m.Down()
	if err == migrate.ErrNoChange {
		logger.Info("No migrations to rollback")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to run migrations down: %w", err)
	}

	logger.Info("All migrations rolled back")
	return nil
}

// MigrateVersion gets the current migration version
func MigrateVersion(cfg *config.DatabaseConfig) (uint, bool, error) {
	m, err := createMigrator(cfg)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	return m.Version()
}

// MigrateSteps runs a specific number of migration steps
func MigrateSteps(cfg *config.DatabaseConfig, steps int, logger *slog.Logger) error {
	m, err := createMigrator(cfg)
	if err != nil {
		return err
	}
	defer m.Close()

	if steps > 0 {
		logger.Info("Running migrations up", slog.Int("steps", steps))
	} else {
		logger.Warn("Running migrations down", slog.Int("steps", -steps))
	}
	
	err = m.Steps(steps)
	if err == migrate.ErrNoChange {
		logger.Info("No migrations to apply")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		logger.Warn("Failed to get migration version after steps", "error", err)
	} else {
		logger.Info("Migration steps completed",
			slog.Uint64("version", uint64(version)),
			slog.Bool("dirty", dirty),
		)
	}

	return nil
}

// createMigrator creates a migrate instance
func createMigrator(cfg *config.DatabaseConfig) (*migrate.Migrate, error) {
	return migrate.New(
		"file://migrations",
		cfg.DatabaseURL(),
	)
}

// WaitForDatabase waits for the database to be available
func WaitForDatabase(ctx context.Context, cfg *config.DatabaseConfig, logger *slog.Logger) error {
	logger.Info("Waiting for database to be available")
	
	// Create a simple connection to test availability
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			db, err := New(ctx, cfg, logger)
			if err != nil {
				logger.Debug("Database not yet available", "error", err)
				continue
			}
			
			if err := db.Ping(ctx); err != nil {
				db.Close()
				logger.Debug("Database ping failed", "error", err)
				continue
			}
			
			db.Close()
			logger.Info("Database is available")
			return nil
		}
	}
}

// EnsureMigrations ensures that migrations are up to date
// This is useful for applications that should auto-migrate on startup
func EnsureMigrations(cfg *config.DatabaseConfig, logger *slog.Logger) error {
	logger.Info("Ensuring database migrations are up to date")
	
	// Check current version
	version, dirty, err := MigrateVersion(cfg)
	if err != nil {
		if err == migrate.ErrNilVersion {
			logger.Info("No migrations applied yet, running initial migration")
			return MigrateUp(cfg, logger)
		}
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	
	if dirty {
		return fmt.Errorf("database is in dirty state at version %d - manual intervention required", version)
	}
	
	logger.Info("Current migration version", slog.Uint64("version", uint64(version)))
	
	// Try to apply any pending migrations
	err = MigrateUp(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to apply pending migrations: %w", err)
	}
	
	return nil
}