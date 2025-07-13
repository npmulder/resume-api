// Package cache provides caching functionality for the resume API
package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string, dest interface{}) error

	// Set stores a value in the cache with the specified TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Close closes the cache connection
	Close() error
}

// Options defines configuration options for the cache
type Options struct {
	// TTL is the default time-to-live for cache entries
	TTL time.Duration

	// Enabled indicates whether caching is enabled
	Enabled bool
}