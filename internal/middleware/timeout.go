package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware returns a middleware that cancels the context after the specified timeout.
// If the handler doesn't complete within the timeout, a 408 Request Timeout status is returned.
func TimeoutMiddleware(timeout time.Duration, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update the request with the new context
		c.Request = c.Request.WithContext(ctx)

		// Create a channel to signal when the request is complete
		done := make(chan struct{})
		
		// Process the request in a goroutine
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for the request to complete or timeout
		select {
		case <-done:
			// Request completed before timeout
			return
		case <-ctx.Done():
			// Request timed out
			if ctx.Err() == context.DeadlineExceeded {
				logger.Warn("request timed out",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"timeout", timeout,
				)
				c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
					"error": "Request timed out",
				})
			}
		}
	}
}