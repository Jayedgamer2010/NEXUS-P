package admin

import (
	"nexus/backend/repositories"
	"nexus/backend/services"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type StatsController struct {
	userRepo     *repositories.UserRepository
	serverRepo   *repositories.ServerRepository
	nodeRepo     *repositories.NodeRepository
	wingsService *services.WingsService
}

func NewStatsController(
	userRepo *repositories.UserRepository,
	serverRepo *repositories.ServerRepository,
	nodeRepo *repositories.NodeRepository,
	wingsSvc *services.WingsService,
) *StatsController {
	return &StatsController{
		userRepo:     userRepo,
		serverRepo:   serverRepo,
		nodeRepo:     nodeRepo,
		wingsService: wingsSvc,
	}
}

func (sc *StatsController) GetStats(c *fiber.Ctx) error {
	userCount, _ := sc.userRepo.Count()
	serverCount, _ := sc.serverRepo.Count()
	nodeCount, _ := sc.nodeRepo.Count()
	runningCount, _ := sc.serverRepo.CountRunning()

	nodes, _, _ := sc.nodeRepo.All(1, 100)
	nodeStatuses := make([]map[string]interface{}, 0)
	for _, node := range nodes {
		online := node.IsOnline()
		nodeStatuses = append(nodeStatuses, map[string]interface{}{
			"id":     node.ID,
			"name":   node.Name,
			"online": online,
		})
	}

	stats := map[string]interface{}{
		"users":           userCount,
		"nodes":           nodeCount,
		"servers":         serverCount,
		"running_servers": runningCount,
		"node_statuses":   nodeStatuses,
	}

	return utils.Success(c, stats, "Stats retrieved")
}
