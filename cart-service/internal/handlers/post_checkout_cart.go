package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

// POST /checkout/:userId - Checkout (saves to order-service)
func CheckoutCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var checkoutReq models.CheckoutRequest
	if err := c.BodyParser(&checkoutReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid checkout data",
		})
	}

	// Get current cart from MongoDB
	cart, err := db.MongoStorage.GetCart(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get cart",
		})
	}

	if len(cart.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cart is empty",
		})
	}

	// Convert cart items to order items
	orderItems := make([]models.OrderItem, len(cart.Items))
	for i, item := range cart.Items {
		orderItems[i] = models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		}
	}

	// Create order request
	orderReq := models.OrderRequest{
		UserID:          userID,
		Items:           orderItems,
		TotalAmount:     cart.Total,
		PaymentMethod:   checkoutReq.PaymentMethod,
		ShippingAddress: checkoutReq.ShippingAddress,
	}

	// Send order to order-service
	orderJSON, err := json.Marshal(orderReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to prepare order",
		})
	}

	resp, err := http.Post(db.OrderServiceUrl, "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		log.Printf("Order service returned status: %d", resp.StatusCode)
		return c.Status(500).JSON(fiber.Map{
			"error": "Order service failed to create order",
		})
	}

	// Parse order service response
	var orderResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse order response",
		})
	}

	// Clear cart after successful checkout (this will trigger polling notification)
	if err := db.MongoStorage.ClearCart(userID); err != nil {
		log.Printf("Warning: Failed to clear cart: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
		"order":   orderResponse,
	})
}
