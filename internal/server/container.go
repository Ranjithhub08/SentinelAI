package server

import (
	"fmt"

	"github.com/ranjithkumar/sentinelai/internal/auth"
	"github.com/ranjithkumar/sentinelai/internal/monitor"
	"github.com/ranjithkumar/sentinelai/internal/repository"
	"github.com/ranjithkumar/sentinelai/internal/service"
	"github.com/ranjithkumar/sentinelai/pkg/config"
)

// Container holds all application dependencies
type Container struct {
	Repository  repository.Repository
	Service     service.Service
	AuthRepo    auth.Repository
	AuthSvc     auth.Service
	MonitorRepo monitor.Repository
	MonitorSvc  monitor.Service
}

// NewContainer initializes and wires dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	repo := repository.New()
	svc := service.New(repo)

	authRepo := auth.NewRepository()
	authSvc := auth.NewService(authRepo)

	var monitorRepo monitor.Repository
	var err error

	if cfg.DBHost != "" {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		monitorRepo, err = monitor.NewPostgresRepository(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to init postgres repo: %w", err)
		}
	} else {
		monitorRepo = monitor.NewRepository()
	}
	monitorSvc := monitor.NewService(monitorRepo)

	return &Container{
		Repository:  repo,
		Service:     svc,
		AuthRepo:    authRepo,
		AuthSvc:     authSvc,
		MonitorRepo: monitorRepo,
		MonitorSvc:  monitorSvc,
	}, nil
}
