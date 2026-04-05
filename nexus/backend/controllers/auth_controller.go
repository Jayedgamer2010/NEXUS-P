package controllers

import (
	"errors"
	"nexus/backend/models"
	"nexus/backend/services"
	"nexus/backend/requests"
	"nexus/backend/utils"
	"nexus/backend/transformers"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, token, err := ctrl.authService.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailTaken):
			return utils.Error(c, 422, "The email address is already in use")
		case errors.Is(err, services.ErrUsernameTaken):
			return utils.Error(c, 422, "The username is already in use")
		default:
			return utils.Error(c, 500, "Failed to create account")
		}
	}

	return utils.Success(c, fiber.Map{
		"user":  transformers.TransformUserDetail(*user),
		"token": token,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	user, token, err := ctrl.authService.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return utils.Error(c, 401, "Invalid email or password")
		case errors.Is(err, services.ErrAccountSuspended):
			return utils.Error(c, 403, "Account is suspended")
		default:
			return utils.Error(c, 500, "Failed to login")
		}
	}

	return utils.Success(c, fiber.Map{
		"user":  transformers.TransformUserDetail(*user),
		"token": token,
	})
}

func (ctrl *AuthController) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return utils.Success(c, transformers.TransformUserDetail(*user))
}
