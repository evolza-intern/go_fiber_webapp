package dao

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/models"
	"os"
)

func LoadActivityLogs() ([]models.ActivityLog, error) {
	data, err := os.ReadFile(activityFile)
	if err != nil {
		return nil, err
	}
	var logs []models.ActivityLog
	err = json.Unmarshal(data, &logs)
	return logs, err
}
