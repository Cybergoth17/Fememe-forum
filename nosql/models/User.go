package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      *string            `json:"username" validate:"required,max=12" bson:"username" form="name"`
	Password      *string            `json:"Password" validate:"required,min=6" bson:"password" form="password"`
	Email         *string            `json:"email" validate:"email,required" bson:"email" form="email"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
}
