package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/warley004/packing-optimizer-api/internal/api/http/handlers"
)

// RegisterRoutes registers all HTTP routes for the API.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")
	{
		packingHandler := handlers.NewPackingHandler()
		v1.POST("/packing", packingHandler.Pack)
	}
}
