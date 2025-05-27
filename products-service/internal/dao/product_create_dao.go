package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductCreateDAO struct {
	collection *mongo.Collection
}

func NewProductCreateDAO(collection *mongo.Collection) *ProductCreateDAO {
	return &ProductCreateDAO{
		collection: collection,
	}
}

func (dao *ProductCreateDAO) CreateProduct(ctx context.Context, product *models.Product) error {
	// Generate a new ObjectID if not provided
	if product.ProductID.IsZero() {
		product.ProductID = primitive.NewObjectID()
	}

	_, err := dao.collection.InsertOne(ctx, product)
	return err
}
