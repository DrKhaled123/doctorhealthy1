package security

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-key-generator/internal/middleware"
	"api-key-generator/internal/models"
	"api-key-generator/internal/utils"
)

// Test Input Sanitization Functions
func TestSanitizationFunctions(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedLog   string
		expectedHTML  string
		expectedJSON  string
		expectedInput string
	}{
		{
			name:          "Clean Input",
			input:         "normal text",
			expectedLog:   "normal text",
			expectedHTML:  "normal text",
			expectedJSON:  "normal text",
			expectedInput: "normal text",
		},
		{
			name:          "XSS Script Tag",
			input:         "<script>alert('xss')</script>",
			expectedLog:   "",
			expectedHTML:  "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			expectedJSON:  "",
			expectedInput: "",
		},
		{
			name:          "SQL Injection",
			input:         "'; DROP TABLE users; --",
			expectedLog:   "",
			expectedHTML:  "&#39;; DROP TABLE users; --",
			expectedJSON:  "\\'; DROP TABLE users; --",
			expectedInput: "'; DROP TABLE users; --",
		},
		{
			name:          "Control Characters",
			input:         "test\x00\x01\x02data",
			expectedLog:   "testdata",
			expectedHTML:  "testdata",
			expectedJSON:  "testdata",
			expectedInput: "testdata",
		},
		{
			name:          "Unicode XSS",
			input:         "\u003cscript\u003ealert('xss')\u003c/script\u003e",
			expectedLog:   "alert('xss')",
			expectedHTML:  "＜script＞alert(&#39;xss&#39;)＜/script＞",
			expectedJSON:  "alert('xss')",
			expectedInput: "alert('xss')",
		},
		{
			name:          "Long Input",
			input:         strings.Repeat("a", 300),
			expectedLog:   strings.Repeat("a", 200) + "...",
			expectedHTML:  strings.Repeat("a", 300),
			expectedJSON:  strings.Repeat("a", 300),
			expectedInput: strings.Repeat("a", 300),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test SanitizeForLog
			logResult := utils.SanitizeForLog(tc.input)
			assert.Equal(t, tc.expectedLog, logResult, "SanitizeForLog failed")

			// Test SanitizeForHTML
			htmlResult := utils.SanitizeForHTML(tc.input)
			assert.Equal(t, tc.expectedHTML, htmlResult, "SanitizeForHTML failed")

			// Test SanitizeForJSON
			jsonResult := utils.SanitizeForJSON(tc.input)
			assert.Equal(t, tc.expectedJSON, jsonResult, "SanitizeForJSON failed")

			// Test SanitizeInput
			inputResult := utils.SanitizeInput(tc.input)
			assert.Equal(t, tc.expectedInput, inputResult, "SanitizeInput failed")
		})
	}
}

// Test XSS Prevention
func TestXSSPreventionValidation(t *testing.T) {
	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
		"<svg onload=alert('xss')>",
		"javascript:alert('xss')",
		"<iframe src=\"javascript:alert('xss')\"></iframe>",
		"<body onload=alert('xss')>",
		"<input onfocus=alert('xss') autofocus>",
		"<select onfocus=alert('xss') autofocus>",
		"<textarea onfocus=alert('xss') autofocus>",
		"<keygen onfocus=alert('xss') autofocus>",
		"<video><source onerror=\"alert('xss')\">",
		"<audio src=x onerror=alert('xss')>",
		"<details open ontoggle=alert('xss')>",
		"<marquee onstart=alert('xss')>",
		"\"><script>alert('xss')</script>",
		"';alert('xss');//",
		"\";alert('xss');//",
		"<script>alert(String.fromCharCode(88,83,83))</script>",
		"<img src=\"javascript:alert('xss')\">",
		"<div style=\"background-image:url(javascript:alert('xss'))\">",
	}

	for _, payload := range xssPayloads {
		t.Run("XSS: "+payload[:min(len(payload), 30)], func(t *testing.T) {
			// Test HTML sanitization
			sanitized := utils.SanitizeForHTML(payload)
			assert.NotContains(t, sanitized, "<script", "Script tag not properly escaped")
			assert.NotContains(t, sanitized, "javascript:", "JavaScript protocol not removed")
			assert.NotContains(t, sanitized, "onerror=", "Event handler not escaped")
			assert.NotContains(t, sanitized, "onload=", "Event handler not escaped")

			// Test input sanitization
			inputSanitized := utils.SanitizeInput(payload)
			assert.NotContains(t, inputSanitized, "<script", "Script tag not removed from input")
			assert.NotContains(t, inputSanitized, "javascript:", "JavaScript protocol not removed from input")
		})
	}
}

// Test SQL Injection Input Validation
func TestSQLInjectionInputValidation(t *testing.T) {
	sqlPayloads := []string{
		"'; DROP TABLE users; --",
		"' OR '1'='1",
		"' OR 1=1 --",
		"admin'--",
		"admin' /*",
		"' OR 'x'='x",
		"' UNION SELECT NULL,NULL,NULL --",
		"'; INSERT INTO users VALUES ('hacker','password'); --",
		"1; DELETE FROM users WHERE 1=1 --",
		"' OR '1'='1' /*",
		"' OR '1'='1' --",
		"' OR '1'='1' #",
		"') OR ('1'='1",
		"') OR ('1'='1' --",
		"1' OR '1'='1",
		"1' OR '1'='1' --",
		"1' OR '1'='1' /*",
		"1' OR '1'='1' #",
		"admin'; --",
		"admin' #",
	}

	suite := setupSecurityTestSuite(t)

	for _, payload := range sqlPayloads {
		t.Run("SQL Injection: "+payload[:min(len(payload), 30)], func(t *testing.T) {
			// Test that malicious SQL doesn't break the system
			_, err := suite.apiKeyService.GetAPIKeyByKey(context.Background(), payload)

			// Should return not found error, not crash
			assert.Error(t, err)
			assert.Contains(t, strings.ToLower(err.Error()), "not found")

			// Verify database integrity by checking valid key still works
			validKey, err := suite.apiKeyService.GetAPIKeyByKey(context.Background(), suite.testAPIKey)
			assert.NoError(t, err)
			assert.NotNil(t, validKey)
		})
	}
}

// Test API Key Creation Input Validation
func TestAPIKeyCreationInputValidation(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	maliciousInputs := []struct {
		name        string
		request     models.CreateAPIKeyRequest
		expectError bool
	}{
		{
			name: "XSS in Name",
			request: models.CreateAPIKeyRequest{
				Name:        "<script>alert('xss')</script>",
				Permissions: []string{"recipes:read"},
			},
			expectError: false, // Should be sanitized, not rejected
		},
		{
			name: "SQL Injection in Name",
			request: models.CreateAPIKeyRequest{
				Name:        "test'; DROP TABLE api_keys; --",
				Permissions: []string{"recipes:read"},
			},
			expectError: false, // Should be sanitized, not rejected
		},
		{
			name: "Empty Name",
			request: models.CreateAPIKeyRequest{
				Name:        "",
				Permissions: []string{"recipes:read"},
			},
			expectError: true,
		},
		{
			name: "Very Long Name",
			request: models.CreateAPIKeyRequest{
				Name:        strings.Repeat("a", 200),
				Permissions: []string{"recipes:read"},
			},
			expectError: true,
		},
		{
			name: "Invalid Permission",
			request: models.CreateAPIKeyRequest{
				Name:        "test",
				Permissions: []string{"invalid:permission"},
			},
			expectError: true,
		},
		{
			name: "XSS in Permission",
			request: models.CreateAPIKeyRequest{
				Name:        "test",
				Permissions: []string{"<script>alert('xss')</script>"},
			},
			expectError: true,
		},
		{
			name: "Empty Permissions",
			request: models.CreateAPIKeyRequest{
				Name:        "test",
				Permissions: []string{},
			},
			expectError: true,
		},
	}

	for _, tc := range maliciousInputs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := suite.apiKeyService.CreateAPIKey(context.Background(), &tc.request)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test HTTP Header Injection
func TestHTTPHeaderInjection(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	headerInjectionPayloads := []string{
		"test\r\nContent-Type: text/html",
		"test\nSet-Cookie: admin=true",
		"test\r\n\r\n<script>alert('xss')</script>",
		"test%0d%0aContent-Type:%20text/html",
		"test%0aSet-Cookie:%20admin=true",
	}

	for _, payload := range headerInjectionPayloads {
		t.Run("Header Injection: "+payload[:min(len(payload), 30)], func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("X-API-Key", payload)
			rec := httptest.NewRecorder()
			c := suite.e.NewContext(req, rec)

			// Test that header injection doesn't work
			mw := middleware.ScopesAny(suite.apiKeyService, "recipes:read")
			handler := mw(func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"status": "success"})
			})

			err := handler(c)

			// Should return unauthorized, not inject headers
			assert.Error(t, err)

			// Check that no malicious headers were set
			assert.NotContains(t, rec.Header().Get("Content-Type"), "text/html")
			assert.Empty(t, rec.Header().Get("Set-Cookie"))
		})
	}
}

// Test JSON Injection
func TestJSONInjection(t *testing.T) {
	suite := setupSecurityTestSuite(t)

	jsonPayloads := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name: "Object Prototype Pollution",
			payload: map[string]interface{}{
				"__proto__": map[string]interface{}{
					"admin": true,
				},
				"name": "test",
			},
		},
		{
			name: "Constructor Pollution",
			payload: map[string]interface{}{
				"constructor": map[string]interface{}{
					"prototype": map[string]interface{}{
						"admin": true,
					},
				},
				"name": "test",
			},
		},
		{
			name: "Large Number",
			payload: map[string]interface{}{
				"name":   "test",
				"number": int64(9223372036854775807), // Max int64
			},
		},
		{
			name: "Special Characters",
			payload: map[string]interface{}{
				"name": "test\x00\x01\x02",
			},
		},
	}

	for _, tc := range jsonPayloads {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := suite.e.NewContext(req, rec)

			// Test that JSON parsing is safe
			var parsed map[string]interface{}
			err = c.Bind(&parsed)

			// Should either parse safely or reject
			if err == nil {
				// If parsed, ensure no prototype pollution
				assert.Nil(t, parsed["__proto__"])
				assert.Nil(t, parsed["constructor"])
			}
		})
	}
}

// Test Path Traversal Prevention
func TestPathTraversalPrevention(t *testing.T) {
	pathTraversalPayloads := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"....//....//....//etc/passwd",
		"..%2f..%2f..%2fetc%2fpasswd",
		"..%252f..%252f..%252fetc%252fpasswd",
		"..%c0%af..%c0%af..%c0%afetc%c0%afpasswd",
		"..%c1%9c..%c1%9c..%c1%9cetc%c1%9cpasswd",
		"/%2e%2e/%2e%2e/%2e%2e/etc/passwd",
		"/\\.\\.\\/\\.\\.\\/etc/passwd",
		"..///////../../../etc/passwd",
	}

	for _, payload := range pathTraversalPayloads {
		t.Run("Path Traversal: "+payload[:min(len(payload), 30)], func(t *testing.T) {
			sanitized := utils.SanitizeInput(payload)

			// Should not contain path traversal sequences
			assert.NotContains(t, sanitized, "../")
			assert.NotContains(t, sanitized, "..\\")
			assert.NotContains(t, sanitized, "%2e%2e")
			assert.NotContains(t, sanitized, "%252e")
		})
	}
}

// Test NoSQL Injection Prevention
func TestNoSQLInjectionPrevention(t *testing.T) {
	noSQLPayloads := []string{
		`{"$gt": ""}`,
		`{"$ne": null}`,
		`{"$regex": ".*"}`,
		`{"$where": "function() { return true; }"}`,
		`{"$expr": {"$gt": [{"$size": "$items"}, 0]}}`,
		`{"username": {"$ne": null}, "password": {"$ne": null}}`,
		`{"$or": [{"username": "admin"}, {"role": "admin"}]}`,
		`{"$and": [{"$or": [{"username": "admin"}]}]}`,
	}

	for _, payload := range noSQLPayloads {
		t.Run("NoSQL Injection: "+payload[:min(len(payload), 30)], func(t *testing.T) {
			sanitized := utils.SanitizeInput(payload)

			// Should remove or escape NoSQL operators
			assert.NotContains(t, sanitized, "$gt")
			assert.NotContains(t, sanitized, "$ne")
			assert.NotContains(t, sanitized, "$regex")
			assert.NotContains(t, sanitized, "$where")
			assert.NotContains(t, sanitized, "$expr")
			assert.NotContains(t, sanitized, "$or")
			assert.NotContains(t, sanitized, "$and")
		})
	}
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
