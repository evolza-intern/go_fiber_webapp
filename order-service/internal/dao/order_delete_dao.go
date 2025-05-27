package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderDeleteDAO struct {
	collection *mongo.Collection
}

func NewOrderDeleteDAO(collection *mongo.Collection) *OrderDeleteDAO {
	return &OrderDeleteDAO{
		collection: collection,
	}
}

type DeleteResult struct {
	MatchedCount  int64
	ModifiedCount int64
}

func (dao *OrderDeleteDAO) CancelOrder(ctx context.Context, orderID primitive.ObjectID) (*DeleteResult, error) {
	filter := bson.M{
		"order_id": orderID,
	}

	update := bson.M{
		"$set": bson.M{
			"status": "cancelled",
		},
	}

	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &DeleteResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
	}, nil
}
