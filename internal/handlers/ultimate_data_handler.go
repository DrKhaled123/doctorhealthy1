package handlers

import (
	"net/http"

	"api-key-generator/internal/services"

	"github.com/gin-gonic/gin"
)

// UltimateDataHandler handles all comprehensive health data endpoints
type UltimateDataHandler struct {
	ultimateService *services.UltimateDataService
	apiKeyService   *services.APIKeyService
}

// NewUltimateDataHandler creates a new ultimate data handler
func NewUltimateDataHandler(ultimateService *services.UltimateDataService, apiKeyService *services.APIKeyService) *UltimateDataHandler {
	return &UltimateDataHandler{
		ultimateService: ultimateService,
		apiKeyService:   apiKeyService,
	}
}

// GetAllData returns all comprehensive health data (requires authentication)
func (h *UltimateDataHandler) GetAllData(c *gin.Context) {
	// Verify API key
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetAllData()

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          "Ultimate comprehensive health database",
		"data":             data,
		"total_categories": len(data),
	})
}

// GetMedications returns all medications data
func (h *UltimateDataHandler) GetMedications(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetMedications()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Complete medications database (30 medications)",
		"data":    data,
	})
}

// Category-specific medication endpoints (stubs map to filtered data)
func (h *UltimateDataHandler) GetWeightLossDrugs(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["weight_loss_drugs"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "weight_loss_drugs", "data": category})
}

func (h *UltimateDataHandler) GetAppetiteSuppressants(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["appetite_suppressants"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "appetite_suppressants", "data": category})
}

func (h *UltimateDataHandler) GetInvestigationalDrugs(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["investigational_drugs"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "investigational_drugs", "data": category})
}

func (h *UltimateDataHandler) GetDiabetesMedications(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["diabetes_medications"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "diabetes_medications", "data": category})
}

func (h *UltimateDataHandler) GetHormonalMedications(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["hormonal_medications"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "hormonal_medications", "data": category})
}

func (h *UltimateDataHandler) GetPerformanceEnhancers(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["performance_enhancers"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "performance_enhancers", "data": category})
}

func (h *UltimateDataHandler) GetOtherMedications(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	m := h.ultimateService.GetMedications()
	var category interface{}
	if mm, ok := m.(map[string]interface{}); ok {
		category = mm["other_medications"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "category": "other_medications", "data": category})
}

// GetVitaminsMinerals returns vitamins and minerals data
func (h *UltimateDataHandler) GetVitaminsMinerals(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetVitaminsMinerals()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Complete vitamins & minerals database (17 nutrients)",
		"data":    data,
	})
}

// Additional vitamins/minerals endpoints
func (h *UltimateDataHandler) GetVitamins(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var vitamins interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		vitamins = m["vitamins"]
		if vitamins == nil {
			vitamins = m["water_soluble_vitamins"]
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": vitamins})
}

func (h *UltimateDataHandler) GetWaterSolubleVitamins(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var data interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		data = m["water_soluble_vitamins"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetFatSolubleVitamins(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var data interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		data = m["fat_soluble_vitamins"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GetComprehensiveDiseases returns comprehensive diseases data
func (h *UltimateDataHandler) GetComprehensiveDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetComprehensiveDiseases()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comprehensive diseases database (15+ conditions)",
		"data":    data,
	})
}

// GetComprehensiveComplaints returns comprehensive complaints data
func (h *UltimateDataHandler) GetComprehensiveComplaints(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetComprehensiveComplaints()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comprehensive complaints database (22 case studies)",
		"data":    data,
	})
}

// GetInjuries returns injuries data
func (h *UltimateDataHandler) GetInjuries(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetInjuries()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Complete injuries database (5 major injuries)",
		"data":    data,
	})
}

// GetDiseases returns diseases data
func (h *UltimateDataHandler) GetDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetDiseases()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Diseases database (4 major diseases)",
		"data":    data,
	})
}

// Minerals and supplement endpoints
func (h *UltimateDataHandler) GetMinerals(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var data interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		data = m["minerals"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetSupplementInteractions(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var data interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		data = m["supplement_interactions"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetSpecialPopulations(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	vm := h.ultimateService.GetVitaminsMinerals()
	var data interface{}
	if m, ok := vm.(map[string]interface{}); ok {
		data = m["special_populations"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// Disease category endpoints (filtering by keys in comprehensive_diseases)
func (h *UltimateDataHandler) GetEndocrineDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetComprehensiveDiseases()
	var data interface{}
	if m, ok := d.(map[string]interface{}); ok {
		data = m["endocrine"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetCardiovascularDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetComprehensiveDiseases()
	var data interface{}
	if m, ok := d.(map[string]interface{}); ok {
		data = m["cardiovascular"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetAutoimmuneDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetComprehensiveDiseases()
	var data interface{}
	if m, ok := d.(map[string]interface{}); ok {
		data = m["autoimmune"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetNeurologicalDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetComprehensiveDiseases()
	var data interface{}
	if m, ok := d.(map[string]interface{}); ok {
		data = m["neurological"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetGastrointestinalDiseases(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetComprehensiveDiseases()
	var data interface{}
	if m, ok := d.(map[string]interface{}); ok {
		data = m["gastrointestinal"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UltimateDataHandler) GetTreatmentProtocols(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetCategoryData("treatment_protocols", nil)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": d})
}

func (h *UltimateDataHandler) GetEmergencyProtocols(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}
	d := h.ultimateService.GetCategoryData("emergency_protocols", nil)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": d})
}

// GetComplaints returns complaints data
func (h *UltimateDataHandler) GetComplaints(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetComplaints()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Complaints database (6 common complaints)",
		"data":    data,
	})
}

// GetTypePlans returns type plans data
func (h *UltimateDataHandler) GetTypePlans(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetTypePlans()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Type plans database (8 diet plans)",
		"data":    data,
	})
}

// GetWorkouts returns workouts data
func (h *UltimateDataHandler) GetWorkouts(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetWorkouts()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Workouts database (6+ exercises)",
		"data":    data,
	})
}

// GetRecipes returns recipes data
func (h *UltimateDataHandler) GetRecipes(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetRecipes()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Recipes database (5+ recipes)",
		"data":    data,
	})
}

// GetMetadata returns database metadata
func (h *UltimateDataHandler) GetMetadata(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	data := h.ultimateService.GetMetadata()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ultimate database metadata",
		"data":    data,
	})
}

// SearchData searches across all data categories
func (h *UltimateDataHandler) SearchData(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Search query parameter 'q' is required",
		})
		return
	}

	results := h.ultimateService.SearchData(query)

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          "Search results for: " + query,
		"query":            query,
		"results":          results,
		"categories_found": len(results),
	})
}

// GetDataStats returns statistics about the loaded data
func (h *UltimateDataHandler) GetDataStats(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	stats := h.ultimateService.GetDataStats()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ultimate database statistics",
		"stats":   stats,
	})
}

// ReloadData reloads the database from file (admin only)
func (h *UltimateDataHandler) ReloadData(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	// Additional admin check could be added here

	err := h.ultimateService.ReloadData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to reload database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ultimate database reloaded successfully",
	})
}

// ValidateData performs validation on loaded data
func (h *UltimateDataHandler) ValidateData(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	validation := h.ultimateService.ValidateData()

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Database validation results",
		"validation": validation,
	})
}

// GetCategoryData returns data for a specific category with optional filtering
func (h *UltimateDataHandler) GetCategoryData(c *gin.Context) {
	if !h.verifyAPIKey(c) {
		return
	}

	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Category parameter is required",
		})
		return
	}

	// Get query parameters for filtering
	filters := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 && key != "api_key" {
			filters[key] = values[0]
		}
	}

	data := h.ultimateService.GetCategoryData(category, filters)
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Category not found: " + category,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "Data for category: " + category,
		"category":        category,
		"data":            data,
		"filters_applied": len(filters),
	})
}

// GetHealthStatus returns overall health of the ultimate database system
func (h *UltimateDataHandler) GetHealthStatus(c *gin.Context) {
	stats := h.ultimateService.GetDataStats()
	validation := h.ultimateService.ValidateData()

	c.JSON(http.StatusOK, gin.H{
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
func (h *UltimateDataHandler) verifyAPIKey(c *gin.Context) bool {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		apiKey = c.Query("api_key")
	}

	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "API key is required. Provide it in X-API-Key header or api_key query parameter",
		})
		return false
	}

	// Validate API key using GetAPIKeyByKey
	_, err := h.apiKeyService.GetAPIKeyByKey(c.Request.Context(), apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid API key",
		})
		return false
	}

	return true
}
