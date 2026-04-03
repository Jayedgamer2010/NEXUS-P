package controllers

import (
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type NodeController struct{}

// GetAll returns all nodes
func (nc *NodeController) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var nodes []models.Node
	var total int64

	database.DB.Model(&models.Node{}).Count(&total)
	if err := database.DB.Offset(offset).Limit(limit).Find(&nodes).Error; err != nil {
		return utils.InternalError(c, "Failed to fetch nodes")
	}

	// Hide token in response
	response := make([]interface{}, len(nodes))
	for i, node := range nodes {
		response[i] = fiber.Map{
			"id":                node.ID,
			"uuid":              node.UUID,
			"name":              node.Name,
			"fqdn":              node.FQDN,
			"scheme":            node.Scheme,
			"wings_port":        node.WingsPort,
			"memory":            node.Memory,
			"memory_overalloc":  node.MemoryOveralloc,
			"disk":              node.Disk,
			"disk_overalloc":    node.DiskOveralloc,
			"token_id":          node.TokenID,
			"created_at":        node.CreatedAt,
			"updated_at":        node.UpdatedAt,
		}
	}

	return utils.Paginated(c, response, total, page, limit)
}

// GetByID returns a specific node
func (nc *NodeController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	return utils.Success(c, fiber.Map{
		"id":               node.ID,
		"uuid":             node.UUID,
		"name":             node.Name,
		"fqdn":             node.FQDN,
		"scheme":           node.Scheme,
		"wings_port":       node.WingsPort,
		"memory":           node.Memory,
		"memory_overalloc": node.MemoryOveralloc,
		"disk":             node.Disk,
		"disk_overalloc":   node.DiskOveralloc,
		"token_id":         node.TokenID,
		"created_at":       node.CreatedAt,
		"updated_at":       node.UpdatedAt,
	}, "Node retrieved")
}

// Create creates a new node
func (nc *NodeController) Create(c *fiber.Ctx) error {
	var req struct {
		Name            string `json:"name"`
		FQDN            string `json:"fqdn"`
		Scheme          string `json:"scheme"`
		WingsPort       int    `json:"wings_port"`
		Memory          int64  `json:"memory"`
		MemoryOveralloc int    `json:"memory_overalloc"`
		Disk            int64  `json:"disk"`
		DiskOveralloc   int    `json:"disk_overalloc"`
		TokenID         string `json:"token_id"`
		Token           string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate
	if req.Name == "" || req.FQDN == "" || req.TokenID == "" || req.Token == "" {
		return utils.BadRequest(c, "Name, FQDN, token_id, and token are required")
	}

	if req.Scheme != "http" && req.Scheme != "https" {
		req.Scheme = "https"
	}
	if req.WingsPort == 0 {
		req.WingsPort = 8080
	}

	nodeUUID := utils.GenerateUUID()

	node := models.Node{
		UUID:            nodeUUID,
		Name:            req.Name,
		FQDN:            req.FQDN,
		Scheme:          req.Scheme,
		WingsPort:       req.WingsPort,
		Memory:          req.Memory,
		MemoryOveralloc: req.MemoryOveralloc,
		Disk:            req.Disk,
		DiskOveralloc:   req.DiskOveralloc,
		TokenID:         req.TokenID,
		Token:           req.Token,
	}

	if err := database.DB.Create(&node).Error; err != nil {
		return utils.InternalError(c, "Failed to create node")
	}

	return utils.Success(c, fiber.Map{
		"id":     node.ID,
		"uuid":   node.UUID,
		"name":   node.Name,
	}, "Node created")
}

// Update updates a node
func (nc *NodeController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	var req struct {
		Name            *string `json:"name"`
		FQDN            *string `json:"fqdn"`
		Scheme          *string `json:"scheme"`
		WingsPort       *int    `json:"wings_port"`
		Memory          *int64  `json:"memory"`
		MemoryOveralloc *int    `json:"memory_overalloc"`
		Disk            *int64  `json:"disk"`
		DiskOveralloc   *int    `json:"disk_overalloc"`
		TokenID         *string `json:"token_id"`
		Token           *string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Name != nil {
		node.Name = *req.Name
	}
	if req.FQDN != nil {
		node.FQDN = *req.FQDN
	}
	if req.Scheme != nil {
		if *req.Scheme != "http" && *req.Scheme != "https" {
			return utils.BadRequest(c, "Scheme must be 'http' or 'https'")
		}
		node.Scheme = *req.Scheme
	}
	if req.WingsPort != nil {
		node.WingsPort = *req.WingsPort
	}
	if req.Memory != nil {
		node.Memory = *req.Memory
	}
	if req.MemoryOveralloc != nil {
		node.MemoryOveralloc = *req.MemoryOveralloc
	}
	if req.Disk != nil {
		node.Disk = *req.Disk
	}
	if req.DiskOveralloc != nil {
		node.DiskOveralloc = *req.DiskOveralloc
	}
	if req.TokenID != nil {
		node.TokenID = *req.TokenID
	}
	if req.Token != nil {
		node.Token = *req.Token
	}

	if err := database.DB.Save(&node).Error; err != nil {
		return utils.InternalError(c, "Failed to update node")
	}

	return utils.Success(c, fiber.Map{
		"id":   node.ID,
		"uuid": node.UUID,
	}, "Node updated")
}

// Delete deletes a node
func (nc *NodeController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	// Check if node has servers
	var serverCount int64
	database.DB.Model(&models.Server{}).Where("node_id = ?", node.ID).Count(&serverCount)
	if serverCount > 0 {
		return utils.BadRequest(c, "Cannot delete node with active servers")
	}

	if err := database.DB.Delete(&node).Error; err != nil {
		return utils.InternalError(c, "Failed to delete node")
	}

	return utils.Success(c, nil, "Node deleted")
}

// GetNodeStats returns aggregate stats for a node
func (nc *NodeController) GetStats(c *fiber.Ctx) error {
	id := c.Params("id")

	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Node not found")
	}

	// Get server count
	var serverCount int64
	database.DB.Model(&models.Server{}).Where("node_id = ?", node.ID).Count(&serverCount)

	// Get total allocated memory and disk
	type Stats struct {
		TotalMemory  int64 `json:"total_memory"`
		TotalDisk    int64 `json:"total_disk"`
		MemoryLimit  int64 `json:"memory_limit"`
		DiskLimit    int64 `json:"disk_limit"`
		ServerCount  int64 `json:"server_count"`
	}
	stats := Stats{
		MemoryLimit: node.Memory,
		DiskLimit:   node.Disk,
		ServerCount: serverCount,
	}

	database.DB.Model(&models.Server{}).Where("node_id = ?", node.ID).Select("COALESCE(SUM(memory), 0)").Scan(&stats.TotalMemory)
	database.DB.Model(&models.Server{}).Where("node_id = ?", node.ID).Select("COALESCE(SUM(disk), 0)").Scan(&stats.TotalDisk)

	return utils.Success(c, stats, "Node stats retrieved")
}
