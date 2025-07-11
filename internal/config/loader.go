package config

import (
	"fmt"
	"log/slog"
	"os"
)

// MustLoad loads configuration and panics if it fails
// Use this in main.go where configuration failure should stop the application
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	return config
}

// LoadForTesting loads configuration optimized for testing
// Sets sensible defaults for test environment
func LoadForTesting() *Config {
	// Set test environment if not already set
	if os.Getenv("RESUME_API_ENVIRONMENT") == "" {
		os.Setenv("RESUME_API_ENVIRONMENT", "test")
	}
	
	// Override database name for testing if not set
	if os.Getenv("RESUME_API_DATABASE_NAME") == "" {
		os.Setenv("RESUME_API_DATABASE_NAME", "resume_api_test")
	}
	
	// Use text logging for easier test debugging
	if os.Getenv("RESUME_API_LOGGING_FORMAT") == "" {
		os.Setenv("RESUME_API_LOGGING_FORMAT", "text")
	}
	
	// Use debug logging in tests
	if os.Getenv("RESUME_API_LOGGING_LEVEL") == "" {
		os.Setenv("RESUME_API_LOGGING_LEVEL", "debug")
	}
	
	config, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load test configuration: %v", err))
	}
	
	return config
}

// PrintConfig logs the current configuration (with sensitive data masked)
func PrintConfig(config *Config, logger *slog.Logger) {
	logger.Info("Configuration loaded",
		slog.String("environment", config.Environment),
		slog.String("server_address", config.Server.ServerAddress()),
		slog.Duration("server_read_timeout", config.Server.ReadTimeout),
		slog.Duration("server_write_timeout", config.Server.WriteTimeout),
		slog.String("database_host", config.Database.Host),
		slog.Int("database_port", config.Database.Port),
		slog.String("database_name", config.Database.Name),
		slog.String("database_user", config.Database.User),
		slog.String("database_password", maskPassword(config.Database.Password)),
		slog.String("database_ssl_mode", config.Database.SSLMode),
		slog.Int("database_max_connections", config.Database.MaxConnections),
		slog.String("logging_level", config.Logging.Level),
		slog.String("logging_format", config.Logging.Format),
	)
}

// maskPassword masks a password for logging
func maskPassword(password string) string {
	if len(password) <= 2 {
		return "***"
	}
	return password[:1] + "***" + password[len(password)-1:]
}

// ValidateForProduction performs additional validation for production environment
func ValidateForProduction(config *Config) error {
	if !config.IsProduction() {
		return nil
	}
	
	// Production-specific validations
	if config.Database.Password == "devpass" || config.Database.Password == "password" {
		return fmt.Errorf("insecure database password detected in production")
	}
	
	if config.Database.SSLMode == "disable" {
		return fmt.Errorf("SSL must be enabled in production")
	}
	
	if config.Logging.Level == "debug" {
		return fmt.Errorf("debug logging should not be used in production")
	}
	
	if config.Server.Host == "localhost" || config.Server.Host == "127.0.0.1" {
		return fmt.Errorf("server should bind to 0.0.0.0 in production, not localhost")
	}
	
	return nil
}

// GetDatabaseDSN returns a database DSN for external tools (like migration scripts)
// This is a convenience function for backward compatibility
func GetDatabaseDSN() string {
	// Try to get from environment variable first (for migration scripts)
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	
	// Fall back to loading from config
	config, err := Load()
	if err != nil {
		// Return default DSN if config loading fails
		return "postgres://dev:devpass@localhost:5432/resume_api_dev?sslmode=disable"
	}
	
	return config.Database.DatabaseURL()
}