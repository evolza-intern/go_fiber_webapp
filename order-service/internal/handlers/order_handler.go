package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// POST /api/orders - Save order (called by cart)
func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Set default values
	if order.OrderID.IsZero() {
		order.OrderID = primitive.NewObjectID()
	}
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

	result, err := orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	order.OrderID = result.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   order,
	})
}

// GET /api/orders/user/:userId - Get all orders of a user
func GetUserOrders(c *fiber.Ctx) error {
	userID := c.Params("userId")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user OrderID format",
		})
	}

	cursor, err := orderCollection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}
	defer cursor.Close(context.TODO())

	var orders []models.Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to decode orders",
		})
	}

	if orders == nil {
		orders = []models.Order{}
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"count":  len(orders),
	})
}

// GET /api/orders/:orderId/invoice - Get specific order for invoice generation
func GetOrderInvoice(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order OrderID format",
		})
	}

	var order models.Order
	err = orderCollection.FindOne(context.TODO(), bson.M{"_id": orderObjectID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Generate invoice file if not already generated
	if !order.InvoiceGenerated {
		invoiceContent := utils.GenerateInvoiceContent(order)
		filename := fmt.Sprintf("invoice_%s.txt", order.OrderID.Hex())
		filepath := fmt.Sprintf("./invoices/%s", filename)

		if err := os.WriteFile(filepath, []byte(invoiceContent), 0644); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate invoice file",
			})
		}

		// Update order to mark invoice as generated
		_, err = orderCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": orderObjectID},
			bson.M{"$set": bson.M{"invoice_generated": true}},
		)
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

// GET /download/invoice/:filename - Download specific invoice file

// GET /download/orders/:userId - List user's available invoices
func ListUserInvoices(c *fiber.Ctx) error {
	userID := c.Params("userId")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user OrderID format",
		})
	}

	// Get user's orders with generated invoices
	cursor, err := orderCollection.Find(context.TODO(), bson.M{
		"user_id":           userObjectID,
		"invoice_generated": true,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}
	defer cursor.Close(context.TODO())

	var orders []models.Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to decode orders",
		})
	}

	// Build invoice list
	invoices := make([]fiber.Map, 0)
	for _, order := range orders {
		filename := fmt.Sprintf("invoice_%s.txt", order.OrderID.Hex())
		invoices = append(invoices, fiber.Map{
			"order_id":     order.OrderID.Hex(),
			"order_date":   order.OrderDate,
			"total_amount": order.TotalAmount,
			"status":       order.Status,
			"filename":     filename,
			"download_url": fmt.Sprintf("/download/invoice/%s", filename),
		})
	}

	return c.JSON(fiber.Map{
		"user_id":  userID,
		"invoices": invoices,
		"count":    len(invoices),
	})
}
