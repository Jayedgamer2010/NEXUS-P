package middleware

import (
	"strings"

	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthConfig holds JWT configuration
type AuthConfig struct {
	JWTSecret string
	JWTExpire int
}

func Auth(cfg *AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "Authorization header required")
		}

		// Check for Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return utils.Unauthorized(c, "Invalid authorization header format")
		}

		tokenString := parts[1]
		// Create a temporary config for validation
		tempCfg := &config.Config{JWTSecret: cfg.JWTSecret}
		claims, err := utils.ValidateToken(tokenString, tempCfg)
		if err != nil {
			return utils.Unauthorized(c, "Invalid or expired token")
		}

		// Fetch user from database
		var user models.User
		if err := database.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
			return utils.Unauthorized(c, "User not found")
		}

		// Store user in context
		c.Locals("user", user)
		return c.Next()
	}
}

func GetUser(c *fiber.Ctx) *models.User {
	if user, ok := c.Locals("user").(models.User); ok {
		return &user
	}
	return nil
}
