package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`                  //? Unique order ID (MongoDB ObjectID)
	OrderID   string             `json:"order_id"`                       //? Order ID as a string
	TableID   *string            `json:"table_id"`                       //? ID of the table associated with this order
	OrderDate time.Time          `json:"order_date" validate:"required"` //? Time when order was placed
	CreatedAt time.Time          `json:"created_at"`                     //? Timestamp when order was created
	UpdatedAt time.Time          `json:"updated_at"`                     //? Timestamp when order was last updated
}
