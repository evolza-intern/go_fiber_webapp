package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderCreateDAO struct {
	collection *mongo.Collection
}

func NewOrderCreateDAO(collection *mongo.Collection) *OrderCreateDAO {
	return &OrderCreateDAO{
		collection: collection,
	}
}

func (dao *OrderCreateDAO) CreateOrder(ctx context.Context, order *models.Order) error {
	// Generate a new ObjectID if not provided
	if order.OrderID.IsZero() {
		order.OrderID = primitive.NewObjectID()
	}

	_, err := dao.collection.InsertOne(ctx, order)
	return err
}
