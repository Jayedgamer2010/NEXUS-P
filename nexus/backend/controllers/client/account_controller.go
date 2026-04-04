package client

import (
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AccountController struct {
	userRepo *repositories.UserRepository
}

func NewAccountController(userRepo *repositories.UserRepository) *AccountController {
	return &AccountController{userRepo: userRepo}
}

func (ac *AccountController) Get(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	return utils.Success(c, transformers.TransformUser(user.Sanitize()), "Account retrieved")
}

func (ac *AccountController) Update(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	var req requests.UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.NameFirst != nil {
		user.NameFirst = *req.NameFirst
	}
	if req.NameLast != nil {
		user.NameLast = *req.NameLast
	}
	if req.Language != nil {
		user.Language = *req.Language
	}

	if err := ac.userRepo.Update(&user); err != nil {
		return utils.InternalError(c, "Failed to update account")
	}

	return utils.Success(c, transformers.TransformUser(user.Sanitize()), "Account updated")
}
