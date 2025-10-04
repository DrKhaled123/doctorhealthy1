package utils

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"api-key-generator/internal/models"
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EnableCSP           bool
	EnableXSSProtection bool
	EnableHSTS          bool
	EnableCSRF          bool
	AllowedOrigins      []string
	BlockedPatterns     []string
	SanitizeInput       bool
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		EnableCSP:           true,
		EnableXSSProtection: true,
		EnableHSTS:          true,
		EnableCSRF:          true,
		AllowedOrigins:      []string{"https://app.example.com", "https://www.example.com"},
		BlockedPatterns: []string{
			`<script[^>]*>.*?</script>`,
			`javascript:`,
			`vbscript:`,
			`onload=`,
			`onerror=`,
			`onclick=`,
			`onmouseover=`,
		},
		SanitizeInput: true,
	}
}

// InputSanitizer provides comprehensive input sanitization
type InputSanitizer struct {
	config *SecurityConfig
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer(config *SecurityConfig) *InputSanitizer {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	return &InputSanitizer{config: config}
}

// SanitizeString sanitizes a string input
func (is *InputSanitizer) SanitizeString(input string) string {
	if !is.config.SanitizeInput {
		return input
	}

	// HTML escape
	sanitized := html.EscapeString(input)

	// Remove potentially dangerous patterns
	for _, pattern := range is.config.BlockedPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		sanitized = re.ReplaceAllString(sanitized, "")
	}

	return strings.TrimSpace(sanitized)
}

// SanitizeURL sanitizes and validates URLs
func (is *InputSanitizer) SanitizeURL(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	// Parse URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Validate scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}

	// Check against allowed origins
	if !is.isAllowedOrigin(parsedURL.Host) {
		return "", fmt.Errorf("origin not allowed: %s", parsedURL.Host)
	}

	return parsedURL.String(), nil
}

// SanitizeEmail sanitizes and validates email addresses
func (is *InputSanitizer) SanitizeEmail(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input) {
		return "", fmt.Errorf("invalid email format")
	}

	// Check length
	if len(input) > 254 {
		return "", fmt.Errorf("email too long")
	}

	return strings.ToLower(strings.TrimSpace(input)), nil
}

// SanitizeJSON sanitizes JSON input
func (is *InputSanitizer) SanitizeJSON(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	// Remove potential JSON injection
	sanitized := strings.ReplaceAll(input, "<", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	sanitized = strings.ReplaceAll(sanitized, "&", "")

	// Validate JSON structure
	if !is.isValidJSON(sanitized) {
		return "", fmt.Errorf("invalid JSON structure")
	}

	return sanitized, nil
}

// SanitizeForJSON sanitizes input for JSON output
func SanitizeForJSON(input string) string {
	// Escape backslashes first, then quotes to prevent double-escaping
	sanitized := strings.ReplaceAll(input, "\\", "\\\\")
	sanitized = strings.ReplaceAll(sanitized, "\"", "\\\"")

	// Remove control characters
	controlCharsRegex := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	sanitized = controlCharsRegex.ReplaceAllString(sanitized, "")

	return sanitized
}

// SanitizeForLog sanitizes input for logging to prevent log injection
func SanitizeForLog(input string) string {
	// Remove newlines and carriage returns
	sanitized := strings.ReplaceAll(input, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")

	// Remove control characters
	controlCharsRegex := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	sanitized = controlCharsRegex.ReplaceAllString(sanitized, "")

	// Limit length
	if len(sanitized) > 200 {
		sanitized = sanitized[:200] + "..."
	}

	return sanitized
}

// SanitizeForHTML sanitizes input for HTML output to prevent XSS
func SanitizeForHTML(input string) string {
	// Remove JavaScript protocols
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsProtocolRegex.ReplaceAllString(input, "")

	// Remove event handlers
	eventHandlerRegex := regexp.MustCompile(`(?i)on\w+\s*=`)
	input = eventHandlerRegex.ReplaceAllString(input, "")

	// Remove script tags
	scriptTagRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptTagRegex.ReplaceAllString(input, "")

	// Remove Unicode XSS patterns
	unicodeXssRegex := regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)
	input = unicodeXssRegex.ReplaceAllString(input, "")

	// Escape HTML entities
	return html.EscapeString(input)
}

// isAllowedOrigin checks if an origin is in the allowed list
func (is *InputSanitizer) isAllowedOrigin(origin string) bool {
	for _, allowed := range is.config.AllowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// isValidJSON performs basic JSON validation
func (is *InputSanitizer) isValidJSON(input string) bool {
	// Simple check for balanced braces and brackets
	braceCount := 0
	bracketCount := 0

	for _, char := range input {
		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		}
	}

	return braceCount == 0 && bracketCount == 0
}

// CSPHeaders generates Content Security Policy headers
type CSPHeaders struct {
	DefaultSrc string
	ScriptSrc  string
	StyleSrc   string
	ImgSrc     string
	ConnectSrc string
	FontSrc    string
	ObjectSrc  string
	MediaSrc   string
	FrameSrc   string
}

// DefaultCSPHeaders returns secure default CSP headers
func DefaultCSPHeaders() *CSPHeaders {
	return &CSPHeaders{
		DefaultSrc: "'self'",
		ScriptSrc:  "'self' 'unsafe-inline'",
		StyleSrc:   "'self' 'unsafe-inline' https://fonts.googleapis.com",
		ImgSrc:     "'self' data: https:",
		ConnectSrc: "'self' https://api.example.com",
		FontSrc:    "'self' https://fonts.gstatic.com",
		ObjectSrc:  "'none'",
		MediaSrc:   "'self'",
		FrameSrc:   "'none'",
	}
}

// GenerateCSPHeader generates a complete CSP header string
func (csp *CSPHeaders) GenerateCSPHeader() string {
	var directives []string

	if csp.DefaultSrc != "" {
		directives = append(directives, fmt.Sprintf("default-src %s", csp.DefaultSrc))
	}
	if csp.ScriptSrc != "" {
		directives = append(directives, fmt.Sprintf("script-src %s", csp.ScriptSrc))
	}
	if csp.StyleSrc != "" {
		directives = append(directives, fmt.Sprintf("style-src %s", csp.StyleSrc))
	}
	if csp.ImgSrc != "" {
		directives = append(directives, fmt.Sprintf("img-src %s", csp.ImgSrc))
	}
	if csp.ConnectSrc != "" {
		directives = append(directives, fmt.Sprintf("connect-src %s", csp.ConnectSrc))
	}
	if csp.FontSrc != "" {
		directives = append(directives, fmt.Sprintf("font-src %s", csp.FontSrc))
	}
	if csp.ObjectSrc != "" {
		directives = append(directives, fmt.Sprintf("object-src %s", csp.ObjectSrc))
	}
	if csp.MediaSrc != "" {
		directives = append(directives, fmt.Sprintf("media-src %s", csp.MediaSrc))
	}
	if csp.FrameSrc != "" {
		directives = append(directives, fmt.Sprintf("frame-src %s", csp.FrameSrc))
	}

	return strings.Join(directives, "; ")
}

// SecurityValidator validates security aspects of requests and responses
type SecurityValidator struct {
	sanitizer *InputSanitizer
	csp       *CSPHeaders
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		sanitizer: NewInputSanitizer(DefaultSecurityConfig()),
		csp:       DefaultCSPHeaders(),
	}
}

// ValidateInput validates and sanitizes user input
func (sv *SecurityValidator) ValidateInput(input string, inputType string) (string, []models.ValidationError) {
	var errors []models.ValidationError
	var sanitized string

	switch inputType {
	case "string":
		sanitized = sv.sanitizer.SanitizeString(input)
	case "url":
		var err error
		sanitized, err = sv.sanitizer.SanitizeURL(input)
		if err != nil {
			errors = append(errors, models.ValidationError{
				Field:   inputType,
				Value:   input,
				Message: err.Error(),
			})
		}
	case "email":
		var err error
		sanitized, err = sv.sanitizer.SanitizeEmail(input)
		if err != nil {
			errors = append(errors, models.ValidationError{
				Field:   inputType,
				Value:   input,
				Message: err.Error(),
			})
		}
	case "json":
		var err error
		sanitized, err = sv.sanitizer.SanitizeJSON(input)
		if err != nil {
			errors = append(errors, models.ValidationError{
				Field:   inputType,
				Value:   input,
				Message: err.Error(),
			})
		}
	default:
		sanitized = sv.sanitizer.SanitizeString(input)
	}

	// Check for suspicious patterns
	if sv.containsSuspiciousPatterns(sanitized) {
		errors = append(errors, models.ValidationError{
			Field:   inputType,
			Value:   input,
			Message: "Input contains suspicious patterns",
		})
	}

	return sanitized, errors
}

// containsSuspiciousPatterns checks for potentially dangerous patterns
func (sv *SecurityValidator) containsSuspiciousPatterns(input string) bool {
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"eval(",
		"document.cookie",
		"window.location",
		"innerHTML",
		"outerHTML",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// ValidatePassword validates password strength
func (sv *SecurityValidator) ValidatePassword(password string) []models.ValidationError {
	var errors []models.ValidationError

	if len(password) < 8 {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		})
	}

	if len(password) > 128 {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must be less than 128 characters",
		})
	}

	// Check for required character types
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
		})
	}

	if !hasUpper {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
		})
	}

	if !hasDigit {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one digit",
		})
	}

	if !hasSpecial {
		errors = append(errors, models.ValidationError{
			Field:   "password",
			Message: "Password must contain at least one special character",
		})
	}

	return errors
}

// GenerateSecurityHeaders generates all security headers
func (sv *SecurityValidator) GenerateSecurityHeaders() map[string]string {
	headers := make(map[string]string)

	if sv.sanitizer.config.EnableCSP {
		headers["Content-Security-Policy"] = sv.csp.GenerateCSPHeader()
	}

	if sv.sanitizer.config.EnableXSSProtection {
		headers["X-XSS-Protection"] = "1; mode=block"
	}

	if sv.sanitizer.config.EnableHSTS {
		headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	}

	// Additional security headers
	headers["X-Content-Type-Options"] = "nosniff"
	headers["X-Frame-Options"] = "DENY"
	headers["Referrer-Policy"] = "strict-origin-when-cross-origin"
	headers["Permissions-Policy"] = "geolocation=(), microphone=(), camera=()"

	return headers
}

// SecurityScanner scans for common vulnerabilities
type SecurityScanner struct {
	Patterns []VulnerabilityPattern
}

// VulnerabilityPattern represents a vulnerability pattern
type VulnerabilityPattern struct {
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Severity    string
	Solution    string
}

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner() *SecurityScanner {
	scanner := &SecurityScanner{
		Patterns: make([]VulnerabilityPattern, 0),
	}

	// Add common vulnerability patterns
	scanner.AddPattern(VulnerabilityPattern{
		Name:        "SQL Injection",
		Description: "Potential SQL injection vulnerability",
		Pattern:     regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter).*from`),
		Severity:    "high",
		Solution:    "Use parameterized queries or prepared statements",
	})

	scanner.AddPattern(VulnerabilityPattern{
		Name:        "XSS Vulnerability",
		Description: "Potential cross-site scripting vulnerability",
		Pattern:     regexp.MustCompile(`(?i)(<script|javascript:|on\w+\s*=)`),
		Severity:    "high",
		Solution:    "Sanitize input and use CSP headers",
	})

	scanner.AddPattern(VulnerabilityPattern{
		Name:        "Path Traversal",
		Description: "Potential path traversal vulnerability",
		Pattern:     regexp.MustCompile(`(?i)(\.\.\\|\.\./|%2e%2e%2f|%2e%2e/)`),
		Severity:    "high",
		Solution:    "Validate and sanitize file paths",
	})

	scanner.AddPattern(VulnerabilityPattern{
		Name:        "Command Injection",
		Description: "Potential command injection vulnerability",
		Pattern:     regexp.MustCompile(`(?i)(;\s*|\|\||&&)\s*\w+`),
		Severity:    "critical",
		Solution:    "Avoid shell execution or use allowlists",
	})

	return scanner
}

// AddPattern adds a vulnerability pattern to scan for
func (ss *SecurityScanner) AddPattern(pattern VulnerabilityPattern) {
	ss.Patterns = append(ss.Patterns, pattern)
}

// Scan scans input for vulnerabilities
func (ss *SecurityScanner) Scan(input string) []error {
	var errors []error

	for _, pattern := range ss.Patterns {
		if pattern.Pattern.MatchString(input) {
			errors = append(errors, models.NewSecurityError(
				fmt.Sprintf("Detected %s: %s", pattern.Name, pattern.Description),
				pattern.Name,
				pattern.Severity,
			))
		}
	}

	return errors
}

// Global security instances
var GlobalSecurityValidator = NewSecurityValidator()
var GlobalSecurityScanner = NewSecurityScanner()
var GlobalInputSanitizer = NewInputSanitizer(DefaultSecurityConfig())

// Convenience functions for global security

// SanitizeInputGlobal sanitizes input using global sanitizer
func SanitizeInputGlobal(input, inputType string) (string, []models.ValidationError) {
	return GlobalSecurityValidator.ValidateInput(input, inputType)
}

// ScanForVulnerabilitiesGlobal scans input for vulnerabilities
func ScanForVulnerabilitiesGlobal(input string) []error {
	return GlobalSecurityScanner.Scan(input)
}

// GenerateSecurityHeaders generates security headers
func GenerateSecurityHeaders() map[string]string {
	return GlobalSecurityValidator.GenerateSecurityHeaders()
}

// ValidatePassword validates password strength
func ValidatePassword(password string) []models.ValidationError {
	return GlobalSecurityValidator.ValidatePassword(password)
}

// SanitizeInput sanitizes general input by removing dangerous characters
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove control characters
	controlCharsRegex := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	input = controlCharsRegex.ReplaceAllString(input, "")

	// Remove Unicode XSS patterns
	unicodeXssRegex := regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)
	input = unicodeXssRegex.ReplaceAllString(input, "")

	// Remove JavaScript protocols
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsProtocolRegex.ReplaceAllString(input, "")

	// Remove event handlers
	eventHandlerRegex := regexp.MustCompile(`(?i)on\w+\s*=`)
	input = eventHandlerRegex.ReplaceAllString(input, "")

	// Remove script tags
	scriptTagRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptTagRegex.ReplaceAllString(input, "")

	// Remove potentially dangerous characters for SQL/NoSQL injection
	dangerous := []string{"vbscript:", "eval(", "document.cookie", "window.location", "innerHTML", "outerHTML"}
	for _, danger := range dangerous {
		input = strings.ReplaceAll(strings.ToLower(input), strings.ToLower(danger), "")
	}

	return input
}
