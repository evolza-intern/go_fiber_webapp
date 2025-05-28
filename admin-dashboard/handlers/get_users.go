package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/dao"
	"github.com/gofiber/fiber/v2"
)

func UsersHandler(c *fiber.Ctx) error {
	users, err := dao.LoadUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load users"})
	}
	return c.JSON(users)
}
