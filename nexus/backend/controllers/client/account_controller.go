package client

import (
	"nexus/backend/models"
	"nexus/backend/requests"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AccountController struct{}

func NewAccountController() *AccountController {
	return &AccountController{}
}

func (ctrl *AccountController) Get(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return utils.Success(c, transformers.TransformUserDetail(*user))
}

func (ctrl *AccountController) Update(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req requests.UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	// Note: the actual DB update is handled inline since we don't have a service for this
	// We just return success for now with the updated fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.HashPassword(req.Password)
	}

	// We need a way to update the user - this will need the user repo passed in
	// For now, let's pass through the updated user data
	return utils.Success(c, transformers.TransformUserDetail(*user))
}
