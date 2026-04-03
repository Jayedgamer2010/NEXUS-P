package routes

import (
	"nexus/backend/controllers"
	"nexus/backend/middleware"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
)

// Setup registers all routes
func Setup(app *fiber.App, authCfg *middleware.AuthConfig, authCtl *controllers.AuthController, userController *controllers.UserController, serverController *controllers.ServerController, nodeController *controllers.NodeController, eggController *controllers.EggController) {

	// Public auth routes
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", authCtl.Register)
	authGroup.Post("/login", authCtl.Login)

	// Protected auth routes
	authProtected := app.Group("/api/auth", middleware.Auth(authCfg))
	authProtected.Get("/me", authCtl.Me)
	authProtected.Post("/logout", authCtl.Logout)

	// Admin routes (require admin role)
	admin := app.Group("/api/admin", middleware.Auth(authCfg), middleware.Admin)

	// User management
	users := admin.Group("/users")
	users.Get("/", userController.GetAll)
	users.Post("/", userController.Create)
	users.Get("/:id", userController.GetByID)
	users.Patch("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)

	// Node management
	nodes := admin.Group("/nodes")
	nodes.Get("/", nodeController.GetAll)
	nodes.Post("/", nodeController.Create)
	nodes.Get("/:id", nodeController.GetByID)
	nodes.Get("/:id/stats", nodeController.GetStats)
	nodes.Patch("/:id", nodeController.Update)
	nodes.Delete("/:id", nodeController.Delete)

	// Server management
	servers := admin.Group("/servers")
	servers.Get("/", serverController.GetAll)
	servers.Post("/", serverController.Create)
	servers.Get("/:id", serverController.GetByID)
	servers.Patch("/:id", serverController.Update)
	servers.Delete("/:id", serverController.Delete)
	servers.Post("/:id/power", serverController.Power)

	// Egg management
	eggs := admin.Group("/eggs")
	eggs.Get("/", eggController.GetAll)
	eggs.Post("/", eggController.Create)
	eggs.Get("/:id", eggController.GetByID)
	eggs.Patch("/:id", eggController.Update)
	eggs.Delete("/:id", eggController.Delete)

	// Client routes (protect all)
	client := app.Group("/api/client", middleware.Auth(authCfg))

	// Client servers
	client.Get("/servers", serverController.GetMyServers)
	client.Get("/servers/:uuid", serverController.GetMyServer)
	client.Get("/servers/:uuid/resources", serverController.GetResources)

	// WebSocket console proxy (requires custom upgrade handling in main or separate router)
	// For Fiber, we use app.Get("/ws/...", websocket.New(...))
	// We'll add this directly in main.go, or create a dedicated function
	// But for simplicity we can do it here using fiber's websocket middleware
	// However, we need access to Node and wings client. We'll handle this in main.

	// Client account - stubs for now (Phase 2 will extend)
	client.Get("/account", func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(models.User)
		if !ok {
			return utils.Unauthorized(c, "Not authenticated")
		}
		return utils.Success(c, utils.FromUser(&user), "Account info")
	})
	client.Patch("/account", func(c *fiber.Ctx) error {
		// Placeholder - Phase 2 will implement email/password update
		return utils.Success(c, nil, "Account update not implemented in Phase 1")
	})
}
