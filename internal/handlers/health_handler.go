package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles the request to check the health of the service.
// @Summary Health check
// @Description Check if the service is up and running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Service is healthy"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
