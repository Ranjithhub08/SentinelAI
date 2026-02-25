package monitor

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Repository defines data access for monitors
type Repository interface {
	Add(ctx context.Context, m *Monitor) error
	List(ctx context.Context, userID string) ([]*Monitor, error)
	GetAll(ctx context.Context) ([]*Monitor, error)
	UpdateStatus(ctx context.Context, id string, lastChecked time.Time, statusCode int, responseTime time.Duration, isHealthy bool) error
	SetRunning(ctx context.Context, id string, isRunning bool) error
}

type inMemoryRepository struct {
	mu       sync.RWMutex
	monitors map[string]*Monitor
}

// NewRepository creates a new in-memory monitor repository
func NewRepository() Repository {
	return &inMemoryRepository{
		monitors: make(map[string]*Monitor),
	}
}

func cloneMonitor(m *Monitor) *Monitor {
	if m == nil {
		return nil
	}
	clone := *m
	return &clone
}

func (r *inMemoryRepository) Add(ctx context.Context, m *Monitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.monitors[m.ID]; exists {
		return errors.New("monitor already exists")
	}

	r.monitors[m.ID] = cloneMonitor(m)
	return nil
}

func (r *inMemoryRepository) List(ctx context.Context, userID string) ([]*Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Monitor
	for _, m := range r.monitors {
		if m.UserID == userID {
			result = append(result, cloneMonitor(m))
		}
	}
	return result, nil
}

func (r *inMemoryRepository) GetAll(ctx context.Context) ([]*Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Monitor
	for _, m := range r.monitors {
		result = append(result, cloneMonitor(m))
	}
	return result, nil
}

func (r *inMemoryRepository) UpdateStatus(ctx context.Context, id string, lastChecked time.Time, statusCode int, responseTime time.Duration, isHealthy bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, exists := r.monitors[id]
	if !exists {
		return errors.New("monitor not found")
	}

	m.LastChecked = lastChecked
	m.StatusCode = statusCode
	m.ResponseTime = responseTime
	m.IsHealthy = isHealthy

	return nil
}

func (r *inMemoryRepository) SetRunning(ctx context.Context, id string, isRunning bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, exists := r.monitors[id]
	if !exists {
		return errors.New("monitor not found")
	}

	m.IsRunning = isRunning
	return nil
}
