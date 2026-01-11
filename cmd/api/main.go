package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/warley004/packing-optimizer-api/internal/api/http"
)

func main() {
	// Gin mode: debug (default) or release
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		// Keep default debug for local dev; Docker/CI can set GIN_MODE=release
	}

	router := gin.New()
	router.Use(gin.Recovery())

	http.RegisterRoutes(router)

	addr := ":8080"
	log.Printf("starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
