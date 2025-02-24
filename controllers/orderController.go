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

// Initialize order collection
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

// Get all orders
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing orders"})
			return
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allOrders)
	}
}

// Get a single order by ID
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderID := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		var order models.Order

		// Parse JSON body
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		validationErr := validate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if table exists
		if order.TableID != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
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

		// Insert into DB
		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order could not be created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Update an order
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		var order models.Order

		orderID := c.Param("order_id")

		// Parse JSON body
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prepare update object
		var updateObj primitive.D

		if order.TableID != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.TableID})
		}

		// Update timestamp
		order.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", order.UpdatedAt})

		// Update options
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}
		filter := bson.M{"order_id": orderID}

		// Perform update
		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Create order and return OrderID
func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Assign timestamps and IDs
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	// Insert into DB
	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return ""
	}

	return order.OrderID
}
