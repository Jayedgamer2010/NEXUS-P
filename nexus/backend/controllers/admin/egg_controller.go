package admin

import (
	"strconv"

	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type EggController struct {
	eggRepo    *repositories.EggRepository
	serverRepo *repositories.ServerRepository
}

func NewEggController(eggRepo *repositories.EggRepository, serverRepo *repositories.ServerRepository) *EggController {
	return &EggController{eggRepo: eggRepo, serverRepo: serverRepo}
}

func (ctrl *EggController) GetAll(c *fiber.Ctx) error {
	eggs, err := ctrl.eggRepo.FindAll()
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch eggs")
	}

	items := make([]transformers.EggItem, 0, len(eggs))
	for _, e := range eggs {
		items = append(items, transformers.TransformEgg(e))
	}

	return utils.Success(c, items)
}

func (ctrl *EggController) Create(c *fiber.Ctx) error {
	var req requests.CreateEggRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	egg := &models.Egg{
		Author:      req.Author,
		Name:        req.Name,
		Description: req.Description,
		DockerImage: req.DockerImage,
		Startup:     req.Startup,
		ConfigStop:  req.ConfigStop,
	}

	if err := ctrl.eggRepo.Create(egg); err != nil {
		return utils.Error(c, 500, "Failed to create egg")
	}

	return utils.Success(c, transformers.TransformEggDetail(*egg))
}

func (ctrl *EggController) GetOne(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid egg ID")
	}

	egg, err := ctrl.eggRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Egg not found")
	}

	return utils.Success(c, transformers.TransformEggDetail(*egg))
}

func (ctrl *EggController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid egg ID")
	}

	egg, err := ctrl.eggRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Egg not found")
	}

	var req requests.UpdateEggRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if req.Name != "" {
		egg.Name = req.Name
	}
	if req.Author != "" {
		egg.Author = req.Author
	}
	if req.Description != "" {
		egg.Description = req.Description
	}
	if req.DockerImage != "" {
		egg.DockerImage = req.DockerImage
	}
	if req.Startup != "" {
		egg.Startup = req.Startup
	}
	if req.ConfigStop != "" {
		egg.ConfigStop = req.ConfigStop
	}

	if err := ctrl.eggRepo.Update(egg); err != nil {
		return utils.Error(c, 500, "Failed to update egg")
	}

	return utils.Success(c, transformers.TransformEggDetail(*egg))
}

func (ctrl *EggController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid egg ID")
	}

	count := ctrl.eggRepo.CountByEggID(uint(id))
	if count > 0 {
		return utils.Error(c, 422, "Cannot delete egg with active servers")
	}

	if err := ctrl.eggRepo.Delete(uint(id)); err != nil {
		return utils.Error(c, 500, "Failed to delete egg")
	}

	return utils.SuccessMessage(c, "Egg deleted successfully", nil)
}

func (ctrl *EggController) DeleteWithServers(c *fiber.Ctx) error {
	return utils.Error(c, 422, "Cannot delete egg with active servers")
}
