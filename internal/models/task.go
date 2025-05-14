package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Completed   bool               `json:"completed" bson:"completed" default:"false"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
}
