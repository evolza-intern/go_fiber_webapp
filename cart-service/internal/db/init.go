package db

import (
	"log"
)

const OrderServiceUrl = "http://localhost:3001"

var MongoStorage *MongoCartStorage

func InitHandlers() error {
	var err error
	MongoStorage, err = NewMongoCartStorage(
		"mongodb+srv://admin:OohHyQBJW4d43Ter@cluster0.y1nym.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", // MongoDB connection string
		"go_fiber_eCommerce", // Database name
	)
	if err != nil {
		return err
	}

	log.Println("Cart handlers initialized with MongoDB")
	return nil
}

// Cleanup function to call when shutting down
func Cleanup() {
	if MongoStorage != nil {
		err := MongoStorage.Close()
		if err != nil {
			return
		}
	}
}
