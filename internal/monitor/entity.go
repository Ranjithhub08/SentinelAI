package monitor

import "time"

// Monitor represents a health check target
type Monitor struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	URL          string        `json:"url"`
	Interval     time.Duration `json:"interval"`
	LastChecked  time.Time     `json:"last_checked"`
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	IsHealthy    bool          `json:"is_healthy"`
	IsRunning    bool          `json:"is_running"`
}
