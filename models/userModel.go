package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_name    *string            `json:"first_name"`
	Last_name     *string            `json:"last_name"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
	Medications	  []Medication		 `json:"medications" bson:"medications"`
}


type Medication struct {
	Medication_id          primitive.ObjectID `json:"medication_id" bson:"_id"` // ID created by GO so we dont have to pass it help of bson pkg Golang understands this type
	Name        		   *string            `json:"name" bson:"name"`
	Dosage      		   *string            `json:"dosage" bson:"dosage"`
	Description			   *string            `json:"description" bson:"description"`
	// UserID		primitive.ObjectID `json:"user_id"`
}