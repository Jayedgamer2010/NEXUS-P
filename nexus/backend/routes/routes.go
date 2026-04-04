package routes

import (
	"nexus/backend/controllers"
	"nexus/backend/controllers/admin"
	"nexus/backend/controllers/client"
	"nexus/backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	authCfg *middleware.AuthConfig,
	authCtl *controllers.AuthController,
	userCtl *admin.UserController,
	serverCtl *admin.ServerController,
	nodeCtl *admin.NodeController,
	eggCtl *admin.EggController,
	statsCtl *admin.StatsController,
	clientSrvCtl *client.ServerController,
	accCtl *client.AccountController,
) {
	// No auth
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "NEXUS", "version": "1.0.0"})
	})

	// Auth routes
	authGroup := app.Group("/api/auth")
	authGroup.Post("/login", authCtl.Login)
	authGroup.Post("/register", authCtl.Register)

	// Protected auth routes
	authProtected := app.Group("/api/auth", middleware.Auth(authCfg))
	authProtected.Get("/me", authCtl.Me)
	authProtected.Post("/logout", authCtl.Logout)

	// Client routes (authenticated)
	clientGroup := app.Group("/api/client", middleware.Auth(authCfg))
	clientGroup.Get("/account", accCtl.Get)
	clientGroup.Patch("/account", accCtl.Update)
	clientGroup.Get("/servers", clientSrvCtl.GetMyServers)
	clientGroup.Get("/servers/:uuid", clientSrvCtl.GetMyServer)
	clientGroup.Get("/servers/:uuid/resources", clientSrvCtl.GetResources)
	clientGroup.Post("/servers/:uuid/power", clientSrvCtl.Power)

	// Admin routes
	adminGroup := app.Group("/api/admin", middleware.Auth(authCfg), middleware.Admin)

	// Users
	users := adminGroup.Group("/users")
	users.Get("/", userCtl.GetAll)
	users.Post("/", userCtl.Create)
	users.Get("/:id", userCtl.GetByID)
	users.Patch("/:id", userCtl.Update)
	users.Delete("/:id", userCtl.Delete)

	// Servers
	servers := adminGroup.Group("/servers")
	servers.Get("/", serverCtl.GetAll)
	servers.Post("/", serverCtl.Create)
	servers.Get("/:id", serverCtl.GetByID)
	servers.Patch("/:id", serverCtl.Update)
	servers.Delete("/:uuid", serverCtl.Delete)
	servers.Post("/:uuid/power", serverCtl.Power)
	servers.Post("/:uuid/suspend", serverCtl.Suspend)
	servers.Post("/:uuid/unsuspend", serverCtl.Unsuspend)
	servers.Post("/:uuid/reinstall", serverCtl.Reinstall)

	// Nodes
	nodes := adminGroup.Group("/nodes")
	nodes.Get("/", nodeCtl.GetAll)
	nodes.Post("/", nodeCtl.Create)
	nodes.Get("/:id", nodeCtl.GetByID)
	nodes.Patch("/:id", nodeCtl.Update)
	nodes.Delete("/:id", nodeCtl.Delete)
	nodes.Get("/:id/allocations", nodeCtl.GetAllocations)
	nodes.Post("/:id/allocations", nodeCtl.CreateAllocation)

	// Allocation deletion must be registered before "/:id" to match correctly,
	// but since we're using groups we handle it separately
	adminGroup.Delete("/allocations/:allocID", nodeCtl.DeleteAllocation)

	// Eggs
	eggs := adminGroup.Group("/eggs")
	eggs.Get("/", eggCtl.GetAll)
	eggs.Post("/", eggCtl.Create)
	eggs.Get("/:id", eggCtl.GetByID)
	eggs.Patch("/:id", eggCtl.Update)
	eggs.Delete("/:id", eggCtl.Delete)

	// Stats
	adminGroup.Get("/stats", statsCtl.GetStats)
}