package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`                                                //? Unique invoice ID (MongoDB ObjectID)
	InvoiceID      string             `json:"invoice_id" validate:"required"`                               //? Invoice ID as a string
	OrderID        string             `json:"order_id" validate:"required"`                                 //? Associated order ID
	PaymentMethod  string             `json:"payment_method" validate:"required,oneof=cash card upi"`       //? Payment method used
	PaymentStatus  string             `json:"payment_status" validate:"required,oneof=pending paid failed"` //? Status of the payment
	PaymentDueDate time.Time          `json:"payment_due_date"`                                             //? Due date for the payment
	CreatedAt      time.Time          `json:"created_at"`                                                   //? Timestamp when the invoice was created
	UpdatedAt      time.Time          `json:"updated_at"`                                                   //? Timestamp when the invoice was last updated
}
