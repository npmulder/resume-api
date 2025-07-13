package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns a handler that exposes Prometheus metrics
// @Summary Prometheus metrics
// @Description Expose Prometheus metrics for monitoring
// @Tags metrics
// @Accept json
// @Produce text/plain
// @Success 200 {string} string "Prometheus metrics in text format"
// @Router /metrics [get]
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
