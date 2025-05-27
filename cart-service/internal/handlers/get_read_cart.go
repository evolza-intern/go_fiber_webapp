package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/db"
	"github.com/gofiber/fiber/v2"
)

// GET /:userId - Get current cart
func GetCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	cart, err := db.MongoStorage.GetCart(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get cart",
		})
	}

	return c.JSON(cart)
}
