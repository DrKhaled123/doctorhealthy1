package security

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"api-key-generator/internal/config"
	"api-key-generator/internal/middleware"
)

// Test Rate Limiting Middleware
func TestRateLimitingMiddleware(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			RateLimitRequests: 5,           // 5 requests
			RateLimitWindow:   time.Second, // per second
		},
	}

	e := echo.New()
	e.Use(middleware.RateLimit(cfg))

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	e.GET("/test", handler)

	// Test rate limiting
	successCount := 0
	rateLimitedCount := 0

	// Make requests rapidly to trigger rate limiting
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1") // Same IP
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			successCount++
		} else if rec.Code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// Should have some successful requests and some rate-limited
	assert.Greater(t, successCount, 0, "Should have some successful requests")
	assert.Greater(t, rateLimitedCount, 0, "Should have some rate-limited requests")
	assert.Equal(t, 10, successCount+rateLimitedCount, "All requests should be accounted for")
}

// Test Rate Limiting with Different IPs
func TestRateLimitingDifferentIPs(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			RateLimitRequests: 2, // Very low limit
			RateLimitWindow:   time.Second,
		},
	}

	e := echo.New()
	e.Use(middleware.RateLimit(cfg))

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	e.GET("/test", handler)

	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for _, ip := range ips {
		t.Run("IP: "+ip, func(t *testing.T) {
			successCount := 0
			rateLimitedCount := 0

			// Each IP should get its own rate limit bucket
			for i := 0; i < 5; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.Header.Set("X-Forwarded-For", ip)
				rec := httptest.NewRecorder()

				e.ServeHTTP(rec, req)

				if rec.Code == http.StatusOK {
					successCount++
				} else if rec.Code == http.StatusTooManyRequests {
					rateLimitedCount++
				}
			}

			// Each IP should have its own limit
			assert.Greater(t, successCount, 0, "Should have some successful requests for IP "+ip)
		})
	}
}

// Test User-Based Rate Limiting
func TestUserRateLimit(t *testing.T) {
	e := echo.New()

	// Use a stricter user rate limit
	e.Use(middleware.UserRateLimit(rate.Limit(2), 2)) // 2 requests per second, burst of 2

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}

	tests := []struct {
		name   string
		userID string
	}{
		{"User1", "user-123"},
		{"User2", "user-456"},
		{"NoUser", ""}, // Should skip rate limiting
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Test multiple requests
			successCount := 0
			rateLimitedCount := 0

			for i := 0; i < 5; i++ {
				err := handler(c)

				// Reset recorder for each iteration
				rec = httptest.NewRecorder()
				c = e.NewContext(req, rec)
				if tt.userID != "" {
					c.Set("user_id", tt.userID)
				}

				mw := middleware.UserRateLimit(rate.Limit(2), 2)
				wrappedHandler := mw(handler)
				err = wrappedHandler(c)

				if err == nil {
					successCount++
				} else {
					if httpErr, ok := err.(*echo.HTTPError); ok && httpErr.Code == http.StatusTooManyRequests {
						rateLimitedCount++
					}
				}
			}

			if tt.userID == "" {
				// No user ID, should not be rate limited
				assert.Equal(t, 5, successCount, "Requests without user_id should not be rate limited")
			} else {
				// With user ID, should be rate limited
				assert.Greater(t, rateLimitedCount, 0, "User requests should be rate limited")
			}
		})
	}
}

// Test Monthly Quota Middleware
func TestMonthlyQuotaMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(middleware.UserQuotaMiddleware())

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "generated"})
	}

	// Test endpoints that should trigger quota
	e.POST("/api/v1/enhanced/diet/generate", handler)
	e.POST("/api/v1/enhanced/workout/generate", handler)
	e.POST("/api/v1/enhanced/lifestyle/generate", handler)

	// Test endpoint that should NOT trigger quota
	e.GET("/api/v1/recipes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	tests := []struct {
		name          string
		endpoint      string
		method        string
		shouldQuota   bool
		plan          string
		expectedLimit int
	}{
		{
			name:          "Generate endpoint - Free plan",
			endpoint:      "/api/v1/enhanced/diet/generate",
			method:        http.MethodPost,
			shouldQuota:   true,
			plan:          "free",
			expectedLimit: 3,
		},
		{
			name:          "Generate endpoint - Pro plan",
			endpoint:      "/api/v1/enhanced/workout/generate",
			method:        http.MethodPost,
			shouldQuota:   true,
			plan:          "pro",
			expectedLimit: 50,
		},
		{
			name:          "Generate endpoint - Lifetime plan",
			endpoint:      "/api/v1/enhanced/lifestyle/generate",
			method:        http.MethodPost,
			shouldQuota:   true,
			plan:          "lifetime",
			expectedLimit: 1000000, // Effectively unlimited
		},
		{
			name:          "Non-generate endpoint",
			endpoint:      "/api/v1/recipes",
			method:        http.MethodGet,
			shouldQuota:   false,
			plan:          "free",
			expectedLimit: 0, // Not applicable
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			successCount := 0
			quotaExceededCount := 0
			testLimit := minInt(tt.expectedLimit+2, 10) // Test a reasonable number

			for i := 0; i < testLimit; i++ {
				req := httptest.NewRequest(tt.method, tt.endpoint, nil)

				// Set plan cookie
				if tt.plan != "" {
					cookie := &http.Cookie{
						Name:  "plan",
						Value: tt.plan,
					}
					req.AddCookie(cookie)
				}

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				if rec.Code == http.StatusOK {
					successCount++
				} else if rec.Code == http.StatusTooManyRequests {
					quotaExceededCount++
				}
			}

			if !tt.shouldQuota {
				// Non-quota endpoints should always succeed
				assert.Equal(t, testLimit, successCount, "Non-quota endpoints should not be limited")
				assert.Equal(t, 0, quotaExceededCount, "Non-quota endpoints should not be quota limited")
			} else {
				// Quota endpoints should eventually hit limits (except lifetime)
				if tt.plan == "lifetime" {
					assert.Equal(t, testLimit, successCount, "Lifetime plan should not hit quota")
				} else if tt.expectedLimit < testLimit {
					assert.Greater(t, quotaExceededCount, 0, "Quota should be enforced for plan: "+tt.plan)
				}
			}
		})
	}
}

// Test Quota with Shared Plan Bonus
func TestQuotaWithSharedBonus(t *testing.T) {
	e := echo.New()
	e.Use(middleware.UserQuotaMiddleware())

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "generated"})
	}
	e.POST("/api/v1/enhanced/diet/generate", handler)

	tests := []struct {
		name          string
		shared        string
		expectedLimit int
	}{
		{
			name:          "Free plan without shared",
			shared:        "",
			expectedLimit: 3,
		},
		{
			name:          "Free plan with shared=yes",
			shared:        "yes",
			expectedLimit: 11,
		},
		{
			name:          "Free plan with shared=no",
			shared:        "no",
			expectedLimit: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			successCount := 0
			quotaExceededCount := 0

			// Test up to the limit + 2 to ensure quota is enforced
			testRequests := tt.expectedLimit + 2

			for i := 0; i < testRequests; i++ {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/enhanced/diet/generate", nil)

				// Set plan to free
				planCookie := &http.Cookie{
					Name:  "plan",
					Value: "free",
				}
				req.AddCookie(planCookie)

				// Set shared cookie if specified
				if tt.shared != "" {
					sharedCookie := &http.Cookie{
						Name:  "shared",
						Value: tt.shared,
					}
					req.AddCookie(sharedCookie)
				}

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				if rec.Code == http.StatusOK {
					successCount++
				} else if rec.Code == http.StatusTooManyRequests {
					quotaExceededCount++
				}
			}

			// Should succeed up to the limit, then fail
			assert.LessOrEqual(t, successCount, tt.expectedLimit, "Should not exceed quota limit")
			assert.Greater(t, quotaExceededCount, 0, "Should hit quota limit")
		})
	}
}

// Test Concurrent Rate Limiting
func TestConcurrentRateLimit(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			RateLimitRequests: 10,
			RateLimitWindow:   time.Second,
		},
	}

	e := echo.New()
	e.Use(middleware.RateLimit(cfg))

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	e.GET("/test", handler)

	// Test concurrent requests
	const numGoroutines = 5
	const requestsPerGoroutine = 10

	var wg sync.WaitGroup
	results := make(chan int, numGoroutines*requestsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.Header.Set("X-Forwarded-For", "127.0.0.1")
				rec := httptest.NewRecorder()

				e.ServeHTTP(rec, req)
				results <- rec.Code
			}
		}(i)
	}

	wg.Wait()
	close(results)

	// Count results
	successCount := 0
	rateLimitedCount := 0

	for code := range results {
		if code == http.StatusOK {
			successCount++
		} else if code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// Should have both successful and rate-limited requests
	assert.Greater(t, successCount, 0, "Should have some successful requests")
	assert.Greater(t, rateLimitedCount, 0, "Should have some rate-limited requests")

	totalRequests := numGoroutines * requestsPerGoroutine
	assert.Equal(t, totalRequests, successCount+rateLimitedCount, "All requests should be accounted for")
}

// Test Rate Limit Reset After Time Window
func TestRateLimitReset(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			RateLimitRequests: 2,                      // Very low limit
			RateLimitWindow:   100 * time.Millisecond, // Short window
		},
	}

	e := echo.New()
	e.Use(middleware.RateLimit(cfg))

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
	e.GET("/test", handler)

	// First batch of requests - should hit rate limit
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}

	// Wait for rate limit window to reset
	time.Sleep(200 * time.Millisecond)

	// Second batch of requests - should succeed again
	successCount := 0
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			successCount++
		}
	}

	// Should have some successful requests after reset
	assert.Greater(t, successCount, 0, "Should have successful requests after rate limit reset")
}

// Helper function - using minInt to avoid conflicts
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
