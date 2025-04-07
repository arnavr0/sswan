package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arnavr0/sswan/internal/jsonlog"
	"github.com/arnavr0/sswan/ws"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	wsHandler := ws.NewWsHandler(app.logger)

	http.HandleFunc("/ws", wsHandler.ServeWS)
	// Simple root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("WebSocket Server is Running"))
		app.logger.PrintInfo("Handled root request", map[string]string{
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
		})
	})

	server := &http.Server{
		Addr:              ":4000",
		Handler:           http.DefaultServeMux,   // Use the default mux where we registered handlers
		ErrorLog:          log.New(logger, "", 0), // Use custom logger for server errors!
		ReadHeaderTimeout: 3 * time.Second,
		// Add other timeouts
	}

	err := server.ListenAndServe()
	app.logger.PrintInfo("server starting on", map[string]string{})
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		// Use PrintFatal which includes stack trace and exits
		app.logger.PrintFatal(err, nil)
	}
	app.logger.PrintInfo("Server stopped", nil)
}
