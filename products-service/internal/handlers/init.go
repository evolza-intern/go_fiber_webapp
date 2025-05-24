package handlers

import "go.mongodb.org/mongo-driver/v2/mongo"

var productCollection *mongo.Collection

func InitProductHandler(collection *mongo.Collection) {
	productCollection = collection
}
