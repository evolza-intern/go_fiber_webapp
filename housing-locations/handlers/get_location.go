package handlers

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/models"
	"github.com/gofiber/fiber/v2"
	"os"
)

func GetLocationById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	// Read the JSON file
	data, err := os.ReadFile("db/db.json")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read JSON file")
	}

	// Unmarshal JSON into a slice
	var locations []models.HousingLocation
	if err := json.Unmarshal(data, &locations); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON")
	}

	// Search for the location by ID
	for _, location := range locations {
		if location.ID == id {
			return ctx.JSON(location)
		}
	}

	// Return 404 if not found
	return ctx.Status(fiber.StatusNotFound).SendString("Location not found")
}
