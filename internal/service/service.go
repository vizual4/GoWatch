package service

import (
	"GoWatch/internal/domain"
	"GoWatch/internal/repository"
	"fmt"
	"net/url"
	"time"
)

type MonitorService struct {
	repo repository.PostgresRepository
}

func NewMonitorService(repo repository.PostgresRepository) *MonitorService {
	return &MonitorService{
		repo: repo,
	}
}

func (s *MonitorService) CreateMon(data *domain.Monitor) error {
	if err := check(data); err != nil {
		return err
	}

	data.CreatedAt = time.Now()
	return s.repo.Create(data)
}

func (s *MonitorService) AllMons() ([]*domain.Monitor, error) {
	return s.repo.GetAllMon()
}

func (s *MonitorService) GetMon(id int64) (*domain.Monitor, error) {
	return s.repo.GetMon(id)
}

func (s *MonitorService) DelMon(id int64) error {
	return s.repo.DelMon(id)
}

func check(data *domain.Monitor) error {
	if len(data.Name) == 0 {
		return fmt.Errorf("name is empty")
	}

	u, err := url.ParseRequestURI(data.URL)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("wrong protocol")
	}

	if data.IntervalSec < 10 {
		return fmt.Errorf("wrong interval")
	}

	return nil
}
