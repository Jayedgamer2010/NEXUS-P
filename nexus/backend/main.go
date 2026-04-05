package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/middleware"
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/services"
	"nexus/backend/utils"
	"nexus/backend/wings"

	adminCtrl "nexus/backend/controllers/admin"
	clientCtrl "nexus/backend/controllers/client"
	"nexus/backend/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Connect database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// 3. Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	serverRepo := repositories.NewServerRepository(database.DB)
	nodeRepo := repositories.NewNodeRepository(database.DB)
	eggRepo := repositories.NewEggRepository(database.DB)
	allocRepo := repositories.NewAllocationRepository(database.DB)

	// 4. Initialize wings client
	wingsClient := wings.NewWingsClient()

	// 5. Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	serverSvc := services.NewServerService(serverRepo, nodeRepo, eggRepo, userRepo, allocRepo, wingsClient)
	nodeSvc := services.NewNodeService(nodeRepo, serverRepo, wingsClient)
	userSvc := services.NewUserService(userRepo, serverRepo)

	// 6. Initialize controllers
	authController := controllers.NewAuthController(authService)
	statsController := adminCtrl.NewStatsController(userRepo, serverRepo, nodeRepo)
	userController := adminCtrl.NewUserController(userRepo, userSvc)
	serverController := adminCtrl.NewServerController(serverRepo, serverSvc)
	nodeController := adminCtrl.NewNodeController(nodeRepo, serverRepo, allocRepo, nodeSvc)
	eggController := adminCtrl.NewEggController(eggRepo, serverRepo)
	clientServerController := clientCtrl.NewClientServerController(serverRepo, serverSvc, nodeRepo)
	clientAccountController := clientCtrl.NewAccountController()

	// 7. Seed admin user
	seedAdmin(userRepo, cfg)

	// 8. Seed default eggs
	seedEggs(eggRepo)

	// 9. Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:   cfg.AppName,
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// CORS first
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Preflight
	app.Options("/*", func(c *fiber.Ctx) error { return c.SendStatus(204) })

	// Health (no auth)
	app.Get("/health", healthHandler)

	// Auth (no auth)
	app.Post("/api/auth/register", authController.Register)
	app.Post("/api/auth/login", authController.Login)

	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Protected routes
	api := app.Group("/api", authMiddleware)
	api.Get("/auth/me", authController.GetMe)

	// Client routes
	client := api.Group("/client")
	client.Get("/account", clientAccountController.Get)
	client.Patch("/account", clientAccountController.Update)
	client.Get("/servers", clientServerController.GetAll)
	client.Get("/servers/:uuid", clientServerController.GetOne)
	client.Get("/servers/:uuid/resources", clientServerController.GetResources)
	client.Post("/servers/:uuid/power", clientServerController.PowerAction)

	// Admin middleware
	adminMiddleware := middleware.AdminMiddleware()

	// Admin routes
	admin := api.Group("/admin", adminMiddleware)
	admin.Get("/stats", statsController.GetStats)

	admin.Get("/users", userController.GetAll)
	admin.Post("/users", userController.Create)
	admin.Get("/users/:id", userController.GetOne)
	admin.Patch("/users/:id", userController.Update)
	admin.Delete("/users/:id", userController.Delete)

	admin.Get("/servers", serverController.GetAll)
	admin.Post("/servers", serverController.Create)
	admin.Get("/servers/:id", serverController.GetOne)
	admin.Patch("/servers/:id", serverController.Update)
	admin.Delete("/servers/:id", serverController.Delete)
	admin.Post("/servers/:id/power", serverController.PowerAction)
	admin.Post("/servers/:id/suspend", serverController.Suspend)
	admin.Post("/servers/:id/unsuspend", serverController.Unsuspend)

	admin.Get("/nodes", nodeController.GetAll)
	admin.Post("/nodes", nodeController.Create)
	admin.Get("/nodes/:id", nodeController.GetOne)
	admin.Patch("/nodes/:id", nodeController.Update)
	admin.Delete("/nodes/:id", nodeController.Delete)
	admin.Get("/nodes/:id/allocations", nodeController.GetAllocations)
	admin.Post("/nodes/:id/allocations", nodeController.AddAllocation)
	admin.Delete("/allocations/:id", nodeController.DeleteAllocation)

	admin.Get("/eggs", eggController.GetAll)
	admin.Post("/eggs", eggController.Create)
	admin.Get("/eggs/:id", eggController.GetOne)
	admin.Patch("/eggs/:id", eggController.Update)
	admin.Delete("/eggs/:id", eggController.Delete)

	// WebSocket console - handled by ws_handler.go
	registerWSConsole(app)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
	}()

	// Start server
	log.Printf("NEXUS panel starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var startTime = time.Now()

func healthHandler(c *fiber.Ctx) error {
	return utils.Success(c, fiber.Map{
		"status": "ok",
		"uptime": time.Since(startTime).String(),
	})
}

func seedAdmin(userRepo *repositories.UserRepository, cfg *config.Config) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		return
	}

	if _, err := userRepo.FindByEmail(cfg.AdminEmail); err == nil {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	username := cfg.AdminUsername
	if username == "" {
		username = "admin"
	}

	adminUser := &models.User{
		Username:  username,
		Email:     cfg.AdminEmail,
		Role:      "admin",
		RootAdmin: true,
	}

	if err := adminUser.HashPassword(cfg.AdminPassword); err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	if err := userRepo.Create(adminUser); err != nil {
		log.Printf("Failed to seed admin user: %v", err)
		return
	}

	log.Printf("Admin user seeded: %s", adminUser.Email)
}

func seedEggs(eggRepo *repositories.EggRepository) {
	eggs, err := eggRepo.FindAll()
	if err == nil && len(eggs) > 0 {
		log.Println("Eggs already exist, skipping seed")
		return
	}

	defaultEggs := []models.Egg{
		{
			Author:      "NEXUS",
			Name:        "Vanilla Minecraft",
			Description: "Standard Minecraft server with vanilla gameplay",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -XX:MaxRAMPercentage=95.0 -jar {{SERVER_JARFILE}}",
			ConfigStop:  "stop",
		},
		{
			Author:      "NEXUS",
			Name:        "Paper MC",
			Description: "Paper Minecraft server - high performance fork",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -XX:MaxRAMPercentage=95.0 -jar {{SERVER_JARFILE}}",
			ConfigStop:  "stop",
		},
		{
			Author:      "NEXUS",
			Name:        "BungeeCord",
			Description: "BungeeCord proxy for server networking",
			DockerImage: "ghcr.io/pterodactyl/yolks:java_17",
			Startup:     "java -Xms128M -XX:MaxRAMPercentage=95.0 -jar {{SERVER_JARFILE}}",
			ConfigStop:  "end",
		},
	}

	for _, egg := range defaultEggs {
		if err := eggRepo.Create(&egg); err != nil {
			log.Printf("Failed to seed egg %s: %v", egg.Name, err)
		} else {
			log.Printf("Egg seeded: %s", egg.Name)
		}
	}
}
