package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var orderCollection *mongo.Collection

func InitOrderHandler(collection *mongo.Collection) {
	orderCollection = collection
}

// POST /api/orders - Save order (called by cart)
func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Set default values
	if order.ID.IsZero() {
		order.ID = primitive.NewObjectID()
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

	order.ID = result.InsertedID.(primitive.ObjectID)
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
			"error": "Invalid user ID format",
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
			"error": "Invalid order ID format",
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
		filename := fmt.Sprintf("invoice_%s.txt", order.ID.Hex())
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
		"invoice_url": fmt.Sprintf("/download/invoice/invoice_%s.txt", order.ID.Hex()),
	})
}

// GET /download/invoice/:filename - Download specific invoice file
func DownloadInvoice(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Security: Prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	// Validate filename format (should be invoice_<orderid>.txt)
	if !strings.HasPrefix(filename, "invoice_") || !strings.HasSuffix(filename, ".txt") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid invoice filename format",
		})
	}

	filepath := fmt.Sprintf("./invoices/%s", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "Invoice file not found",
		})
	}

	// Set headers for file download
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Set("Content-Type", "text/plain")

	return c.SendFile(filepath)
}

// GET /download/orders/:userId - List user's available invoices
func ListUserInvoices(c *fiber.Ctx) error {
	userID := c.Params("userId")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID format",
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
		filename := fmt.Sprintf("invoice_%s.txt", order.ID.Hex())
		invoices = append(invoices, fiber.Map{
			"order_id":     order.ID.Hex(),
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
