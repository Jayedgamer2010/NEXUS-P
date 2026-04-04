package wings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"nexus/backend/models"

	gorillaws "github.com/gorilla/websocket"
	fiberws "github.com/gofiber/contrib/websocket"
)

type ConsoleConnection struct {
	Node       *models.Node
	ServerUUID string
	WingsConn  *gorillaws.Conn
	ClientConn *fiberws.Conn
	Done       chan bool
}

func NewConsoleConnection(node *models.Node, serverUUID string) *ConsoleConnection {
	return &ConsoleConnection{
		Node:       node,
		ServerUUID: serverUUID,
		Done:       make(chan bool, 1),
	}
}

func (cc *ConsoleConnection) ConnectConsole(clientConn *fiberws.Conn) error {
	cc.ClientConn = clientConn

	wingsURL := fmt.Sprintf("%s://%s:%d/api/servers/%s/ws",
		cc.Node.Scheme, cc.Node.FQDN, cc.Node.DaemonListen)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+cc.Node.DaemonToken)
	header.Set("Origin", fmt.Sprintf("%s://%s", cc.Node.Scheme, cc.Node.FQDN))

	conn, _, err := gorillaws.DefaultDialer.Dial(wingsURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to wings: %w", err)
	}
	cc.WingsConn = conn

	go cc.readFromWings()
	go cc.readFromClient()

	<-cc.Done
	return nil
}

func (cc *ConsoleConnection) readFromWings() {
	defer func() {
		cc.Done <- true
		cc.closeAll()
	}()

	for {
		messageType, p, err := cc.WingsConn.ReadMessage()
		if err != nil {
			return
		}

		if err := cc.ClientConn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func (cc *ConsoleConnection) readFromClient() {
	defer func() {
		cc.Done <- true
		cc.closeAll()
	}()

	for {
		messageType, p, err := cc.ClientConn.ReadMessage()
		if err != nil {
			return
		}

		var msg map[string]interface{}
		if bytes.HasPrefix(p, []byte("{")) {
			if err := json.Unmarshal(p, &msg); err == nil {
				cc.WingsConn.WriteMessage(messageType, p)
				continue
			}
		}

		commandMsg := map[string]string{"cmd": string(p)}
		cmdJSON, _ := json.Marshal(commandMsg)

		if err := cc.WingsConn.WriteMessage(gorillaws.TextMessage, cmdJSON); err != nil {
			return
		}
	}
}

func (cc *ConsoleConnection) closeAll() {
	if cc.WingsConn != nil {
		cc.WingsConn.Close()
	}
	if cc.ClientConn != nil {
		cc.ClientConn.Close()
	}
}
