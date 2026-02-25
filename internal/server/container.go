package server

import (
	"github.com/ranjithkumar/sentinelai/internal/repository"
	"github.com/ranjithkumar/sentinelai/internal/service"
)

// Container holds all application dependencies
type Container struct {
	Repository repository.Repository
	Service    service.Service
}

// NewContainer initializes and wires dependencies
func NewContainer() (*Container, error) {
	repo := repository.New()
	svc := service.New(repo)

	return &Container{
		Repository: repo,
		Service:    svc,
	}, nil
}
