package controllers

import (
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

// GetAll returns all users (admin only)
func (uc *UserController) GetAll(c *fiber.Ctx) error {
	// Pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	database.DB.Model(&models.User{}).Count(&total)
	if err := database.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return utils.InternalError(c, "Failed to fetch users")
	}

	// Convert to response (exclude passwords)
	responseUsers := make([]interface{}, len(users))
	for i, user := range users {
		responseUsers[i] = utils.FromUser(&user)
	}

	return utils.Paginated(c, responseUsers, total, page, limit)
}

// GetByID returns a specific user
func (uc *UserController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	return utils.Success(c, utils.FromUser(&user), "User retrieved")
}

// Create creates a new user (admin only)
func (uc *UserController) Create(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Coins    int    `json:"coins"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Username, email, and password are required")
	}

	// Validate role
	if req.Role != "admin" && req.Role != "client" {
		return utils.BadRequest(c, "Role must be 'admin' or 'client'")
	}

	userUUID := utils.GenerateUUID()

	user := models.User{
		UUID:      userUUID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		Coins:     req.Coins,
		RootAdmin: req.Role == "admin",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.InternalError(c, "Failed to create user")
	}

	return utils.Success(c, utils.FromUser(&user), "User created")
}

// Update updates a user
func (uc *UserController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"`
		Coins     int    `json:"coins"`
		RootAdmin bool   `json:"root_admin"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Role != "" {
		if req.Role != "admin" && req.Role != "client" {
			return utils.BadRequest(c, "Role must be 'admin' or 'client'")
		}
		user.Role = req.Role
		user.RootAdmin = req.Role == "admin"
	}
	if req.Coins != 0 {
		user.Coins = req.Coins
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return utils.InternalError(c, "Failed to update user")
	}

	return utils.Success(c, utils.FromUser(&user), "User updated")
}

// Delete deletes a user
func (uc *UserController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "User not found")
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return utils.InternalError(c, "Failed to delete user")
	}

	return utils.Success(c, nil, "User deleted")
}
