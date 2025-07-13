package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Global validator instance
var validate = validator.New()

// RateLimiterConfig holds configuration for the rate limiter middleware
type RateLimiterConfig struct {
	// RequestsPerSecond defines the maximum rate of requests per second
	RequestsPerSecond int
	// BurstSize defines the maximum burst size
	BurstSize int
	// TTL defines how long to keep client entries in the limiter map
	TTL time.Duration
}

// DefaultRateLimiterConfig returns a default configuration for the rate limiter
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,    // 10 requests per second
		BurstSize:         20,    // Allow bursts of up to 20 requests
		TTL:               time.Hour, // Clean up client entries after 1 hour
	}
}

// client represents a client in the rate limiter
type client struct {
	tokens     int       // Current token count
	lastAccess time.Time // Last time tokens were added
	lastSeen   time.Time // Last time client was seen
}

// RateLimiterMiddleware returns a middleware that limits the number of requests per client IP
func RateLimiterMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)

	// Start a goroutine to clean up old clients
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > config.TTL {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()

		// Create new client if not exists
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				tokens:     config.BurstSize, // Start with full tokens
				lastAccess: now,
				lastSeen:   now,
			}
		}

		// Update last seen time
		clients[ip].lastSeen = now

		// Calculate tokens to add based on time elapsed
		elapsed := now.Sub(clients[ip].lastAccess).Seconds()
		tokensToAdd := int(elapsed * float64(config.RequestsPerSecond))

		// Update tokens and last access time
		if tokensToAdd > 0 {
			// Use if statement instead of min function
			newTokens := clients[ip].tokens + tokensToAdd
			if newTokens > config.BurstSize {
				newTokens = config.BurstSize
			}
			clients[ip].tokens = newTokens
			clients[ip].lastAccess = now
		}

		// Check if request can be allowed
		if clients[ip].tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		// Consume a token
		clients[ip].tokens--

		mu.Unlock()
		c.Next()
	}
}

// SecurityHeadersMiddleware adds security-related headers to all responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Clickjacking protection
		c.Header("X-Frame-Options", "DENY")

		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// HTTP Strict Transport Security (HSTS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// QueryParams represents common query parameters used across the API
type QueryParams struct {
	Limit  string `form:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Offset string `form:"offset" validate:"omitempty,numeric,min=0"`
}

// InputValidationMiddleware validates request parameters and bodies
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate query parameters
		if err := validateQueryParams(c); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Next()
	}
}

// validateQueryParams performs validation on common query parameters
func validateQueryParams(c *gin.Context) error {
	var params QueryParams

	// Bind query parameters to struct
	if err := c.ShouldBindQuery(&params); err != nil {
		return err
	}

	// Validate using validator
	if err := validate.Struct(params); err != nil {
		return err
	}

	return nil
}
