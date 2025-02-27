package database

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client instance
var Client *mongo.Client

// DbInstance initializes the MongoDB client
func DbInstance() (*mongo.Client, error) {
    if Client != nil {
        return Client, nil
    }
    // Rely on main.go's init() for godotenv.Load(), or keep it here and remove from main.go
    MongoDbURI := os.Getenv("MONGODB_URI")
    if MongoDbURI == "" {
        log.Fatal("MONGODB_URI is not set in .env file")
    }
    log.Printf("Connecting to MongoDB with URI: %s", MongoDbURI) // Debug
    clientOptions := options.Client().ApplyURI(MongoDbURI)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
    }
    if client == nil {
        return nil, fmt.Errorf("mongo.Connect returned nil client")
    }
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