package handlers

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/housing-locations/models"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
)

func UpdateListing(ctx *fiber.Ctx) error {
	var (
		updatedListing models.HousingLocation
		dbPath         = "db/db.json"
	)

	// Get the ID from URL parameters
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid listing ID",
		})
	}

	// Parse request body into the updated listing
	if err := ctx.BodyParser(&updatedListing); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Load existing listings
	var listings []models.HousingLocation
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read listings database",
		})
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &listings); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse existing listings",
			})
		}
	}

	// Find and update the listing
	found := false
	for i, listing := range listings {
		if listing.ID == id {
			// Preserve the original ID
			updatedListing.ID = id
			listings[i] = updatedListing
			found = true
			break
		}
	}

	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Listing not found",
		})
	}

	// Save updated listings
	if err := SaveJSONToFile(dbPath, listings); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save updated listing",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(updatedListing)
}
