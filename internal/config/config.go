// Package config provides configuration management for the resume API
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration
type Config struct {
	Environment string         `mapstructure:"environment" validate:"required,oneof=development production test"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Logging     LoggingConfig  `mapstructure:"logging"`
	Redis       RedisConfig    `mapstructure:"redis"`
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

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	// Set up Viper
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure Viper to read from environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("RESUME_API")

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
	}

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
