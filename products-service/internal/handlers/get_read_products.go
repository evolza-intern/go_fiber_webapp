package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

// GetProducts READ ALL
func (h *ProductHandlers) GetProducts(c *fiber.Ctx) error {
	ctx := context.TODO()
	products, err := h.readDAO.GetAllProducts(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Find failed",
		})
	}

	return c.JSON(products)
}
