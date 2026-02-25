package server

import (
	"github.com/ranjithkumar/sentinelai/internal/auth"
	"github.com/ranjithkumar/sentinelai/internal/repository"
	"github.com/ranjithkumar/sentinelai/internal/service"
)

// Container holds all application dependencies
type Container struct {
	Repository repository.Repository
	Service    service.Service
	AuthRepo   auth.Repository
	AuthSvc    auth.Service
}

// NewContainer initializes and wires dependencies
func NewContainer() (*Container, error) {
	repo := repository.New()
	svc := service.New(repo)

	authRepo := auth.NewRepository()
	authSvc := auth.NewService(authRepo)

	return &Container{
		Repository: repo,
		Service:    svc,
		AuthRepo:   authRepo,
		AuthSvc:    authSvc,
	}, nil
}
