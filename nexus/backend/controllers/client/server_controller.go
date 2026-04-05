package client

import (
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type ServerController struct {
	serverRepo   *repositories.ServerRepository
	serverSvc    *services.ServerService
	nodeRepo     *repositories.NodeRepository
}

func NewClientServerController(
	serverRepo *repositories.ServerRepository,
	serverSvc *services.ServerService,
	nodeRepo *repositories.NodeRepository,
) *ServerController {
	return &ServerController{
		serverRepo: serverRepo,
		serverSvc:  serverSvc,
		nodeRepo:   nodeRepo,
	}
}

func (ctrl *ServerController) GetAll(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	servers, err := ctrl.serverRepo.FindByUserID(user.ID)
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch servers")
	}

	items := make([]transformers.ServerItem, 0, len(servers))
	for _, s := range servers {
		items = append(items, transformers.TransformServer(s))
	}

	return utils.Success(c, items)
}

func (ctrl *ServerController) GetOne(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	uid := c.Params("uuid")

	server, err := ctrl.serverRepo.FindByUUID(uid)
	if err != nil {
		return utils.Error(c, 404, "Server not found")
	}

	if server.UserID != user.ID {
		return utils.Error(c, 403, "Not authorized to access this server")
	}

	return utils.Success(c, transformers.TransformServerDetail(*server))
}

func (ctrl *ServerController) GetResources(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	uid := c.Params("uuid")

	server, err := ctrl.serverRepo.FindByUUID(uid)
	if err != nil {
		return utils.Error(c, 404, "Server not found")
	}

	if server.UserID != user.ID {
		return utils.Error(c, 403, "Not authorized to access this server")
	}

	resources, err := ctrl.serverSvc.GetServerResources(server.ID)
	if err != nil {
		return utils.Error(c, 502, "Failed to fetch resources from Wings")
	}

	return utils.Success(c, resources)
}

func (ctrl *ServerController) PowerAction(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	uid := c.Params("uuid")

	server, err := ctrl.serverRepo.FindByUUID(uid)
	if err != nil {
		return utils.Error(c, 404, "Server not found")
	}

	if server.UserID != user.ID {
		return utils.Error(c, 403, "Not authorized to access this server")
	}

	var req requests.PowerActionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if err := ctrl.serverSvc.PowerAction(server.ID, req.Action); err != nil {
		return utils.Error(c, 502, "Failed to send power action: "+err.Error())
	}

	return utils.SuccessMessage(c, "Power action sent: "+req.Action, nil)
}
