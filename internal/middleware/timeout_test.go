package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a logger for testing
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	t.Run("request completes before timeout", func(t *testing.T) {
		// Create a new Gin router
		router := gin.New()

		// Add the timeout middleware with a 500ms timeout
		router.Use(TimeoutMiddleware(500*time.Millisecond, logger))

		// Add a handler that completes quickly (100ms)
		router.GET("/quick", func(c *gin.Context) {
			time.Sleep(100 * time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/quick", nil)
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("request times out", func(t *testing.T) {
		// Create a new Gin router
		router := gin.New()

		// Add the timeout middleware with a 100ms timeout
		router.Use(TimeoutMiddleware(100*time.Millisecond, logger))

		// Add a handler that takes too long (300ms)
		router.GET("/slow", func(c *gin.Context) {
			time.Sleep(300 * time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusRequestTimeout, w.Code)
		assert.Contains(t, w.Body.String(), "timed out")
	})
}