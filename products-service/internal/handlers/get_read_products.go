package handlers

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetProducts READ ALL
func GetProducts(c *fiber.Ctx) error {
	cursor, err := productCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Find failed",
		})
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Cursor decode failed",
		})
	}

	// Handle empty results
	if products == nil {
		products = []models.Product{}
	}

	return c.JSON(products)
}
