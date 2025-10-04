package middleware

import (
	"html"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// SecurityHeaders middleware adds security headers to responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Security headers
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			return next(c)
		}
	}
}

// ConfigurableCORS middleware with configurable allowed origins
func ConfigurableCORS(allowedOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				c.Response().Header().Set("Access-Control-Max-Age", "3600")
			}

			// Handle preflight requests
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}

// Input sanitization functions
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	// HTML escape
	input = html.EscapeString(input)
	// Trim whitespace
	input = strings.TrimSpace(input)
	return input
}

// ValidateInput validates input length and basic security
func ValidateInput(input string, maxLength int) bool {
	if len(input) > maxLength {
		return false
	}
	// Check for potentially dangerous patterns
	dangerous := []string{"<script", "javascript:", "vbscript:", "onload=", "onerror="}
	inputLower := strings.ToLower(input)
	for _, pattern := range dangerous {
		if strings.Contains(inputLower, pattern) {
			return false
		}
	}
	return true
}

// SQL injection prevention helper (use with prepared statements)
func ValidateSQLInput(input string, maxLength int) bool {
	if len(input) > maxLength {
		return false
	}
	// Basic SQL injection pattern detection
	sqlPatterns := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_", "exec", "union", "select", "drop", "delete", "insert", "update"}
	inputLower := strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(inputLower, pattern) {
			return false
		}
	}
	return true
}
