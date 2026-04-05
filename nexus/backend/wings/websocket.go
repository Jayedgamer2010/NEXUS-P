package wings

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"nexus/backend/models"

	"github.com/gorilla/websocket"
)

type WSProxy struct {
	node      models.Node
	uuid      string
	upgrader  websocket.Upgrader
	mu        sync.Mutex
	wsWings   *websocket.Conn
}

func NewWSProxy(node models.Node, uuid string) *WSProxy {
	return &WSProxy{
		node: node,
		uuid: uuid,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (p *WSProxy) BuildWingsWebSocketURL() string {
	scheme := "wss"
	if p.node.Scheme == "http" {
		scheme = "ws"
	}
	return fmt.Sprintf("%s://%s:%d/api/servers/%s/ws", scheme, p.node.FQDN, p.node.DaemonListen, p.uuid)
}

func (p *WSProxy) BuildAuthHeader() string {
	return fmt.Sprintf("Bearer %s.%s", p.node.DaemonTokenID, p.node.DaemonToken)
}

func (p *WSProxy) ConnectToWings() (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	header := make(http.Header)
	header.Set("Authorization", p.BuildAuthHeader())
	header.Set("Origin", p.node.GetConnectionAddress())

	conn, _, err := dialer.Dial(p.BuildWingsWebSocketURL(), header)
	if err != nil {
		return nil, fmt.Errorf("dial wings websocket: %w", err)
	}

	p.mu.Lock()
	p.wsWings = conn
	p.mu.Unlock()

	return conn, nil
}

func (p *WSProxy) CloseWings() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.wsWings != nil {
		p.wsWings.Close()
		p.wsWings = nil
	}
}

func (p *WSProxy) GetUpgrader() *websocket.Upgrader {
	return &p.upgrader
}
