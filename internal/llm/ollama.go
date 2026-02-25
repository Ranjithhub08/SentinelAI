package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ollamaReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaRes struct {
	Response string `json:"response"`
}

type ollamaProvider struct {
	client *http.Client
	url    string
	model  string
}

// NewOllamaProvider creates a new local Ollama LLM provider
func NewOllamaProvider(url, model string) Provider {
	return &ollamaProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		url:    url,
		model:  model,
	}
}

func (p *ollamaProvider) AnalyzeFailure(ctx context.Context, input FailureInput) (string, error) {
	prompt := fmt.Sprintf(
		"Analyze this monitoring failure. URL: %s, Status Code: %d, Response Time: %s, Timestamp: %s. Briefly explain what might have gone wrong.",
		input.URL, input.StatusCode, input.ResponseTime.String(), input.Timestamp.Format(time.RFC3339),
	)

	reqBody := ollamaReq{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status: %d", res.StatusCode)
	}

	var parsedRes ollamaRes
	if err := json.NewDecoder(res.Body).Decode(&parsedRes); err != nil {
		return "", err
	}

	return parsedRes.Response, nil
}
