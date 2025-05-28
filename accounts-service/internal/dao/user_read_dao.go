package dao

import (
	"fmt"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/models"
)

func FindUserByUsername(username string) (*models.User, error) {
	users, err := LoadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func FindUserByID(id int) (*models.User, error) {
	users, err := LoadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
