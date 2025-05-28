package middleware

import (
	"github.com/evolza-intern/go_fiber_webapp/accounts-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Missing authorization token",
		})
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Store user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)

	return c.Next()
}
