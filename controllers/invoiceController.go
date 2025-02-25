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
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := database.InvoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing invoice items"})
			return
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

// Get a single invoice
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceID := c.Param("invoice_id")
		var invoice models.Invoice

		err := database.InvoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching invoice item"})
			return
		}

		// Format invoice response
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderID)
		invoiceView.OrderID = invoice.OrderID
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.InvoiceID = invoice.InvoiceID

		// Handle nil pointers safely
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		} else {
			invoiceView.PaymentMethod = "null"
		}

		if invoice.PaymentStatus != nil {
			invoiceView.PaymentStatus = invoice.PaymentStatus
		}

		// Handle missing order details safely
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order was not found"})
			return
		}

		// Assign default status
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		// Assign timestamps
		now := time.Now()
		invoice.PaymentDueDate = now.AddDate(0, 0, 1)
		invoice.CreatedAt = now
		invoice.UpdatedAt = now
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()

		// Insert into DB
		result, insertErr := database.InvoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invoice item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Update an invoice
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"invoice_id": invoiceID}
		var updateObj primitive.D

		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{"payment_method", *invoice.PaymentMethod})
		}

		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{"payment_status", *invoice.PaymentStatus})
		}

		// Update timestamp
		invoice.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updated_at", invoice.UpdatedAt})

		// Perform update
		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}

		result, err := database.InvoiceCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invoice update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
