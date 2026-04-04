package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
	}
}

func ErrorResponse(message string) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": message,
	}
}

func ValidationErrorResponse(errors map[string]string) fiber.Map {
	return fiber.Map{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	}
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func Paginated(c *fiber.Ctx, data interface{}, total int64, page, limit int) error {
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta": fiber.Map{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    lastPage,
		},
	})
}
