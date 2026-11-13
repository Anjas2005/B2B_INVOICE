package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Anjas2005/bookmarks-api/db"
	"github.com/Anjas2005/bookmarks-api/modules"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {

	// Load .env if available
	if err := godotenv.Load(); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("loading environment: %w", err)
	}

	// Create startup context
	dbContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	// Initialize database pool
	dbPool, err := db.DBConnect(dbContext)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer dbPool.Close()

	// Initialize Gin
	router := gin.Default()

	modules.RegisterInvoice(router, dbPool)

	// Read port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	log.Printf("Starting server on %s", addr)

	// Start HTTP server
	if err := router.Run(addr); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
