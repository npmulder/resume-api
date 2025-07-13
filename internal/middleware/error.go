package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/utils"
)

// ErrorHandlerMiddleware returns a middleware that handles errors and panics
// with a consistent error response format
func ErrorHandlerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from any panics and return a 500 error
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
					models.WithDetails(map[string]interface{}{
						"error": err,
					}),
				)
			}
		}()

		// Process the request
		c.Next()

		// If there were any errors during the request handling,
		// ensure they are returned in the standardized format
		if len(c.Errors) > 0 {
			// Log the errors
			for _, e := range c.Errors {
				logger.Error("request error",
					"error", e.Err,
					"meta", e.Meta,
					"type", e.Type,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)
			}

			// If no response has been written yet, return a 500 error
			if !c.Writer.Written() {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error",
					models.WithCode(models.ErrCodeInternalError),
					models.WithDetails(map[string]interface{}{
						"errors": c.Errors.Errors(),
					}),
				)
			}
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
// This is useful for tracing requests across logs and error responses
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID from header or generate a new one
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}

		// Set the request ID in the context
		c.Set("RequestID", requestID)

		// Add the request ID to the response headers
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}