package admin

import (
	"strconv"

	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type EggController struct {
	svc *services.EggService
}

func NewEggController(svc *services.EggService) *EggController {
	return &EggController{svc: svc}
}

func (ec *EggController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	eggs, total, err := ec.svc.All(page, limit)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch eggs")
	}

	return utils.Paginated(c, transformers.TransformEggs(eggs), total, page, limit)
}

func (ec *EggController) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid egg ID")
	}

	egg, err := ec.svc.FindByID(uint(id))
	if err != nil || egg == nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	return utils.Success(c, transformers.TransformEgg(*egg), "Egg retrieved")
}

func (ec *EggController) Create(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		DockerImage string `json:"docker_image" validate:"required"`
		Startup     string `json:"startup" validate:"required"`
		Author      string `json:"author"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name == "" || req.DockerImage == "" || req.Startup == "" {
		return utils.BadRequest(c, "name, docker_image, and startup are required")
	}

	egg, err := ec.svc.Create(req.Name, req.Description, req.DockerImage, req.Startup, req.Author)
	if err != nil {
		return utils.InternalError(c, "Failed to create egg")
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(transformers.TransformEgg(*egg)))
}

func (ec *EggController) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid egg ID")
	}

	egg, err := ec.svc.FindByID(uint(id))
	if err != nil || egg == nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		DockerImage string `json:"docker_image"`
		Startup     string `json:"startup"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := ec.svc.Update(egg, req.Name, req.Description, req.DockerImage, req.Startup); err != nil {
		return utils.InternalError(c, "Failed to update egg")
	}

	return utils.Success(c, transformers.TransformEgg(*egg), "Egg updated")
}

func (ec *EggController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid egg ID")
	}

	egg, err := ec.svc.FindByID(uint(id))
	if err != nil || egg == nil {
		return utils.Error(c, fiber.StatusNotFound, "Egg not found")
	}

	if err := ec.svc.Delete(egg); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Egg deleted")
}
