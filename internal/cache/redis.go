// Package cache provides caching functionality for the resume API
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/npmulder/resume-api/internal/config"
)

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = errors.New("cache miss")

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	if !cfg.Enabled {
		return nil, errors.New("redis cache is disabled")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Set stores a value in the cache with the specified TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for cache: %w", err)
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// Close closes the Redis client connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// NoOpCache is a cache implementation that does nothing
// Used when caching is disabled
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns ErrCacheMiss
func (c *NoOpCache) Get(ctx context.Context, key string, dest interface{}) error {
	return ErrCacheMiss
}

// Set does nothing and returns nil
func (c *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete does nothing and returns nil
func (c *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

// Close does nothing and returns nil
func (c *NoOpCache) Close() error {
	return nil
}

// New creates a new cache based on the configuration
func New(cfg *config.RedisConfig) (Cache, error) {
	if !cfg.Enabled {
		return NewNoOpCache(), nil
	}
	return NewRedisCache(cfg)
}