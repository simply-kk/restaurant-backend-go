package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initialize table collection
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

// Get all tables
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing tables"})
			return
		}

		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTables)
	}
}

// Get a single table by ID
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

// Create a new table
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		// Parse JSON body
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Assign timestamps and IDs
		table.CreatedAt = time.Now()
		table.UpdatedAt = time.Now()
		table.ID = primitive.NewObjectID()
		table.TableID = table.ID.Hex()

		// Insert into DB
		result, insertErr := tableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table could not be created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Update a table
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		tableId := c.Param("table_id")

		// Parse JSON body
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prepare update object
		var updateObj primitive.D

		if table.NumberOfGuests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", table.NumberOfGuests})
		}

		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{"table_number", table.TableNumber})
		}

		// Update timestamp
		table.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", table.UpdatedAt})

		// Update options
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}
		filter := bson.M{"table_id": tableId}

		// Perform update
		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
