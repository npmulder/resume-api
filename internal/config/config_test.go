package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Clean environment
	clearEnv()
	
	t.Run("loads default configuration", func(t *testing.T) {
		config, err := Load()
		require.NoError(t, err)
		require.NotNil(t, config)
		
		// Check defaults
		assert.Equal(t, "development", config.Environment)
		assert.Equal(t, "localhost", config.Server.Host)
		assert.Equal(t, 8080, config.Server.Port)
		assert.Equal(t, 15*time.Second, config.Server.ReadTimeout)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, 5432, config.Database.Port)
		assert.Equal(t, "resume_api_dev", config.Database.Name)
		assert.Equal(t, "info", config.Logging.Level)
		assert.Equal(t, "json", config.Logging.Format)
	})
	
	t.Run("loads from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("RESUME_API_ENVIRONMENT", "production")
		os.Setenv("RESUME_API_SERVER_PORT", "9000")
		os.Setenv("RESUME_API_DATABASE_NAME", "resume_api_prod")
		os.Setenv("RESUME_API_LOGGING_LEVEL", "error")
		defer clearEnv()
		
		config, err := Load()
		require.NoError(t, err)
		
		assert.Equal(t, "production", config.Environment)
		assert.Equal(t, 9000, config.Server.Port)
		assert.Equal(t, "resume_api_prod", config.Database.Name)
		assert.Equal(t, "error", config.Logging.Level)
	})
	
	t.Run("validates configuration", func(t *testing.T) {
		os.Setenv("RESUME_API_ENVIRONMENT", "invalid")
		defer clearEnv()
		
		_, err := Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid environment")
	})
}

func TestDatabaseURL(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}
	
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expected, config.DatabaseURL())
}

func TestServerAddress(t *testing.T) {
	config := &ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}
	
	assert.Equal(t, "0.0.0.0:8080", config.ServerAddress())
}

func TestEnvironmentHelpers(t *testing.T) {
	tests := []struct {
		env         string
		isDev       bool
		isProd      bool
		isTest      bool
	}{
		{"development", true, false, false},
		{"production", false, true, false},
		{"test", false, false, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := &Config{Environment: tt.env}
			assert.Equal(t, tt.isDev, config.IsDevelopment())
			assert.Equal(t, tt.isProd, config.IsProduction())
			assert.Equal(t, tt.isTest, config.IsTest())
		})
	}
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		config := &Config{
			Environment: "development",
			Server: ServerConfig{
				Port: 8080,
			},
			Database: DatabaseConfig{
				Port:               5432,
				SSLMode:            "disable",
				MaxConnections:     10,
				MaxIdleConnections: 5,
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "json",
			},
		}
		
		err := validateConfig(config)
		assert.NoError(t, err)
	})
	
	t.Run("invalid server port", func(t *testing.T) {
		config := &Config{
			Environment: "development",
			Server: ServerConfig{
				Port: 0,
			},
			Database: DatabaseConfig{
				Port:               5432,
				SSLMode:            "disable",
				MaxConnections:     10,
				MaxIdleConnections: 5,
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "json",
			},
		}
		
		err := validateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server port")
	})
	
	t.Run("invalid idle connections", func(t *testing.T) {
		config := &Config{
			Environment: "development",
			Server: ServerConfig{
				Port: 8080,
			},
			Database: DatabaseConfig{
				Port:               5432,
				SSLMode:            "disable",
				MaxConnections:     5,
				MaxIdleConnections: 10, // More than max connections
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "json",
			},
		}
		
		err := validateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max_idle_connections cannot be greater than max_connections")
	})
}

// Helper function to clear environment variables
func clearEnv() {
	envVars := []string{
		"RESUME_API_ENVIRONMENT",
		"RESUME_API_SERVER_HOST",
		"RESUME_API_SERVER_PORT",
		"RESUME_API_SERVER_READ_TIMEOUT",
		"RESUME_API_SERVER_WRITE_TIMEOUT",
		"RESUME_API_SERVER_IDLE_TIMEOUT",
		"RESUME_API_SERVER_GRACEFUL_STOP",
		"RESUME_API_DATABASE_HOST",
		"RESUME_API_DATABASE_PORT",
		"RESUME_API_DATABASE_NAME",
		"RESUME_API_DATABASE_USER",
		"RESUME_API_DATABASE_PASSWORD",
		"RESUME_API_DATABASE_SSL_MODE",
		"RESUME_API_DATABASE_MAX_CONNECTIONS",
		"RESUME_API_DATABASE_MAX_IDLE_CONNECTIONS",
		"RESUME_API_DATABASE_CONN_MAX_LIFETIME",
		"RESUME_API_DATABASE_CONN_MAX_IDLE_TIME",
		"RESUME_API_LOGGING_LEVEL",
		"RESUME_API_LOGGING_FORMAT",
	}
	
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}