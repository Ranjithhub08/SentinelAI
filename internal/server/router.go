package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjithkumar/sentinelai/internal/handler"
	"github.com/ranjithkumar/sentinelai/internal/middleware"
	"go.uber.org/zap"
)

// SetupRouter configures the HTTP router and registers routes
func SetupRouter(logger *zap.Logger, container *Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.Logger(logger), gin.Recovery())

	healthHandler := handler.NewHealthHandler()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
	}

	return r
}
