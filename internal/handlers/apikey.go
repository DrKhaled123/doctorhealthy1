package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"api-key-generator/internal/models"
	"api-key-generator/internal/services"
	"api-key-generator/internal/utils"

	"github.com/labstack/echo/v4"
)

// APIKeyHandler handles API key related HTTP requests
type APIKeyHandler struct {
	service *services.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(service *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		service: service,
	}
}

// CreateAPIKey creates a new API key
func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	// Extract and validate input
	var req models.CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Call service layer
	apiKey, err := h.service.CreateAPIKey(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	// Sanitize response data
	apiKey.Name = utils.SanitizeForJSON(apiKey.Name)
	if apiKey.Description != nil {
		sanitized := utils.SanitizeForJSON(*apiKey.Description)
		apiKey.Description = &sanitized
	}

	// Return response
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "API key created successfully",
		"api_key": apiKey,
	})
}

// GetAPIKey retrieves an API key by ID
func (h *APIKeyHandler) GetAPIKey(c echo.Context) error {
	// Extract ID from path
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key ID is required")
	}

	// Call service layer
	apiKey, err := h.service.GetAPIKey(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(err)
	}

	// Mask the key for security
	maskAPIKey(apiKey)

	// Sanitize response data
	apiKey.Name = utils.SanitizeForJSON(apiKey.Name)
	if apiKey.Description != nil {
		sanitized := utils.SanitizeForJSON(*apiKey.Description)
		apiKey.Description = &sanitized
	}

	// Return response
	return c.JSON(http.StatusOK, apiKey)
}

// ListAPIKeys lists API keys with pagination and filters
func (h *APIKeyHandler) ListAPIKeys(c echo.Context) error {
	// Parse query parameters
	var params models.ListAPIKeysParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	// Parse is_active parameter
	if isActiveStr := c.QueryParam("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	// Call service layer
	result, err := h.service.ListAPIKeys(c.Request().Context(), &params)
	if err != nil {
		return handleServiceError(err)
	}

	// Return response
	return c.JSON(http.StatusOK, result)
}

// UpdateAPIKey updates an API key
func (h *APIKeyHandler) UpdateAPIKey(c echo.Context) error {
	// Extract ID from path
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key ID is required")
	}

	// Extract and validate input
	var req models.UpdateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Call service layer
	apiKey, err := h.service.UpdateAPIKey(c.Request().Context(), id, &req)
	if err != nil {
		return handleServiceError(err)
	}

	// Mask the key for security
	maskAPIKey(apiKey)

	// Sanitize response data
	apiKey.Name = utils.SanitizeForJSON(apiKey.Name)
	if apiKey.Description != nil {
		sanitized := utils.SanitizeForJSON(*apiKey.Description)
		apiKey.Description = &sanitized
	}

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "API key updated successfully",
		"api_key": apiKey,
	})
}

// DeleteAPIKey deletes an API key
func (h *APIKeyHandler) DeleteAPIKey(c echo.Context) error {
	// Extract ID from path
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key ID is required")
	}

	// Call service layer
	err := h.service.DeleteAPIKey(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(err)
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"message": "API key deleted successfully",
	})
}

// GetPermissions returns available permissions
func (h *APIKeyHandler) GetPermissions(c echo.Context) error {
	permissions := h.service.GetAvailablePermissions()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"permissions": permissions,
	})
}

// RenewAPIKey extends the API key expiry
func (h *APIKeyHandler) RenewAPIKey(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key ID is required")
	}

	// Optional extend_days query param; when absent or <=0, default policy applies
	extendDays := 0
	if v := c.QueryParam("extend_days"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			extendDays = parsed
		}
	}

	updated, err := h.service.RenewAPIKey(c.Request().Context(), id, extendDays)
	if err != nil {
		return handleServiceError(err)
	}

	maskAPIKey(updated)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "API key renewed",
		"api_key":     updated,
		"extend_days": extendDays,
	})
}

// ValidateAPIKey validates an API key (for testing purposes)
func (h *APIKeyHandler) ValidateAPIKey(c echo.Context) error {
	// Extract API key from header
	apiKey := c.Request().Header.Get("X-API-Key")
	if apiKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "API key is required")
	}

	// Validate API key
	key, err := h.service.GetAPIKeyByKey(c.Request().Context(), apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
	}

	// Record usage
	err = h.service.RecordUsage(
		c.Request().Context(),
		key.ID,
		c.Request().URL.Path,
		c.Request().Method,
		http.StatusOK,
		c.RealIP(),
		c.Request().UserAgent(),
	)
	if err != nil {
		// Log error but don't fail the request (sanitized)
		c.Logger().Error("Failed to record API key usage:", utils.SanitizeForLog(err.Error()))
	}

	// Return validation result
	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid":       true,
		"api_key_id":  key.ID,
		"name":        utils.SanitizeForJSON(key.Name),
		"permissions": key.Permissions,
		"expires_at":  key.ExpiresAt,
	})
}

// RegisterAPIKeyRoutes registers all API key routes
func RegisterAPIKeyRoutes(g *echo.Group, service *services.APIKeyService) {
	handler := NewAPIKeyHandler(service)

	// Bootstrap creation is allowed without existing admin key but requires X-Bootstrap-Token
	g.POST("/api-keys", func(c echo.Context) error {
		// If any key exists, defer to normal handler (requires auth elsewhere)
		has, err := service.HasAnyKeys(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
		if !has {
			// Validate bootstrap token
			token := c.Request().Header.Get("X-Bootstrap-Token")
			if token == "" || token != utils.SanitizeForLog(utils.GetEnvOrDefault("BOOTSTRAP_TOKEN", "")) {
				return echo.NewHTTPError(http.StatusUnauthorized, "bootstrap token required")
			}
		}
		return handler.CreateAPIKey(c)
	})
	g.GET("/api-keys", handler.ListAPIKeys)
	g.GET("/api-keys/:id", handler.GetAPIKey)
	g.PUT("/api-keys/:id", handler.UpdateAPIKey)
	g.DELETE("/api-keys/:id", handler.DeleteAPIKey)
	g.POST("/api-keys/:id/renew", handler.RenewAPIKey)

	// Utility routes
	g.GET("/permissions", handler.GetPermissions)
	g.POST("/validate", handler.ValidateAPIKey)
}

// handleServiceError converts service errors to HTTP errors
func handleServiceError(err error) error {
	errMsg := err.Error()

	switch {
	case contains(errMsg, "not found"):
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	case contains(errMsg, "validation failed"):
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	case contains(errMsg, "invalid"):
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	case contains(errMsg, "expired"):
		return echo.NewHTTPError(http.StatusUnauthorized, "API key has expired")
	case contains(errMsg, "inactive"):
		return echo.NewHTTPError(http.StatusUnauthorized, "API key is inactive")
	default:
		// Log internal error but don't expose details
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}

// maskAPIKey masks an API key for security
func maskAPIKey(apiKey *models.APIKey) {
	if len(apiKey.Key) > 12 {
		apiKey.Key = apiKey.Key[:8] + "..." + apiKey.Key[len(apiKey.Key)-4:]
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
