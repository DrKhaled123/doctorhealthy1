package handlers

import (
	"net/http"
	"strings"

	"api-key-generator/internal/utils"

	"github.com/labstack/echo/v4"
)

// (validateAPIKey) removed; authorization is handled inline in handlers now

// hasPermission is retained for compatibility (unused in revised handler)
// func (h *RecipeHandler) hasPermission(apiKey *models.APIKey, permission string) bool {
//     for _, perm := range apiKey.Permissions {
//         if perm == permission || perm == "admin" {
//             return true
//         }
//     }
//     return false
// }

// recordUsage records API key usage for analytics
func (h *RecipeHandler) recordUsage(c echo.Context, apiKeyID string, status int) {
	err := h.apiKeyService.RecordUsage(
		c.Request().Context(),
		apiKeyID,
		c.Request().URL.Path,
		c.Request().Method,
		status,
		c.RealIP(),
		c.Request().UserAgent(),
	)
	if err != nil {
		c.Logger().Error("Failed to record usage:", utils.SanitizeForLog(err.Error()))
	}
}

// handleServiceError converts service errors to proper HTTP responses
func (h *RecipeHandler) handleServiceError(err error) error {
	errMsg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errMsg, "not found"):
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	case strings.Contains(errMsg, "validation"):
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed")
	case strings.Contains(errMsg, "invalid"):
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	case strings.Contains(errMsg, "timeout"):
		return echo.NewHTTPError(http.StatusRequestTimeout, "request timeout")
	case strings.Contains(errMsg, "rate limit"):
		return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}

// sanitizeStringSlice sanitizes a slice of strings
func (h *RecipeHandler) sanitizeStringSlice(slice []string) []string {
	var sanitized []string
	for _, s := range slice {
		if clean := utils.SanitizeInput(strings.TrimSpace(s)); clean != "" {
			sanitized = append(sanitized, clean)
		}
	}
	return sanitized
}

// (validatePagination) removed; pagination handled inline where needed

// (createSuccessResponse) removed; standard responses used directly

// (createErrorResponse) removed; centralized HTTPError handler used
