package repository

import (
	"GoWatch/internal/domain"
	"fmt"
	"sync"
)

type MonitorRepository interface {
	Create(data *domain.Monitor) error
	GetAllMon() ([]*domain.Monitor, error)
	GetMon(id int64) (*domain.Monitor, error)
	DelMon(id int64) error
}

type MemoryRepository struct {
	mu     sync.RWMutex
	memory map[int64]*domain.Monitor
	index  int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		memory: make(map[int64]*domain.Monitor),
		index:  1,
	}
}

func (r *MemoryRepository) Create(data *domain.Monitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data.ID = r.index
	r.memory[r.index] = data
	r.index++

	return nil
}

func (r *MemoryRepository) GetAllMon() ([]*domain.Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.memory) == 0 {
		return nil, fmt.Errorf("no values in memory")
	}

	res := make([]*domain.Monitor, 0, len(r.memory))

	for _, v := range r.memory {
		res = append(res, v)
	}

	return res, nil
}

func (r *MemoryRepository) GetMon(id int64) (*domain.Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res, ok := r.memory[id]
	if !ok {
		return nil, NF
	}

	return res, nil
}

func (r *MemoryRepository) DelMon(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.memory[id]; !ok {
		return NF
	}

	delete(r.memory, id)

	return nil
}
