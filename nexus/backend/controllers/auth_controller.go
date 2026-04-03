package controllers

import (
	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

var cfg *config.Config

func Init(c *config.Config) {
	cfg = c
}

type AuthController struct{}

// Register handles user registration
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // admin or client, default client
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Username, email, and password are required")
	}

	// Validate role (only allow admin if it's the first user or explicitly set)
	role := "client"
	if req.Role == "admin" {
		// Check if any admin exists
		var adminCount int64
		database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount == 0 {
			role = "admin" // First user can be admin
		} else {
			return utils.BadRequest(c, "Cannot self-assign admin role")
		}
	}

	// Generate UUID
	userUUID := utils.GenerateUUID()

	// Create user
	user := models.User{
		UUID:     userUUID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
		Coins:    0,
		RootAdmin: role == "admin",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.InternalError(c, "Failed to create user")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.UUID, user.Role, cfg)
	if err != nil {
		return utils.InternalError(c, "Failed to generate token")
	}

	return utils.Success(c, fiber.Map{
		"token": token,
		"user":  utils.FromUser(&user),
	}, "Registration successful")
}

// Login handles user login
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Email and password are required")
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return utils.Unauthorized(c, "Invalid credentials")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return utils.Unauthorized(c, "Invalid credentials")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.UUID, user.Role, cfg)
	if err != nil {
		return utils.InternalError(c, "Failed to generate token")
	}

	return utils.Success(c, fiber.Map{
		"token": token,
		"user":  utils.FromUser(&user),
	}, "Login successful")
}

// Me returns the current user's info
func (ac *AuthController) Me(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	return utils.Success(c, utils.FromUser(&user), "User retrieved")
}

// Logout (optional - JWT is stateless, but included for completeness)
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	// JWT is stateless, so we just return success
	// In a production system, you might want to blacklist the token
	return utils.Success(c, nil, "Logged out successfully")
}
