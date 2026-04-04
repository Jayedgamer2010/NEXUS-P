package admin

import (
	"strconv"

	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type ServerController struct {
	svc *services.ServerService
}

func NewServerController(svc *services.ServerService) *ServerController {
	return &ServerController{svc: svc}
}

func (sc *ServerController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	servers, total, err := sc.svc.All(page, limit)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch servers")
	}

	return utils.Paginated(c, transformers.TransformServers(servers), total, page, limit)
}

func (sc *ServerController) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid server ID")
	}

	server, err := sc.svc.FindByID(uint(id))
	if err != nil || server == nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	return utils.Success(c, transformers.TransformServer(*server), "Server retrieved")
}

func (sc *ServerController) Create(c *fiber.Ctx) error {
	var req requests.CreateServerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	server, err := sc.svc.Create(req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(transformers.TransformServer(*server)))
}

func (sc *ServerController) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid server ID")
	}

	server, err := sc.svc.FindByID(uint(id))
	if err != nil || server == nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	var req requests.UpdateServerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := sc.svc.Update(server, req); err != nil {
		return utils.InternalError(c, "Failed to update server")
	}

	return utils.Success(c, transformers.TransformServer(*server), "Server updated")
}

func (sc *ServerController) Delete(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		if id == 0 {
			return utils.BadRequest(c, "Invalid server ID")
		}
		server, err := sc.svc.FindByID(uint(id))
		if err != nil || server == nil {
			return utils.Error(c, fiber.StatusNotFound, "Server not found")
		}
		uuid = server.UUID
	}

	if err := sc.svc.Delete(uuid); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Server deleted")
}

func (sc *ServerController) Power(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		if id == 0 {
			return utils.BadRequest(c, "Invalid server ID")
		}
		server, err := sc.svc.FindByID(uint(id))
		if err != nil || server == nil {
			return utils.Error(c, fiber.StatusNotFound, "Server not found")
		}
		uuid = server.UUID
	}

	var req requests.PowerActionRequest
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

func (sc *ServerController) Suspend(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		if id == 0 {
			return utils.BadRequest(c, "Invalid server ID")
		}
		server, err := sc.svc.FindByID(uint(id))
		if err != nil || server == nil {
			return utils.Error(c, fiber.StatusNotFound, "Server not found")
		}
		uuid = server.UUID
	}

	if err := sc.svc.Suspend(uuid); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Server suspended")
}

func (sc *ServerController) Unsuspend(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		if id == 0 {
			return utils.BadRequest(c, "Invalid server ID")
		}
		server, err := sc.svc.FindByID(uint(id))
		if err != nil || server == nil {
			return utils.Error(c, fiber.StatusNotFound, "Server not found")
		}
		uuid = server.UUID
	}

	if err := sc.svc.Unsuspend(uuid); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Server unsuspended")
}

func (sc *ServerController) Reinstall(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
		if id == 0 {
			return utils.BadRequest(c, "Invalid server ID")
		}
		server, err := sc.svc.FindByID(uint(id))
		if err != nil || server == nil {
			return utils.Error(c, fiber.StatusNotFound, "Server not found")
		}
		uuid = server.UUID
	}

	if err := sc.svc.Reinstall(uuid); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, nil, "Server reinstall initiated")
}
