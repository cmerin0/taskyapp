package db

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client      *mongo.Client
	clientOnce  sync.Once
	isConnected bool
)

var Client *mongo.Client

func ConnectDB(uri string) {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatal("Failed to connect to MongoDB: ", err)
		}

		// Ping the primary
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatal("Failed to ping MongoDB:", err)
		}

		// Set the client to connected
		isConnected = true
		log.Println("MongoDB connection established")

		// Assign the client to the global variable
		Client = client
	})
}

func GetCollection(name string) *mongo.Collection {
	// Load environment variables from .env file
	// This is useful for loading the database name
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}
	dbname := os.Getenv("MONGO_DBNAME")
	if !isConnected {
		log.Fatal("MongoDB client not initialized. Call ConnectDB() first")
	}
	return client.Database(dbname).Collection(name)
}
