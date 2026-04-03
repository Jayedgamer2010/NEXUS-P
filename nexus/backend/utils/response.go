package utils

import (
	"nexus/backend/models"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return Error(c, fiber.StatusUnauthorized, message)
}

func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return Error(c, fiber.StatusForbidden, message)
}

func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

// PaginatedResponse for list endpoints
type PaginatedResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
	Total   int64         `json:"total"`
	Page    int           `json:"page"`
	Limit   int           `json:"limit"`
}

func Paginated(c *fiber.Ctx, data []interface{}, total int64, page, limit int) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: "",
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// Standard JSON structure for user responses (hide sensitive fields)
type UserResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Coins     int    `json:"coins"`
	RootAdmin bool   `json:"root_admin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FromUser(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Coins:     user.Coins,
		RootAdmin: user.RootAdmin,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
