package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Text     string             `json:"text" validate:"required" bson:"text"`
	Date     time.Time          `json:"date" validate:"required" bson:"date"`
	Username string             `json:"username"`
	PostID   primitive.ObjectID `json:"post_id" validate:"required" bson:"post_id"`
	Avatar   string             `json:"avatar" validate:"required" bson:"avatar"`
}
