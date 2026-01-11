package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all HTTP routes for the API.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
