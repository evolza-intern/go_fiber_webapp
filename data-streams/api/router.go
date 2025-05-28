package api

import (
	"github.com/evolza-intern/go_fiber_webapp/data-streams/handlers"
	"github.com/evolza-intern/go_fiber_webapp/data-streams/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, processor *services.DataProcessor) {
	// API group
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", healthCheck)

	// Data processing routes
	api.Post("/data", handlers.StreamDataHandler(processor))
	api.Post("/data/batch", handlers.BatchDataHandler(processor))
	api.Get("/data/status", handlers.GetProcessingStatus(processor))
	api.Get("/data/results/:id", handlers.GetProcessingResults(processor))
}

// healthCheck returns API health status
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "data-processor-api",
		"version": "1.0.0",
	})
}
