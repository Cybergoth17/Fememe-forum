package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `json:"username" bson:"username"`
	Title    string             `json:"title" `
	Text     string             `json:"text" `
	Comment  []Comment
	Tags     []string  `json:"tags"`
	Date     time.Time `json:"date" `
}
