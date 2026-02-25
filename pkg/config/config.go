package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port          int
	Env           string
	JwtSecret     string
	JwtExpiration int
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-for-local-dev-only"
	}

	jwtExp := 24
	if e := os.Getenv("JWT_EXPIRATION_HOURS"); e != "" {
		if parsed, err := strconv.Atoi(e); err == nil {
			jwtExp = parsed
		}
	}

	return &Config{
		Port:          port,
		Env:           env,
		JwtSecret:     jwtSecret,
		JwtExpiration: jwtExp,
	}, nil
}
