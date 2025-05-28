package main

import (
	"github.com/evolza-intern/go_fiber_webapp/data-streams/api"
	"github.com/evolza-intern/go_fiber_webapp/data-streams/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Data-Processor-API",
		AppName:       "Data Processing API v1.0.0",
		BodyLimit:     50 * 1024 * 1024, // 50MB for large data streams
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Initialize data processor service
	processor := services.NewDataProcessor()

	// Setup routes
	api.SetupRoutes(app, processor)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		err := app.Shutdown()
		if err != nil {
			return
		}
		processor.Shutdown()
	}()

	// Start server
	log.Println("🚀 Server starting on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
