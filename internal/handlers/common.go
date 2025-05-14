package handlers

import (
	"github.com/cmerin0/tasky/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// Declare collections at package level
// to avoid multiple calls to GetCollection
// and to ensure they are initialized only once
var (
	userCollection *mongo.Collection
	taskCollection *mongo.Collection
)

// getUserCollection returns the user collection
// from the database. It initializes it if not already done.
func getUserCollection() *mongo.Collection {
	if userCollection == nil {
		userCollection = db.GetCollection("users")
	}
	return userCollection
}

// getTaskCollection returns the task collection
// from the database. It initializes it if not already done.
func getTaskCollection() *mongo.Collection {
	if taskCollection == nil {
		taskCollection = db.GetCollection("tasks")
	}
	return taskCollection
}
