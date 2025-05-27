package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteProduct DELETE
func (h *ProductHandlers) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert string ProductID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ProductID format",
		})
	}

	ctx := context.TODO()
	result, err := h.deleteDAO.DeleteProduct(ctx, objectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Delete failed",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Product not found or already deleted",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
