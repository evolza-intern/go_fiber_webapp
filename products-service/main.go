package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	database "github.com/evolza-intern/go_fiber_webapp/products-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/handlers"

	"github.com/gofiber/fiber/v2"

	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/api"
)

func main() {
	app := fiber.New()

	// Connect to MongoDB
	database.ConnectMongo()

	// Get collection
	productCollection := database.GetCollection("go_fiber_eCommerce", "products")

	// Inject collection into handlers
	handlers.InitProductHandler(productCollection)

	// Setup routes
	api.SetupRoutes(app)

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		database.DisconnectMongo()
		app.Shutdown()
	}()

	log.Println("Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
