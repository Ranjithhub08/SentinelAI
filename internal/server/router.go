package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjithkumar/sentinelai/internal/auth"
	"github.com/ranjithkumar/sentinelai/internal/handler"
	"github.com/ranjithkumar/sentinelai/internal/middleware"
	"github.com/ranjithkumar/sentinelai/internal/monitor"
	"github.com/ranjithkumar/sentinelai/pkg/config"
	"go.uber.org/zap"
)

// SetupRouter configures the HTTP router and registers routes
func SetupRouter(cfg *config.Config, logger *zap.Logger, container *Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.Logger(logger), gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	authHandler := auth.NewHandler(container.AuthSvc, cfg)
	monitorHandler := monitor.NewHandler(container.MonitorSvc)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		monitorGroup := v1.Group("/monitor")
		monitorGroup.Use(auth.Middleware(cfg.JwtSecret))
		{
			monitorGroup.POST("/add", monitorHandler.Add)
			monitorGroup.GET("/list", monitorHandler.List)
		}
	}

	return r
}
