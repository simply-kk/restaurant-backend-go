package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct to hold order items
type OrderItemPack struct {
	TableID    *string
	OrderItems []models.OrderItem
}

// Get all order items
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := database.OrderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing order items"})
			return
		}

		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// Get order items by Order ID
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing order items"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// Get a single order item by ID
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemID := c.Param("order_item_id")
		var orderItem models.OrderItem

		err := database.OrderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemID}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

// Update an order item
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		orderItemID := c.Param("order_item_id")

		// Parse JSON body
		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prepare update object
		var updateObj primitive.D

		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{"unit_price", *orderItem.UnitPrice})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
		}
		if orderItem.FoodID != nil {
			updateObj = append(updateObj, bson.E{"food_id", *orderItem.FoodID})
		}

		// Update timestamp
		orderItem.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.UpdatedAt})

		// Update options
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}
		filter := bson.M{"order_item_id": orderItemID}

		// Perform update
		result, err := database.OrderItemCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order item update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Create order items
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order models.Order

		// Parse JSON body
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Assign timestamps
		order.OrderDate = time.Now()
		order.TableID = orderItemPack.TableID
		orderID := OrderItemOrderCreator(order)

		// Prepare order items for insertion
		orderItemsToBeInserted := []interface{}{}

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderID = orderID

			// Validate input
			validationErr := helpers.Validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			// Assign IDs and timestamps
			orderItem.ID = primitive.NewObjectID()
			orderItem.OrderItemID = orderItem.ID.Hex()
			orderItem.CreatedAt = time.Now()
			orderItem.UpdatedAt = time.Now()

			// Round unit price
			var num = toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		// Insert into DB
		insertedOrderItems, err := database.OrderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order items"})
			return
		}

		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

// Items by OrderID aggregation pipeline
func ItemsByOrder(orderID string) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.D{{"order_id", orderID}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	// Execute aggregation
	result, err := database.OrderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, lookupStage, unwindStage,
	})

	if err != nil {
		return nil, err
	}

	var orderItems []bson.M
	if err = result.All(ctx, &orderItems); err != nil {
		return nil, err
	}

	return orderItems, nil
}
