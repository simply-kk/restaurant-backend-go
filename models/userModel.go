package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`                                       //? Unique user ID (MongoDB ObjectID)
	FirstName string             `json:"first_name" validate:"required"`                      //? First name of the user
	LastName  string             `json:"last_name" validate:"required"`                       //? Last name of the user
	Email     string             `json:"email" validate:"required,email"`                     //? User email (must be valid)
	Password  string             `json:"password" validate:"required,min=6"`                  //? Hashed password
	Phone     string             `json:"phone" validate:"required"`                           //? Contact phone number
	Role      string             `json:"role" validate:"required,oneof=admin staff customer"` //? User role (admin, staff, customer)
	CreatedAt time.Time          `json:"created_at"`                                          //? Timestamp when the user was created
	UpdatedAt time.Time          `json:"updated_at"`                                          //? Timestamp when the user was last updated
}
