package utils

import "os"

// GetEnvOrDefault returns env var value or default
func GetEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

