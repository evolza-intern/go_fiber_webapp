package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderReadDAO struct {
	collection *mongo.Collection
}

func NewOrderReadDAO(collection *mongo.Collection) *OrderReadDAO {
	return &OrderReadDAO{
		collection: collection,
	}
}

// GetAllOrdersByUser returns all orders belonging to a specific user
func (dao *OrderReadDAO) GetAllOrdersByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Order, error) {
	filter := bson.M{
		"user_id": userID,
	}
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []models.Order{}
	}

	return orders, nil
}

// GetOrderByIDForUser fetches a specific order by order_id and user_id
func (dao *OrderReadDAO) GetOrderByIDForUser(ctx context.Context, orderID, userID primitive.ObjectID) (*models.Order, error) {
	filter := bson.M{
		"order_id": orderID,
		"user_id":  userID,
	}

	var order models.Order
	err := dao.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderByID fetches a specific order by order_id only (for admin/system operations)
func (dao *OrderReadDAO) GetOrderByID(ctx context.Context, orderID primitive.ObjectID) (*models.Order, error) {
	filter := bson.M{
		"order_id": orderID,
	}

	var order models.Order
	err := dao.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
