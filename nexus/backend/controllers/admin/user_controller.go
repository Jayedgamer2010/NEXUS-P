package admin

import (
	"strconv"

	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	svc      *services.UserService
	allocRepo *repositories.AllocationRepository
}

func NewUserController(svc *services.UserService, allocRepo *repositories.AllocationRepository) *UserController {
	return &UserController{svc: svc, allocRepo: allocRepo}
}

func (uc *UserController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	users, total, err := uc.svc.All(page, limit)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch users")
	}

	return utils.Paginated(c, transformers.TransformUsers(users), total, page, limit)
}

func (uc *UserController) Create(c *fiber.Ctx) error {
	var req requests.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	user, err := uc.svc.Create(req)
	if err != nil {
		if err == services.ErrEmailTaken {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse("Email already in use"))
		}
		if err == services.ErrUsernameTaken {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse("Username already in use"))
		}
		return utils.InternalError(c, "Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(transformers.TransformUser(user.Sanitize())))
}

func (uc *UserController) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid user ID")
	}

	user, err := uc.svc.FindByID(uint(id))
	if err != nil {
		return utils.InternalError(c, "Failed to fetch user")
	}
	if user == nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	return utils.Success(c, transformers.TransformUser(user.Sanitize()), "User retrieved")
}

func (uc *UserController) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid user ID")
	}

	user, err := uc.svc.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	var req requests.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := uc.svc.Update(user, req); err != nil {
		return utils.InternalError(c, "Failed to update user")
	}

	return utils.Success(c, transformers.TransformUser(user.Sanitize()), "User updated")
}

func (uc *UserController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid user ID")
	}

	user, err := uc.svc.FindByID(uint(id))
	if err != nil || user == nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	if err := uc.svc.Delete(user); err != nil {
		return utils.InternalError(c, "Failed to delete user")
	}

	return utils.Success(c, nil, "User deleted")
}

func (uc *UserController) GetMyServers(c *fiber.Ctx) error {
	return nil // placeholder - client server logic in client controller
}

func (uc *UserController) GetMyServer(c *fiber.Ctx) error {
	return nil
}

func (uc *UserController) GetResources(c *fiber.Ctx) error {
	return nil
}

func (uc *UserController) GetMyServerByUUID(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	_ = user // used in final implementation
	return utils.Success(c, nil, "OK")
}
