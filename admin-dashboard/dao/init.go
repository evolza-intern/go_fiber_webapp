package dao

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/models"
	"os"
	"time"
)

var (
	usersFile      = "data/users.json"
	activityFile   = "data/activity.json"
	lastActivityID = 0
)

// Initialize sample data
func InitializeData() error {

	// Create data directory
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Initialize users if file doesn't exist
	if _, err := os.Stat(usersFile); os.IsNotExist(err) {
		sampleUsers := []models.User{
			{ID: 1, Username: "admin", Email: "admin@example.com", Role: "Admin", Status: "Active", CreatedAt: time.Now().AddDate(0, -2, 0), LastLogin: time.Now().Add(-1 * time.Hour)},
			{ID: 2, Username: "john_doe", Email: "john@example.com", Role: "User", Status: "Active", CreatedAt: time.Now().AddDate(0, -1, -15), LastLogin: time.Now().Add(-2 * time.Hour)},
			{ID: 3, Username: "jane_smith", Email: "jane@example.com", Role: "Moderator", Status: "Active", CreatedAt: time.Now().AddDate(0, -1, -10), LastLogin: time.Now().Add(-30 * time.Minute)},
			{ID: 4, Username: "inactive_user", Email: "inactive@example.com", Role: "User", Status: "Inactive", CreatedAt: time.Now().AddDate(0, -3, 0), LastLogin: time.Now().AddDate(0, 0, -7)},
		}
		err := SaveUsers(sampleUsers)
		if err != nil {
			return err
		}
	}

	// Initialize activity logs if file doesn't exist
	if _, err := os.Stat(activityFile); os.IsNotExist(err) {
		sampleLogs := []models.ActivityLog{
			{ID: 1, Action: "User Login", UserID: 1, Username: "admin", Timestamp: time.Now().Add(-2 * time.Hour), Details: "Admin user logged in"},
			{ID: 2, Action: "User Created", UserID: 2, Username: "system", Timestamp: time.Now().Add(-1 * time.Hour), Details: "New user john_doe created"},
			{ID: 3, Action: "Profile Updated", UserID: 3, Username: "jane_smith", Timestamp: time.Now().Add(-30 * time.Minute), Details: "Profile information updated"},
		}
		lastActivityID = 3
		err := saveActivityLogs(sampleLogs)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveUsers(users []models.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFile, data, 0644)
}

func saveActivityLogs(logs []models.ActivityLog) error {
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(activityFile, data, 0644)
}
