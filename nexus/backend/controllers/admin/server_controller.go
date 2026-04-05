package admin

import (
	"strconv"

	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type ServerController struct {
	serverRepo  *repositories.ServerRepository
	serverSvc   *services.ServerService
}

func NewServerController(
	serverRepo *repositories.ServerRepository,
	serverSvc *services.ServerService,
) *ServerController {
	return &ServerController{serverRepo: serverRepo, serverSvc: serverSvc}
}

func (ctrl *ServerController) GetAll(c *fiber.Ctx) error {
	page := utils.GetPage(c)
	perPage := utils.GetPerPage(c)

	servers, total, err := ctrl.serverRepo.FindAll(page, perPage)
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch servers")
	}

	return utils.PaginatedResponse(c, transformers.TransformServers(servers), utils.BuildMeta(total, page, perPage))
}

func (ctrl *ServerController) Create(c *fiber.Ctx) error {
	var req requests.CreateServerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	server, err := ctrl.serverSvc.Create(req)
	if err != nil {
		return utils.Error(c, 422, err.Error())
	}

	return utils.Success(c, transformers.TransformServerDetail(*server))
}

func (ctrl *ServerController) GetOne(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	server, err := ctrl.serverRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Server not found")
	}

	return utils.Success(c, transformers.TransformServerDetail(*server))
}

func (ctrl *ServerController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	server, err := ctrl.serverRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Server not found")
	}

	var req requests.UpdateServerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if req.Name != "" {
		server.Name = req.Name
	}
	if req.Description != "" {
		server.Description = req.Description
	}
	if req.Memory != nil {
		server.Memory = *req.Memory
	}
	if req.Disk != nil {
		server.Disk = *req.Disk
	}
	if req.CPU != nil {
		server.CPU = *req.CPU
	}
	if req.Swap != nil {
		server.Swap = *req.Swap
	}

	if err := ctrl.serverRepo.Update(server); err != nil {
		return utils.Error(c, 500, "Failed to update server")
	}

	return utils.Success(c, transformers.TransformServerDetail(*server))
}

func (ctrl *ServerController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	if err := ctrl.serverSvc.Delete(uint(id)); err != nil {
		return utils.Error(c, 500, "Failed to delete server")
	}

	return utils.SuccessMessage(c, "Server deleted successfully", nil)
}

func (ctrl *ServerController) PowerAction(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	var req requests.PowerActionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if err := ctrl.serverSvc.PowerAction(uint(id), req.Action); err != nil {
		return utils.Error(c, 502, "Failed to send power action: "+err.Error())
	}

	return utils.SuccessMessage(c, "Power action sent: "+req.Action, nil)
}

func (ctrl *ServerController) Suspend(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	if err := ctrl.serverSvc.Suspend(uint(id)); err != nil {
		return utils.Error(c, 500, "Failed to suspend server")
	}

	return utils.SuccessMessage(c, "Server suspended", nil)
}

func (ctrl *ServerController) Unsuspend(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid server ID")
	}

	if err := ctrl.serverSvc.Unsuspend(uint(id)); err != nil {
		return utils.Error(c, 500, "Failed to unsuspend server")
	}

	return utils.SuccessMessage(c, "Server unsuspended", nil)
}
