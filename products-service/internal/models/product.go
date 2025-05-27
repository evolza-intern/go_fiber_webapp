package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ProductID   primitive.ObjectID `bson:"product_id,omitempty" json:"productid"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	IsDeleted   bool               `bson:"is_deleted" json:"is_deleted"`
}
