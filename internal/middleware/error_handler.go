package middleware

import (
	"net/http"
	"time"

	"api-key-generator/internal/models"

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

		// Map status code
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok && he != nil {
			if he.Code > 0 {
				code = he.Code
			}
		}

		// Build standardized error response
		resp := models.ErrorResponse{
			Error:     http.StatusText(code),
			Message:   sanitizeHTTPErrorMessage(err),
			Code:      httpStatusCodeString(code),
			Timestamp: time.Now().UTC(),
		}

		// Attempt JSON response; on failure, fall back to default handler
		if jsonErr := c.JSON(code, resp); jsonErr != nil {
			c.Echo().DefaultHTTPErrorHandler(err, c)
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

