package monitor

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitorResponse is the DTO used to shape the API response
type MonitorResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	URL           string    `json:"url"`
	Interval      int64     `json:"interval"`
	LastChecked   time.Time `json:"last_checked"`
	StatusCode    int       `json:"status_code"`
	ResponseTime  int64     `json:"response_time"`
	IsHealthy     bool      `json:"is_healthy"`
	IsRunning     bool      `json:"is_running"`
	AIExplanation string    `json:"ai_explanation,omitempty"`
}

func mapToResponse(m *Monitor) MonitorResponse {
	return MonitorResponse{
		ID:            m.ID,
		UserID:        m.UserID,
		URL:           m.URL,
		Interval:      int64(m.Interval),
		LastChecked:   m.LastChecked,
		StatusCode:    m.StatusCode,
		ResponseTime:  m.ResponseTime.Milliseconds(),
		IsHealthy:     m.IsHealthy,
		IsRunning:     m.IsRunning,
		AIExplanation: m.AIExplanation,
	}
}

// Handler processes HTTP monitoring actions
type Handler struct {
	svc Service
}

// NewHandler generates a dependency-resolved Handler
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Add handles POST payloads to register a new URL for interval checking
func (h *Handler) Add(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized", "data": nil})
		return
	}

	var req AddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request data", "data": nil})
		return
	}

	m, err := h.svc.Add(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to add monitor", "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "monitor added successfully",
		"data":    mapToResponse(m),
	})
}

// List handles returning active user monitors and their recent states
func (h *Handler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized", "data": nil})
		return
	}

	monitors, err := h.svc.List(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list monitors", "data": nil})
		return
	}

	var responseData []MonitorResponse
	for _, m := range monitors {
		responseData = append(responseData, mapToResponse(m))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "monitors retrieved",
		"data":    responseData,
	})
}
