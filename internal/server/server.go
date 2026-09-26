package server

import (
	"GoWatch/internal/handler"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(logger *log.Logger) *http.Server {
	mux := http.NewServeMux()
	handler.Init(mux)

	port := os.Getenv("GOWATCH_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	server := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
		IdleTimeout:  10 * time.Millisecond,
	}

	return &server
}
