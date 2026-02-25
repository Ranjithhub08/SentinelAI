package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterReq defines registration payload
type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginReq defines login payload
type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Service defines auth business logic
type Service interface {
	Register(ctx context.Context, req RegisterReq) (*User, error)
	Login(ctx context.Context, req LoginReq, secret string, expHours int) (string, error)
}

type serviceImpl struct {
	repo Repository
}

// NewService creates a new auth service
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) Register(ctx context.Context, req RegisterReq) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           generateID(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *serviceImpl) Login(ctx context.Context, req LoginReq, secret string, expHours int) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return GenerateToken(user.ID, secret, expHours)
}

func generateID() string {
	return time.Now().Format("20060102150405")
}
