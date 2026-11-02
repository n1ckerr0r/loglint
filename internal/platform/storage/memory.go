package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/n1ckerr0r/loglint/internal/service"
)

var ErrRunNotFound = errors.New("run not found")

type MemoryRepository struct {
	mu   sync.RWMutex
	runs map[string]service.Run
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{runs: make(map[string]service.Run)}
}

func (repo *MemoryRepository) CreateRun(_ context.Context, run service.Run) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.runs[run.ID] = run
	return nil
}

func (repo *MemoryRepository) GetRun(_ context.Context, id string) (service.Run, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	run, ok := repo.runs[id]
	if !ok {
		return service.Run{}, ErrRunNotFound
	}
	return run, nil
}

func (repo *MemoryRepository) UpdateRun(_ context.Context, id string, status service.RunStatus, report string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	run, ok := repo.runs[id]
	if !ok {
		return ErrRunNotFound
	}
	run.Status = status
	run.Report = report
	run.UpdatedAt = time.Now().UTC()
	repo.runs[id] = run
	return nil
}

func (repo *MemoryRepository) Close() error {
	return nil
}
