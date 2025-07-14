package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryConfig(t *testing.T) {
	// Clean environment
	clearTelemetryEnv()

	t.Run("loads telemetry defaults", func(t *testing.T) {
		config, err := Load()
		require.NoError(t, err)
		require.NotNil(t, config)

		// Check telemetry defaults
		assert.False(t, config.Telemetry.Enabled)
		assert.Equal(t, "resume-api", config.Telemetry.ServiceName)
		assert.Equal(t, "stdout", config.Telemetry.ExporterType)
		assert.Equal(t, "", config.Telemetry.ExporterEndpoint)
		assert.Equal(t, 1.0, config.Telemetry.SamplingRate)
	})

	t.Run("loads telemetry from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("RESUME_API_TELEMETRY_ENABLED", "true")
		os.Setenv("RESUME_API_TELEMETRY_SERVICE_NAME", "test-service")
		os.Setenv("RESUME_API_TELEMETRY_EXPORTER_TYPE", "otlp")
		os.Setenv("RESUME_API_TELEMETRY_EXPORTER_ENDPOINT", "localhost:4317")
		os.Setenv("RESUME_API_TELEMETRY_SAMPLING_RATE", "0.5")
		defer clearTelemetryEnv()

		config, err := Load()
		require.NoError(t, err)

		assert.True(t, config.Telemetry.Enabled)
		assert.Equal(t, "test-service", config.Telemetry.ServiceName)
		assert.Equal(t, "otlp", config.Telemetry.ExporterType)
		assert.Equal(t, "localhost:4317", config.Telemetry.ExporterEndpoint)
		assert.Equal(t, 0.5, config.Telemetry.SamplingRate)
	})

	t.Run("loads redis from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("RESUME_API_REDIS_ENABLED", "true")
		os.Setenv("RESUME_API_REDIS_HOST", "redis-host")
		os.Setenv("RESUME_API_REDIS_PORT", "6380")
		os.Setenv("RESUME_API_REDIS_PASSWORD", "redis-password")
		os.Setenv("RESUME_API_REDIS_DB", "1")
		os.Setenv("RESUME_API_REDIS_TTL", "30m")
		defer clearRedisEnv()

		config, err := Load()
		require.NoError(t, err)

		assert.True(t, config.Redis.Enabled)
		assert.Equal(t, "redis-host", config.Redis.Host)
		assert.Equal(t, 6380, config.Redis.Port)
		assert.Equal(t, "redis-password", config.Redis.Password)
		assert.Equal(t, 1, config.Redis.DB)
		assert.Equal(t, 30*time.Minute, config.Redis.TTL)
	})
}

// Helper function to clear telemetry environment variables
func clearTelemetryEnv() {
	envVars := []string{
		"RESUME_API_TELEMETRY_ENABLED",
		"RESUME_API_TELEMETRY_SERVICE_NAME",
		"RESUME_API_TELEMETRY_EXPORTER_TYPE",
		"RESUME_API_TELEMETRY_EXPORTER_ENDPOINT",
		"RESUME_API_TELEMETRY_SAMPLING_RATE",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}

// Helper function to clear redis environment variables
func clearRedisEnv() {
	envVars := []string{
		"RESUME_API_REDIS_ENABLED",
		"RESUME_API_REDIS_HOST",
		"RESUME_API_REDIS_PORT",
		"RESUME_API_REDIS_PASSWORD",
		"RESUME_API_REDIS_DB",
		"RESUME_API_REDIS_TTL",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
