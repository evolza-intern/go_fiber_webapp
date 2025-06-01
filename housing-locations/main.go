package main

import (
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api.SetupRoutes(app)

	err := app.Listen(":3003")
	if err != nil {
		return
	}
}
