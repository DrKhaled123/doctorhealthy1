package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	APIKey   APIKeyConfig
	Security SecurityConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string
}

// APIKeyConfig holds API key configuration
type APIKeyConfig struct {
	Length         int
	ExpiryDuration time.Duration
	Prefix         string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "data/apikeys.db"),
		},
		APIKey: APIKeyConfig{
			Length:         getEnvAsInt("API_KEY_LENGTH", 32),
			ExpiryDuration: getEnvAsDuration("API_KEY_EXPIRY", "365d"),
			Prefix:         getEnv("API_KEY_PREFIX", "ak_"),
		},
		Security: SecurityConfig{
			RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),
		},
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Printf("Warning: failed to parse %s=%s as int, using default %d: %v", key, value, defaultValue, err)
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := parseDuration(value); err == nil {
			return duration
		} else {
			log.Printf("Warning: failed to parse %s=%s as duration, using default %s: %v", key, value, defaultValue, err)
		}
	}
	duration, _ := parseDuration(defaultValue)
	return duration
}

// parseDuration parses duration strings like "365d", "24h", "30m"
func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return time.ParseDuration(s)
	}

	unit := s[len(s)-1:]
	value := s[:len(s)-1]

	switch unit {
	case "d":
		if days, err := strconv.Atoi(value); err == nil {
			return time.Duration(days) * 24 * time.Hour, nil
		}
	case "w":
		if weeks, err := strconv.Atoi(value); err == nil {
			return time.Duration(weeks) * 7 * 24 * time.Hour, nil
		}
	}

	return time.ParseDuration(s)
}
