package admin

import (
	"encoding/base64"
	"crypto/rand"
	"strconv"
	"strings"

	"nexus/backend/models"
	"nexus/backend/requests"
	"nexus/backend/services"
	"nexus/backend/repositories"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type NodeController struct {
	nodeRepo    *repositories.NodeRepository
	serverRepo  *repositories.ServerRepository
	allocRepo   *repositories.AllocationRepository
	nodeSvc     *services.NodeService
}

func NewNodeController(
	nodeRepo *repositories.NodeRepository,
	serverRepo *repositories.ServerRepository,
	allocRepo *repositories.AllocationRepository,
	nodeSvc *services.NodeService,
) *NodeController {
	return &NodeController{
		nodeRepo:   nodeRepo,
		serverRepo: serverRepo,
		allocRepo:  allocRepo,
		nodeSvc:    nodeSvc,
	}
}

func (ctrl *NodeController) GetAll(c *fiber.Ctx) error {
	nodes, err := ctrl.nodeRepo.FindAll()
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch nodes")
	}

	items := make([]transformers.NodeItem, 0, len(nodes))
	for _, node := range nodes {
		item := transformers.TransformNode(node)
		item.ServerCount = int(ctrl.nodeSvc.CountServersOnNode(node.ID))
		items = append(items, item)
	}

	return utils.Success(c, items)
}

func (ctrl *NodeController) Create(c *fiber.Ctx) error {
	var req requests.CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	scheme := req.Scheme
	if scheme == "" {
		scheme = "https"
	}
	daemonListen := req.DaemonListen
	if daemonListen == 0 {
		daemonListen = 8080
	}
	daemonSFTP := req.DaemonSFTP
	if daemonSFTP == 0 {
		daemonSFTP = 2022
	}
	daemonBase := req.DaemonBase
	if daemonBase == "" {
		daemonBase = "/var/lib/pterodactyl"
	}

	node := &models.Node{
		Name:              req.Name,
		Description:       req.Description,
		FQDN:              req.FQDN,
		Scheme:            scheme,
		Memory:            req.Memory,
		MemoryOverallocate: req.MemoryOveralloc,
		Disk:              req.Disk,
		DiskOverallocate:  req.DiskOveralloc,
		DaemonListen:      daemonListen,
		DaemonSFTP:        daemonSFTP,
		DaemonBase:        daemonBase,
		DaemonTokenID:     req.DaemonTokenID,
		DaemonToken:       req.DaemonToken,
		BehindProxy:       req.BehindProxy,
		MaintenanceMode:   req.MaintenanceMode,
	}

	if err := ctrl.nodeRepo.Create(node); err != nil {
		return utils.Error(c, 500, "Failed to create node")
	}

	return utils.Success(c, transformers.TransformNodeDetail(*node))
}

func (ctrl *NodeController) GetOne(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid node ID")
	}

	node, err := ctrl.nodeRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Node not found")
	}

	return utils.Success(c, transformers.TransformNodeDetail(*node))
}

func (ctrl *NodeController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid node ID")
	}

	node, err := ctrl.nodeRepo.FindByID(uint(id))
	if err != nil {
		return utils.Error(c, 404, "Node not found")
	}

	var req requests.UpdateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	if req.Name != "" {
		node.Name = req.Name
	}
	if req.Description != "" {
		node.Description = req.Description
	}
	if req.FQDN != "" {
		node.FQDN = req.FQDN
	}
	if req.Scheme != "" {
		node.Scheme = req.Scheme
	}
	if req.Memory > 0 {
		node.Memory = req.Memory
	}
	if req.Disk > 0 {
		node.Disk = req.Disk
	}
	if req.DaemonListen > 0 {
		node.DaemonListen = req.DaemonListen
	}
	if req.DaemonSFTP > 0 {
		node.DaemonSFTP = req.DaemonSFTP
	}
	if req.DaemonTokenID != "" {
		node.DaemonTokenID = req.DaemonTokenID
	}
	if req.DaemonToken != "" {
		node.DaemonToken = req.DaemonToken
	}
	if req.DaemonBase != "" {
		node.DaemonBase = req.DaemonBase
	}

	if err := ctrl.nodeRepo.Update(node); err != nil {
		return utils.Error(c, 500, "Failed to update node")
	}

	return utils.Success(c, transformers.TransformNodeDetail(*node))
}

func (ctrl *NodeController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid node ID")
	}

	if err := ctrl.nodeSvc.Delete(uint(id)); err != nil {
		if err == services.ErrNodeHasServers {
			return utils.Error(c, 422, "Cannot delete node with existing servers")
		}
		return utils.Error(c, 500, "Failed to delete node")
	}

	return utils.SuccessMessage(c, "Node deleted successfully", nil)
}

func (ctrl *NodeController) GetAllocations(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid node ID")
	}

	allocations, err := ctrl.allocRepo.FindByNodeID(uint(id))
	if err != nil {
		return utils.Error(c, 500, "Failed to fetch allocations")
	}

	type allocItem struct {
		ID       uint   `json:"id"`
		NodeID   uint   `json:"node_id"`
		IP       string `json:"ip"`
		IPAlias  string `json:"ip_alias"`
		Port     int    `json:"port"`
		Notes    string `json:"notes"`
		Assigned bool   `json:"assigned"`
		ServerID *uint  `json:"server_id"`
	}

	items := make([]allocItem, 0, len(allocations))
	for _, a := range allocations {
		items = append(items, allocItem{
			ID:       a.ID,
			NodeID:   a.NodeID,
			IP:       a.IP,
			IPAlias:  a.IPAlias,
			Port:     a.Port,
			Notes:    a.Notes,
			Assigned: a.IsAssigned(),
			ServerID: a.ServerID,
		})
	}

	return utils.Success(c, items)
}

func (ctrl *NodeController) AddAllocation(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid node ID")
	}

	var req requests.CreateAllocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "Invalid request body")
	}

	if errs := utils.Validate(req); errs != nil {
		return utils.ValidationError(c, errs)
	}

	alloc := &models.Allocation{
		NodeID:  uint(id),
		IP:      req.IP,
		IPAlias: req.IPAlias,
		Port:    req.Port,
		Notes:   req.Notes,
	}

	if err := ctrl.allocRepo.Create(alloc); err != nil {
		return utils.Error(c, 500, "Failed to create allocation")
	}

	return utils.Success(c, alloc)
}

func (ctrl *NodeController) DeleteAllocation(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.Error(c, 400, "Invalid allocation ID")
	}

	if err := ctrl.allocRepo.Delete(uint(id)); err != nil {
		if err == repositories.ErrAllocationAssigned {
			return utils.Error(c, 422, "Cannot delete allocation assigned to a server")
		}
		return utils.Error(c, 500, "Failed to delete allocation")
	}

	return utils.SuccessMessage(c, "Allocation deleted successfully", nil)
}

func generateDaemonToken() (string, string) {
	b := make([]byte, 32)
	rand.Read(b)
	tokenID := strings.ToLower(generateShortID())
	token := base64.RawStdEncoding.EncodeToString(b)
	return tokenID, token
}

func generateShortID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:8]
}
