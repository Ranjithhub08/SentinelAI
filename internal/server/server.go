package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ranjithkumar/sentinelai/pkg/config"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// New creates a new HTTP server instance
func New(cfg *config.Config, logger *zap.Logger, container *Container) *Server {
	router := SetupRouter(cfg, logger, container)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

// Start runs the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
