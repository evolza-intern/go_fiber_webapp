package api

import (
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupRoutes(app *fiber.App, orderCollection *mongo.Collection) {
	orderHandlers := handlers.NewOrderHandlers(orderCollection)

	app.Post("/", orderHandlers.SaveOrder)                      // Save order (called by cart)
	app.Get("/user/:userId", orderHandlers.GetReadOrder)        // Get all orders of a user
	app.Get("/:orderId/invoice", orderHandlers.GetInvoiceOrder) // Get specific order for invoice generation

	// File Server Routes for Invoice Downloads
	app.Static("/invoices", "./invoices")                            // Static file serving
	app.Get("/download/invoice/:filename", handlers.DownloadInvoice) // Download specific invoice

}
