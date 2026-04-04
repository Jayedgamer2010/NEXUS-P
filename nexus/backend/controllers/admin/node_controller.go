package admin

import (
	"strconv"

	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type NodeController struct {
	svc *services.NodeService
}

func NewNodeController(svc *services.NodeService) *NodeController {
	return &NodeController{svc: svc}
}

func (nc *NodeController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	nodes, total, err := nc.svc.All(page, limit)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch nodes")
	}

	result := make([]transformers.NodeTransformed, len(nodes))
	for i, n := range nodes {
		online := nc.svc.PingNode(&n)
		usedMem := n.GetUsedMemory(nil)
		usedDisk := n.GetUsedDisk(nil)
		result[i] = transformers.NodeTransformed{
			ID:              n.ID,
			UUID:            n.UUID,
			Public:          n.Public,
			Name:            n.Name,
			FQDN:            n.FQDN,
			Scheme:          n.Scheme,
			Memory:          n.Memory,
			Disk:            n.Disk,
			DaemonListen:    n.DaemonListen,
			DaemonSFTP:      n.DaemonSFTP,
			MaintenanceMode: n.MaintenanceMode,
			UsedMemory:      usedMem,
			UsedDisk:        usedDisk,
			CreatedAt:       n.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		_ = online
	}

	return utils.Paginated(c, result, total, page, limit)
}

func (nc *NodeController) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid node ID")
	}

	node, err := nc.svc.FindByID(uint(id))
	if err != nil || node == nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	result := transformers.TransformNode(*node, node.GetUsedMemory(nil), node.GetUsedDisk(nil), 0)
	return utils.Success(c, result, "Node retrieved")
}

func (nc *NodeController) Create(c *fiber.Ctx) error {
	var req requests.CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	node, err := nc.svc.Create(req)
	if err != nil {
		return utils.InternalError(c, "Failed to create node")
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(transformers.TransformNode(*node, 0, 0, 0)))
}

func (nc *NodeController) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid node ID")
	}

	node, err := nc.svc.FindByID(uint(id))
	if err != nil || node == nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	var req requests.UpdateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := nc.svc.Update(node, req); err != nil {
		return utils.InternalError(c, "Failed to update node")
	}

	return utils.Success(c, transformers.TransformNode(*node, 0, 0, 0), "Node updated")
}

func (nc *NodeController) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid node ID")
	}

	node, err := nc.svc.FindByID(uint(id))
	if err != nil || node == nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	if err := nc.svc.Delete(node); err != nil {
		return utils.InternalError(c, "Failed to delete node")
	}

	return utils.Success(c, nil, "Node deleted")
}

func (nc *NodeController) GetAllocations(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid node ID")
	}

	allocs, err := nc.svc.GetAllocations(uint(id))
	if err != nil {
		return utils.InternalError(c, "Failed to fetch allocations")
	}

	return utils.Success(c, transformers.TransformAllocations(allocs), "Allocations retrieved")
}

func (nc *NodeController) CreateAllocation(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	if id == 0 {
		return utils.BadRequest(c, "Invalid node ID")
	}

	var req requests.CreateAllocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if errors := utils.ValidateRequest(req); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ValidationErrorResponse(errors))
	}

	alloc, err := nc.svc.CreateAllocation(uint(id), req)
	if err != nil {
		return utils.InternalError(c, "Failed to create allocation")
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(transformers.TransformAllocation(*alloc)))
}

func (nc *NodeController) DeleteAllocation(c *fiber.Ctx) error {
	allocID, _ := strconv.ParseUint(c.Params("allocID"), 10, 64)
	if allocID == 0 {
		return utils.BadRequest(c, "Invalid allocation ID")
	}

	if err := nc.svc.DeleteAllocation(uint(allocID)); err != nil {
		return utils.InternalError(c, "Failed to delete allocation")
	}

	return utils.Success(c, nil, "Allocation deleted")
}
