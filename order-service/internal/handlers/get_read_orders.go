package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetReadOrder GET /api/orders/user/:userId - Get all orders of a user
func (h *OrderHandlers) GetReadOrder(c *fiber.Ctx) error {
	userID := c.Params("userId")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx := context.TODO()
	orders, err := h.readDAO.GetAllOrdersByUser(ctx, userObjectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"count":  len(orders),
	})
}
