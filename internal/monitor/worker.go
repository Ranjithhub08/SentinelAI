package monitor

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Job represents a single health check execution
type Job struct {
	Monitor *Monitor
}

// WorkerPool manages concurrent health checks
type WorkerPool struct {
	numWorkers int
	jobChan    chan Job
	repo       Repository
	logger     *zap.Logger
}

// NewWorkerPool creates a new monitor worker pool
func NewWorkerPool(numWorkers int, repo Repository, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobChan:    make(chan Job, 1000), // Buffer jobs
		repo:       repo,
		logger:     logger,
	}
}

// Start spawns the configured number of workers
func (wp *WorkerPool) Start(ctx context.Context) {
	wp.logger.Info("Starting monitor worker pool", zap.Int("workers", wp.numWorkers))
	for i := 0; i < wp.numWorkers; i++ {
		go wp.worker(ctx)
	}
}

// Submit queues a health check job
func (wp *WorkerPool) Submit(job Job) {
	select {
	case wp.jobChan <- job:
	default:
		wp.logger.Warn("Worker pool job channel is full, dropping health check job", zap.String("monitor_id", job.Monitor.ID))
	}
}

func (wp *WorkerPool) worker(ctx context.Context) {
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-wp.jobChan:
			wp.safeProcessJob(ctx, client, job)
		}
	}
}

func (wp *WorkerPool) safeProcessJob(ctx context.Context, client *http.Client, job Job) {
	defer func() {
		if r := recover(); r != nil {
			wp.logger.Error("Job panic recovered", zap.Any("panic", r), zap.String("monitor_id", job.Monitor.ID))
		}
		_ = wp.repo.SetRunning(ctx, job.Monitor.ID, false)
	}()
	wp.processJob(ctx, client, job)
}

func (wp *WorkerPool) processJob(ctx context.Context, client *http.Client, job Job) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.Monitor.URL, nil)
	if err != nil {
		_ = wp.repo.UpdateStatus(ctx, job.Monitor.ID, time.Now(), 0, 0, false)
		wp.logger.Error("Failed to create request", zap.Error(err), zap.String("url", job.Monitor.URL))
		return
	}

	res, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		_ = wp.repo.UpdateStatus(ctx, job.Monitor.ID, time.Now(), 0, duration, false)
		wp.logger.Warn("Health check unreachable", zap.Error(err), zap.String("url", job.Monitor.URL))
		return
	}
	defer res.Body.Close()

	isHealthy := res.StatusCode >= 200 && res.StatusCode < 500
	_ = wp.repo.UpdateStatus(ctx, job.Monitor.ID, time.Now(), res.StatusCode, duration, isHealthy)

	wp.logger.Info("Health check executed",
		zap.String("monitor_id", job.Monitor.ID),
		zap.String("url", job.Monitor.URL),
		zap.Int("status", res.StatusCode),
		zap.Duration("latency", duration),
		zap.Bool("healthy", isHealthy),
	)
}
