package server

import (
	"GoWatch/internal/handler"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(logger *log.Logger, h *handler.MonitorHandler) *http.Server {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	port := os.Getenv("GOWATCH_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	server := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return &server
}
