package ws

import "github.com/coder/websocket"

// Client represents a single connected user/Websocket connection
type Client struct {
	ID   string // Unique identifier for the client (e.g., generated UUID or derived from connection)
	Conn *websocket.Conn
	Room string
	// hub *Hub // Optional: A reference back to the hub if needed for direct calls
}
