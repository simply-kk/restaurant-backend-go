package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get all orders with pagination
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
		matchStage := bson.D{{"$match", bson.D{}}}
		skipStage := bson.D{{"$skip", startIndex}}
		limitStage := bson.D{{"$limit", recordPerPage}}

		// Fetch orders with pagination
		result, err := database.OrderCollection.Aggregate(ctx, mongo.Pipeline{matchStage, skipStage, limitStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing orders"})
			return
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
			return
		}

		c.JSON(http.StatusOK, allOrders)
	}
}

// Get a single order by ID
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderID := c.Param("order_id")
		var order models.Order

		err := database.OrderCollection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

// Create a new order
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		var order models.Order

		// Parse JSON body
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if order.OrderDate.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order date is required"})
			return
		}

		validationErr := helpers.Validate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if table exists
		if order.TableID != nil {
			err := database.TableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			}
		}

		// Assign timestamps and IDs
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		order.ID = primitive.NewObjectID()
		order.OrderID = order.ID.Hex()

		// Insert into database
		result, insertErr := database.OrderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order could not be created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Update an existing order
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		var order models.Order

		orderID := c.Param("order_id")

		// Parse JSON body
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate request
		validationErr := helpers.Validate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Prepare update object
		var updateObj primitive.D

		// Validate if table exists before updating
		if order.TableID != nil {
			err := database.TableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.TableID})
		}

		// Update timestamp
		order.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", order.UpdatedAt})

		// Perform update
		filter := bson.M{"order_id": orderID}
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}

		result, err := database.OrderCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// OrderItemOrderCreator - Creates an order and returns its OrderID
func OrderItemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Assign timestamps and IDs
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	// Insert into database
	_, err := database.OrderCollection.InsertOne(ctx, order)
	if err != nil {
		return ""
	}

	return order.OrderID
}
