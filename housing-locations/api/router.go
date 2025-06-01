package api

import (
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/locations", handlers.GetAllLocations)
	app.Get("/locations/:id", handlers.GetLocationById)
	app.Post("/listing", handlers.CreateListing)
	app.Put("/listing/:id", handlers.UpdateListing)
}
