package ws

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/arnavr0/sswan/internal/jsonlog"
	"github.com/coder/websocket"
)

type WsHandler struct {
	logger *jsonlog.Logger
}

func NewWsHandler(logger *jsonlog.Logger) *WsHandler {
	return &WsHandler{
		logger: logger,
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

	h.logger.PrintInfo("Websocket connection established", map[string]string{
		"remote_addr": r.RemoteAddr,
	})

	// Start goroutine for the client
	go func() {
		logger := h.logger
		remoteAddr := r.RemoteAddr

		defer func() {
			logger.PrintInfo("Closing connection in defer", map[string]string{"remote_addr": remoteAddr})
			conn.Close(websocket.StatusInternalError, "Internal server error occurred")
		}()

		ctx := context.Background()

		logger.PrintInfo("Client reader goroutine started", map[string]string{"remote_addr": remoteAddr})

		for {
			readCtx, cancelRead := context.WithTimeout(ctx, 10*time.Second)
			msgType, p, err := conn.Read(readCtx)
			cancelRead()

			if err != nil {
				closeStatus := websocket.CloseStatus(err)
				properties := map[string]string{"remote_addr": remoteAddr, "status_code": closeStatus.String()}

				if closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway {
					logger.PrintInfo("Client disconnected normally", properties)
				} else if errors.Is(err, context.DeadlineExceeded) {
					// This DeadlineExceeded is from our readCtx timeout
					logger.PrintInfo("Client read timeout", properties)
					conn.Close(websocket.StatusPolicyViolation, "Idle timeout")
				} else if errors.Is(err, context.Canceled) {
					// This Canceled likely means the goroutine's base ctx was canceled (e.g., server shutdown)
					logger.PrintInfo("Context canceled for client", properties)
					// Close normally if cancellation was the reason
					conn.Close(websocket.StatusNormalClosure, "Context canceled")
				} else {
					properties["component"] = "conn.Read"
					logger.PrintError(err, properties)
					if closeStatus == -1 {
						conn.Close(websocket.StatusUnsupportedData, "Read error")
					}
				}
				break // Exit loop on any error
			}

			logger.PrintInfo("Received message", map[string]string{
				"remote_addr": remoteAddr,
				"type":        msgType.String(),
				"content":     string(p),
			})

		}

		logger.PrintInfo("Client reader goroutine finished", map[string]string{"remote_addr": remoteAddr})
	}() // End of goroutine

}
