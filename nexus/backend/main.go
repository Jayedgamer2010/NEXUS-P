package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"nexus/backend/config"
	"nexus/backend/controllers"
	"nexus/backend/controllers/admin"
	"nexus/backend/controllers/client"
	"nexus/backend/database"
	"nexus/backend/middleware"
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/routes"
	"nexus/backend/services"
	"nexus/backend/wings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := autoMigrate(); err != nil {
		log.Printf("Warning: migration issues: %v", err)
	}

	seedAdminFromEnv(database.DB)

	db := database.DB
	wingsSvc := services.NewWingsService()

	userRepo := repositories.NewUserRepository(db)
	serverRepo := repositories.NewServerRepository(db)
	nodeRepo := repositories.NewNodeRepository(db)
	eggRepo := repositories.NewEggRepository(db)
	allocRepo := repositories.NewAllocationRepository(db)

	userSvc := services.NewUserService(userRepo, cfg)
	serverSvc := services.NewServerService(serverRepo, nodeRepo, eggRepo, userRepo, allocRepo, wingsSvc)
	nodeSvc := services.NewNodeService(nodeRepo, allocRepo, wingsSvc)
	eggSvc := services.NewEggService(eggRepo, serverRepo, db)

	// Seed default eggs
	if err := eggSvc.SeedDefaults(); err != nil {
		log.Printf("Warning: egg seeding issues: %v", err)
	}

	authCtl := controllers.NewAuthController(userSvc)
	userCtl := admin.NewUserController(userSvc, allocRepo)
	serverCtl := admin.NewServerController(serverSvc)
	nodeCtl := admin.NewNodeController(nodeSvc)
	eggCtl := admin.NewEggController(eggSvc)
	statsCtl := admin.NewStatsController(userRepo, serverRepo, nodeRepo, wingsSvc)
	clientSrvCtl := client.NewServerController(serverSvc, wingsSvc)
	accCtl := client.NewAccountController(userRepo)

	app := fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
		ReadTimeout:           30 * 1e9,
		WriteTimeout:          30 * 1e9,
		IdleTimeout:           120 * 1e9,
	})

	// CORS - must be first
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc:     func(origin string) bool { return true },
		AllowMethods:         "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",
		AllowHeaders:         "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials:     false,
		MaxAge:               300,
	}))
	app.Options("/*", func(c *fiber.Ctx) error { return c.SendStatus(204) })

	app.Use(recover.New())

	authCfg := &middleware.AuthConfig{JWTSecret: cfg.JWTSecret}
	routes.Setup(app, authCfg, authCtl, userCtl, serverCtl, nodeCtl, eggCtl, statsCtl, clientSrvCtl, accCtl)

	// WebSocket console
	app.Get("/ws/console", middleware.Auth(authCfg), websocket.New(func(c *websocket.Conn) {
		serverUUID := c.Query("server_uuid")
		if serverUUID == "" {
			c.WriteJSON(fiber.Map{"error": "server_uuid query parameter required"})
			c.Close()
			return
		}

		user, ok := c.Locals("user").(models.User)
		if !ok {
			c.WriteJSON(fiber.Map{"error": "Unauthorized"})
			c.Close()
			return
		}

		var server models.Server
		if err := db.Where("uuid = ? AND user_id = ?", serverUUID, user.ID).
			Preload("Node").First(&server).Error; err != nil {
			c.WriteJSON(fiber.Map{"error": "Server not found or access denied"})
			c.Close()
			return
		}

		conn := wings.NewConsoleConnection(&server.Node, server.UUID)
		if err := conn.ConnectConsole(c); err != nil {
			log.Printf("Console connection error: %v", err)
		}
	}))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("%s starting on port %s", cfg.AppName, cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func autoMigrate() error {
	db := database.DB
	db.Exec("SET session_replication_role = replica")
	defer db.Exec("SET session_replication_role = DEFAULT")

	return db.AutoMigrate(
		&models.User{},
		&models.Node{},
		&models.Egg{},
		&models.Allocation{},
		&models.Server{},
		&models.Ticket{},
		&models.CoinTransaction{},
	)
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"data":    nil,
	})
}

func seedAdminFromEnv(db *gorm.DB) {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	username := os.Getenv("ADMIN_USERNAME")
	if email == "" || password == "" {
		log.Println("Warning: ADMIN_EMAIL or ADMIN_PASSWORD not set, skipping admin seed")
		return
	}
	var user models.User
	if db.Where("email = ?", email).First(&user).Error == nil {
		return // already exists
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}
	admin := models.User{
		UUID:      uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  string(hashed),
		Role:      "admin",
		RootAdmin: true,
		Coins:     0,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin: %v", err)
		return
	}
	log.Println("Admin account created from environment variables")
}
