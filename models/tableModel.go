package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`                        //? Unique table ID (MongoDB ObjectID)
	NumberOfGuests *int               `json:"number_of_guests" validate:"required"` //? Number of guests at the table
	TableNumber    *int               `json:"table_number" validate:"required"`     //? Table number
	CreatedAt      time.Time          `json:"created_at"`                           //? Timestamp when the table was created
	UpdatedAt      time.Time          `json:"updated_at"`                           //? Timestamp when the table was last updated
	TableID        string             `json:"table_id"`                             //? Unique table identifier as a string
}
