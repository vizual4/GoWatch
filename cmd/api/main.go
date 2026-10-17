package main

import (
	"GoWatch/internal/handler"
	"GoWatch/internal/repository"
	"GoWatch/internal/server"
	"GoWatch/internal/service"
	"GoWatch/internal/worker"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "INFO:", log.LstdFlags)

	dsn := os.Getenv("DATABASE_URL")
	if len(dsn) == 0 {
		dsn = `postgres://gowatch_user:gowatch_pass@localhost:5432/gowatch_db?sslmode=disable`
	}

	repo, err := repository.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewMonitorService(repo)

	h := handler.NewMonitorHandler(svc)

	ctx := context.Background()

	wp := worker.NewWorkerPool(5, repo)
	wp.Start(ctx)

	scheduler := worker.NewScheduler(repo, wp, 10*time.Second)

	go scheduler.Start(ctx)

	serv := server.NewServer(logger, h)
	if err = serv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
