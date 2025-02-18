package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang-restaurant-management/database"
	"golang-restaurant-management/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Initialize MongoDB collections
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

// GetFoods retrieves all food items (to be implemented)
func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// ! GetFood retrieves a single food item by its food_id
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set timeout for the database operation
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodId := c.Param("food_id") // Corrected parameter

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

// ! CreateFood adds a new food item
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		// Parse JSON body
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if the referenced menu exists
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuID}).Decode(&menu)
		if err != nil {
			msg := "Menu was not found"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Set timestamps
		food.CreatedAt = time.Now()
		food.UpdatedAt = time.Now()

		// Generate unique ID
		food.ID = primitive.NewObjectID()
		food.FoodID = food.ID.Hex()

		// Round price to 2 decimal places
		num := toFixed(*food.Price, 2)
		food.Price = &num

		// Insert into database
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := "Food item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, result)
	}
}

// toFixed rounds a float64 to a given number of decimal places
func toFixed(num float64, precision int) float64 {
	output := fmt.Sprintf("%.*f", precision, num)
	var rounded float64
	fmt.Sscanf(output, "%f", &rounded)
	return rounded
}
