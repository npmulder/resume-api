package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware returns a new logging middleware.
func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"ip", c.ClientIP(),
		)
	}
}
