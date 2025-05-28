package dao

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/models"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/utils"
	"os"
	"path/filepath"
)

var (
	dbPath = "db/users.json"
)

func LoadUsers() ([]models.User, error) {
	// Create db directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create empty users file with default admin user
		defaultUsers := []models.User{
			{
				ID:       1,
				Username: "admin",
				Password: utils.HashPassword("admin123"), // Default password
			},
		}
		return defaultUsers, SaveUsers(defaultUsers)
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = json.Unmarshal(data, &users)
	return users, err
}

func SaveUsers(users []models.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbPath, data, 0644)
}
