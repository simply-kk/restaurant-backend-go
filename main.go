package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv" 
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client instance
var Client *mongo.Client

// DbInstance initializes the MongoDB client only once
func DbInstance() (*mongo.Client, error) {
	if Client != nil {
		return Client, nil // Return existing connection
	}

	// ✅ Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default values")
	}

	// ✅ Read MongoDB URI from .env
	MongoDbURI := os.Getenv("MONGODB_URI")
	if MongoDbURI == "" {
		MongoDbURI = "mongodb://localhost:27017" // Default URI
	}
	fmt.Println("Connecting to MongoDB at:", MongoDbURI)

	// Set client options
	clientOptions := options.Client().ApplyURI(MongoDbURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping MongoDB to check connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	Client = client
	return Client, nil
}

// OpenCollection returns a reference to a MongoDB collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("Cannot open collection: MongoDB client is nil")
	}
	return client.Database("restaurant").Collection(collectionName)
}

// CloseDB gracefully closes the MongoDB connection
func CloseDB() {
	if Client != nil {
		err := Client.Disconnect(context.TODO())
		if err != nil {
			log.Println("Error disconnecting MongoDB:", err)
		} else {
			fmt.Println("MongoDB connection closed successfully.")
		}
		Client = nil
	}
}
