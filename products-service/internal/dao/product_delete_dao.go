package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductDeleteDAO struct {
	collection *mongo.Collection
}

type DeleteResult struct {
	MatchedCount int64
}

func NewProductDeleteDAO(collection *mongo.Collection) *ProductDeleteDAO {
	return &ProductDeleteDAO{
		collection: collection,
	}
}

// DeleteProduct performs soft delete by setting isDeleted to true
func (dao *ProductDeleteDAO) DeleteProduct(ctx context.Context, id primitive.ObjectID) (*DeleteResult, error) {
	// Only update documents that are not already deleted
	filter := bson.M{
		"_id":        id,
		"is_deleted": bson.M{"$ne": true}, // Only match if not already deleted
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
		},
	}

	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &DeleteResult{
		MatchedCount: result.MatchedCount,
	}, nil
}
