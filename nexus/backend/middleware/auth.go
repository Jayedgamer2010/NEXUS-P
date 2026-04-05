package middleware

import (
	"strings"

	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Error(c, 401, "Missing authorization header")
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return utils.Error(c, 401, "Invalid authorization format")
		}

		claims, err := utils.ValidateToken(tokenParts[1], cfg)
		if err != nil {
			return utils.Error(c, 401, "Invalid or expired token")
		}

		var user models.User
		if err := database.DB.Where("uuid = ?", claims.UUID).First(&user).Error; err != nil {
			return utils.Error(c, 401, "User not found")
		}

		if user.Suspended {
			return utils.Error(c, 403, "Account is suspended")
		}

		c.Locals("user", &user)
		c.Locals("user_id", user.ID)
		c.Locals("user_uuid", user.UUID)
		return c.Next()
	}
}

func GetUser(c *fiber.Ctx) *models.User {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}
