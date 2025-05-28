package handlers

import (
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/dao"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/models"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"time"
)

func StatsHandler(c *fiber.Ctx) error {
	stats := getSystemStats()
	return c.JSON(stats)
}

// Get system statistics
func getSystemStats() models.SystemStats {
	users, _ := dao.LoadUsers()
	totalUsers := len(users)
	activeUsers := 0
	inactiveUsers := 0

	for _, user := range users {
		if user.Status == "Active" {
			activeUsers++
		} else {
			inactiveUsers++
		}
	}

	// Get memory stats
	var startTime = time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsage := fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024)

	// Calculate uptime
	uptime := time.Since(startTime).Round(time.Second).String()

	return models.SystemStats{
		TotalUsers:    totalUsers,
		ActiveUsers:   activeUsers,
		InactiveUsers: inactiveUsers,
		MemoryUsage:   memoryUsage,
		CPUUsage:      0.0, // Simplified - would need additional package for real CPU usage
		Uptime:        uptime,
	}
}
