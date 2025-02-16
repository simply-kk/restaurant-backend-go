package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`               //? Unique note ID (MongoDB ObjectID)
	Title     string             `json:"title" validate:"required"`   //? Title of the note
	NoteID    string             `json:"note_id" validate:"required"` //? Unique note identifier
	Text      string             `json:"text"`                        //? Content of the note
	CreatedAt time.Time          `json:"created_at"`                  //? Timestamp when the note was created
	UpdatedAt time.Time          `json:"updated_at"`                  //? Timestamp when the note was last updated
}
