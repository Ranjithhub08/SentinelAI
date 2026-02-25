package service

import "github.com/ranjithkumar/sentinelai/internal/repository"

// Service defines the interface for business logic
type Service interface {
}

type serviceImpl struct {
	repo repository.Repository
}

// New returns a new Service implementation
func New(repo repository.Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}
