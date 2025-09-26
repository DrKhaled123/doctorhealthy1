package utils

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regex patterns for better performance
var (
	controlCharsRegex = regexp.MustCompile(`[\x00-\x1f\x7f]`)
)

// SanitizeForLog sanitizes input for logging to prevent log injection
func SanitizeForLog(input string) string {
	// Remove newlines and carriage returns
	sanitized := strings.ReplaceAll(input, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")

	// Remove control characters
	sanitized = controlCharsRegex.ReplaceAllString(sanitized, "")

	// Limit length
	if len(sanitized) > 200 {
		sanitized = sanitized[:200] + "..."
	}

	return sanitized
}

// SafeKeyPrefix returns a bounded, sanitized prefix of a sensitive key for logs
// Example: SafeKeyPrefix("abcdef", 8) -> "abcdef"; SafeKeyPrefix("abcdefghijk", 8) -> "abcdefgh"
func SafeKeyPrefix(key string, prefixLen int) string {
	if prefixLen <= 0 {
		return ""
	}
	if len(key) > prefixLen {
		key = key[:prefixLen]
	}
	return SanitizeForLog(key)
}

// SanitizeForHTML sanitizes input for HTML output to prevent XSS
func SanitizeForHTML(input string) string {
	return html.EscapeString(input)
}

// SanitizeForJSON sanitizes input for JSON output
func SanitizeForJSON(input string) string {
	// Escape backslashes first, then quotes to prevent double-escaping
	sanitized := strings.ReplaceAll(input, "\\", "\\\\")
	sanitized = strings.ReplaceAll(sanitized, "\"", "\\\"")

	// Remove control characters
	sanitized = controlCharsRegex.ReplaceAllString(sanitized, "")

	return sanitized
}

// SanitizeInput sanitizes general input by removing dangerous characters
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove control characters
	input = controlCharsRegex.ReplaceAllString(input, "")

	// Truncate at first HTML tag to remove scripts/injections quickly
	if idx := strings.Index(input, "<"); idx >= 0 {
		input = input[:idx]
	}

	// Remove potentially dangerous characters for SQL/NoSQL injection
	dangerous := []string{"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror="}
	for _, danger := range dangerous {
		input = strings.ReplaceAll(strings.ToLower(input), danger, "")
	}

	return input
}

// ParsePositiveInt parses a string to positive integer
func ParsePositiveInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if val <= 0 {
		return 0, strconv.ErrRange
	}
	return val, nil
}

// ParseNonNegativeInt parses a string to non-negative integer
func ParseNonNegativeInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, strconv.ErrRange
	}
	return val, nil
}
