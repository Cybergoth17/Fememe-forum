package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID       primitive.ObjectID `bson:"_id"`
	Text     string            `json:"text" validate:"required" bson:"text"`
	Date     time.Time         `json:"date" validate:"required" bson:"date"`
	Username string
}
