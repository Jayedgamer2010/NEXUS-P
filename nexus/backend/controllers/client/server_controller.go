package client

import (
	"nexus/backend/models"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type ServerController struct {
	svc          *services.ServerService
	wingsService *services.WingsService
}

func NewServerController(svc *services.ServerService, wingsSvc *services.WingsService) *ServerController {
	return &ServerController{svc: svc, wingsService: wingsSvc}
}

func (sc *ServerController) GetMyServers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	servers, err := sc.svc.FindByUserID(user.ID)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch servers")
	}

	result := transformers.TransformServers(servers)
	return utils.Success(c, result, "Servers retrieved")
}

func (sc *ServerController) GetMyServer(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	uuid := c.Params("uuid")
	server, err := sc.svc.FindByUUID(uuid)
	if err != nil || server == nil || server.UserID != user.ID {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	return utils.Success(c, transformers.TransformServer(*server), "Server retrieved")
}

func (sc *ServerController) GetResources(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	uuid := c.Params("uuid")
	server, err := sc.svc.FindByUUID(uuid)
	if err != nil || server == nil || server.UserID != user.ID {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	node, err := sc.svc.FindByUUID(uuid)
	_ = node
	resources, err := sc.wingsService.GetServerResources(server.Node, server.UUID)
	if err != nil {
		return utils.Error(c, fiber.StatusBadGateway, "Node is offline")
	}

	return utils.Success(c, resources, "Resources retrieved")
}

func (sc *ServerController) Power(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "Not authenticated")
	}

	uuid := c.Params("uuid")
	server, err := sc.svc.FindByUUID(uuid)
	if err != nil || server == nil || server.UserID != user.ID {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	var req struct {
		Action string `json:"action" validate:"required,oneof=start stop restart kill"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	if err := sc.svc.PowerAction(uuid, req.Action); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Power action sent")
}
