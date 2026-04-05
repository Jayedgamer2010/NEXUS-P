package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Data: data})
}

func SuccessMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(APIResponse{Success: false, Message: message})
}

func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return c.Status(422).JSON(APIResponse{Success: false, Message: "Validation failed", Errors: errors})
}

func PaginatedResponse(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.JSON(APIResponse{
		Success: true,
		Data: struct {
			Data interface{} `json:"data"`
			Meta interface{} `json:"meta"`
		}{Data: data, Meta: meta},
	})
}
