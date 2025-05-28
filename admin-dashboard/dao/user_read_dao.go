package dao

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/admin-dashboard/models"
	"os"
)

func LoadUsers() ([]models.User, error) {
	data, err := os.ReadFile(usersFile)
	if err != nil {
		return nil, err
	}
	var users []models.User
	err = json.Unmarshal(data, &users)
	return users, err
}
