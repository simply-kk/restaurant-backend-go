package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`                  //? Unique food ID (MongoDB ObjectID)
	Name      *string            `json:"name" validate:"required"`       //? Name of the food item
	Price     *float64           `json:"price" validate:"required,gt=0"` //? Price of the food item (must be greater than 0)
	FoodImage *string            `json:"image"`                          //? Image URL of the food item
	MenuID    *string            `json:"menu_id" validate:"required"`    //? Associated menu ID
	FoodID    string             `json:"food_id"`                        //? Unique food identifier
	CreatedAt time.Time          `json:"created_at"`                     //? Timestamp when the food item was created
	UpdatedAt time.Time          `json:"updated_at"`                     //? Timestamp when the food item was last updated
}
