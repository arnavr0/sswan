package ws

import (
	"context"
	"sync"
	"time"

	"github.com/arnavr0/sswan/internal/jsonlog"
	"github.com/coder/websocket/wsjson"
)

// Hub manages the set of active clients and broadcasts messages.
// It ensures concurrent access to client data is safe.
type Hub struct {
	clients map[*Client]bool
	mu      sync.Mutex
	logger  *jsonlog.Logger
}

func NewHub(logger *jsonlog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		logger:  logger,
	}
}

// Register adds a new client to the Hub
func (h *Hub) Register(client *Client) {
	h.mu.Lock() // Lock mutex before accessing shared resrouce (the clients map)
	defer h.mu.Unlock()

	h.clients[client] = true
	h.logger.PrintInfo("Client registered", map[string]string{"clientID": client.ID})
}
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		h.logger.PrintInfo("Client unregistered", map[string]string{"clientID": client.ID})
	}
}

// Broadcast sends a message to ALL connected clients EXCEPT the sender.
// TODO: This needs significant enhancement for rooms and targeted sending later.
// This simple broadcast is just for initial testing
func (h *Hub) Broadcast(ctx context.Context, msg SignalMessage, sender *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.logger.PrintInfo("Broadcasting message", map[string]string{
		"senderID": sender.ID,
		"type":     msg.Type,
		// "target": msg.Target, // Add later when targetting exists
		// "payload": fmt.Sprintf("%+v", msg.Payload), // Be careful logging large payloads
	})

	// failedClients := []*Clients{}

	for client := range h.clients {
		if client != sender {
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

			err := wsjson.Write(writeCtx, client.Conn, msg)
			if err != nil {
				h.logger.PrintError(err, map[string]string{"clientID": client.ID, "component": "hub.broadcast"})
			}
			cancel() // release timeout context
		}
	}
}
