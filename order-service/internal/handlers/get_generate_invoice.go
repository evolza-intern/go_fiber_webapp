package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"os"
)

// GetInvoiceOrder GET /api/orders/:orderId/invoice - Get specific order for invoice generation
func (h *OrderHandlers) GetInvoiceOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	ctx := context.TODO()
	// For invoice generation, we need to get order by ID regardless of user
	// You might want to add a GetOrderByID method to readDAO
	order, err := h.readDAO.GetOrderByID(ctx, orderObjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Generate the invoice file if not already generated
	if !order.InvoiceGenerated {
		invoiceContent := utils.GenerateInvoiceContent(*order)
		filename := fmt.Sprintf("invoice_%s.txt", order.OrderID.Hex())
		filepath := fmt.Sprintf("./invoices/%s", filename)

		if err := os.WriteFile(filepath, []byte(invoiceContent), 0644); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate invoice file",
			})
		}

		// Update order to mark the invoice as generated
		orderUpdate := *order
		orderUpdate.InvoiceGenerated = true
		_, err = h.updateDAO.UpdateInvoiceGeneratedField(ctx, orderObjectID)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to update invoice_generated flag: %v\n", err)
		}
	}

	return c.JSON(fiber.Map{
		"order":       order,
		"invoice_url": fmt.Sprintf("/download/invoice/invoice_%s.txt", order.OrderID.Hex()),
	})
}
