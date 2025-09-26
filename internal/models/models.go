package models

import (
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
