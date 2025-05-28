package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ProtectedHandler(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Hello, %s! This is a protected route.", username),
		"data":    "This is sensitive data that only authenticated users can see",
	})
}
