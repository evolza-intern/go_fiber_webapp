package main

import (
	"errors"
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/api"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/dao"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize database (load users to create default admin if needed)
	_, err := dao.LoadUsers()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	fmt.Println("Database initialized. Default admin user: username='admin', password='admin123'")

	// Setup routes
	api.SetupRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
