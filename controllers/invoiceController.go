package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct to hold invoice view format
type InvoiceViewFormat struct {
	InvoiceID      string
	PaymentMethod  string
	OrderID        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

// Get all invoices
func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := database.InvoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing invoice items"})
			return
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
			return
		}

		c.JSON(http.StatusOK, allInvoices)
	}
}

// Get a single invoice
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceID := c.Param("invoice_id")
		var invoice models.Invoice

		err := database.InvoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}

		// Format invoice response
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}

		invoiceView.OrderID = invoice.OrderID
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.InvoiceID = invoice.InvoiceID

		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		} else {
			invoiceView.PaymentMethod = "UNKNOWN"
		}

		if invoice.PaymentStatus != nil {
			invoiceView.PaymentStatus = invoice.PaymentStatus
		}

		if len(allOrderItems) > 0 {
			invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
			invoiceView.TableNumber = allOrderItems[0]["table_number"]
			invoiceView.OrderDetails = allOrderItems[0]["order_items"]
		}

		c.JSON(http.StatusOK, invoiceView)
	}
}

// Create a new invoice
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if order exists
		var order models.Order
		err := database.OrderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order was not found"})
			return
		}

		// Assign default status
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		now := time.Now()
		invoice.PaymentDueDate = now.AddDate(0, 0, 1)
		invoice.CreatedAt = now
		invoice.UpdatedAt = now
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()

		result, insertErr := database.InvoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invoice item was not created",
				"details": insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Invoice created successfully", "result": result})
	}
}

// Update an invoice
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")

		// Parse JSON request body
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Validate invoice_id format
		objID, err := primitive.ObjectIDFromHex(invoiceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format"})
			return
		}

		// Prepare update object
		filter := bson.M{"_id": objID} // Ensure querying by ObjectID
		var updateObj bson.D

		// Update fields if not nil
		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: *invoice.PaymentMethod})
		}

		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_status", Value: *invoice.PaymentStatus})
		}

		// Always update timestamp
		invoice.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

		// Ensure there's at least one update field
		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
			return
		}

		// Perform update operation
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}
		result, err := database.InvoiceCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invoice update failed",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully", "result": result})
	}
}
