package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/config"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/database"
	"github.com/Mohammed-Yasin-Mulla/easyplots-whtasapp.git/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize database connection
	dbpool, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer database.Close(dbpool)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, dbpool)

	// Start server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
