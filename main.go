package main

import (
	"fmt"
	"log"
	"os"

	"golang-restaurant-management/database"
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection

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
	foodCollection = database.OpenCollection(client, "food")

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

	// Start the server
	fmt.Println("Server running on port:", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Ensure MongoDB closes when the server shuts down
	defer database.CloseDB()
}
