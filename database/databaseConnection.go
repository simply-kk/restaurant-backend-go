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

//  Load environment variables
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

//  DbInstance initializes the MongoDB client only once
func DbInstance() (*mongo.Client, error) {
	if Client != nil {
		return Client, nil
	}

	// Load .env variables before getting MONGODB_URI
	loadEnv()

	MongoDbURI := os.Getenv("MONGODB_URI")
	if MongoDbURI == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(MongoDbURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	//  Ping MongoDB to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("could not ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	Client = client
	return Client, nil
}

//  OpenCollection returns a reference to a MongoDB collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("Cannot open collection: MongoDB client is nil")
	}
	return client.Database("restaurant").Collection(collectionName)
}

//  CloseDB gracefully closes the MongoDB connection
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
