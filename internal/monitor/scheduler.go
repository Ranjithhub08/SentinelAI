package monitor

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Scheduler runs periodic health checks against designated URLs
type Scheduler struct {
	repo       Repository
	workerPool *WorkerPool
	logger     *zap.Logger
	interval   int
}

// NewScheduler creates a new monitor scheduler
func NewScheduler(repo Repository, workerPool *WorkerPool, logger *zap.Logger, interval int) *Scheduler {
	return &Scheduler{
		repo:       repo,
		workerPool: workerPool,
		logger:     logger,
		interval:   interval,
	}
}

// Start begins ticking processing routines that loop over monitors to queue due jobs
func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("Starting monitor scheduler", zap.Int("interval_seconds", s.interval))
	ticker := time.NewTicker(time.Duration(s.interval) * time.Second)

	go func() {
		defer ticker.Stop()
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Scheduler panic recovered", zap.Any("panic", r))
			}
		}()
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Stopping monitor scheduler")
				return
			case <-ticker.C:
				s.queueDueMonitors(ctx)
			}
		}
	}()
}

func (s *Scheduler) queueDueMonitors(ctx context.Context) {
	monitors, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get monitors for scheduling", zap.Error(err))
		return
	}

	now := time.Now()
	for _, m := range monitors {
		if !m.IsRunning && now.Sub(m.LastChecked) >= m.Interval {
			if err := s.repo.SetRunning(ctx, m.ID, true); err == nil {
				s.workerPool.Submit(Job{Monitor: m})
			}
		}
	}
}
