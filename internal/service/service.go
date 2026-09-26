package service

import (
	"GoWatch/internal/domain"
	"GoWatch/internal/repository"
	"fmt"
	"net/url"
	"time"
)

func CreateMon(data *domain.Monitor) (int64, error) {
	err := check(data)
	if err != nil {
		return 0, err
	}

	data.CreatedAt = time.Now()

	err = repository.Create(data)
	if err != nil {
		return 0, err
	}

	return data.ID, nil
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
