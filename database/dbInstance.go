package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MenuCollection    *mongo.Collection
	OrderCollection   *mongo.Collection
	TableCollection   *mongo.Collection
	FoodCollection    *mongo.Collection
	InvoiceCollection *mongo.Collection
	OrderItemCollection *mongo.Collection
)

func InitCollections(client *mongo.Client) {
	MenuCollection = client.Database("restaurant").Collection("menu")
	OrderCollection = client.Database("restaurant").Collection("order")
	TableCollection = client.Database("restaurant").Collection("table")
	FoodCollection = client.Database("restaurant").Collection("food")
	InvoiceCollection = client.Database("restaurant").Collection("invoice")
	OrderItemCollection = client.Database("restaurant").Collection("orderItem")
}
