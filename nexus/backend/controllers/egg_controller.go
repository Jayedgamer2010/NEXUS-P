package controllers

import (
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type EggController struct{}

// GetAll returns all eggs
func (ec *EggController) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var eggs []models.Egg
	var total int64

	database.DB.Model(&models.Egg{}).Count(&total)
	if err := database.DB.Offset(offset).Limit(limit).Find(&eggs).Error; err != nil {
		return utils.InternalError(c, "Failed to fetch eggs")
	}

	// Convert []models.Egg to []interface{}
	data := make([]interface{}, len(eggs))
	for i, egg := range eggs {
		data[i] = egg
	}

	return utils.Paginated(c, data, total, page, limit)
}

// GetByID returns a specific egg
func (ec *EggController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var egg models.Egg
	if err := database.DB.First(&egg, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	return utils.Success(c, egg, "Egg retrieved")
}

// Create creates a new egg (admin only)
func (ec *EggController) Create(c *fiber.Ctx) error {
	var req struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		DockerImage    string `json:"docker_image"`
		StartupCommand string `json:"startup_command"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" || req.DockerImage == "" || req.StartupCommand == "" {
		return utils.BadRequest(c, "Name, docker_image, and startup_command are required")
	}

	eggUUID := utils.GenerateUUID()

	egg := models.Egg{
		UUID:           eggUUID,
		Name:           req.Name,
		Description:    req.Description,
		DockerImage:    req.DockerImage,
		StartupCommand: req.StartupCommand,
	}

	if err := database.DB.Create(&egg).Error; err != nil {
		return utils.InternalError(c, "Failed to create egg")
	}

	return utils.Success(c, egg, "Egg created")
}

// Update updates an egg
func (ec *EggController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var egg models.Egg
	if err := database.DB.First(&egg, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	var req struct {
		Name           *string `json:"name"`
		Description    *string `json:"description"`
		DockerImage    *string `json:"docker_image"`
		StartupCommand *string `json:"startup_command"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name != nil {
		egg.Name = *req.Name
	}
	if req.Description != nil {
		egg.Description = *req.Description
	}
	if req.DockerImage != nil {
		egg.DockerImage = *req.DockerImage
	}
	if req.StartupCommand != nil {
		egg.StartupCommand = *req.StartupCommand
	}

	if err := database.DB.Save(&egg).Error; err != nil {
		return utils.InternalError(c, "Failed to update egg")
	}

	return utils.Success(c, egg, "Egg updated")
}

// Delete deletes an egg
func (ec *EggController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var egg models.Egg
	if err := database.DB.First(&egg, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	// Check if any servers use this egg
	var serverCount int64
	database.DB.Model(&models.Server{}).Where("egg_id = ?", egg.ID).Count(&serverCount)
	if serverCount > 0 {
		return utils.BadRequest(c, "Cannot delete egg with active servers")
	}

	if err := database.DB.Delete(&egg).Error; err != nil {
		return utils.InternalError(c, "Failed to delete egg")
	}

	return utils.Success(c, nil, "Egg deleted")
}
