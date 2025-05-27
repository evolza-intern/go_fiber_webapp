package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductReadDAO struct {
	collection *mongo.Collection
}

func NewProductReadDAO(collection *mongo.Collection) *ProductReadDAO {
	return &ProductReadDAO{
		collection: collection,
	}
}

func (dao *ProductReadDAO) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	// Only fetch products that are not deleted
	filter := bson.M{
		"is_deleted": bson.M{"$ne": true}, // Exclude deleted products
	}

	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// Log error if needed
		}
	}(cursor, ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	// Handle empty results
	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

func (dao *ProductReadDAO) GetProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	// Only fetch product if it's not deleted
	filter := bson.M{
		"_id":        id,
		"is_deleted": bson.M{"$ne": true}, // Exclude deleted products
	}

	var product models.Product
	err := dao.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
