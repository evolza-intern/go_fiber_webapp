package dao

import (
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/models"
	"time"
)

func AddActivityLog(action string, userID int, username string, details string) {
	lastActivityID++
	log := models.ActivityLog{
		ID:        lastActivityID,
		Action:    action,
		UserID:    userID,
		Username:  username,
		Timestamp: time.Now(),
		Details:   details,
	}

	logs, _ := LoadActivityLogs()
	logs = append(logs, log)
	err := saveActivityLogs(logs)
	if err != nil {
		return
	}
}
