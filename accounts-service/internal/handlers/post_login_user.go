package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/dao"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/models"
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// LoginHandler handles user login
func LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Find user
	user, err := dao.FindUserByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return response
	response := models.LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username

	return c.JSON(response)
}
