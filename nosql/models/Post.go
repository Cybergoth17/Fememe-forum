package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `json:"username" bson:"username"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Title    string             `json:"title" bson:"title"`
	Text     string             `json:"text" `
	Comment  []Comment          `json:"comments" `
	Tags     []string           `json:"tags"`
	Date     time.Time          `json:"date" `
}
