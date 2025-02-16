package controllers

import (
	"context"
	"net/http"
	"time"

	"golang-restaurant-management/database"
	"golang-restaurant-management/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Initialize MongoDB collection
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// GetFoods retrieves all food items (to be implemented)
func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// GetFood retrieves a single food item by its food_id
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set timeout for the database operation
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodId := c.Param("food_id") // Corrected the parameter

		var food models.Food

		// Fetch food item from MongoDB
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching food item"})
			return
		}

		// Return the found food item
		c.JSON(http.StatusOK, food)
	}
}

// CreateFood adds a new food item (to be implemented)
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
	}
}

func Round() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func toFixed() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
