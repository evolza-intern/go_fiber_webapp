package handlers

import (
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/dao"
	"github.com/gofiber/fiber/v2"
)

// ProfileHandler handles user profile requests
func ProfileHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	username := c.Locals("username").(string)

	user, err := dao.FindUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"message":  fmt.Sprintf("Hello, %s! This is your protected profile.", username),
	})
}
