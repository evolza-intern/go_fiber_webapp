package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductUpdateDAO struct {
	collection *mongo.Collection
}

type UpdateResult struct {
	MatchedCount int64
}

func NewProductUpdateDAO(collection *mongo.Collection) *ProductUpdateDAO {
	return &ProductUpdateDAO{
		collection: collection,
	}
}

func (dao *ProductUpdateDAO) UpdateProduct(ctx context.Context, id primitive.ObjectID, product *models.Product) (*UpdateResult, error) {
	// Only update documents that are not deleted
	filter := bson.M{
		"_id":        id,
		"is_deleted": bson.M{"$ne": true}, // Only match if not deleted
	}

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
		},
	}

	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &UpdateResult{
		MatchedCount: result.MatchedCount,
	}, nil
}
