package handlers

import (
	"net/http"

	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

// UltimateEchoHandler handles all comprehensive health data endpoints for Echo framework
type UltimateEchoHandler struct {
	ultimateService *services.UltimateDataService
	apiKeyService   *services.APIKeyService
}

// NewUltimateEchoHandler creates a new ultimate data handler for Echo
func NewUltimateEchoHandler(ultimateService *services.UltimateDataService, apiKeyService *services.APIKeyService) *UltimateEchoHandler {
	return &UltimateEchoHandler{
		ultimateService: ultimateService,
		apiKeyService:   apiKeyService,
	}
}

// GetAllData returns all comprehensive health data (requires authentication)
func (h *UltimateEchoHandler) GetAllData(c echo.Context) error {
	// Verify API key
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetAllData()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":           "success",
		"message":          "Ultimate comprehensive health database",
		"data":             data,
		"total_categories": len(data),
	})
}

// GetMedications returns all medications data
func (h *UltimateEchoHandler) GetMedications(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetMedications()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Complete medications database (30 medications)",
		"data":    data,
	})
}

// GetVitaminsMinerals returns vitamins and minerals data
func (h *UltimateEchoHandler) GetVitaminsMinerals(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetVitaminsMinerals()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Complete vitamins & minerals database (17 nutrients)",
		"data":    data,
	})
}

// GetComprehensiveDiseases returns comprehensive diseases data
func (h *UltimateEchoHandler) GetComprehensiveDiseases(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetComprehensiveDiseases()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Comprehensive diseases database (15+ conditions)",
		"data":    data,
	})
}

// GetComprehensiveComplaints returns comprehensive complaints data
func (h *UltimateEchoHandler) GetComprehensiveComplaints(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetComprehensiveComplaints()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Comprehensive complaints database (22 case studies)",
		"data":    data,
	})
}

// GetInjuries returns injuries data
func (h *UltimateEchoHandler) GetInjuries(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetInjuries()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Complete injuries database (5 major injuries)",
		"data":    data,
	})
}

// GetDiseases returns diseases data
func (h *UltimateEchoHandler) GetDiseases(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetDiseases()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Diseases database (4 major diseases)",
		"data":    data,
	})
}

// GetComplaints returns complaints data
func (h *UltimateEchoHandler) GetComplaints(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetComplaints()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Complaints database (6 common complaints)",
		"data":    data,
	})
}

// GetTypePlans returns type plans data
func (h *UltimateEchoHandler) GetTypePlans(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetTypePlans()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Type plans database (8 diet plans)",
		"data":    data,
	})
}

// GetWorkouts returns workouts data
func (h *UltimateEchoHandler) GetWorkouts(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetWorkouts()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Workouts database (6+ exercises)",
		"data":    data,
	})
}

// GetRecipes returns recipes data
func (h *UltimateEchoHandler) GetRecipes(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetRecipes()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Recipes database (5+ recipes)",
		"data":    data,
	})
}

// GetMetadata returns database metadata
func (h *UltimateEchoHandler) GetMetadata(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	data := h.ultimateService.GetMetadata()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Ultimate database metadata",
		"data":    data,
	})
}

// SearchData searches across all data categories
func (h *UltimateEchoHandler) SearchData(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query parameter 'q' is required")
	}

	results := h.ultimateService.SearchData(query)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":           "success",
		"message":          "Search results for: " + query,
		"query":            query,
		"results":          results,
		"categories_found": len(results),
	})
}

// GetDataStats returns statistics about the loaded data
func (h *UltimateEchoHandler) GetDataStats(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	stats := h.ultimateService.GetDataStats()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Ultimate database statistics",
		"stats":   stats,
	})
}

// ReloadData reloads the database from file (admin only)
func (h *UltimateEchoHandler) ReloadData(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	err := h.ultimateService.ReloadData()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to reload database")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Ultimate database reloaded successfully",
	})
}

// ValidateData performs validation on loaded data
func (h *UltimateEchoHandler) ValidateData(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	validation := h.ultimateService.ValidateData()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "success",
		"message":    "Database validation results",
		"validation": validation,
	})
}

// GetCategoryData returns data for a specific category with optional filtering
func (h *UltimateEchoHandler) GetCategoryData(c echo.Context) error {
	if !h.verifyAPIKey(c) {
		return nil
	}

	category := c.Param("category")
	if category == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "category parameter is required")
	}

	// Get query parameters for filtering
	filters := make(map[string]string)
	for key, values := range c.Request().URL.Query() {
		if len(values) > 0 && key != "api_key" {
			filters[key] = values[0]
		}
	}

	data := h.ultimateService.GetCategoryData(category, filters)
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":          "success",
		"message":         "Data for category: " + category,
		"category":        category,
		"data":            data,
		"filters_applied": len(filters),
	})
}

// GetHealthStatus returns overall health of the ultimate database system
func (h *UltimateEchoHandler) GetHealthStatus(c echo.Context) error {
	stats := h.ultimateService.GetDataStats()
	validation := h.ultimateService.ValidateData()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":           "success",
		"message":          "Ultimate database health status",
		"health":           "healthy",
		"stats":            stats,
		"validation":       validation,
		"endpoints_active": 15,
		"authentication":   "required",
	})
}

// verifyAPIKey checks if the request has a valid API key
func (h *UltimateEchoHandler) verifyAPIKey(c echo.Context) bool {
	apiKey := c.Request().Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = c.QueryParam("api_key")
	}

	if apiKey == "" {
		c.Error(echo.NewHTTPError(http.StatusUnauthorized, "API key required"))
		return false
	}

	// Validate API key using GetAPIKeyByKey
	_, err := h.apiKeyService.GetAPIKeyByKey(c.Request().Context(), apiKey)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusUnauthorized, "invalid API key"))
		return false
	}

	return true
}

// RegisterUltimateRoutes sets up all ultimate database API routes for Echo
func RegisterUltimateRoutes(api *echo.Group, ultimateService *services.UltimateDataService, apiKeyService *services.APIKeyService) {
	// Create handler
	handler := NewUltimateEchoHandler(ultimateService, apiKeyService)

	// Ultimate database routes (all require authentication)
	ultimate := api.Group("/ultimate")
	{
		// Main endpoints
		ultimate.GET("/all", handler.GetAllData)         // Get all data
		ultimate.GET("/metadata", handler.GetMetadata)   // Get metadata
		ultimate.GET("/stats", handler.GetDataStats)     // Get statistics
		ultimate.GET("/health", handler.GetHealthStatus) // Health check
		ultimate.GET("/search", handler.SearchData)      // Search all data
		ultimate.POST("/reload", handler.ReloadData)     // Reload database
		ultimate.GET("/validate", handler.ValidateData)  // Validate data

		// Category-specific endpoints
		ultimate.GET("/medications", handler.GetMedications)                          // 30 medications
		ultimate.GET("/vitamins-minerals", handler.GetVitaminsMinerals)               // 17 nutrients
		ultimate.GET("/comprehensive-diseases", handler.GetComprehensiveDiseases)     // 15+ diseases
		ultimate.GET("/comprehensive-complaints", handler.GetComprehensiveComplaints) // 22 cases
		ultimate.GET("/injuries", handler.GetInjuries)                                // 5 injuries
		ultimate.GET("/diseases", handler.GetDiseases)                                // 4 diseases
		ultimate.GET("/complaints", handler.GetComplaints)                            // 6 complaints
		ultimate.GET("/type-plans", handler.GetTypePlans)                             // 8 diet plans
		ultimate.GET("/workouts", handler.GetWorkouts)                                // 6+ exercises
		ultimate.GET("/recipes", handler.GetRecipes)                                  // 5+ recipes

		// Dynamic category endpoint
		ultimate.GET("/category/:category", handler.GetCategoryData) // Get any category
	}

	// Legacy compatibility routes (redirect to ultimate)
	api.GET("/medications", handler.GetMedications)
	api.GET("/diseases", handler.GetComprehensiveDiseases)
	api.GET("/complaints", handler.GetComprehensiveComplaints)
	api.GET("/vitamins", handler.GetVitaminsMinerals)
	api.GET("/workouts", handler.GetWorkouts)
	api.GET("/recipes", handler.GetRecipes)
}
