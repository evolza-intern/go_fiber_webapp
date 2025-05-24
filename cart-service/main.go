package main

import (
	"log"

	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/handlers"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS for microservice communication
	app.Use(cors.New())

	// Initialize in-memory storage
	storage.InitCartStorage()

	// Initialize handlers
	handlers.InitHandlers()

	// REST API Routes
	app.Get("/:userId", handlers.GetCart)                // Get current cart (fallback)
	app.Post("/checkout/:userId", handlers.CheckoutCart) // Checkout (saves to order-service)

	// WebSocket Route
	app.Get("/ws/:userId", handlers.HandleWebSocket) // WebSocket connection

	log.Println("Cart Service running on http://localhost:3002")
	log.Fatal(app.Listen(":3002"))
}
