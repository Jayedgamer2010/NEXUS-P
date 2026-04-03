package controllers

import (
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/wings"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type ServerController struct{}

// Admin: Get all servers
func (sc *ServerController) GetAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var servers []models.Server
	var total int64

	database.DB.Model(&models.Server{}).Count(&total)
	if err := database.DB.Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return utils.InternalError(c, "Failed to fetch servers")
	}

	// For simplicity, return IDs - in production you'd want proper DTOs
	response := make([]interface{}, len(servers))
	for i, server := range servers {
		response[i] = fiber.Map{
			"id":              server.ID,
			"uuid":            server.UUID,
			"name":            server.Name,
			"status":          server.Status,
			"suspended":       server.Suspended,
			"memory":          server.Memory,
			"disk":            server.Disk,
			"cpu":             server.CPU,
			"user_id":         server.UserID,
			"node_id":         server.NodeID,
			"egg_id":          server.EggID,
			"allocation_id":   server.AllocationID,
			"created_at":      server.CreatedAt,
			"updated_at":      server.UpdatedAt,
		}
	}

	return utils.Paginated(c, response, total, page, limit)
}

// Admin: Get server by ID
func (sc *ServerController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var server models.Server
	if err := database.DB.Preload("Node").Preload("Egg").Preload("Allocation").First(&server, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	response := fiber.Map{
		"id":            server.ID,
		"uuid":          server.UUID,
		"name":          server.Name,
		"status":        server.Status,
		"suspended":     server.Suspended,
		"memory":        server.Memory,
		"disk":          server.Disk,
		"cpu":           server.CPU,
		"user_id":       server.UserID,
		"node_id":       server.NodeID,
		"egg_id":        server.EggID,
		"allocation_id": server.AllocationID,
		"node": fiber.Map{
			"id":   server.Node.ID,
			"uuid": server.Node.UUID,
			"name": server.Node.Name,
			"fqdn": server.Node.FQDN,
		},
		"egg": fiber.Map{
			"id":   server.Egg.ID,
			"uuid": server.Egg.UUID,
			"name": server.Egg.Name,
		},
		"allocation": fiber.Map{
			"ip":   server.Allocation.IP,
			"port": server.Allocation.Port,
		},
		"created_at": server.CreatedAt,
		"updated_at": server.UpdatedAt,
	}

	return utils.Success(c, response, "Server retrieved")
}

// Admin: Create server
func (sc *ServerController) Create(c *fiber.Ctx) error {
	var req struct {
		Name          string `json:"name"`
		UserID        uint   `json:"user_id"`
		NodeID        uint   `json:"node_id"`
		EggID         uint   `json:"egg_id"`
		AllocationID  uint   `json:"allocation_id"`
		Memory        int64  `json:"memory"`
		Disk          int64  `json:"disk"`
		CPU           int    `json:"cpu"`
		Startup       string `json:"startup"`
		Image         string `json:"image"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if req.Name == "" || req.UserID == 0 || req.NodeID == 0 || req.EggID == 0 ||
		req.AllocationID == 0 || req.Memory == 0 || req.Disk == 0 || req.CPU == 0 {
		return utils.BadRequest(c, "Missing required fields")
	}

	// Check user exists
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		return utils.BadRequest(c, "User not found")
	}

	// Check node exists
	var node models.Node
	if err := database.DB.First(&node, req.NodeID).Error; err != nil {
		return utils.BadRequest(c, "Node not found")
	}

	// Check egg exists
	var egg models.Egg
	if err := database.DB.First(&egg, req.EggID).Error; err != nil {
		return utils.BadRequest(c, "Egg not found")
	}

	// Check allocation exists and is not assigned
	var allocation models.Allocation
	if err := database.DB.First(&allocation, req.AllocationID).Error; err != nil {
		return utils.BadRequest(c, "Allocation not found")
	}
	if allocation.Assigned {
		return utils.BadRequest(c, "Allocation already assigned")
	}

	// Generate UUID
	serverUUID := utils.GenerateUUID()

	// Create server record
	server := models.Server{
		UUID:          serverUUID,
		Name:          req.Name,
		UserID:        req.UserID,
		NodeID:        req.NodeID,
		EggID:         req.EggID,
		AllocationID:  req.AllocationID,
		Memory:        req.Memory,
		Disk:          req.Disk,
		CPU:           req.CPU,
		Status:        "installing",
		Suspended:     false,
	}

	if err := database.DB.Create(&server).Error; err != nil {
		return utils.InternalError(c, "Failed to create server")
	}

	// Mark allocation as assigned
	allocation.Assigned = true
	allocation.ServerID = &server.ID
	database.DB.Save(&allocation)

	// Notify Wings to create server (async in production)
	go func() {
		wingsClient := wings.NewClient(&node)
		if err := wingsClient.CreateServer(server.UUID, true); err != nil {
			database.DB.Model(&server).Update("status", "error")
		} else {
			database.DB.Model(&server).Update("status", "running")
		}
	}()

	return utils.Success(c, fiber.Map{
		"id":     server.ID,
		"uuid":   server.UUID,
		"status": server.Status,
	}, "Server creation initiated")
}

// Admin: Update server
func (sc *ServerController) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var server models.Server
	if err := database.DB.First(&server, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	var req struct {
		Name          *string `json:"name"`
		Memory        *int64  `json:"memory"`
		Disk          *int64  `json:"disk"`
		CPU           *int    `json:"cpu"`
		Status        *string `json:"status"`
		Suspended     *bool   `json:"suspended"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Update fields if provided
	if req.Name != nil {
		server.Name = *req.Name
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
	if req.Status != nil {
		server.Status = *req.Status
	}
	if req.Suspended != nil {
		server.Suspended = *req.Suspended
	}

	if err := database.DB.Save(&server).Error; err != nil {
		return utils.InternalError(c, "Failed to update server")
	}

	return utils.Success(c, fiber.Map{
		"id":   server.ID,
		"uuid": server.UUID,
	}, "Server updated")
}

// Admin: Delete server
func (sc *ServerController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var server models.Server
	if err := database.DB.Preload("Node").First(&server, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	// Free allocation
	if err := database.DB.Model(&models.Allocation{}).Where("server_id = ?", server.ID).Updates(map[string]interface{}{
		"assigned": false,
		"server_id": nil,
	}).Error; err != nil {
		return utils.InternalError(c, "Failed to free allocation")
	}

	// Notify Wings to delete server
	go func() {
		wingsClient := wings.NewClient(&server.Node)
		wingsClient.DeleteServer(server.UUID)
	}()

	if err := database.DB.Delete(&server).Error; err != nil {
		return utils.InternalError(c, "Failed to delete server")
	}

	return utils.Success(c, nil, "Server deleted")
}

// Admin: Send power action
func (sc *ServerController) Power(c *fiber.Ctx) error {
	id := c.Params("id")

	var server models.Server
	if err := database.DB.Preload("Node").First(&server, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	var req struct {
		Action string `json:"action"` // start, stop, restart, kill
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if req.Action == "" {
		return utils.BadRequest(c, "Action is required")
	}

	// Validate action
	validActions := map[string]bool{"start": true, "stop": true, "restart": true, "kill": true}
	if !validActions[req.Action] {
		return utils.BadRequest(c, "Invalid action. Must be one of: start, stop, restart, kill")
	}

	// Send power action to Wings (async)
	go func() {
		wingsClient := wings.NewClient(&server.Node)
		if err := wingsClient.SendPowerAction(server.UUID, req.Action); err != nil {
			// Log error but don't fail the request
			database.DB.Model(&server).Update("status", "error")
		}
	}()

	return utils.Success(c, nil, "Power action sent")
}

// Client: Get user's own servers
func (sc *ServerController) GetMyServers(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	var servers []models.Server
	if err := database.DB.Where("user_id = ?", user.ID).Find(&servers).Error; err != nil {
		return utils.InternalError(c, "Failed to fetch servers")
	}

	response := make([]interface{}, len(servers))
	for i, server := range servers {
		response[i] = fiber.Map{
			"id":            server.ID,
			"uuid":          server.UUID,
			"name":          server.Name,
			"status":        server.Status,
			"suspended":     server.Suspended,
			"memory":        server.Memory,
			"disk":          server.Disk,
			"cpu":           server.CPU,
			"node_id":       server.NodeID,
			"egg_id":        server.EggID,
			"allocation_id": server.AllocationID,
			"created_at":    server.CreatedAt,
		}
	}

	return utils.Success(c, response, "Servers retrieved")
}

// Client: Get specific server by UUID (client's own)
func (sc *ServerController) GetMyServer(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}
	uuid := c.Params("uuid")

	var server models.Server
	if err := database.DB.Where("uuid = ? AND user_id = ?", uuid, user.ID).
		Preload("Node").Preload("Egg").Preload("Allocation").First(&server).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	response := fiber.Map{
		"id":            server.ID,
		"uuid":          server.UUID,
		"name":          server.Name,
		"status":        server.Status,
		"suspended":     server.Suspended,
		"memory":        server.Memory,
		"disk":          server.Disk,
		"cpu":           server.CPU,
		"node_id":       server.NodeID,
		"egg_id":        server.EggID,
		"allocation_id": server.AllocationID,
		"node": fiber.Map{
			"id":   server.Node.ID,
			"uuid": server.Node.UUID,
			"name": server.Node.Name,
			"fqdn": server.Node.FQDN,
		},
		"egg": fiber.Map{
			"id":   server.Egg.ID,
			"uuid": server.Egg.UUID,
			"name": server.Egg.Name,
		},
		"allocation": fiber.Map{
			"ip":   server.Allocation.IP,
			"port": server.Allocation.Port,
		},
		"created_at": server.CreatedAt,
		"updated_at": server.UpdatedAt,
	}

	return utils.Success(c, response, "Server retrieved")
}

// Client: Get server resources
func (sc *ServerController) GetResources(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}
	uuid := c.Params("uuid")

	var server models.Server
	if err := database.DB.Where("uuid = ? AND user_id = ?", uuid, user.ID).
		Preload("Node").First(&server); err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Server not found")
	}

	// Fetch from Wings
	wingsClient := wings.NewClient(&server.Node)
	resources, err := wingsClient.GetServerResources(server.UUID)
	if err != nil {
		return utils.InternalError(c, "Failed to fetch resources from Wings")
	}

	return utils.Success(c, resources, "Resources retrieved")
}
