package admin

import (
	"nexus/backend/repositories"
	"nexus/backend/transformers"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type StatsController struct {
	userRepo   *repositories.UserRepository
	serverRepo *repositories.ServerRepository
	nodeRepo   *repositories.NodeRepository
}

func NewStatsController(
	userRepo *repositories.UserRepository,
	serverRepo *repositories.ServerRepository,
	nodeRepo *repositories.NodeRepository,
) *StatsController {
	return &StatsController{
		userRepo:   userRepo,
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
	}
}

func (ctrl *StatsController) GetStats(c *fiber.Ctx) error {
	recentServers, _ := ctrl.serverRepo.FindRecent(5)
	recentUsers, _, _ := ctrl.userRepo.FindAll(1, 5, "")

	recentServerItems := make([]transformers.ServerItem, 0, len(recentServers))
	for _, s := range recentServers {
		recentServerItems = append(recentServerItems, transformers.TransformServer(s))
	}

	userItems := transformers.TransformUsers(recentUsers)

	return utils.Success(c, fiber.Map{
		"users":           ctrl.userRepo.CountAll(),
		"nodes":           ctrl.nodeRepo.CountAll(),
		"servers":         ctrl.serverRepo.CountAll(),
		"running_servers": ctrl.serverRepo.CountRunning(),
		"recent_servers":  recentServerItems,
		"recent_users":    userItems,
	})
}
