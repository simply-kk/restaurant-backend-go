package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-restaurant-management/database"
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not set
	}

	// Initialize Database Connection
	client, err := database.DbInstance()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize collection
	database.InitCollections(client)

	// Initialize Gin Router
	router := gin.Default() // Includes default logging and recovery middleware

	// Public Routes (No Authentication Required)
	routes.UserRoutes(router)

	// Apply Authentication Middleware (Protected Routes)
	router.Use(middleware.Authentication())

	// Protected Routes
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	// Run server in a goroutine so it doesn’t block shutdown handling
	go func() {
		fmt.Println("Server running on port:", port)
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// **Graceful Shutdown Handling**
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // Catch termination signals
	<-quit // Block until a signal is received

	fmt.Println("\nShutting down server...")

	// Gracefully close the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Println("Error disconnecting MongoDB:", err)
	} else {
		fmt.Println("MongoDB connection closed successfully.")
	}

	fmt.Println("Server shutdown completed.")
}
