package handlers

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

// POST /api/orders - Save order (called by cart)
func (h *OrderHandlers) SaveOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Set default values
	order.OrderDate = time.Now()
	if order.Status == "" {
		order.Status = "pending"
	}
	if order.PaymentStatus == "" {
		order.PaymentStatus = "pending"
	}
	order.InvoiceGenerated = false

	// Calculate total if not provided
	if order.TotalAmount == 0 {
		for _, item := range order.Items {
			order.TotalAmount += item.Subtotal
		}
	}

	ctx := context.TODO()
	err := h.createDAO.CreateOrder(ctx, &order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   order,
	})
}
