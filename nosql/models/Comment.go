package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID   primitive.ObjectID `bson:"_id"`
	Text *string            `json:"text" validate:"required" bson:"text"`
	Date *time.Time         `json:"date" validate:"required" bson:"date"`
}
