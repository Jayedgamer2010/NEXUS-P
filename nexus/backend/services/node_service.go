package services

import (
	"errors"
	"fmt"

	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/wings"

	"github.com/google/uuid"
)

type NodeService struct {
	repo         *repositories.NodeRepository
	allocRepo    *repositories.AllocationRepository
	wingsService *WingsService
}

func NewNodeService(
	nodeRepo *repositories.NodeRepository,
	allocRepo *repositories.AllocationRepository,
	wingsSvc *WingsService,
) *NodeService {
	return &NodeService{repo: nodeRepo, allocRepo: allocRepo, wingsService: wingsSvc}
}

func (s *NodeService) All(page, perPage int) ([]models.Node, int64, error) {
	return s.repo.All(page, perPage)
}

func (s *NodeService) FindByID(id uint) (*models.Node, error) {
	return s.repo.FindByID(id)
}

func (s *NodeService) Create(req requests.CreateNodeRequest) (*models.Node, error) {
	node := &models.Node{
		UUID:          uuid.New().String(),
		Name:          req.Name,
		FQDN:          req.FQDN,
		Scheme:        req.Scheme,
		Memory:        req.Memory,
		Disk:          req.Disk,
		DaemonTokenID: req.DaemonTokenID,
		DaemonToken:   req.DaemonToken,
		DaemonListen:  req.DaemonListen,
		DaemonSFTP:    req.DaemonSFTP,
		DaemonBase:    req.DaemonBase,
		Public:        req.Public,
		Description:   req.Description,
		BehindProxy:   req.BehindProxy,
		LocationID:    req.LocationID,
	}

	if err := s.repo.Create(node); err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	return node, nil
}

func (s *NodeService) Update(node *models.Node, req requests.UpdateNodeRequest) error {
	if req.Name != nil {
		node.Name = *req.Name
	}
	if req.FQDN != nil {
		node.FQDN = *req.FQDN
	}
	if req.Scheme != nil {
		node.Scheme = *req.Scheme
	}
	if req.Memory != nil {
		node.Memory = *req.Memory
	}
	if req.Disk != nil {
		node.Disk = *req.Disk
	}
	if req.DaemonToken != nil {
		node.DaemonToken = *req.DaemonToken
	}
	if req.DaemonTokenID != nil {
		node.DaemonTokenID = *req.DaemonTokenID
	}
	if req.DaemonListen != nil {
		node.DaemonListen = *req.DaemonListen
	}
	if req.DaemonSFTP != nil {
		node.DaemonSFTP = *req.DaemonSFTP
	}
	if req.DaemonBase != nil {
		node.DaemonBase = *req.DaemonBase
	}
	if req.Public != nil {
		node.Public = *req.Public
	}
	if req.Description != nil {
		node.Description = *req.Description
	}
	if req.MaintenanceMode != nil {
		node.MaintenanceMode = *req.MaintenanceMode
	}

	return s.repo.Update(node)
}

func (s *NodeService) Delete(node *models.Node) error {
	return s.repo.Delete(node)
}

func (s *NodeService) IsOnline(node *models.Node) bool {
	return node.IsOnline()
}

func (s *NodeService) CreateAllocation(nodeID uint, req requests.CreateAllocationRequest) (*models.Allocation, error) {
	allocation := &models.Allocation{
		NodeID:  nodeID,
		IP:      req.IP,
		IPAlias: req.IPAlias,
		Port:    req.Port,
		Notes:   req.Notes,
	}

	if err := s.allocRepo.Create(allocation); err != nil {
		return nil, fmt.Errorf("failed to create allocation: %w", err)
	}

	return allocation, nil
}

func (s *NodeService) DeleteAllocation(id uint) error {
	return s.allocRepo.Delete(id)
}

func (s *NodeService) GetAllocations(nodeID uint) ([]models.Allocation, error) {
	return s.allocRepo.ByNodeID(nodeID)
}

func (s *NodeService) GetNodeResources(node *models.Node) (*wings.ServerResources, error) {
	info, err := s.wingsService.GetSystemInfo(*node)
	if err != nil {
		return nil, err
	}
	_ = info
	return nil, nil
}

func (s *NodeService) PingNode(node *models.Node) bool {
	_, err := s.wingsService.GetSystemInfo(*node)
	return err == nil
}

func (s *NodeService) GetNodeServerCount(nodeID uint) (int64, error) {
	var count int64
	err := nodeIDToDB(nodeID, &count)
	return count, err
}

func nodeIDToDB(_ uint, _ *int64) error {
	return errors.New("use repository for server count")
}
