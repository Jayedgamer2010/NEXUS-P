package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/wings"

	"github.com/google/uuid"
)

type ServerService struct {
	repo         *repositories.ServerRepository
	nodeRepo     *repositories.NodeRepository
	eggRepo      *repositories.EggRepository
	userRepo     *repositories.UserRepository
	allocRepo    *repositories.AllocationRepository
	wingsService *WingsService
}

func NewServerService(
	serverRepo *repositories.ServerRepository,
	nodeRepo *repositories.NodeRepository,
	eggRepo *repositories.EggRepository,
	userRepo *repositories.UserRepository,
	allocRepo *repositories.AllocationRepository,
	wingsSvc *WingsService,
) *ServerService {
	return &ServerService{
		repo:         serverRepo,
		nodeRepo:     nodeRepo,
		eggRepo:      eggRepo,
		userRepo:     userRepo,
		allocRepo:    allocRepo,
		wingsService: wingsSvc,
	}
}

func (s *ServerService) All(page, perPage int) ([]models.Server, int64, error) {
	return s.repo.All(page, perPage)
}

func (s *ServerService) FindByID(id uint) (*models.Server, error) {
	return s.repo.FindByID(id)
}

func (s *ServerService) FindByUUID(uuid string) (*models.Server, error) {
	return s.repo.FindByUUID(uuid)
}

func (s *ServerService) FindByUserID(userID uint) ([]models.Server, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ServerService) Create(req requests.CreateServerRequest) (*models.Server, error) {
	node, err := s.nodeRepo.FindByID(req.NodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find node: %w", err)
	}
	if node == nil {
		return nil, errors.New("node not found")
	}

	egg, err := s.eggRepo.FindByID(req.EggID)
	if err != nil {
		return nil, fmt.Errorf("failed to find egg: %w", err)
	}
	if egg == nil {
		return nil, errors.New("egg not found")
	}

	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	allocation, err := s.allocRepo.FindUnassigned(req.NodeID)
	if err != nil {
		return nil, errors.New("no available allocations on this node")
	}

	serverUUID := uuid.New().String()
	serverUUIDShort := uuid.New().String()[:8]

	envBytes, _ := json.Marshal(req.Environment)

	dockerImage := req.DockerImage
	if dockerImage == "" {
		dockerImage = egg.DockerImage
	}

	startup := req.StartupCmd
	if startup == "" {
		startup = egg.Startup
	}

	server := &models.Server{
		UUID:         serverUUID,
		UUIDShort:    serverUUIDShort,
		Name:         req.Name,
		Description:  req.Description,
		UserID:       req.UserID,
		NodeID:       req.NodeID,
		EggID:        req.EggID,
		AllocationID: allocation.ID,
		Memory:       req.Memory,
		Disk:         req.Disk,
		CPU:          req.CPU,
		Image:        dockerImage,
		Startup:      startup,
		EnvVariables: string(envBytes),
		Status:       models.StatusInstalling,
	}

	if err := s.repo.Create(server); err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	allocation.Assign(server.ID)
	_ = s.allocRepo.Update(allocation)

	wingsPayload := wings.CreateServerPayload{
		UUID:              serverUUID,
		StartOnCompletion: true,
		Build: wings.BuildConfig{
			MemoryLimit: req.Memory,
			Swap:        0,
			Disk:        req.Disk,
			IOWeight:    500,
			CPU:         req.CPU,
		},
		Container: wings.ContainerConfig{
			Image:       dockerImage,
			Startup:     startup,
			Environment: req.Environment,
		},
		Allocation: wings.AllocationConfig{
			Default:    allocation.Port,
			Additional: []interface{}{},
		},
	}

	if err := s.wingsService.CreateServer(*node, wingsPayload); err != nil {
		server.Status = models.StatusInstallFailed
		_ = s.repo.Update(server)
		return server, fmt.Errorf("server created in db but wings failed: %w", err)
	}

	server.Status = models.StatusRunning
	_ = s.repo.Update(server)

	return server, nil
}

func (s *ServerService) Update(server *models.Server, req requests.UpdateServerRequest) error {
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
	if req.Description != nil {
		server.Description = *req.Description
	}

	return s.repo.Update(server)
}

func (s *ServerService) Delete(serverUUID string) error {
	server, err := s.repo.FindByUUID(serverUUID)
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	if server == nil {
		return errors.New("server not found")
	}

	node, err := s.nodeRepo.FindByID(server.NodeID)
	if err == nil && node != nil {
		_ = s.wingsService.DeleteServer(*node, serverUUID)
	}

	allocation, albErr := s.allocRepo.FindByID(server.AllocationID)
	if albErr == nil {
		allocation.Unassign()
		_ = s.allocRepo.Update(allocation)
	}

	return s.repo.Delete(server)
}

func (s *ServerService) PowerAction(serverUUID string, action string) error {
	validActions := map[string]bool{models.PowerStart: true, models.PowerStop: true, models.PowerRestart: true, models.PowerKill: true}
	if !validActions[action] {
		return errors.New("invalid action")
	}

	server, err := s.repo.FindByUUID(serverUUID)
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	if server == nil {
		return errors.New("server not found")
	}

	node, err := s.nodeRepo.FindByID(server.NodeID)
	if err != nil || node == nil {
		return errors.New("node not found")
	}

	return s.wingsService.SendPowerAction(*node, serverUUID, action)
}

func (s *ServerService) Suspend(serverUUID string) error {
	server, err := s.repo.FindByUUID(serverUUID)
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	if server == nil {
		return errors.New("server not found")
	}

	server.Suspended = true
	_ = s.wingsService.SendPowerAction(*s.mustGetNode(server.NodeID), serverUUID, "stop")
	return s.repo.Update(server)
}

func (s *ServerService) Unsuspend(serverUUID string) error {
	server, err := s.repo.FindByUUID(serverUUID)
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	if server == nil {
		return errors.New("server not found")
	}

	server.Suspended = false
	return s.repo.Update(server)
}

func (s *ServerService) Reinstall(serverUUID string) error {
	server, err := s.repo.FindByUUID(serverUUID)
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	if server == nil {
		return errors.New("server not found")
	}

	server.Status = models.StatusInstalling
	return s.repo.Update(server)
}

func (s *ServerService) GetResources(server *models.Server, node *models.Node) (*wings.ServerResources, error) {
	return s.wingsService.GetServerResources(*node, server.UUID)
}

func (s *ServerService) mustGetNode(nodeID uint) *models.Node {
	node, _ := s.nodeRepo.FindByID(nodeID)
	return node
}
