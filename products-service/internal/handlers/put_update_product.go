package handlers

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateProduct UPDATE
func (h *ProductHandlers) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert string ProductID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ProductID format",
		})
	}

	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	ctx := context.TODO()
	result, err := h.updateDAO.UpdateProduct(ctx, objectID, &product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Update failed",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}
