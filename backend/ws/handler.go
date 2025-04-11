package ws

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/arnavr0/sswan/internal/jsonlog"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type WsHandler struct {
	logger *jsonlog.Logger
	hub    *Hub
}

func NewWsHandler(logger *jsonlog.Logger, hub *Hub) *WsHandler {
	return &WsHandler{
		logger: logger,
		hub:    hub,
	}
}

func (h *WsHandler) ServeWS(w http.ResponseWriter, r *http.Request) {

	//Log info
	h.logger.PrintInfo("Received websocket upgrade request", map[string]string{
		"remote_addr": r.RemoteAddr,
		"path":        r.URL.Path,
	})

	//Accept the connection
	// NOTE: InsecureSkipVerify is for DEVELOPMENT ONLY to allow connections
	// from localhost without TLS. REMOVE THIS and configure OriginPatterns
	// properly for production deployments with HTTPS/WSS.
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		// OriginPatterns: []string{"localhost:4000", "yourdomain.com"}, // Production
	})
	if err != nil {
		h.logger.PrintError(err, map[string]string{
			"remote_addr": r.RemoteAddr,
			"component":   "websocket.Accept",
		})
		return
	}

	client := &Client{
		ID:   uuid.NewString(),
		Conn: conn,
		Room: "", // initially no room
	}
	h.hub.Register(client)
	h.logger.PrintInfo("Websocket connection established and client created", map[string]string{
		"remote_addr": r.RemoteAddr,
		"clientID":    client.ID,
	})

	// Start goroutine for the client
	go func() {
		logger := h.logger
		remoteAddr := r.RemoteAddr
		clientID := client.ID

		// Seperate defer for conn.Close to manage closedIntentionally
		defer func() {
			logger.PrintInfo("Goroutine exiting, unregistering client", map[string]string{"clientID": client.ID})
			h.hub.Unregister(client)
			logger.PrintInfo("Closing connection in defer", map[string]string{"remote_addr": remoteAddr})
			conn.Close(websocket.StatusInternalError, "Internal server error occurred")
		}()

		closedIntentionally := false // we want it to close intentionally
		defer func() {
			if !closedIntentionally {
				logger.PrintInfo("Closing connection via inner defer (abnormal exit)", map[string]string{"clientID": clientID})
			} else {
				logger.PrintInfo("Skipping close in inner defer (already handled)", map[string]string{"clientID": clientID})
			}
		}()

		ctx := context.Background()

		logger.PrintInfo("Client reader goroutine started", map[string]string{"remote_addr": remoteAddr})

		for {
			readCtx, cancelRead := context.WithTimeout(ctx, 10*time.Second)
			// Read structured json message
			var msg SignalMessage
			err := wsjson.Read(readCtx, conn, &msg)
			cancelRead()

			if err != nil {
				closeStatus := websocket.CloseStatus(err)
				properties := map[string]string{"remote_addr": remoteAddr, "status_code": closeStatus.String()}

				if closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway {
					logger.PrintInfo("Client disconnected normally", properties)
				} else if errors.Is(err, context.DeadlineExceeded) {
					// This DeadlineExceeded is from our readCtx timeout
					logger.PrintInfo("Client read timeout", properties)
					errClose := conn.Close(websocket.StatusPolicyViolation, "Idle timeout")
					closedIntentionally = (errClose == nil) // Set flag to true if close initiated successfully
				} else if errors.Is(err, context.Canceled) {
					// This Canceled likely means the goroutine's base ctx was canceled (e.g., server shutdown)
					logger.PrintInfo("Context canceled for client", properties)
					// Close normally if cancellation was the reason
					conn.Close(websocket.StatusNormalClosure, "Context canceled")
				} else {
					properties["component"] = "wsjson.Read"
					logger.PrintError(err, properties)
					if closeStatus == -1 {
						errClose := conn.Close(websocket.StatusUnsupportedData, "Read error")
						closedIntentionally = (errClose == nil)
						if !closedIntentionally {
							logger.PrintError(errClose, map[string]string{"clientID": clientID, "component": "conn.Close(read_error)"})
						}
					} else {
						closedIntentionally = true // Assume peer closed uncleanly, we handled it
					}

				}
				break // Exit loop on any error
			}

			msg.Sender = client.ID

			logger.PrintInfo("Received message", map[string]string{
				"clientID": clientID,
				"type":     msg.Type,
				"sender":   msg.Sender,
				"payload":  msg.Payload.(string),
			})

			h.hub.Broadcast(ctx, msg, client) // Broadcast the message

		}

		logger.PrintInfo("Client reader goroutine finished", map[string]string{"remote_addr": remoteAddr})
	}() // End of goroutine

	// Log that the main handler function has finished setup
	h.logger.PrintInfo("ServeWS handler finished setup for client", map[string]string{
		"clientID": client.ID, // Log client ID here too
	})
}
