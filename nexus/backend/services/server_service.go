package services

import (
	"errors"
	"fmt"
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/wings"
)

type ServerService struct {
	serverRepo   *repositories.ServerRepository
	nodeRepo     *repositories.NodeRepository
	eggRepo      *repositories.EggRepository
	userRepo     *repositories.UserRepository
	allocRepo    *repositories.AllocationRepository
	wingsClient  *wings.WingsClient
}

func NewServerService(
	serverRepo *repositories.ServerRepository,
	nodeRepo *repositories.NodeRepository,
	eggRepo *repositories.EggRepository,
	userRepo *repositories.UserRepository,
	allocRepo *repositories.AllocationRepository,
	wingsClient *wings.WingsClient,
) *ServerService {
	return &ServerService{
		serverRepo:  serverRepo,
		nodeRepo:    nodeRepo,
		eggRepo:     eggRepo,
		userRepo:    userRepo,
		allocRepo:   allocRepo,
		wingsClient: wingsClient,
	}
}

func (s *ServerService) Create(req requests.CreateServerRequest) (*models.Server, error) {
	// Verify node exists
	node, err := s.nodeRepo.FindByID(req.NodeID)
	if err != nil {
		return nil, errors.New("node not found")
	}

	// Verify egg exists
	egg, err := s.eggRepo.FindByID(req.EggID)
	if err != nil {
		return nil, errors.New("egg not found")
	}

	// Verify user exists
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Find available allocation
	alloc, err := s.allocRepo.FindAvailable(req.NodeID)
	if err != nil {
		return nil, errors.New("no available allocations on this node")
	}

	// Build environment map from egg startup
	envMap := map[string]string{
		"SERVER_JARFILE": "server.jar",
		"STARTUP":        egg.Startup,
	}

	// Build Wings payload
	payload := wings.CreateServerPayload{
		UUID:              "", // Will be set by GORM BeforeCreate
		StartOnCompletion: true,
		Image:             egg.DockerImage,
		Startup: wings.StartupConfig{
			Done:            "",
			UserInteraction: nil,
			StripAnsi:       nil,
		},
		Environment: envMap,
		Limits: wings.ServerLimits{
			Memory:  req.Memory,
			Swap:    req.Swap,
			Disk:    req.Disk,
			IO:      500,
			CPU:     req.CPU,
			Threads: "",
		},
		FeatureLimits: wings.FeatureLimits{
			Databases:   0,
			Allocations: 0,
			Backups:     0,
		},
		Allocations: wings.AllocationConfig{
			Default:    alloc.Port,
			Additional: []int{},
		},
	}

	// Create server in DB first (so GORM generates UUID)
	server := &models.Server{
		Name:         req.Name,
		Description:  req.Description,
		UserID:       req.UserID,
		NodeID:       req.NodeID,
		EggID:        req.EggID,
		AllocationID: alloc.ID,
		Memory:       req.Memory,
		Disk:         req.Disk,
		CPU:          req.CPU,
		Swap:         req.Swap,
		IO:           500,
		Image:        egg.DockerImage,
		Startup:      egg.Startup,
		Environment:  "{}",
		Status:       models.StatusInstalling,
	}

	// Try to create on Wings
	if err := s.wingsClient.CreateServer(*node, payload); err != nil {
		server.Status = models.StatusInstallFailed
	}

	// Save to DB
	if err := s.serverRepo.Create(server); err != nil {
		return nil, fmt.Errorf("create server in database: %w", err)
	}

	// Assign allocation
	if err := s.allocRepo.Assign(alloc.ID, server.ID); err != nil {
		return nil, fmt.Errorf("assign allocation: %w", err)
	}

	// Refresh server with preloads
	result, _ := s.serverRepo.FindByID(server.ID)
	if result != nil {
		return result, nil
	}
	return server, nil
}

func (s *ServerService) Delete(id uint) error {
	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return err
	}

	node, _ := s.nodeRepo.FindByID(server.NodeID)
	if node != nil {
		// Ignore Wings errors — we still delete from DB
		s.wingsClient.DeleteServer(*node, server.UUID, true)
	}

	// Unassign allocation
	s.allocRepo.Unassign(server.AllocationID)

	return s.serverRepo.Delete(id)
}

func (s *ServerService) PowerAction(id uint, action string) error {
	switch action {
	case "start", "stop", "restart", "kill":
	default:
		return errors.New("invalid power action: " + action)
	}

	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return err
	}

	node, _ := s.nodeRepo.FindByID(server.NodeID)
	if node == nil {
		return errors.New("node not found for server")
	}

	return s.wingsClient.SendPowerAction(*node, server.UUID, action)
}

func (s *ServerService) Suspend(id uint) error {
	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return err
	}

	server.Suspended = true
	server.Status = models.StatusSuspended
	return s.serverRepo.Update(server)
}

func (s *ServerService) Unsuspend(id uint) error {
	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return err
	}

	server.Suspended = false
	server.Status = models.StatusOffline
	return s.serverRepo.Update(server)
}

func (s *ServerService) GetServerResources(id uint) (*wings.ServerResources, error) {
	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	node, _ := s.nodeRepo.FindByID(server.NodeID)
	if node == nil {
		return nil, errors.New("node not found for server")
	}

	return s.wingsClient.GetServerResources(*node, server.UUID)
}

func (s *ServerService) GetConsoleToken(server *models.Server) (*wings.ServerConsoleToken, error) {
	node, err := s.nodeRepo.FindByID(server.NodeID)
	if err != nil {
		return nil, errors.New("node not found")
	}
	return s.wingsClient.GetConsoleToken(*node, server.UUID)
}

// ParseEnv parses environment JSON string to map
func ParseEnv(envJSON string) map[string]string {
	result := map[string]string{}
	if envJSON == "" {
		return result
	}
	// Simple: if it starts with {, try to parse
	if len(envJSON) > 0 && envJSON[0] == '{' {
		result["STARTUP"] = envJSON
	} else {
		result["STARTUP"] = envJSON
	}
	return result
}
