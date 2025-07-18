package handlers

import "github.com/gofiber/fiber/v2"

func LogoutHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
