package main

import (
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/api"
	"log"
	"os"

	database "github.com/evolza-intern/go_fiber_webapp/order-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS for microservice communication
	app.Use(cors.New())

	// Connect to MongoDB
	database.ConnectMongo()
	orderCollection := database.GetCollection("go_fiber_eCommerce", "orders")

	// Inject collection into handlers
	handlers.InitOrderHandler(orderCollection)

	// Create invoices directory if it doesn't exist
	if err := os.MkdirAll("./invoices", 0755); err != nil {
		log.Fatal("Failed to create invoices directory:", err)
	}

	api.SetupRoutes(app, orderCollection)

	log.Println("Order Service running on http://localhost:3001")
	log.Fatal(app.Listen(":3001"))
}
