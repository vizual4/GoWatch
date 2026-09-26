package repository

import (
	"GoWatch/internal/domain"
	"sync"
)

var memory = make(map[int64]domain.Monitor)
var index int64 = 1
var mu sync.RWMutex

func Create(data *domain.Monitor) error {
	mu.Lock()
	defer mu.Unlock()

	data.ID = index
	memory[index] = *data
	index++

	return nil
}
