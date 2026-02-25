package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port              int
	Env               string
	JwtSecret         string
	JwtExpiration     int
	SchedulerInterval int
	OllamaURL         string
	LLMModel          string
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		return nil, errors.New("PORT environment variable is required")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("PORT must be a valid integer")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	jwtExpStr := os.Getenv("TOKEN_EXPIRATION")
	if jwtExpStr == "" {
		return nil, errors.New("TOKEN_EXPIRATION environment variable is required")
	}
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		return nil, errors.New("TOKEN_EXPIRATION must be a valid integer")
	}

	schedulerInterval := 1
	if sIntStr := os.Getenv("SCHEDULER_INTERVAL"); sIntStr != "" {
		if parsed, err := strconv.Atoi(sIntStr); err == nil && parsed > 0 {
			schedulerInterval = parsed
		}
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434/api/generate"
	}

	llmModel := os.Getenv("LLM_MODEL")
	if llmModel == "" {
		llmModel = "llama3"
	}

	return &Config{
		Port:              port,
		Env:               env,
		JwtSecret:         jwtSecret,
		JwtExpiration:     jwtExp,
		SchedulerInterval: schedulerInterval,
		OllamaURL:         ollamaURL,
		LLMModel:          llmModel,
	}, nil
}
