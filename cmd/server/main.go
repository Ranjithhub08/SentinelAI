package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ranjithkumar/sentinelai/internal/server"
	"github.com/ranjithkumar/sentinelai/pkg/config"
	"github.com/ranjithkumar/sentinelai/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. Configuration loader (reads from .env with validation)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Logger initialization using zap (production config)
	zlog, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		_ = zlog.Sync()
	}()

	// 7. Dependency container struct to wire dependencies
	container, err := server.NewContainer()
	if err != nil {
		zlog.Fatal("Failed to initialize dependency container", zap.Error(err))
	}

	// 3, 4, 10. HTTP server setup using Gin & Clear separation of router and server bootstrap
	srv := server.New(cfg, zlog, container)

	// Start server in a background goroutine
	go func() {
		if err := srv.Start(); err != nil {
			zlog.Fatal("Server start failed", zap.Error(err))
		}
	}()

	// 8. Graceful shutdown with runtime signal capturing
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until signal is received

	zlog.Info("Shutdown signal received")

	// Shutdown timeout (10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		zlog.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zlog.Info("Server stopped cleanly")
}
