package main

import (
	"errors"
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/api"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/dao"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func main() {
	if err := dao.InitializeData(); err != nil {
		log.Fatal("Failed to initialize data:", err)
	}

	// Create Fiber app
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

	api.SetupRoutes(app)

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Admin Dashboard starting on port %s\n", port)
	fmt.Printf("Open http://localhost:%s to view the dashboard\n", port)
	log.Fatal(app.Listen(":" + port))

}
