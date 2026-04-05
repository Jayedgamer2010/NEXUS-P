package services

import (
	"errors"
	"nexus/backend/repositories"
	"nexus/backend/wings"
)

var (
	ErrNodeHasServers = errors.New("cannot delete node with existing servers")
)

type NodeService struct {
	nodeRepo   *repositories.NodeRepository
	serverRepo *repositories.ServerRepository
	wingsClient *wings.WingsClient
}

func NewNodeService(
	nodeRepo *repositories.NodeRepository,
	serverRepo *repositories.ServerRepository,
	wingsClient *wings.WingsClient,
) *NodeService {
	return &NodeService{
		nodeRepo:    nodeRepo,
		serverRepo:  serverRepo,
		wingsClient: wingsClient,
	}
}

func (s *NodeService) CountServersOnNode(nodeID uint) int64 {
	return s.serverRepo.CountByNodeID(nodeID)
}

func (s *NodeService) Delete(id uint) error {
	count := s.CountServersOnNode(id)
	if count > 0 {
		return ErrNodeHasServers
	}
	return s.nodeRepo.Delete(id)
}

func (s *NodeService) TestConnection(id uint) error {
	node, err := s.nodeRepo.FindByID(id)
	if err != nil {
		return err
	}
	_, err = s.wingsClient.GetSystemInfo(*node)
	return err
}
