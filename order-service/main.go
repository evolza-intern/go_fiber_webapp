package main

import (
	"log"
	"os"

	database "github.com/evolza-intern/go_fiber_webapp/order-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS for microservice communication
	app.Use(cors.New())

	// Connect to MongoDB
	database.ConnectMongo()
	orderCollection := database.GetCollection("go_fiber_eCommerce", "orders")

	// Inject collection into handlers
	handlers.InitOrderHandler(orderCollection)

	// Create invoices directory if it doesn't exist
	if err := os.MkdirAll("./invoices", 0755); err != nil {
		log.Fatal("Failed to create invoices directory:", err)
	}

	// API Routes
	api := app.Group("/api/orders")
	api.Post("/", handlers.CreateOrder)                    // Save order (called by cart)
	api.Get("/user/:userId", handlers.GetUserOrders)       // Get all orders of a user
	api.Get("/:orderId/invoice", handlers.GetOrderInvoice) // Get specific order for invoice generation

	// File Server Routes for Invoice Downloads
	app.Static("/invoices", "./invoices")                            // Static file serving
	app.Get("/download/invoice/:filename", handlers.DownloadInvoice) // Download specific invoice
	app.Get("/download/orders/:userId", handlers.ListUserInvoices)   // List user's available invoices

	log.Println("Order Service running on http://localhost:3001")
	log.Fatal(app.Listen(":3001"))
}
