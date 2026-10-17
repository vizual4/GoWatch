package worker

import (
	"GoWatch/internal/repository"
	"context"
	"log"
	"time"
)

type Scheduler struct {
	repo     repository.PostgresRepository
	work     *WorkerPool
	interval time.Duration
}

func NewScheduler(repo repository.PostgresRepository, work *WorkerPool, t time.Duration) *Scheduler {
	return &Scheduler{
		repo:     repo,
		work:     work,
		interval: t,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			res, err := s.repo.GetAllMon()
			if err != nil {
				log.Printf("Scheduler can't get monitors, error: %v", err)
				continue
			}

			for _, v := range res {
				s.work.AddJob(*v)
			}
		}
	}
}
