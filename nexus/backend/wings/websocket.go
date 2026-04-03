package wings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"nexus/backend/models"

	gorillaws "github.com/gorilla/websocket"
	fiberws "github.com/gofiber/contrib/websocket"
)

// ConsoleConnection manages a WebSocket connection to Wings console
type ConsoleConnection struct {
	Node         *models.Node
	ServerUUID   string
	WingsConn    *gorillaws.Conn      // gorilla websocket (upstream to Wings)
	ClientConn   *fiberws.Conn        // fiber websocket (downstream to client)
	Done         chan bool
}

// NewConsoleConnection creates a new console connection manager
func NewConsoleConnection(node *models.Node, serverUUID string) *ConsoleConnection {
	return &ConsoleConnection{
		Node:       node,
		ServerUUID: serverUUID,
		Done:       make(chan bool, 1),
	}
}

// ConnectConsole establishes WebSocket to Wings and starts message forwarding
func (cc *ConsoleConnection) ConnectConsole(clientConn *fiberws.Conn) error {
	cc.ClientConn = clientConn

	// Connect to Wings WebSocket
	wingsURL := fmt.Sprintf("%s://%s:%d/api/servers/%s/ws",
		cc.Node.Scheme, cc.Node.FQDN, cc.Node.WingsPort)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+cc.Node.Token)
	header.Set("Origin", fmt.Sprintf("%s://%s", cc.Node.Scheme, cc.Node.FQDN))

	conn, _, err := gorillaws.DefaultDialer.Dial(wingsURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to wings: %w", err)
	}
	cc.WingsConn = conn

	// Start goroutines for bidirectional message forwarding
	go cc.readFromWings()
	go cc.readFromClient()

	// Wait for completion
	<-cc.Done
	return nil
}

// readFromWings reads messages from Wings and forwards to client
func (cc *ConsoleConnection) readFromWings() {
	defer func() {
		cc.Done <- true
		cc.closeAll()
	}()

	for {
		messageType, p, err := cc.WingsConn.ReadMessage()
		if err != nil {
			log.Printf("Wings read error: %v", err)
			return
		}

		if err := cc.ClientConn.WriteMessage(messageType, p); err != nil {
			log.Printf("Client write error: %v", err)
			return
		}
	}
}

// readFromClient reads messages from client and forwards to Wings
func (cc *ConsoleConnection) readFromClient() {
	defer func() {
		cc.Done <- true
		cc.closeAll()
	}()

	for {
		messageType, p, err := cc.ClientConn.ReadMessage()
		if err != nil {
			log.Printf("Client read error: %v", err)
			return
		}

		// Parse message to determine if it's a command
		var msg map[string]interface{}
		if bytes.HasPrefix(p, []byte("{")) {
			if json.Unmarshal(p, &msg) == nil {
				// It's a JSON message - pass through
				cc.WingsConn.WriteMessage(messageType, p)
				return
			}
		}

		// Plain text command - send as command
		commandMsg := map[string]string{
			"cmd": string(p),
		}
		cmdJSON, _ := json.Marshal(commandMsg)

		if err := cc.WingsConn.WriteMessage(gorillaws.TextMessage, cmdJSON); err != nil {
			log.Printf("Wings write error: %v", err)
			return
		}
	}
}

// closeAll closes both connections
func (cc *ConsoleConnection) closeAll() {
	if cc.WingsConn != nil {
		cc.WingsConn.Close()
	}
	if cc.ClientConn != nil {
		cc.ClientConn.Close()
	}
}

// HandleConsoleWebSocket upgrades the HTTP connection and handles console proxying
func HandleConsoleWebSocket(node *models.Node, serverUUID string) func(c *fiberws.Conn) {
	return func(c *fiberws.Conn) {
		conn := NewConsoleConnection(node, serverUUID)
		if err := conn.ConnectConsole(c); err != nil {
			log.Printf("Console connection failed: %v", err)
			c.WriteJSON(map[string]string{
				"error": fmt.Sprintf("Failed to connect: %v", err),
			})
			c.Close()
		}
	}
}
