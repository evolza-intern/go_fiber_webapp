package api

import (
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

//func SetupRoutes(app *fiber.App, productCollection *mongo.Collection) {
//	// Initialize handlers with the collection
//	productHandlers := handlers.NewProductHandlers(productCollection)
//
//	// Setup routes using the initialized handlers
//	app.Post("/products", productHandlers.CreateProduct)
//	app.Get("/products", productHandlers.GetProducts)
//	app.Get("/products/:id", productHandlers.GetProduct)
//	app.Put("/products/:id", productHandlers.UpdateProduct)
//	app.Delete("/products/:id", productHandlers.DeleteProduct)
//}

func SetupRoutes(app *fiber.App, productCollection *mongo.Collection) {
	// Initialize handlers with the collection
	productHandlers := handlers.NewProductHandlers(productCollection)

	// Global rate limiter - applies to all routes
	globalLimiter := limiter.New(limiter.Config{
		Max:        100,             // Maximum 100 requests
		Expiration: 1 * time.Minute, // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"limit":   "100 requests per minute",
			})
		},
	})

	// Apply global rate limiter to all routes
	app.Use(globalLimiter)

	// Stricter rate limiter for write operations (POST, PUT, DELETE)
	writeOperationsLimiter := limiter.New(limiter.Config{
		Max:        20,              // Maximum 20 requests
		Expiration: 1 * time.Minute, // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Write operation rate limit exceeded",
				"message": "Too many write operations. Please try again later.",
				"limit":   "20 write operations per minute",
			})
		},
	})

	// More lenient rate limiter for read operations (GET)
	readOperationsLimiter := limiter.New(limiter.Config{
		Max:        50,              // Maximum 50 requests
		Expiration: 1 * time.Minute, // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Read operation rate limit exceeded",
				"message": "Too many read operations. Please try again later.",
				"limit":   "50 read operations per minute",
			})
		},
	})

	// Product routes with specific rate limiting
	api := app.Group("/products")

	// Read operations with lenient rate limiting
	api.Get("/", readOperationsLimiter, productHandlers.GetProducts)
	api.Get("/:id", readOperationsLimiter, productHandlers.GetProduct)

	// Write operations with stricter rate limiting
	api.Post("/", writeOperationsLimiter, productHandlers.CreateProduct)
	api.Put("/:id", writeOperationsLimiter, productHandlers.UpdateProduct)
	api.Delete("/:id", writeOperationsLimiter, productHandlers.DeleteProduct)
}
