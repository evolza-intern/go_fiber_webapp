package handlers

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/models"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

func CreateListing(ctx *fiber.Ctx) error {
	var (
		newListing models.HousingLocation
		dbPath     = "db/db.json"
	)

	// Ensure db directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}

	// Parse request body into the new listing
	if err := ctx.BodyParser(&newListing); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Load existing listings
	var listings []models.HousingLocation
	if data, err := os.ReadFile(dbPath); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &listings); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse existing listings",
			})
		}
	}

	// Check if ID already exists (if provided manually)
	for _, listing := range listings {
		if listing.ID == newListing.ID && newListing.ID != 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Listing with this ID already exists",
			})
		}
	}

	// If ID is zero or not provided, auto-increment
	if newListing.ID == 0 {
		maxID := 0
		for _, listing := range listings {
			if listing.ID > maxID {
				maxID = listing.ID
			}
		}
		newListing.ID = maxID + 1
	}

	// Append and save
	listings = append(listings, newListing)

	if err := SaveJSONToFile(dbPath, listings); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save listing",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(newListing)
}

func SaveJSONToFile(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}
