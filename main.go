package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Anjas2005/bookmarks-api/db"
	invoiceRoutes "github.com/Anjas2005/bookmarks-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Loading the environment
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Create an empty context
	dbContext, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	// Will end the cancel fxn when the main fxn ends
	defer cancel()

	// Get a DB connection pool back
	dbPool, err := db.DBConnect(dbContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting a connection pool: %v\n", err)
		return
	}
	// Close the pool after the main server ends
	// NOTE: defer calls run LIFO so cancel calls before dbpool.
	defer dbPool.Close()

	router := gin.Default()

	invoiceRoutes.RegisterInvoice(router, dbPool)

	if err := router.Run(":8080"); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
	}

}
