package wsocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// In memory for testing, TODO

// Hub manages all connections
type Hub struct {
	clients map[string]*Client // address -> Client
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Client is one WebSocket connection
type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Address string // User ID
	Send    chan []byte
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client.Address] = client
	h.mu.Unlock()
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	delete(h.clients, client.Address)
	h.mu.Unlock()
}

func (h *Hub) GetClient(address string) (*Client, bool) {
	h.mu.RLock()
	client, ok := h.clients[address]
	h.mu.RUnlock()
	return client, ok
}
