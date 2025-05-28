package handlers

import (
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/dao"
	"github.com/gofiber/fiber/v2"
)

func ToggleUserStatusHandler(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	users, err := dao.LoadUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load users"})
	}

	// Find and toggle user status
	for i, user := range users {
		if user.ID == userID {
			if user.Status == "Active" {
				users[i].Status = "Inactive"
				dao.AddActivityLog("User Deactivated", userID, user.Username, fmt.Sprintf("User %s was deactivated", user.Username))
			} else {
				users[i].Status = "Active"
				dao.AddActivityLog("User Activated", userID, user.Username, fmt.Sprintf("User %s was activated", user.Username))
			}

			err = dao.SaveUsers(users)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to save users"})
			}

			return c.JSON(fiber.Map{"message": "User status updated successfully"})
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}
