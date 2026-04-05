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
	userRepo    *repositories.UserRepository
	userService *services.UserService
}

func NewUserController(
	userRepo *repositories.UserRepository,
	userService *services.UserService,
) *UserController {
	return &UserController{userRepo: userRepo, userService: userService}
}

func (ctrl *UserController) GetAll(c *fiber.Ctx) error {
	page := utils.GetPage(c)
	perPage := utils.GetPerPage(c)
	search := c.Query("search", "")

	users, total, err := ctrl.userRepo.FindAll(page, perPage, search)
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch users")
	}

	return utils.PaginatedResponse(c, transformers.TransformUsers(users), utils.BuildMeta(total, page, perPage))
}

func (ctrl *UserController) Create(c *fiber.Ctx) error {
	var req requests.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if _, err := ctrl.userRepo.FindByEmail(req.Email); err == nil {
		return utils.Error(c, 422, "Email already in use")
	}
	if _, err := ctrl.userRepo.FindByUsername(req.Username); err == nil {
		return utils.Error(c, 422, "Username already in use")
	}

	role := "client"
	if req.Role == "admin" {
		role = "admin"
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      role,
		RootAdmin: req.RootAdmin,
		Coins:     req.Coins,
	}
	if err := user.HashPassword(req.Password); err != nil {
		return utils.Error(c, 500, "Failed to hash password")
	}

	if err := ctrl.userRepo.Create(user); err != nil {
		return utils.Error(c, 500, "Failed to create user")
	}

	return utils.Success(c, transformers.TransformUserDetail(*user))
}

func (ctrl *UserController) GetOne(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID")
	}

	user, err := ctrl.userRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "User not found")
	}

	return utils.Success(c, transformers.TransformUserDetail(*user))
}

func (ctrl *UserController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID")
	}

	user, err := ctrl.userRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "User not found")
	}

	var req requests.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		// Check email uniqueness if changed
		existing, err := ctrl.userRepo.FindByEmail(req.Email)
		if err == nil && existing.ID != user.ID {
			return utils.Error(c, 422, "Email already in use")
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.RootAdmin != nil {
		user.RootAdmin = *req.RootAdmin
	}
	if req.Coins != nil {
		user.Coins = *req.Coins
	}
	if req.Suspended != nil {
		user.Suspended = *req.Suspended
	}

	if err := ctrl.userService.Update(user, req.Password); err != nil {
		return utils.Error(c, 500, "Failed to update user")
	}

	return utils.Success(c, transformers.TransformUserDetail(*user))
}

func (ctrl *UserController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid user ID")
	}

	if err := ctrl.userService.Delete(uint(id)); err != nil {
		if err == services.ErrUserHasServers {
			return utils.Error(c, 422, "Cannot delete user with existing servers")
		}
		return utils.Error(c, 500, "Failed to delete user")
	}

	return utils.SuccessMessage(c, "User deleted successfully", nil)
}
