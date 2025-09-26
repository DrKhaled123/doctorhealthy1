package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"api-key-generator/internal/services"
	"api-key-generator/internal/utils"
)

// ScopesAny enforces that the request's API key has at least one of the required permissions.
// It expects the raw API key in header "X-API-Key". On success, it stores the resolved API key
// record in context key "api_key_record" for downstream handlers.
func ScopesAny(apiKeyService *services.APIKeyService, requiredPermissions ...string) echo.MiddlewareFunc {
	sanitizedRequired := make([]string, 0, len(requiredPermissions))
	for _, p := range requiredPermissions {
		sanitizedRequired = append(sanitizedRequired, utils.SanitizeInput(p))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rawKey := c.Request().Header.Get("X-API-Key")
			if rawKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
			}

			keyRecord, ok, err := apiKeyService.AuthorizeAny(c.Request().Context(), rawKey, sanitizedRequired)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			c.Set("api_key_record", keyRecord)
			return next(c)
		}
	}
}

// ScopesAll enforces that the request's API key has all required permissions.
// It expects the raw API key in header "X-API-Key". On success, it stores the resolved API key
// record in context key "api_key_record" for downstream handlers.
func ScopesAll(apiKeyService *services.APIKeyService, requiredPermissions ...string) echo.MiddlewareFunc {
	sanitizedRequired := make([]string, 0, len(requiredPermissions))
	for _, p := range requiredPermissions {
		sanitizedRequired = append(sanitizedRequired, utils.SanitizeInput(p))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rawKey := c.Request().Header.Get("X-API-Key")
			if rawKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
			}

			keyRecord, ok, err := apiKeyService.AuthorizeAll(c.Request().Context(), rawKey, sanitizedRequired)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
			}
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			c.Set("api_key_record", keyRecord)
			return next(c)
		}
	}
}
