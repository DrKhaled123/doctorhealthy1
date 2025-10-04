package models

import (
	"fmt"
	"time"
)

// APIKey represents an API key in the system
type APIKey struct {
	ID          string     `json:"id" db:"id"`
	Key         string     `json:"key" db:"key"`
	Name        string     `json:"name" db:"name" validate:"required,min=2,max=100"`
	Description *string    `json:"description,omitempty" db:"description"`
	UserID      *string    `json:"user_id,omitempty" db:"user_id"`
	Permissions []string   `json:"permissions" db:"permissions"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Usage statistics
	UsageCount int64 `json:"usage_count" db:"usage_count"`

	// Rate limiting
	RateLimit     *int `json:"rate_limit,omitempty" db:"rate_limit"`
	RateLimitUsed int  `json:"rate_limit_used" db:"rate_limit_used"`
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
	ExpiryDays  *int     `json:"expiry_days,omitempty" validate:"omitempty,min=1,max=3650"`
	RateLimit   *int     `json:"rate_limit,omitempty" validate:"omitempty,min=1,max=10000"`
}

// UpdateAPIKeyRequest represents the request to update an API key
type UpdateAPIKeyRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Permissions []string `json:"permissions,omitempty" validate:"omitempty,min=1"`
	IsActive    *bool    `json:"is_active,omitempty"`
	RateLimit   *int     `json:"rate_limit,omitempty" validate:"omitempty,min=1,max=10000"`
}

// ListAPIKeysParams represents parameters for listing API keys
type ListAPIKeysParams struct {
	Page     int    `query:"page" validate:"min=1"`
	Limit    int    `query:"limit" validate:"min=1,max=100"`
	Search   string `query:"search"`
	IsActive *bool  `query:"is_active"`
	UserID   string `query:"user_id"`
}

// ListAPIKeysResponse represents the response for listing API keys
type ListAPIKeysResponse struct {
	APIKeys    []APIKey           `json:"api_keys"`
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// APIKeyUsage represents API key usage statistics
type APIKeyUsage struct {
	APIKeyID  string    `json:"api_key_id" db:"api_key_id"`
	Endpoint  string    `json:"endpoint" db:"endpoint"`
	Method    string    `json:"method" db:"method"`
	Status    int       `json:"status" db:"status"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
}

// APIKeyStats represents API key statistics
type APIKeyStats struct {
	TotalRequests     int64            `json:"total_requests"`
	RequestsToday     int64            `json:"requests_today"`
	RequestsThisWeek  int64            `json:"requests_this_week"`
	RequestsThisMonth int64            `json:"requests_this_month"`
	LastUsed          *time.Time       `json:"last_used"`
	TopEndpoints      []EndpointStat   `json:"top_endpoints"`
	StatusCodes       map[string]int64 `json:"status_codes"`
}

// EndpointStat represents endpoint usage statistics
type EndpointStat struct {
	Endpoint string `json:"endpoint"`
	Count    int64  `json:"count"`
}

// Permission represents available permissions
type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// GetAvailablePermissions returns available permissions (prevents modification)
func GetAvailablePermissions() []Permission {
	return []Permission{
		{Name: "read", Description: "Read access to resources", Category: "basic"},
		{Name: "write", Description: "Write access to resources", Category: "basic"},
		{Name: "delete", Description: "Delete access to resources", Category: "basic"},
		{Name: "admin", Description: "Administrative access", Category: "advanced"},
		{Name: "admin:all", Description: "Full administrative access", Category: "advanced"},
		{Name: "users:read", Description: "Read user data", Category: "users"},
		{Name: "users:write", Description: "Modify user data", Category: "users"},
		{Name: "meals:read", Description: "Read meal data", Category: "meals"},
		{Name: "meals:write", Description: "Modify meal data", Category: "meals"},
		{Name: "workouts:read", Description: "Read workout data", Category: "workouts"},
		{Name: "workouts:write", Description: "Modify workout data", Category: "workouts"},
		{Name: "health:read", Description: "Read health data", Category: "health"},
		{Name: "health:write", Description: "Modify health data", Category: "health"},
		{Name: "recipes:read", Description: "Read recipe data", Category: "recipes"},
		{Name: "recipes:write", Description: "Modify recipe data", Category: "recipes"},
		{Name: "nutrition:generate", Description: "Generate nutrition plans", Category: "generation"},
		{Name: "workout:generate", Description: "Generate workout plans", Category: "generation"},
		{Name: "health:generate", Description: "Generate health plans", Category: "generation"},
		{Name: "recipe:generate", Description: "Generate recipes", Category: "generation"},
	}
}

// AvailablePermissions for backward compatibility
var AvailablePermissions = GetAvailablePermissions()

// EnhancedErrorResponse extends the basic ErrorResponse with additional context
type EnhancedErrorResponse struct {
	*ErrorResponse
	Context     *ErrorContext `json:"context,omitempty"`
	Retryable   bool          `json:"retryable,omitempty"`
	Suggestions []string      `json:"suggestions,omitempty"`
	Category    string        `json:"category,omitempty"`
	Severity    string        `json:"severity,omitempty"`
}

// ErrorContext provides additional context for error handling
type ErrorContext struct {
	UserID    string            `json:"user_id,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Endpoint  string            `json:"endpoint,omitempty"`
	Method    string            `json:"method,omitempty"`
	UserAgent string            `json:"user_agent,omitempty"`
	IPAddress string            `json:"ip_address,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value,omitempty"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
}

// ValidationErrorResponse represents validation error responses
type ValidationErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Code      string            `json:"code"`
	Timestamp time.Time         `json:"timestamp"`
	Errors    []ValidationError `json:"errors"`
}

// WasmError represents WebAssembly specific errors
type WasmError struct {
	ErrorType   string                 `json:"error"`
	Message     string                 `json:"message"`
	Code        string                 `json:"code"`
	Timestamp   time.Time              `json:"timestamp"`
	Module      string                 `json:"module,omitempty"`
	Function    string                 `json:"function,omitempty"`
	MemoryUsage int64                  `json:"memory_usage,omitempty"`
	InputSize   int64                  `json:"input_size,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// PerformanceError represents performance-related errors
type PerformanceError struct {
	ErrorType string        `json:"error"`
	Message   string        `json:"message"`
	Code      string        `json:"code"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	Threshold time.Duration `json:"threshold"`
	Operation string        `json:"operation"`
	Resource  string        `json:"resource"`
}

// SecurityError represents security-related errors
type SecurityError struct {
	ErrorType string                 `json:"error"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code"`
	Timestamp time.Time              `json:"timestamp"`
	Threat    string                 `json:"threat"`
	Severity  string                 `json:"severity"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxAttempts     int           `json:"max_attempts"`
	InitialDelay    time.Duration `json:"initial_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	RetryableErrors []string      `json:"retryable_errors"`
}

// ErrorMetrics represents error tracking metrics
type ErrorMetrics struct {
	ErrorType    string            `json:"error_type"`
	Count        int64             `json:"count"`
	LastOccurred time.Time         `json:"last_occurred"`
	Context      map[string]string `json:"context,omitempty"`
	Resolution   string            `json:"resolution,omitempty"`
}

// ErrorPattern represents recurring error patterns
type ErrorPattern struct {
	Pattern     string    `json:"pattern"`
	Description string    `json:"description"`
	Solution    string    `json:"solution"`
	Prevention  string    `json:"prevention"`
	Category    string    `json:"category"`
	Frequency   int       `json:"frequency"`
	LastSeen    time.Time `json:"last_seen"`
}

// NewEnhancedErrorResponse creates a new enhanced error response
func NewEnhancedErrorResponse(errorType, message, code string, retryable bool) *EnhancedErrorResponse {
	return &EnhancedErrorResponse{
		ErrorResponse: &ErrorResponse{
			Error:     errorType,
			Message:   message,
			Code:      code,
			Timestamp: time.Now().UTC(),
		},
		Retryable: retryable,
		Context:   &ErrorContext{},
	}
}

// WithContext adds context to the enhanced error response
func (e *EnhancedErrorResponse) WithContext(key string, value interface{}) *EnhancedErrorResponse {
	if e.Context == nil {
		e.Context = &ErrorContext{}
	}
	switch key {
	case "user_id":
		if v, ok := value.(string); ok {
			e.Context.UserID = v
		}
	case "request_id":
		if v, ok := value.(string); ok {
			e.Context.RequestID = v
		}
	case "endpoint":
		if v, ok := value.(string); ok {
			e.Context.Endpoint = v
		}
	case "method":
		if v, ok := value.(string); ok {
			e.Context.Method = v
		}
	case "user_agent":
		if v, ok := value.(string); ok {
			e.Context.UserAgent = v
		}
	case "ip_address":
		if v, ok := value.(string); ok {
			e.Context.IPAddress = v
		}
	default:
		if e.Context.Metadata == nil {
			e.Context.Metadata = make(map[string]string)
		}
		if v, ok := value.(string); ok {
			e.Context.Metadata[key] = v
		}
	}
	return e
}

// WithSuggestions adds suggestions to the error response
func (e *EnhancedErrorResponse) WithSuggestions(suggestions ...string) *EnhancedErrorResponse {
	e.Suggestions = suggestions
	return e
}

// WithCategory sets the error category
func (e *EnhancedErrorResponse) WithCategory(category string) *EnhancedErrorResponse {
	e.Category = category
	return e
}

// WithSeverity sets the error severity
func (e *EnhancedErrorResponse) WithSeverity(severity string) *EnhancedErrorResponse {
	e.Severity = severity
	return e
}

// WithTraceID adds a trace ID to the error response
func (e *EnhancedErrorResponse) WithTraceID(traceID string) *EnhancedErrorResponse {
	// Since we're extending the basic ErrorResponse, we need to handle this differently
	// For now, we'll store it in the Context map
	if e.Context == nil {
		e.Context = &ErrorContext{}
	}
	if e.Context.Metadata == nil {
		e.Context.Metadata = make(map[string]string)
	}
	e.Context.Metadata["trace_id"] = traceID
	return e
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(message string, errors []ValidationError) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Error:     "validation_error",
		Message:   message,
		Code:      "validation_failed",
		Timestamp: time.Now().UTC(),
		Errors:    errors,
	}
}

// Error implements the error interface for WasmError
func (e *WasmError) Error() string {
	return fmt.Sprintf("[%s] %s (Module: %s, Function: %s)", e.Code, e.Message, e.Module, e.Function)
}

// NewWasmError creates a new WebAssembly error
func NewWasmError(message, module, function string, memoryUsage, inputSize int64) *WasmError {
	return &WasmError{
		ErrorType:   "wasm_error",
		Message:     message,
		Code:        "wasm_execution_failed",
		Timestamp:   time.Now().UTC(),
		Module:      module,
		Function:    function,
		MemoryUsage: memoryUsage,
		InputSize:   inputSize,
		Context:     make(map[string]interface{}),
	}
}

// Error implements the error interface for PerformanceError
func (e *PerformanceError) Error() string {
	return fmt.Sprintf("[%s] %s (Operation: %s, Duration: %v, Threshold: %v)",
		e.Code, e.Message, e.Operation, e.Duration, e.Threshold)
}

// NewPerformanceError creates a new performance error
func NewPerformanceError(message, operation, resource string, duration, threshold time.Duration) *PerformanceError {
	return &PerformanceError{
		ErrorType: "performance_error",
		Message:   message,
		Code:      "performance_threshold_exceeded",
		Timestamp: time.Now().UTC(),
		Duration:  duration,
		Threshold: threshold,
		Operation: operation,
		Resource:  resource,
	}
}

// Error implements the error interface for SecurityError
func (e *SecurityError) Error() string {
	return fmt.Sprintf("[%s] %s (Threat: %s, Severity: %s)", e.Code, e.Message, e.Threat, e.Severity)
}

// NewSecurityError creates a new security error
func NewSecurityError(message, threat, severity string) *SecurityError {
	return &SecurityError{
		ErrorType: "security_error",
		Message:   message,
		Code:      "security_violation",
		Timestamp: time.Now().UTC(),
		Threat:    threat,
		Severity:  severity,
		Context:   make(map[string]interface{}),
	}
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:     3,
		InitialDelay:    time.Second,
		MaxDelay:        30 * time.Second,
		BackoffFactor:   2.0,
		RetryableErrors: []string{"network_error", "timeout", "temporary_failure", "rate_limit"},
	}
}

// IsRetryableError checks if an error is retryable based on configuration
func (r *RetryConfig) IsRetryableError(errorCode string) bool {
	for _, retryable := range r.RetryableErrors {
		if retryable == errorCode {
			return true
		}
	}
	return false
}
