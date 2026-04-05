package middleware

import (
	"github.com/gofiber/fiber/v2"
	"nexus/backend/utils"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUser(c)
		if user == nil {
			return utils.Error(c, 401, "Unauthorized")
		}

		if !user.IsAdmin() {
			return utils.Error(c, 403, "Admin access required")
		}

		return c.Next()
	}
}
