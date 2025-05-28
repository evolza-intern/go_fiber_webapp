package api

import (
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.DashboardHandler)
	app.Get("/api/stats", handlers.StatsHandler)
	app.Get("/api/users", handlers.UsersHandler)
	app.Get("/api/activity", handlers.ActivityHandler)
	app.Post("/api/users/:id/toggle", handlers.ToggleUserStatusHandler)
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No Content
	})
}
