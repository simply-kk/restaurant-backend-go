package database

import (
	"context"
	"fmt"
	// "log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbInstance initializes a MongoDB client and returns it with an error if any.
func DbInstance() (*mongo.Client, error) {
	MongoDbURI := "mongodb://localhost:27017"
	fmt.Println("Connecting to MongoDB at:", MongoDbURI)

	// Set a 10-second timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(MongoDbURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping to check if the connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client, nil
}

// OpenCollection returns a MongoDB collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("restaurant").Collection(collectionName)
}
