package monitor

import (
	"context"
	"time"
)

// AddReq defines the payload for adding a new monitor
type AddReq struct {
	URL      string `json:"url" binding:"required,url"`
	Interval int    `json:"interval" binding:"required,min=10"` // in seconds
}

// Service defines business logic for monitors
type Service interface {
	Add(ctx context.Context, userID string, req AddReq) (*Monitor, error)
	List(ctx context.Context, userID string) ([]*Monitor, error)
}

type serviceImpl struct {
	repo Repository
}

// NewService creates a new monitor service
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) Add(ctx context.Context, userID string, req AddReq) (*Monitor, error) {
	m := &Monitor{
		ID:        generateID(),
		UserID:    userID,
		URL:       req.URL,
		Interval:  time.Duration(req.Interval) * time.Second,
		IsHealthy: false,
	}

	if err := s.repo.Add(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *serviceImpl) List(ctx context.Context, userID string) ([]*Monitor, error) {
	return s.repo.List(ctx, userID)
}

func generateID() string {
	return time.Now().Format("20060102150405000") // simple mock ID generator
}
