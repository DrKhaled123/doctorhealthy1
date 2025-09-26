package handlers

import (
	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

// RegisterEnhancedHealthRoutes registers all enhanced health routes
func RegisterEnhancedHealthRoutes(e *echo.Echo, enhancedService *services.EnhancedHealthService, apiKeyService *services.APIKeyService) {
	handler := NewEnhancedHealthHandler(enhancedService, apiKeyService)

	// Create API group with authentication middleware
	api := e.Group("/api/v1/enhanced")

	// Enhanced Health Data Routes
	api.GET("/all", handler.GetAllData)
	api.GET("/status", handler.GetHealthStatus)

	// Diet and Meals Routes
	api.POST("/diet-plan/generate", handler.GenerateDietPlan)
	api.POST("/diet-plan/pdf", handler.GenerateDietPlanPDF)

	// Workout Routes
	api.POST("/workout-plan/generate", handler.GenerateWorkoutPlan)
	api.POST("/workout-plan/pdf", handler.GenerateWorkoutPlanPDF)

	// Lifestyle Management Routes
	api.POST("/lifestyle-plan/generate", handler.GenerateLifestylePlan)
	api.POST("/lifestyle-plan/pdf", handler.GenerateLifestylePlanPDF)

	// Recipes and Recommendations Routes
	api.POST("/recipes/generate", handler.GenerateRecipes)
	api.POST("/recipes/pdf", handler.GenerateRecipesPDF)
	api.GET("/cooking-skills", handler.GetCookingSkills)
	api.GET("/meal-prep-tips", handler.GetMealPrepTips)

	// Injury Management Routes
	api.GET("/injuries/:type", handler.GetInjuryManagement)
	api.GET("/injuries", handler.GetInjuryManagement)

	// Supplement Guidance Routes
	api.GET("/supplements", handler.GetSupplementsGuidance)

	// Additional utility routes
	api.GET("/help", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Enhanced Health API Help",
			"endpoints": map[string]interface{}{
				"diet_planning": map[string]interface{}{
					"generate":    "POST /api/v1/enhanced/diet-plan/generate",
					"pdf":         "POST /api/v1/enhanced/diet-plan/pdf",
					"description": "Generate personalized diet plans with meal recommendations",
				},
				"workout_planning": map[string]interface{}{
					"generate":    "POST /api/v1/enhanced/workout-plan/generate",
					"pdf":         "POST /api/v1/enhanced/workout-plan/pdf",
					"description": "Generate personalized workout plans with exercise recommendations",
				},
				"lifestyle_management": map[string]interface{}{
					"generate":    "POST /api/v1/enhanced/lifestyle-plan/generate",
					"pdf":         "POST /api/v1/enhanced/lifestyle-plan/pdf",
					"description": "Generate lifestyle management plans for diseases and health conditions",
				},
				"recipes_recommendations": map[string]interface{}{
					"generate":       "POST /api/v1/enhanced/recipes/generate",
					"pdf":            "POST /api/v1/enhanced/recipes/pdf",
					"cooking_skills": "GET /api/v1/enhanced/cooking-skills",
					"meal_prep_tips": "GET /api/v1/enhanced/meal-prep-tips",
					"description":    "Generate personalized recipes with cooking guidance",
				},
				"injury_management": map[string]interface{}{
					"all_injuries":    "GET /api/v1/enhanced/injuries",
					"specific_injury": "GET /api/v1/enhanced/injuries/{type}",
					"description":     "Get injury management and treatment protocols",
				},
				"supplement_guidance": map[string]interface{}{
					"guidance":    "GET /api/v1/enhanced/supplements",
					"description": "Get personalized supplement recommendations",
				},
			},
			"authentication": map[string]interface{}{
				"method": "API Key required",
				"header": "X-API-Key: your_api_key",
				"query":  "?api_key=your_api_key",
			},
			"languages": []string{"en", "ar"},
			"features": []string{
				"Personalized diet plan generation",
				"Custom workout plan creation",
				"Lifestyle and disease management",
				"Recipe recommendations with cooking skills",
				"Injury management protocols",
				"Supplement guidance",
				"PDF generation for all plans",
				"Bilingual support (English/Arabic)",
			},
		})
	})
}
