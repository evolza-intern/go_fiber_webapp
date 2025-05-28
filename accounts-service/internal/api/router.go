package api

import (
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/handlers"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Go Fiber JWT Authentication API",
			"endpoints": fiber.Map{
				"POST /login":    "Login with username/password",
				"POST /logout":   "Logout (requires auth)",
				"GET /profile":   "Get user profile (requires auth)",
				"GET /protected": "Access protected resource (requires auth)",
			},
		})
	})

	app.Post("/login", handlers.LoginHandler)

	// Protected routes group
	protected := app.Group("/", middleware.AuthMiddleware)
	protected.Post("/logout", handlers.LogoutHandler)
	protected.Get("/profile", handlers.ProfileHandler)
	protected.Get("/protected", handlers.ProtectedHandler)
}
