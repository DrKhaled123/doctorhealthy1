package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"api-key-generator/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SetupErrorHandler installs a centralized HTTP error handler with consistent JSON responses
func SetupErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// If response already committed, delegate to default
		if c.Response().Committed {
			c.Echo().DefaultHTTPErrorHandler(err, c)
			return
		}

		// Generate trace ID for error tracking
		traceID := uuid.New().String()

		// Map status code
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok && he != nil {
			if he.Code > 0 {
				code = he.Code
			}
		}

		// Build enhanced error response with context
		errorResp := models.NewEnhancedErrorResponse(
			http.StatusText(code),
			sanitizeHTTPErrorMessage(err),
			httpStatusCodeString(code),
			isRetryableError(code),
		).WithTraceID(traceID)

		// Add request context
		errorResp.WithContext("endpoint", c.Request().URL.Path)
		errorResp.WithContext("method", c.Request().Method)
		errorResp.WithContext("user_agent", c.Request().UserAgent())
		errorResp.WithContext("ip_address", c.RealIP())

		// Add user ID if available from JWT middleware
		if userID := c.Get("user_id"); userID != nil {
			errorResp.WithContext("user_id", userID.(string))
		}

		// Add request ID if available
		if requestID := c.Get("request_id"); requestID != nil {
			errorResp.WithContext("request_id", requestID.(string))
		}

		// Categorize error and add suggestions
		categorizeError(errorResp, err, code)

		// Log error with stack trace for debugging
		logError(err, c, traceID)

		// Attempt JSON response; on failure, fall back to default handler
		if jsonErr := c.JSON(code, errorResp); jsonErr != nil {
			c.Echo().DefaultHTTPErrorHandler(err, c)
		}
	}
}

// SetupAdvancedErrorHandler sets up advanced error handling with error boundaries
func SetupAdvancedErrorHandler(e *echo.Echo) {
	// Set up the enhanced error handler
	SetupErrorHandler(e)

	// Add error recovery middleware
	e.Use(ErrorBoundaryMiddleware())
}

// ErrorBoundaryMiddleware provides error boundary functionality
func ErrorBoundaryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic with stack trace
					traceID := uuid.New().String()
					log.Printf("PANIC RECOVERED [%s]: %v\nStack: %s", traceID, r, debug.Stack())

					// Create panic error response
					panicErr := models.NewEnhancedErrorResponse(
						"Internal Server Error",
						"An unexpected error occurred",
						"internal_server_error",
						false,
					).WithTraceID(traceID)

					panicErr.WithContext("endpoint", c.Request().URL.Path)
					panicErr.WithContext("method", c.Request().Method)
					panicErr.WithContext("panic", "true")
					panicErr.WithCategory("runtime")
					panicErr.WithSeverity("critical")
					panicErr.WithSuggestions(
						"Please try again later",
						"Contact support if the problem persists",
					)

					// Send error response
					c.JSON(http.StatusInternalServerError, panicErr)
				}
			}()

			return next(c)
		}
	}
}

func httpStatusCodeString(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusTooManyRequests:
		return "too_many_requests"
	case http.StatusUnprocessableEntity:
		return "unprocessable_entity"
	case http.StatusInternalServerError:
		fallthrough
	default:
		return "internal_server_error"
	}
}

// isRetryableError determines if an error should be retryable
func isRetryableError(code int) bool {
	switch code {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// categorizeError categorizes the error and adds appropriate suggestions
func categorizeError(errorResp *models.EnhancedErrorResponse, err error, code int) {
	switch code {
	case http.StatusBadRequest:
		errorResp.WithCategory("validation")
		errorResp.WithSeverity("low")
		errorResp.WithSuggestions(
			"Check your request parameters",
			"Ensure all required fields are provided",
		)
	case http.StatusUnauthorized:
		errorResp.WithCategory("authentication")
		errorResp.WithSeverity("medium")
		errorResp.WithSuggestions(
			"Check your API key or authentication token",
			"Ensure you have the required permissions",
		)
	case http.StatusForbidden:
		errorResp.WithCategory("authorization")
		errorResp.WithSeverity("medium")
		errorResp.WithSuggestions(
			"Verify your permissions for this resource",
			"Contact administrator if you need access",
		)
	case http.StatusNotFound:
		errorResp.WithCategory("resource")
		errorResp.WithSeverity("low")
		errorResp.WithSuggestions(
			"Check the resource ID or path",
			"Ensure the resource exists",
		)
	case http.StatusTooManyRequests:
		errorResp.WithCategory("rate_limit")
		errorResp.WithSeverity("medium")
		errorResp.WithSuggestions(
			"Wait before making more requests",
			"Check rate limits for your API key",
		)
	case http.StatusUnprocessableEntity:
		errorResp.WithCategory("validation")
		errorResp.WithSeverity("medium")
		errorResp.WithSuggestions(
			"Check data format and validation rules",
			"Ensure data types are correct",
		)
	case http.StatusInternalServerError:
		errorResp.WithCategory("server")
		errorResp.WithSeverity("high")
		errorResp.WithSuggestions(
			"Try again in a few moments",
			"Contact support if the issue persists",
		)
	default:
		errorResp.WithCategory("unknown")
		errorResp.WithSeverity("medium")
	}

	// Add specific error type detection
	if strings.Contains(err.Error(), "wasm") {
		errorResp.WithCategory("wasm")
		errorResp.WithSuggestions(
			"Check WebAssembly module compatibility",
			"Verify input data size limits",
		)
	}

	if strings.Contains(err.Error(), "database") || strings.Contains(err.Error(), "sql") {
		errorResp.WithCategory("database")
		errorResp.WithSuggestions(
			"Check database connectivity",
			"Verify data integrity",
		)
	}
}

// logError logs the error with appropriate level and context
func logError(err error, c echo.Context, traceID string) {
	// In production, you would use a proper logging framework
	// For now, we'll use the standard log package
	log.Printf("ERROR [%s]: %v | Method: %s | Path: %s | UserAgent: %s | IP: %s",
		traceID,
		err,
		c.Request().Method,
		c.Request().URL.Path,
		c.Request().UserAgent(),
		c.RealIP(),
	)
}

// sanitizeHTTPErrorMessage extracts a safe message from error without leaking internals
func sanitizeHTTPErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if he, ok := err.(*echo.HTTPError); ok && he != nil {
		if msg, ok := he.Message.(string); ok {
			return msg
		}
		// For non-string messages, avoid serializing full struct
		return http.StatusText(he.Code)
	}
	// Generic fallback text
	return http.StatusText(http.StatusInternalServerError)
}
