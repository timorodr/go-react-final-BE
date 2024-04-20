package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// model in DB
type Medication struct {
	ID          primitive.ObjectID `bson:"id"` // ID created by GO so we dont have to pass it help of bson pkg Golang understands this type
	Name        *string            `json: "name"`
	Dosage      *string            `json: "dosage"`
	Description *string            `json: "description"`
}
