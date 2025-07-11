package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware returns a new recovery middleware.
func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", "error", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
