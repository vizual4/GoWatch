package main

import (
	"GoWatch/internal/handler"
	"GoWatch/internal/repository"
	"GoWatch/internal/server"
	"GoWatch/internal/service"
	"GoWatch/internal/worker"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := worker.NewWorkerPool(5, repo)
	wp.Start(ctx)

	scheduler := worker.NewScheduler(repo, wp)

	go scheduler.Start(ctx)

	serv := server.NewServer(logger, h)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	sig := <-sigChan
	logger.Printf("Received signal: %v", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := serv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Server shutdown error: %v", err)
	} else {
		logger.Println("Server stopped gracefully")
	}

	cancel()
	logger.Println("Scheduler and workers stopped")

	if err := repo.Close(); err != nil {
		logger.Printf("Database close error: %v", err)
	} else {
		logger.Println("Databese closed")
	}
}
