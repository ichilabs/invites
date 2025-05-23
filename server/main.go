package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"invites.cc/database"
	"invites.cc/routes"
	"invites.cc/utils"
)

func main() {
	// Load database configuration
	config := database.DWHConfig()

	// Connect to database with retry
	db, err := database.ConnectDB(config.DBConnectionString(), 3)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Verify database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	defer sqlDB.Close()

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Gin router with mode from environment
	gin.SetMode(utils.GetEnv("GIN_MODE", "release"))
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, db)

	// Get port from environment variable or use default
	port := utils.GetEnv("PORT", "9090")

	// Create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()
	
	// TODO : Generated by an IA to make the let the server run indefinitely. Need to check how it should be done.
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
