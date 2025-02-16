package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`                      //? Unique order item ID (MongoDB ObjectID)
	Quantity    *int               `json:"quantity" validate:"required,min=1"` //? Quantity of the item
	UnitPrice   *float64           `json:"unit_price" validate:"required"`     //? Price per unit of the item
	OrderID     string             `json:"order_id" validate:"required"`       //? Associated order ID
	OrderItemID string             `json:"order_item_id" validate:"required"`  //? Unique order item identifier
	FoodID      string             `json:"food_id" validate:"required"`        //? Associated food item ID
	CreatedAt   time.Time          `json:"created_at"`                         //? Timestamp when the order item was created
	UpdatedAt   time.Time          `json:"updated_at"`                         //? Timestamp when the order item was last updated
}
