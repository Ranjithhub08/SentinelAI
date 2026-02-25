package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjithkumar/sentinelai/internal/handler"
	"github.com/ranjithkumar/sentinelai/internal/middleware"
	"go.uber.org/zap"
)

// SetupRouter configures the HTTP router and registers routes
func SetupRouter(logger *zap.Logger, container *Container) *gin.Engine {
	// Set gin to release mode by default, standard for production skeletons
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(middleware.Logger(logger), gin.Recovery())

	// Handlers
	healthHandler := handler.NewHealthHandler()

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
		// Add future routes here, injecting services from container
		// e.g., authHandler := handler.NewAuthHandler(container.Service)
		// v1.POST("/login", authHandler.Login)
	}

	return r
}
