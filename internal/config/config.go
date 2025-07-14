// Package config provides configuration management for the resume API
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration
type Config struct {
	Environment string          `mapstructure:"environment" validate:"required,oneof=development production test"`
	Server      ServerConfig    `mapstructure:"server"`
	Database    DatabaseConfig  `mapstructure:"database"`
	Logging     LoggingConfig   `mapstructure:"logging"`
	Redis       RedisConfig     `mapstructure:"redis"`
	Telemetry   TelemetryConfig `mapstructure:"telemetry"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port" validate:"min=1,max=65535"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	GracefulStop   time.Duration `mapstructure:"graceful_stop"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host               string        `mapstructure:"host" validate:"required"`
	Port               int           `mapstructure:"port" validate:"min=1,max=65535"`
	Name               string        `mapstructure:"name" validate:"required"`
	User               string        `mapstructure:"user" validate:"required"`
	Password           string        `mapstructure:"password" validate:"required"`
	SSLMode            string        `mapstructure:"ssl_mode" validate:"oneof=disable require verify-ca verify-full"`
	MaxConnections     int           `mapstructure:"max_connections" validate:"min=1"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections" validate:"min=1"`
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `mapstructure:"conn_max_idle_time"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level" validate:"oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"oneof=json text"`
}

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host" validate:"required"`
	Port     int           `mapstructure:"port" validate:"min=1,max=65535"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db" validate:"min=0"`
	TTL      time.Duration `mapstructure:"ttl"`
	Enabled  bool          `mapstructure:"enabled"`
}

// TelemetryConfig contains OpenTelemetry configuration
type TelemetryConfig struct {
	Enabled          bool    `mapstructure:"enabled"`
	ServiceName      string  `mapstructure:"service_name" validate:"required_if=Enabled true"`
	ExporterType     string  `mapstructure:"exporter_type" validate:"required_if=Enabled true,oneof=stdout otlp"`
	ExporterEndpoint string  `mapstructure:"exporter_endpoint"`
	SamplingRate     float64 `mapstructure:"sampling_rate" validate:"min=0,max=1"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Set up Viper
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Try to read from .env file (optional)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Read config file if it exists (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Load environment variables from the .env file
		// This is a workaround for viper not correctly reading environment variables from the .env file
		loadEnvFromFile(v.ConfigFileUsed())
	}

	// Configure Viper to read from environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("RESUME_API")

	// Explicitly bind environment variables for telemetry
	// This ensures that environment variables are properly mapped to the configuration struct
	bindEnvVariables(v)

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// loadEnvFromFile loads environment variables from a .env file
func loadEnvFromFile(filePath string) {
	// Read the .env file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// Parse the .env file line by line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Skip comments and empty lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove comments from the value
		if idx := strings.Index(value, "#"); idx != -1 {
			value = strings.TrimSpace(value[:idx])
		}

		// Set the environment variable
		os.Setenv(key, value)
	}
}

// bindEnvVariables explicitly binds environment variables to viper keys
func bindEnvVariables(v *viper.Viper) {
	// Bind telemetry environment variables
	v.BindEnv("telemetry.enabled", "RESUME_API_TELEMETRY_ENABLED")
	v.BindEnv("telemetry.service_name", "RESUME_API_TELEMETRY_SERVICE_NAME")
	v.BindEnv("telemetry.exporter_type", "RESUME_API_TELEMETRY_EXPORTER_TYPE")
	v.BindEnv("telemetry.exporter_endpoint", "RESUME_API_TELEMETRY_EXPORTER_ENDPOINT")
	v.BindEnv("telemetry.sampling_rate", "RESUME_API_TELEMETRY_SAMPLING_RATE")
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Environment
	v.SetDefault("environment", "development")

	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.graceful_stop", "30s")
	v.SetDefault("server.request_timeout", "10s")

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.name", "resume_api_dev")
	v.SetDefault("database.user", "dev")
	v.SetDefault("database.password", "devpass")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_connections", 25)
	v.SetDefault("database.max_idle_connections", 5)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "30m")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.ttl", "15m")
	v.SetDefault("redis.enabled", true)

	// Telemetry defaults
	v.SetDefault("telemetry.enabled", false)
	v.SetDefault("telemetry.service_name", "resume-api")
	v.SetDefault("telemetry.exporter_type", "stdout")
	v.SetDefault("telemetry.exporter_endpoint", "")
	v.SetDefault("telemetry.sampling_rate", 1.0) // 100% sampling by default
}

// validateConfig performs basic validation on the configuration
func validateConfig(config *Config) error {
	// Validate environment
	validEnvs := map[string]bool{
		"development": true,
		"production":  true,
		"test":        true,
	}
	if !validEnvs[config.Environment] {
		return fmt.Errorf("invalid environment: %s (must be development, production, or test)", config.Environment)
	}

	// Validate server port
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be between 1 and 65535)", config.Server.Port)
	}

	// Validate database port
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d (must be between 1 and 65535)", config.Database.Port)
	}

	// Validate SSL mode
	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validSSLModes[config.Database.SSLMode] {
		return fmt.Errorf("invalid SSL mode: %s", config.Database.SSLMode)
	}

	// Validate logging level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[config.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	// Validate logging format
	validLogFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validLogFormats[config.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", config.Logging.Format)
	}

	// Validate database connection settings
	if config.Database.MaxConnections < 1 {
		return fmt.Errorf("max_connections must be at least 1")
	}
	if config.Database.MaxIdleConnections < 1 {
		return fmt.Errorf("max_idle_connections must be at least 1")
	}
	if config.Database.MaxIdleConnections > config.Database.MaxConnections {
		return fmt.Errorf("max_idle_connections cannot be greater than max_connections")
	}

	// Validate Redis configuration if enabled
	if config.Redis.Enabled {
		if config.Redis.Port < 1 || config.Redis.Port > 65535 {
			return fmt.Errorf("invalid redis port: %d (must be between 1 and 65535)", config.Redis.Port)
		}
		if config.Redis.DB < 0 {
			return fmt.Errorf("redis db must be non-negative")
		}
		if config.Redis.TTL < time.Second {
			return fmt.Errorf("redis ttl must be at least 1 second")
		}
	}

	// Validate Telemetry configuration if enabled
	if config.Telemetry.Enabled {
		if config.Telemetry.ServiceName == "" {
			return fmt.Errorf("telemetry service_name is required when telemetry is enabled")
		}

		validExporterTypes := map[string]bool{
			"stdout": true,
			"otlp":   true,
		}
		if !validExporterTypes[config.Telemetry.ExporterType] {
			return fmt.Errorf("invalid telemetry exporter_type: %s (must be one of: stdout, otlp)", config.Telemetry.ExporterType)
		}

		// For exporters other than stdout, endpoint is required
		if config.Telemetry.ExporterType != "stdout" && config.Telemetry.ExporterEndpoint == "" {
			return fmt.Errorf("telemetry exporter_endpoint is required for exporter type: %s", config.Telemetry.ExporterType)
		}

		if config.Telemetry.SamplingRate < 0 || config.Telemetry.SamplingRate > 1 {
			return fmt.Errorf("telemetry sampling_rate must be between 0 and 1, got: %f", config.Telemetry.SamplingRate)
		}
	}

	return nil
}

// DatabaseURL returns a formatted PostgreSQL connection string
func (c *DatabaseConfig) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

// ServerAddress returns the formatted server address
func (c *ServerConfig) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTest returns true if running in test mode
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

// RedisURL returns a formatted Redis connection string
func (c *RedisConfig) RedisURL() string {
	if c.Password == "" {
		return fmt.Sprintf("redis://%s:%d/%d", c.Host, c.Port, c.DB)
	}
	return fmt.Sprintf("redis://:%s@%s:%d/%d", c.Password, c.Host, c.Port, c.DB)
}
