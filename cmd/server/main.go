package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ranjithkumar/sentinelai/internal/llm"
	"github.com/ranjithkumar/sentinelai/internal/monitor"
	"github.com/ranjithkumar/sentinelai/internal/server"
	"github.com/ranjithkumar/sentinelai/pkg/config"
	"github.com/ranjithkumar/sentinelai/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	zlog, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		_ = zlog.Sync()
	}()

	container, err := server.NewContainer(cfg)
	if err != nil {
		zlog.Fatal("Failed to initialize dependency container", zap.Error(err))
	}

	engineCtx, engineCancel := context.WithCancel(context.Background())
	defer engineCancel()

	llmProvider := llm.NewOllamaProvider(cfg.OllamaURL, cfg.LLMModel)

	workerPool := monitor.NewWorkerPool(10, container.MonitorRepo, zlog, llmProvider)
	workerPool.Start(engineCtx)

	scheduler := monitor.NewScheduler(container.MonitorRepo, workerPool, zlog, cfg.SchedulerInterval)
	scheduler.Start(engineCtx)

	srv := server.New(cfg, zlog, container)

	go func() {
		if err := srv.Start(); err != nil {
			zlog.Fatal("Server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	zlog.Info("Shutdown signal received")
	engineCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		zlog.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zlog.Info("Server stopped cleanly")
}
