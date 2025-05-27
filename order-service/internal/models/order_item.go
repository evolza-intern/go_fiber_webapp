package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName string             `bson:"product_name" json:"product_name"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Price       float64            `bson:"price" json:"price"`
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}
