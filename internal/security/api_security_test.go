package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-key-generator/internal/config"
	"api-key-generator/internal/database"
	"api-key-generator/internal/middleware"
	"api-key-generator/internal/models"
	"api-key-generator/internal/services"
	"api-key-generator/internal/utils"
)

// TestSuite for API Security Tests
type APISecurityTestSuite struct {
	e             *echo.Echo
	apiKeyService *services.APIKeyService
	db            *sql.DB
	testAPIKey    string
	testAPIKeyID  string
}

func setupSecurityTestSuite(t *testing.T) *APISecurityTestSuite {
	t.Helper()

	// Create in-memory database
	dbPath := ":memory:"
	db, err := database.Initialize(dbPath)
	require.NoError(t, err)

	// Setup configuration
	cfg := &config.Config{
		APIKey: config.APIKeyConfig{
			Length:         16,
			ExpiryDuration: 24 * time.Hour,
			Prefix:         "ak_",
		},
		Security: config.SecurityConfig{
			RateLimitRequests: 100,
			RateLimitWindow:   time.Minute,
		},
	}

	// Initialize services
	apiKeyService, err := services.NewAPIKeyService(db, cfg)
	require.NoError(t, err, "Failed to create APIKeyService")

	// Create test API key
	req := &models.CreateAPIKeyRequest{
		Name:        "security-test-key",
		Permissions: []string{"recipes:read", "recipes:write", "users:read", "users:write", "admin:all"},
	}

	apiKey, err := apiKeyService.CreateAPIKey(context.Background(), req)
	require.NoError(t, err)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Security())
	e.Use(middleware.RateLimit(cfg))

	t.Cleanup(func() {
		apiKeyService.Close()
		db.Close()
	})

	return &APISecurityTestSuite{
		e:             e,
		apiKeyService: apiKeyService,
		db:            db,
		testAPIKey:    apiKey.Key,
		testAPIKeyID:  apiKey.ID,
	}
}

// Test API Key Authentication Security
func TestAPIKeyAuthenticationSecurity(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid API Key",
			apiKey:         suite.testAPIKey,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Empty API Key",
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "Invalid API Key",
			apiKey:         "invalid-key-123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "Malformed API Key",
			apiKey:         "ak_malformed",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "SQL Injection in API Key",
			apiKey:         "ak_'; DROP TABLE api_keys; --",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name:           "XSS Attempt in API Key",
			apiKey:         "<script>alert('xss')</script>",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test endpoint with API key middleware
			mw := middleware.ScopesAny(suite.apiKeyService, "recipes:read")
			handler := mw(func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"status": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()
			c := suite.e.NewContext(req, rec)

			err := handler(c)

			if tt.expectedError {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

// Test JWT Token Security
func TestJWTTokenSecurity(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	// Generate valid JWT
	validToken, err := utils.GenerateJWT("test-user-123", time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedUserID string
		shouldError    bool
	}{
		{
			name:           "Valid JWT Token",
			authHeader:     "Bearer " + validToken,
			expectedUserID: "test-user-123",
			shouldError:    false,
		},
		{
			name:           "Empty Authorization Header",
			authHeader:     "",
			expectedUserID: "",
			shouldError:    false, // OptionalJWT should not error
		},
		{
			name:           "Invalid Bearer Format",
			authHeader:     "Bearer",
			expectedUserID: "",
			shouldError:    false,
		},
		{
			name:           "Invalid JWT Token",
			authHeader:     "Bearer invalid.jwt.token",
			expectedUserID: "",
			shouldError:    false,
		},
		{
			name:           "Malicious JWT Token",
			authHeader:     "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiJhZG1pbiIsImlhdCI6MTUxNjIzOTAyMn0.",
			expectedUserID: "",
			shouldError:    false,
		},
		{
			name:           "SQL Injection in JWT",
			authHeader:     "Bearer '; DROP TABLE users; --",
			expectedUserID: "",
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := middleware.OptionalJWT()
			handler := mw(func(c echo.Context) error {
				userID, _ := c.Get("user_id").(string)
				return c.JSON(http.StatusOK, map[string]string{"user_id": userID})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			c := suite.e.NewContext(req, rec)

			err := handler(c)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedUserID, response["user_id"])
			}
		})
	}
}

// Test Input Validation and Sanitization
func TestInputValidationSecurity(t *testing.T) {
	maliciousInputs := []struct {
		name  string
		input string
	}{
		{"SQL Injection", "'; DROP TABLE api_keys; --"},
		{"XSS Script Tag", "<script>alert('xss')</script>"},
		{"XSS Event Handler", "<img src=x onerror=alert('xss')>"},
		{"Command Injection", "; rm -rf /"},
		{"Path Traversal", "../../../etc/passwd"},
		{"LDAP Injection", "admin)(&(password=*)"},
		{"XML External Entity", "<!DOCTYPE foo [<!ENTITY xxe SYSTEM \"file:///etc/passwd\">]>"},
		{"NoSQL Injection", "{\"$gt\": \"\"}"},
		{"Control Characters", "\x00\x01\x02\x03"},
		{"Unicode Bypass", "\u003cscript\u003ealert('xss')\u003c/script\u003e"},
	}

	for _, malicious := range maliciousInputs {
		t.Run(malicious.name, func(t *testing.T) {
			// Test sanitization functions
			sanitizedLog := utils.SanitizeForLog(malicious.input)
			sanitizedHTML := utils.SanitizeForHTML(malicious.input)
			sanitizedJSON := utils.SanitizeForJSON(malicious.input)
			sanitizedInput := utils.SanitizeInput(malicious.input)

			// Ensure no malicious content passes through
			assert.NotContains(t, sanitizedLog, "<script")
			assert.NotContains(t, sanitizedHTML, "<script")
			assert.NotContains(t, sanitizedJSON, "<script")
			assert.NotContains(t, sanitizedInput, "<script")

			// Ensure SQL injection patterns are handled
			assert.NotContains(t, sanitizedInput, "DROP TABLE")
			assert.NotContains(t, sanitizedInput, "--")

			// Ensure control characters are removed
			assert.NotContains(t, sanitizedLog, "\x00")
			assert.NotContains(t, sanitizedInput, "\x00")
		})
	}
}

// Test Rate Limiting Security
func TestRateLimitingSecurity(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	// Create a simple handler
	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}

	suite.e.GET("/rate-test", handler)

	// Test rate limiting by making rapid requests
	successCount := 0
	rateLimitedCount := 0

	for i := 0; i < 150; i++ { // Exceed the rate limit
		req := httptest.NewRequest(http.MethodGet, "/rate-test", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1") // Same IP
		rec := httptest.NewRecorder()

		suite.e.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK {
			successCount++
		} else if rec.Code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// Should have some successful requests and some rate-limited
	assert.Greater(t, successCount, 0, "Should have some successful requests")
	assert.Greater(t, rateLimitedCount, 0, "Should have some rate-limited requests")
}

// Test Security Headers
func TestSecurityHeaders(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}

	suite.e.GET("/security-headers-test", handler)

	req := httptest.NewRequest(http.MethodGet, "/security-headers-test", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	// Check security headers
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rec.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", rec.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	assert.Equal(t, "default-src 'self'", rec.Header().Get("Content-Security-Policy"))
}

// Test Permission-Based Access Control
func TestPermissionBasedAccessControl(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	// Create API key with limited permissions
	limitedReq := &models.CreateAPIKeyRequest{
		Name:        "limited-key",
		Permissions: []string{"recipes:read"}, // Only read permission
	}

	limitedKey, err := suite.apiKeyService.CreateAPIKey(context.Background(), limitedReq)
	require.NoError(t, err)

	tests := []struct {
		name               string
		apiKey             string
		requiredPermission string
		expectAccess       bool
	}{
		{
			name:               "Full Access Key - Read Permission",
			apiKey:             suite.testAPIKey,
			requiredPermission: "recipes:read",
			expectAccess:       true,
		},
		{
			name:               "Full Access Key - Write Permission",
			apiKey:             suite.testAPIKey,
			requiredPermission: "recipes:write",
			expectAccess:       true,
		},
		{
			name:               "Limited Key - Allowed Permission",
			apiKey:             limitedKey.Key,
			requiredPermission: "recipes:read",
			expectAccess:       true,
		},
		{
			name:               "Limited Key - Forbidden Permission",
			apiKey:             limitedKey.Key,
			requiredPermission: "recipes:write",
			expectAccess:       false,
		},
		{
			name:               "Limited Key - Admin Permission",
			apiKey:             limitedKey.Key,
			requiredPermission: "admin:all",
			expectAccess:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := middleware.ScopesAny(suite.apiKeyService, tt.requiredPermission)
			handler := mw(func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"status": "authorized"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			rec := httptest.NewRecorder()
			c := suite.e.NewContext(req, rec)

			err := handler(c)

			if tt.expectAccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, http.StatusForbidden, httpErr.Code)
				}
			}
		})
	}
}

// Test SQL Injection Prevention
func TestSQLInjectionPrevention(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	maliciousSQLInputs := []string{
		"'; DROP TABLE api_keys; --",
		"' OR '1'='1",
		"' UNION SELECT * FROM api_keys --",
		"'; INSERT INTO api_keys (name) VALUES ('hacked'); --",
		"' OR 1=1 --",
		"admin'--",
		"admin' /*",
		"' OR 'x'='x",
		"'; EXEC xp_cmdshell('dir'); --",
		"1; DELETE FROM api_keys WHERE 1=1 --",
	}

	for _, maliciousSQL := range maliciousSQLInputs {
		t.Run(fmt.Sprintf("SQL Injection: %s", maliciousSQL), func(t *testing.T) {
			// Test getting API key with malicious input (should not cause errors or expose data)
			_, err := suite.apiKeyService.GetAPIKeyByKey(context.Background(), maliciousSQL)

			// Should return "not found" error, not crash or expose data
			assert.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), "not found")

			// Database should still be intact - verify by getting valid key
			validKey, err := suite.apiKeyService.GetAPIKeyByKey(context.Background(), suite.testAPIKey)
			assert.NoError(t, err)
			assert.NotNil(t, validKey)
		})
	}
}

// Test API Key Expiration Security
func TestAPIKeyExpirationSecurity(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	// Create an already expired API key
	expiredReq := &models.CreateAPIKeyRequest{
		Name:        "expired-key",
		Permissions: []string{"recipes:read"},
		ExpiryDays:  func() *int { days := -1; return &days }(), // Already expired
	}

	expiredKey, err := suite.apiKeyService.CreateAPIKey(context.Background(), expiredReq)
	require.NoError(t, err)

	// Try to use expired key
	mw := middleware.ScopesAny(suite.apiKeyService, "recipes:read")
	handler := mw(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", expiredKey.Key)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	err = handler(c)

	// Should be unauthorized due to expiration
	assert.Error(t, err)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	}
}

// Test Concurrent Access Security
func TestConcurrentAccessSecurity(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	// Test concurrent access to API key validation
	const numGoroutines = 10
	const requestsPerGoroutine = 10

	results := make(chan error, numGoroutines*requestsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				_, err := suite.apiKeyService.GetAPIKeyByKey(context.Background(), suite.testAPIKey)
				results <- err
			}
		}()
	}

	// Collect results
	successCount := 0
	errorCount := 0

	for i := 0; i < numGoroutines*requestsPerGoroutine; i++ {
		err := <-results
		if err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	// All requests should succeed (no race conditions)
	assert.Equal(t, numGoroutines*requestsPerGoroutine, successCount)
	assert.Equal(t, 0, errorCount)
}
