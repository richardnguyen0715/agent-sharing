// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds all application configuration values.
type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	CORSOrigins   []string
	LogLevel      string
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if any required variable is missing or invalid.
func Load() (*Config, error) {
	port := getEnv("PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}

	corsOrigins := strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		Port:          port,
		DatabaseURL:   databaseURL,
		RedisURL:      redisURL,
		JWTSecret:     jwtSecret,
		JWTAccessTTL:  accessTTL,
		JWTRefreshTTL: refreshTTL,
		CORSOrigins:   corsOrigins,
		LogLevel:      logLevel,
	}, nil
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
