package handlers

import (
	"context"
	"net/http"
	"strings"

	"api-key-generator/internal/models"
	"api-key-generator/internal/utils"

	"github.com/labstack/echo/v4"
)

const (
	headerAPIKey               = "X-API-Key"                // nosec G101 - This is a standard HTTP header name, not a credential
	msgInvalidAPIKey           = "invalid API key"          // nosec G101 - This is an error message, not a credential
	msgInsufficientPermissions = "insufficient permissions" // nosec G101 - This is an error message, not a credential
	permSupplementsWrite       = "supplements:write"
	permSupplementsRead        = "supplements:read"
)

// APIKeyAuthorizer defines the minimal contract required from API key service
type APIKeyAuthorizer interface {
	AuthorizeAny(ctx context.Context, rawKey string, requiredPermissions []string) (*models.APIKey, bool, error)
	RecordUsage(ctx context.Context, apiKeyID, endpoint, method string, status int, ipAddress, userAgent string) error
}

// SupplementHandler handles supplement protocol HTTP requests
type SupplementHandler struct {
	supplementService SupplementServiceContract
	apiKeyService     APIKeyAuthorizer
}

// NewSupplementHandler creates a new supplement handler
func NewSupplementHandler(supplementService SupplementServiceContract, apiKeyService APIKeyAuthorizer) *SupplementHandler {
	return &SupplementHandler{
		supplementService: supplementService,
		apiKeyService:     apiKeyService,
	}
}

// SupplementServiceContract defines the minimal contract required from supplement service
type SupplementServiceContract interface {
	GenerateSupplementProtocol(ctx context.Context, req *models.SupplementRequest) (*models.SupplementProtocol, error)
	GetSupplementProtocol(ctx context.Context, id string) (*models.SupplementProtocol, error)
}

// GenerateSupplementProtocol generates personalized supplement recommendations
func (h *SupplementHandler) GenerateSupplementProtocol(c echo.Context) error {
	// Authorize API key with required permissions
	apiKey, ok, err := h.apiKeyService.AuthorizeAny(c.Request().Context(), c.Request().Header.Get(headerAPIKey), []string{permSupplementsWrite})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, msgInvalidAPIKey)
	}
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, msgInsufficientPermissions)
	}

	// Parse and validate request
	var req models.SupplementRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Sanitize inputs
	req.UserID = utils.SanitizeInput(req.UserID)
	req.WorkoutGoal = utils.SanitizeInput(req.WorkoutGoal)
	req.WorkoutType = utils.SanitizeInput(req.WorkoutType)
	req.Intensity = utils.SanitizeInput(req.Intensity)

	// Generate supplement protocol
	protocol, err := h.supplementService.GenerateSupplementProtocol(c.Request().Context(), &req)
	if err != nil {
		c.Logger().Error("Failed to generate supplement protocol:", utils.SanitizeForLog(err.Error()))
		return h.handleServiceError(err)
	}

	// Record usage
	h.recordUsage(c, apiKey.ID, http.StatusOK)

	return c.JSON(http.StatusOK, h.createSuccessResponse(protocol, "Supplement protocol generated successfully"))
}

// GetSupplementProtocol retrieves a supplement protocol by ID
func (h *SupplementHandler) GetSupplementProtocol(c echo.Context) error {
	// Authorize API key with required permissions
	apiKey, ok, err := h.apiKeyService.AuthorizeAny(c.Request().Context(), c.Request().Header.Get(headerAPIKey), []string{permSupplementsRead})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, msgInvalidAPIKey)
	}
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, msgInsufficientPermissions)
	}

	// Get protocol ID from path
	protocolID := utils.SanitizeInput(c.Param("id"))
	if protocolID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "protocol ID is required")
	}

	// Get supplement protocol
	protocol, err := h.supplementService.GetSupplementProtocol(c.Request().Context(), protocolID)
	if err != nil {
		c.Logger().Error("Failed to get supplement protocol:", utils.SanitizeForLog(err.Error()))
		return h.handleServiceError(err)
	}

	// Record usage
	h.recordUsage(c, apiKey.ID, http.StatusOK)

	return c.JSON(http.StatusOK, protocol)
}

// GetSupplementCategories returns available supplement categories
func (h *SupplementHandler) GetSupplementCategories(c echo.Context) error {
	// Authorize API key with required permissions
	apiKey, ok, err := h.apiKeyService.AuthorizeAny(c.Request().Context(), c.Request().Header.Get(headerAPIKey), []string{permSupplementsRead})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, msgInvalidAPIKey)
	}
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, msgInsufficientPermissions)
	}

	// Record usage
	h.recordUsage(c, apiKey.ID, http.StatusOK)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": models.SupplementCategories,
	})
}

// GetSupplementSafetyInfo returns safety information for supplements
func (h *SupplementHandler) GetSupplementSafetyInfo(c echo.Context) error {
	// Authorize API key with required permissions
	apiKey, ok, err := h.apiKeyService.AuthorizeAny(c.Request().Context(), c.Request().Header.Get("X-API-Key"), []string{"supplements:read"})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
	}
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	supplementName := utils.SanitizeInput(c.Param("name"))
	if supplementName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "supplement name is required")
	}

	safetyInfo := h.getSupplementSafetyInfo(supplementName)

	// Record usage
	h.recordUsage(c, apiKey.ID, http.StatusOK)

	return c.JSON(http.StatusOK, safetyInfo)
}

// getSupplementSafetyInfo returns safety information for specific supplements
func (h *SupplementHandler) getSupplementSafetyInfo(name string) map[string]interface{} {
	safetyData := map[string]map[string]interface{}{
		"caffeine": {
			"max_daily_dose":      "400mg",
			"contraindications":   []string{"Pregnancy", "Heart conditions", "Anxiety disorders"},
			"side_effects":        []string{"Jitters", "Insomnia", "Increased heart rate"},
			"interactions":        []string{"Blood thinners", "Stimulant medications"},
			"timing_restrictions": "Avoid 6 hours before bedtime",
		},
		"creatine": {
			"max_daily_dose":      "10g",
			"contraindications":   []string{"Kidney disease", "Liver disease"},
			"side_effects":        []string{"Water retention", "Digestive issues"},
			"interactions":        []string{"Nephrotoxic drugs"},
			"timing_restrictions": "Take with plenty of water",
		},
		"whey_protein": {
			"max_daily_dose":      "50g per serving",
			"contraindications":   []string{"Milk allergy", "Lactose intolerance"},
			"side_effects":        []string{"Digestive upset", "Bloating"},
			"interactions":        []string{"None significant"},
			"timing_restrictions": "Best within 30 minutes post-workout",
		},
	}

	if info, exists := safetyData[name]; exists {
		return info
	}

	return map[string]interface{}{
		"message": "Safety information not available for this supplement",
	}
}

// recordUsage records API key usage
func (h *SupplementHandler) recordUsage(c echo.Context, apiKeyID string, status int) {
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

// handleServiceError converts service errors to HTTP errors
func (h *SupplementHandler) handleServiceError(err error) error {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "not found"):
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	case strings.Contains(errMsg, "validation"):
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed")
	case strings.Contains(errMsg, "invalid"):
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}

// createSuccessResponse creates standardized success response
func (h *SupplementHandler) createSuccessResponse(data interface{}, message string) models.SuccessResponse {
	return models.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}
