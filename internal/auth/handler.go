package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranjithkumar/sentinelai/pkg/config"
)

// Handler processes HTTP auth endpoints
type Handler struct {
	svc Service
	cfg *config.Config
}

// NewHandler creates a new auth handler
func NewHandler(svc Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request data", "data": nil})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user already exists", "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to register user", "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "registration successful",
		"data":    user,
	})
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request data", "data": nil})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req, h.cfg.JwtSecret, h.cfg.JwtExpiration)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid credentials", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
		"data":    gin.H{"token": token},
	})
}
