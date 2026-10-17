package worker

import (
	"GoWatch/internal/domain"
	"context"
	"log"
	"net/http"
	"time"
)

type CheckRepo interface {
	SaveCheck(ctx context.Context, check domain.CheckResult) error
}

type WorkerPool struct {
	numWorkers int
	jobs       chan domain.Monitor
	repo       CheckRepo
	client     *http.Client
}

func NewWorkerPool(num int, c CheckRepo) *WorkerPool {
	return &WorkerPool{
		numWorkers: num,
		jobs:       make(chan domain.Monitor, 100),
		repo:       c,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WorkerPool) Start(ctx context.Context) {

	log.Printf("Запуск %d воркеров", w.numWorkers)

	for i := 1; i <= w.numWorkers; i++ {
		go w.worker(ctx, i)
	}

}

func (w *WorkerPool) worker(ctx context.Context, id int) {
	log.Printf("Воркер %d готов к работе", id)

	for {
		select {
		case <-ctx.Done():
			return
		case monitor, ok := <-w.jobs:
			if !ok {
				return
			}

			w.process(ctx, monitor, id)
		}
	}
}

func (w *WorkerPool) process(ctx context.Context, monitor domain.Monitor, id int) {
	var statusCode int
	var errMsg *string

	start := time.Now()

	resp, err := w.client.Get(monitor.URL)
	durationMs := int(time.Since(start).Milliseconds())
	if err != nil {
		errMsgStr := err.Error()
		errMsg = &errMsgStr
		statusCode = 0
	} else {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	res := domain.CheckResult{
		MonitorID:      monitor.ID,
		StatusCode:     statusCode,
		ResponseTimeMs: durationMs,
		ErrorMsg:       errMsg,
		CheckedAt:      start,
	}

	if w.repo != nil {
		err = w.repo.SaveCheck(ctx, res)
		if err != nil {
			log.Printf("worker %d error save check in DB for monitor %d: %v", id, monitor.ID, err)
		}
	}
}

func (w *WorkerPool) AddJob(m domain.Monitor) {
	w.jobs <- m
}
