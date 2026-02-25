package llm

import (
	"context"
	"time"
)

// FailureInput holds context about a health check failure
type FailureInput struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Timestamp    time.Time
}

// Provider defines the interface for AI-powered log/metrics analysis
type Provider interface {
	AnalyzeFailure(ctx context.Context, input FailureInput) (string, error)
}
