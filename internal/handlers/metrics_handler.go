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
// @Response 200 {string} string "Example response" "# HELP http_requests_total Total number of HTTP requests\n# TYPE http_requests_total counter\nhttp_requests_total{method=\"get\",path=\"/api/v1/profile\",status=\"200\"} 42\nhttp_requests_total{method=\"get\",path=\"/api/v1/experiences\",status=\"200\"} 18\n# HELP http_request_duration_seconds HTTP request duration in seconds\n# TYPE http_request_duration_seconds histogram\nhttp_request_duration_seconds_bucket{method=\"get\",path=\"/api/v1/profile\",le=\"0.1\"} 38\nhttp_request_duration_seconds_bucket{method=\"get\",path=\"/api/v1/profile\",le=\"0.5\"} 42\nhttp_request_duration_seconds_sum{method=\"get\",path=\"/api/v1/profile\"} 1.8\nhttp_request_duration_seconds_count{method=\"get\",path=\"/api/v1/profile\"} 42"
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
