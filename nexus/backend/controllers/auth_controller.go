package controllers

import (
	"nexus/backend/models"
	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController(userSvc *services.UserService) *AuthController {
	return &AuthController{userService: userSvc}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	user, token, err := ac.userService.Register(req)
	if err != nil {
		if err == services.ErrEmailTaken {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse("Email already in use"))
		}
		if err == services.ErrUsernameTaken {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse("Username already in use"))
		}
		return utils.InternalError(c, "Failed to create account")
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.Map{
		"user":  transformers.TransformUser(user.Sanitize()),
		"token": token,
	}))
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	user, token, err := ac.userService.Login(req)
	if err != nil {
		return utils.Unauthorized(c, "Invalid credentials")
	}

	return utils.Success(c, fiber.Map{
		"user":  transformers.TransformUser(user.Sanitize()),
		"token": token,
	}, "Login successful")
}

func (ac *AuthController) Me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	return utils.Success(c, transformers.TransformUser(user.Sanitize()), "User retrieved")
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	return utils.Success(c, nil, "Logged out")
}
