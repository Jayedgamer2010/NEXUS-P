package main

import (
	"log"
	"sync"
	"time"

	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"
	"nexus/backend/wings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gorilla "github.com/gorilla/websocket"
)

func registerWSConsole(app *fiber.App) {
	app.Get("/ws/console", websocket.New(func(c *websocket.Conn) {
		cfg := config.Load()

		// Extract token from query param (standard for WS auth)
		tokenStr := c.Query("token")
		if tokenStr == "" {
			return
		}

		claims, err := utils.ValidateToken(tokenStr, cfg)
		if err != nil {
			return
		}

		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			return
		}

		if user.Suspended {
			return
		}

		serverUUID := c.Query("server")
		if serverUUID == "" {
			return
		}

		var server models.Server
		if err := database.DB.Where("uuid = ? AND user_id = ?", serverUUID, user.ID).First(&server).Error; err != nil {
			return
		}

		var node models.Node
		if err := database.DB.First(&node, server.NodeID).Error; err != nil {
			return
		}

		wsProxy := wings.NewWSProxy(node, server.UUID)
		wingsConn, err := wsProxy.ConnectToWings()
		if err != nil {
			log.Printf("Wings WS connection failed: %v", err)
			return
		}
		defer wsProxy.CloseWings()

		var wg sync.WaitGroup
		done := make(chan struct{})

		// Browser -> Wings
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					return
				}
				if err := wingsConn.WriteMessage(gorilla.TextMessage, msg); err != nil {
					return
				}
			}
		}()

		// Wings -> Browser
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, msg, err := wingsConn.ReadMessage()
				if err != nil {
					return
				}
				if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
		}()

		go func() {
			wg.Wait()
			close(done)
		}()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
					return
				}
			}
		}
	}))
}
