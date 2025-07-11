package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles the request to check the health of the service.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
