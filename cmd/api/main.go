package main

import (
	"GoWatch/internal/server"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "INFO:", log.LstdFlags)
	serv := server.NewServer(logger)
	err := serv.ListenAndServe()
	logger.Fatal(err)
}
