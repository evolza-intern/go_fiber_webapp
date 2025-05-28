package api

import (
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// API v1 group
	v1 := app.Group("/api/v1")

	// Users service routes
	users := v1.Group("/users")
	users.Get("/", handlers.HandleGetUsers)
	users.Get("/:id", handlers.HandleGetUser)
	users.Post("/", handlers.HandleCreateUser)
	users.Put("/:id", handlers.HandleUpdateUser)
	users.Delete("/:id", handlers.HandleDeleteUser)

	// Products service routes
	products := v1.Group("/products")
	products.Get("/", handlers.HandleGetProducts)
	products.Get("/:id", handlers.HandleGetProduct)
	products.Post("/", handlers.HandleCreateProduct)
	products.Put("/:id", handlers.HandleUpdateProduct)
	products.Delete("/:id", handlers.HandleDeleteProduct)

	// Orders service routes
	orders := v1.Group("/orders")
	orders.Get("/", handlers.HandleGetOrders)
	orders.Get("/:id", handlers.HandleGetOrder)
	orders.Post("/", handlers.HandleCreateOrder)
	orders.Put("/:id", handlers.HandleUpdateOrder)
	orders.Delete("/:id", handlers.HandleDeleteOrder)

	// Notifications service routes
	notifications := v1.Group("/notifications")
	notifications.Get("/", handlers.HandleGetNotifications)
	notifications.Post("/", handlers.HandleSendNotification)
}
