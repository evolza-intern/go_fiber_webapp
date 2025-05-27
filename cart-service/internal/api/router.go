package api

import (
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize handlers (this could also be moved to main.go if preferred)
	if err := db.InitHandlers(); err != nil {
		panic("Failed to initialize cart handlers: " + err.Error())
	}

	// REST API Routes
	app.Get("/:userId", handlers.GetCart)                // Get current cart
	app.Post("/checkout/:userId", handlers.CheckoutCart) // Checkout (saves to order-service)

	// WebSocket Route
	app.Get("/ws/:userId", handlers.HandleWebSocket) // WebSocket connection for real-time updates
}
