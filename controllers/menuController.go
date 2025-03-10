package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get all menus
func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := database.MenuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing the menu items"})
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allMenus)
	}
}

// Get a single menu by ID
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := database.MenuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menu"})
			return
		}

		c.JSON(http.StatusOK, menu)
	}
}

// Create a new menu
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		validationErr := helpers.Validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Set timestamps and IDs
		menu.CreatedAt = time.Now()
		menu.UpdatedAt = time.Now()
		menu.ID = primitive.NewObjectID()
		menuID := menu.ID.Hex()
		menu.MenuID = menuID

		// Insert into database
		result, insertErr := database.MenuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Function to check if a time is within a range
func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// Update a menu
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		menuId := c.Param("menu_id")

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"menu_id": menuId}
		var updateObj primitive.D

		// Validate start and end dates
		if menu.StartDate != nil && menu.EndDate != nil {
			if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.EndDate})
		}

		// Update name if provided
		if len(menu.Name) > 0 {
			updateObj = append(updateObj, bson.E{Key:"name", Value: menu.Name})
		}

		// Update category if provided
		if len(menu.Category) > 0 {
			updateObj = append(updateObj, bson.E{Key:"category",Value:  menu.Category})
		}

		// Update timestamp
		menu.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{Key:"updated_at",Value:  menu.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}

		// Perform update
		result, err := database.MenuCollection.UpdateOne(ctx, filter, bson.D{{Key:"$set", Value: updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
