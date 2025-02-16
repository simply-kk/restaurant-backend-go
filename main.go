package main

import (
	"os"

	"golang-restaurant-management/database"
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	// Load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Database
	database.ConnectDatabase()

	// Initialize Router
	router := gin.New()
	router.Use(gin.Logger())

	// Public Routes
	routes.UserRoutes(router)

	// Middleware for Authentication (Apply to Protected Routes)
	router.Use(middleware.Authentication())

	// Protected Routes
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	// Start Server
	router.Run(":" + port)
}
