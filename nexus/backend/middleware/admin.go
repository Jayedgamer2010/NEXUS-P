package middleware

import (
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Admin(c *fiber.Ctx) error {
	user := GetUser(c)
	if user == nil {
		return utils.Unauthorized(c, "User not authenticated")
	}

	if user.Role != "admin" {
		return utils.Forbidden(c, "Admin access required")
	}

	return c.Next()
}
