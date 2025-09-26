package handlers

import (
	"api-key-generator/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupUltimateRoutes sets up all ultimate database API routes
func SetupUltimateRoutes(router *gin.Engine, ultimateService *services.UltimateDataService, apiKeyService *services.APIKeyService) {
	// Create handler
	handler := NewUltimateDataHandler(ultimateService, apiKeyService)

	// API v1 group
	v1 := router.Group("/api/v1")

	// Ultimate database routes (all require authentication)
	ultimate := v1.Group("/ultimate")
	{
		// Main endpoints
		ultimate.GET("/all", handler.GetAllData)         // Get all data
		ultimate.GET("/metadata", handler.GetMetadata)   // Get metadata
		ultimate.GET("/stats", handler.GetDataStats)     // Get statistics
		ultimate.GET("/health", handler.GetHealthStatus) // Health check
		ultimate.GET("/search", handler.SearchData)      // Search all data
		ultimate.POST("/reload", handler.ReloadData)     // Reload database
		ultimate.GET("/validate", handler.ValidateData)  // Validate data

		// Complete medication endpoints
		ultimate.GET("/medications", handler.GetMedications)                                // All 41 medications
		ultimate.GET("/medications/weight-loss", handler.GetWeightLossDrugs)                // 7 weight loss drugs
		ultimate.GET("/medications/appetite-suppressants", handler.GetAppetiteSuppressants) // 6 appetite suppressants
		ultimate.GET("/medications/investigational", handler.GetInvestigationalDrugs)       // 9 investigational drugs
		ultimate.GET("/medications/diabetes", handler.GetDiabetesMedications)               // 4 diabetes medications
		ultimate.GET("/medications/hormonal", handler.GetHormonalMedications)               // 3 hormonal medications
		ultimate.GET("/medications/performance", handler.GetPerformanceEnhancers)           // 3 performance enhancers
		ultimate.GET("/medications/other", handler.GetOtherMedications)                     // 9 other medications

		// Complete vitamins and minerals endpoints
		ultimate.GET("/vitamins", handler.GetVitamins)                              // All 14 vitamins
		ultimate.GET("/vitamins/water-soluble", handler.GetWaterSolubleVitamins)    // 10 water-soluble vitamins
		ultimate.GET("/vitamins/fat-soluble", handler.GetFatSolubleVitamins)        // 4 fat-soluble vitamins
		ultimate.GET("/minerals", handler.GetMinerals)                              // All 7 minerals
		ultimate.GET("/supplement-interactions", handler.GetSupplementInteractions) // All interactions
		ultimate.GET("/special-populations", handler.GetSpecialPopulations)         // Special population data

		// Complete diseases endpoints
		ultimate.GET("/diseases", handler.GetDiseases)                                  // All 85+ diseases
		ultimate.GET("/diseases/endocrine", handler.GetEndocrineDiseases)               // Endocrine diseases
		ultimate.GET("/diseases/cardiovascular", handler.GetCardiovascularDiseases)     // Cardiovascular diseases
		ultimate.GET("/diseases/autoimmune", handler.GetAutoimmuneDiseases)             // Autoimmune diseases
		ultimate.GET("/diseases/neurological", handler.GetNeurologicalDiseases)         // Neurological diseases
		ultimate.GET("/diseases/gastrointestinal", handler.GetGastrointestinalDiseases) // GI diseases
		ultimate.GET("/treatment-protocols", handler.GetTreatmentProtocols)             // Treatment protocols
		ultimate.GET("/emergency-protocols", handler.GetEmergencyProtocols)             // Emergency protocols

		// Legacy endpoints (maintained for compatibility)
		ultimate.GET("/vitamins-minerals", handler.GetVitaminsMinerals)               // Legacy vitamins/minerals
		ultimate.GET("/comprehensive-diseases", handler.GetComprehensiveDiseases)     // Legacy diseases
		ultimate.GET("/comprehensive-complaints", handler.GetComprehensiveComplaints) // Legacy complaints
		ultimate.GET("/injuries", handler.GetInjuries)                                // 5 injuries
		ultimate.GET("/complaints", handler.GetComplaints)                            // 6 complaints
		ultimate.GET("/type-plans", handler.GetTypePlans)                             // 8 diet plans
		ultimate.GET("/workouts", handler.GetWorkouts)                                // 6+ exercises
		ultimate.GET("/recipes", handler.GetRecipes)                                  // 5+ recipes

		// Dynamic category endpoint
		ultimate.GET("/category/:category", handler.GetCategoryData) // Get any category
	}

	// Legacy compatibility routes (redirect to ultimate)
	v1.GET("/medications", handler.GetMedications)
	v1.GET("/diseases", handler.GetDiseases)
	v1.GET("/complaints", handler.GetComprehensiveComplaints)
	v1.GET("/vitamins", handler.GetVitamins)
	v1.GET("/minerals", handler.GetMinerals)
	v1.GET("/workouts", handler.GetWorkouts)
	v1.GET("/recipes", handler.GetRecipes)

	// Public health endpoint (no auth required)
	v1.GET("/health", handler.GetHealthStatus)
}
