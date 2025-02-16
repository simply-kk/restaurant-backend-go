package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`                //? Unique menu ID (MongoDB ObjectID)
	MenuID    string             `json:"menu_id" validate:"required"`  //? Unique menu identifier
	Name      string             `json:"name" validate:"required"`     //? Name of the menu
	Category  string             `json:"category" validate:"required"` //? Category of the menu
	StartDate time.Time          `json:"start_date"`                   //? Start date of menu availability
	EndDate   time.Time          `json:"end_date"`                     //? End date of menu availability
	CreatedAt time.Time          `json:"created_at"`                   //? Timestamp when the menu was created
	UpdatedAt time.Time          `json:"updated_at"`                   //? Timestamp when the menu was last updated
}
