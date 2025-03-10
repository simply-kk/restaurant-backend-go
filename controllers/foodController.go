package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetFoods retrieves all food items with pagination
func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Handle pagination
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		// Aggregation pipeline
		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		skipStage := bson.D{{Key: "$skip", Value: startIndex}}
		limitStage := bson.D{{Key: "$limit", Value: recordPerPage}}

		// Execute aggregation query
		result, err := database.FoodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, skipStage, limitStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing food items"})
			return
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods)
	}
}

// GetFood retrieves a single food item by its food_id
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodID := c.Param("food_id")
		var food models.Food

		err := database.FoodCollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
			return
		}

		c.JSON(http.StatusOK, food)
	}
}

// CreateFood adds a new food item
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

		// Ensure required fields are not nil
		if food.Name == nil || food.Price == nil || food.MenuID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Validate input
		validationErr := helpers.Validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if menu exists
		err := database.MenuCollection.FindOne(ctx, bson.M{"menu_id": *food.MenuID}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
			return
		}

		// Set timestamps
		food.CreatedAt = time.Now()
		food.UpdatedAt = time.Now()
		food.ID = primitive.NewObjectID()
		food.FoodID = food.ID.Hex()

		// Round price safely
		num := toFixed(*food.Price, 2)
		food.Price = &num

		// Insert into database
		result, insertErr := database.FoodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Food item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateFood updates an existing food item
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		foodID := c.Param("food_id")

		// Parse JSON body
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prepare update object
		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}

		if food.Price != nil {
			num := toFixed(*food.Price, 2)
			updateObj = append(updateObj, bson.E{Key: "price", Value: num})
		}

		if food.FoodImage != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.FoodImage})
		}

		if food.MenuID != nil {
			err := database.MenuCollection.FindOne(ctx, bson.M{"menu_id": *food.MenuID}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.MenuID})
		}

		// Update timestamp
		food.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

		// Update options
		upsert := true
		filter := bson.M{"food_id": foodID}
		opt := options.UpdateOptions{Upsert: &upsert}

		// Perform update
		result, err := database.FoodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Food item update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// DeleteFood deletes a food item
func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodID := c.Param("food_id")

		filter := bson.M{"food_id": foodID}
		result, err := database.FoodCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete food item"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Food item deleted successfully"})
	}
}

// toFixed rounds a float to a given number of decimal places
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output+0.5)) / output
}
