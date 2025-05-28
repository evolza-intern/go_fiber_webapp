package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/dao"
	"github.com/gofiber/fiber/v2"
)

func ActivityHandler(c *fiber.Ctx) error {
	logs, err := dao.LoadActivityLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load activity logs"})
	}

	// Return latest 20 logs, sorted by timestamp
	if len(logs) > 20 {
		logs = logs[len(logs)-20:]
	}

	// Reverse to show newest first
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	return c.JSON(logs)
}
