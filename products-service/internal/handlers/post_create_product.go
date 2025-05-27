package handlers

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *ProductHandlers) CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	ctx := context.TODO()
	if err := h.createDAO.CreateProduct(ctx, &product); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Insert failed",
		})
	}

	return c.Status(201).JSON(product)
}
