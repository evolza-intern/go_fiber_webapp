package handlers

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Generate a new ObjectID if not provided
	if product.ID.IsZero() {
		product.ID = primitive.NewObjectID()
	}

	_, err := productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Insert failed",
		})
	}

	// Don't reassign ID - it's already set correctly
	return c.Status(201).JSON(product)
}
