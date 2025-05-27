package main

import (
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/api"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS for microservice communication
	app.Use(cors.New())

	// Initialize in-memory db (if you're still using it alongside MongoDB)
	//db.InitCartStorage()

	// Setup all routes
	api.SetupRoutes(app)

	log.Println("Cart Service running on http://localhost:3002")
	log.Fatal(app.Listen(":3002"))
}
