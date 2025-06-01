package handlers

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/models"
	"github.com/gofiber/fiber/v2"
	"os"
)

func GetAllLocations(c *fiber.Ctx) error {
	data, err := os.ReadFile("db/db.json")
	if err != nil {
		return c.Status(500).SendString("Failed to read JSON file")
	}

	var locations []models.HousingLocation
	if err := json.Unmarshal(data, &locations); err != nil {
		return c.Status(500).SendString("Failed to parse JSON")
	}

	return c.JSON(locations)
}
