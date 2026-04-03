package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"nexus/backend/config"
	"nexus/backend/controllers"
	"nexus/backend/database"
	"nexus/backend/middleware"
	"nexus/backend/models"
	"nexus/backend/routes"
	"nexus/backend/wings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/contrib/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Auto-migrate all models
	if err := autoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize controllers
	authCtl := &controllers.AuthController{}
	userCtl := &controllers.UserController{}
	serverCtl := &controllers.ServerController{}
	nodeCtl := &controllers.NodeController{}
	eggCtl := &controllers.EggController{}
	controllers.Init(cfg)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
		// Production timeouts
		ReadTimeout:  30 * 1000 * 1000 * 1000, // 30s
		WriteTimeout: 30 * 1000 * 1000 * 1000, // 30s
		IdleTimeout:  120 * 1000 * 1000 * 1000, // 120s
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * 60 * 1000 * 1000 * 1000, // 1 minute
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	// Setup routes
	authCfg := &middleware.AuthConfig{JWTSecret: cfg.JWTSecret}
	routes.Setup(app, authCfg, authCtl, userCtl, serverCtl, nodeCtl, eggCtl)

	// WebSocket console route
	// Format: /ws/console?server_uuid=<uuid>
	app.Get("/ws/console", middleware.Auth(authCfg), websocket.New(func(c *websocket.Conn) {
		serverUUID := c.Query("server_uuid")
		if serverUUID == "" {
			c.WriteJSON(fiber.Map{"error": "server_uuid query parameter required"})
			c.Close()
			return
		}

		// Get user from context
		user, ok := c.Locals("user").(models.User)
		if !ok {
			c.WriteJSON(fiber.Map{"error": "Unauthorized"})
			c.Close()
			return
		}

		// Fetch server with node
		var server models.Server
		if err := database.DB.Where("uuid = ? AND user_id = ?", serverUUID, user.ID).
			Preload("Node").First(&server).Error; err != nil {
			c.WriteJSON(fiber.Map{"error": "Server not found or access denied"})
			c.Close()
			return
		}

		// Create console connection to Wings
		conn := wings.NewConsoleConnection(&server.Node, server.UUID)
		if err := conn.ConnectConsole(c); err != nil {
			log.Printf("Console connection error: %v", err)
		}
	}))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("%s starting on port %s", cfg.AppName, cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt
	<-quit
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// autoMigrate runs GORM auto-migration for all models
func autoMigrate() error {
	return database.DB.AutoMigrate(
		&models.User{},
		&models.Node{},
		&models.Server{},
		&models.Egg{},
		&models.Allocation{},
		&models.Ticket{},
		&models.CoinTransaction{},
	)
}

// errorHandler custom error handler
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
