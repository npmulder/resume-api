package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/utils"
)

// RecoveryMiddleware returns a new recovery middleware.
func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				stack := string(debug.Stack())
				logger.Error("panic recovered",
					"error", err,
					"stack", stack,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// Create a standardized error response
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error",
					models.WithCode(models.ErrCodeInternalError),
				)
			}
		}()
		c.Next()
	}
}
