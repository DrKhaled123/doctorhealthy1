package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type RateLimiter struct {
	db     *sql.DB
	mu     sync.RWMutex
	limits map[string]*UserLimit
}

type UserLimit struct {
	Count    int
	ResetAt  time.Time
	LastSeen time.Time
}

func NewRateLimiter(db *sql.DB) *RateLimiter {
	rl := &RateLimiter{
		db:     db,
		limits: make(map[string]*UserLimit),
	}

	// Clean up expired entries every minute
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware(limit int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from context (set by auth middleware)
			userID := c.Get("user_id")
			if userID == nil {
				// Use IP address for anonymous users
				userID = c.RealIP()
			}

			userKey := fmt.Sprintf("%v", userID)

			// Check rate limit
			allowed, resetAt := rl.IsAllowed(userKey, limit)
			if !allowed {
				c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))
				c.Response().Header().Set("Retry-After", fmt.Sprintf("%d", resetAt.Unix()))
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "Rate limit exceeded",
					"reset_at":    resetAt.Unix(),
					"retry_after": resetAt.Unix(),
				})
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limit-1))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))

			return next(c)
		}
	}
}

func (rl *RateLimiter) IsAllowed(userKey string, limit int) (bool, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create user limit
	userLimit, exists := rl.limits[userKey]
	if !exists {
		userLimit = &UserLimit{
			Count:    0,
			ResetAt:  now.Add(time.Minute),
			LastSeen: now,
		}
		rl.limits[userKey] = userLimit
	}

	// Reset if window expired
	if now.After(userLimit.ResetAt) {
		userLimit.Count = 0
		userLimit.ResetAt = now.Add(time.Minute)
	}

	// Check limit
	if userLimit.Count >= limit {
		return false, userLimit.ResetAt
	}

	// Increment counter
	userLimit.Count++
	userLimit.LastSeen = now

	return true, userLimit.ResetAt
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for userKey, userLimit := range rl.limits {
			if now.Sub(userLimit.LastSeen) > time.Minute {
				delete(rl.limits, userKey)
			}
		}
		rl.mu.Unlock()
	}
}
