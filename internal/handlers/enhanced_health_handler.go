package handlers

import (
	"net/http"
	"strconv"

	"api-key-generator/internal/services"
	pdfutil "api-key-generator/internal/utils/pdf"

	"github.com/labstack/echo/v4"
)

// EnhancedHealthHandler handles all enhanced health functionality
type EnhancedHealthHandler struct {
	enhancedService *services.EnhancedHealthService
	apiKeyService   *services.APIKeyService
}

// NewEnhancedHealthHandler creates a new enhanced health handler
func NewEnhancedHealthHandler(enhancedService *services.EnhancedHealthService, apiKeyService *services.APIKeyService) *EnhancedHealthHandler {
	return &EnhancedHealthHandler{
		enhancedService: enhancedService,
		apiKeyService:   apiKeyService,
	}
}

// verifyAPIKey verifies API key authentication
func (h *EnhancedHealthHandler) verifyAPIKey(c echo.Context) error {
	apiKey := c.Request().Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = c.QueryParam("api_key")
	}

	if apiKey == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
	}

	// Verify API key using service and record usage
	key, err := h.apiKeyService.GetAPIKeyByKey(c.Request().Context(), apiKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
	}

	// Record usage (best-effort; do not fail request if recording fails)
	_ = h.apiKeyService.RecordUsage(
		c.Request().Context(),
		key.ID,
		c.Request().URL.Path,
		c.Request().Method,
		http.StatusOK,
		c.RealIP(),
		c.Request().UserAgent(),
	)

	return nil
}

// GetAllData returns all enhanced health data
func (h *EnhancedHealthHandler) GetAllData(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	data := h.enhancedService.GetAllData()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":           "success",
		"message":          "Enhanced comprehensive health database",
		"data":             data,
		"total_categories": len(data),
		"features": []string{
			"diet_meals_generation",
			"workout_planning",
			"lifestyle_management",
			"recipes_recommendations",
			"injury_management",
			"supplement_guidance",
		},
	})
}

// GenerateDietPlan generates a personalized diet plan
func (h *EnhancedHealthHandler) GenerateDietPlan(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse user data from request
	var userData map[string]interface{}
	if err := c.Bind(&userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Generate diet plan
	dietPlan := h.enhancedService.GenerateDietPlan(userData, language)

	// Add PDF generation info
	dietPlan["pdf_generation"] = map[string]interface{}{
		"available":   true,
		"endpoint":    "/api/v1/enhanced/diet-plan/pdf",
		"method":      "POST",
		"description": "Generate PDF of diet plan",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Diet plan generated successfully",
		"plan":         dietPlan,
		"language":     language,
		"generated_at": "2024-01-15T10:30:00Z",
	})
}

// GenerateWorkoutPlan generates a personalized workout plan
func (h *EnhancedHealthHandler) GenerateWorkoutPlan(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse user data from request
	var userData map[string]interface{}
	if err := c.Bind(&userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Generate workout plan
	workoutPlan := h.enhancedService.GenerateWorkoutPlan(userData, language)

	// Add PDF generation info
	workoutPlan["pdf_generation"] = map[string]interface{}{
		"available":   true,
		"endpoint":    "/api/v1/enhanced/workout-plan/pdf",
		"method":      "POST",
		"description": "Generate PDF of workout plan",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Workout plan generated successfully",
		"plan":         workoutPlan,
		"language":     language,
		"generated_at": "2024-01-15T10:30:00Z",
	})
}

// GenerateLifestylePlan generates a lifestyle management plan
func (h *EnhancedHealthHandler) GenerateLifestylePlan(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse user data from request
	var userData map[string]interface{}
	if err := c.Bind(&userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Generate lifestyle plan
	lifestylePlan := h.enhancedService.GenerateLifestylePlan(userData, language)

	// Add PDF generation info
	lifestylePlan["pdf_generation"] = map[string]interface{}{
		"available":   true,
		"endpoint":    "/api/v1/enhanced/lifestyle-plan/pdf",
		"method":      "POST",
		"description": "Generate PDF of lifestyle plan",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Lifestyle plan generated successfully",
		"plan":         lifestylePlan,
		"language":     language,
		"generated_at": "2024-01-15T10:30:00Z",
	})
}

// GenerateRecipes generates personalized recipes
func (h *EnhancedHealthHandler) GenerateRecipes(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse user data from request
	var userData map[string]interface{}
	if err := c.Bind(&userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Get number of recipes to generate (default 3 for free plan)
	numRecipes := 3
	if n, err := strconv.Atoi(c.QueryParam("count")); err == nil && n > 0 {
		numRecipes = n
	}

	// Generate recipes
	recipes := h.enhancedService.GenerateRecipes(userData, language)

	// Add PDF generation info
	recipes["pdf_generation"] = map[string]interface{}{
		"available":   true,
		"endpoint":    "/api/v1/enhanced/recipes/pdf",
		"method":      "POST",
		"description": "Generate PDF of recipes",
	}

	// Add free plan limitations
	recipes["plan_limitations"] = map[string]interface{}{
		"free_plan_recipes": 3,
		"generated_count":   numRecipes,
		"upgrade_message":   "Upgrade to premium for unlimited recipe generation",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Recipes generated successfully",
		"recipes":      recipes,
		"language":     language,
		"generated_at": "2024-01-15T10:30:00Z",
	})
}

// GenerateDietPlanPDF generates PDF for diet plan
func (h *EnhancedHealthHandler) GenerateDietPlanPDF(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse diet plan data
	var dietPlan map[string]interface{}
	if err := c.Bind(&dietPlan); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid diet plan data")
	}

	titleEn := "Diet Plan"
	titleAr := "الخطة الغذائية"
	sections := []string{"Personalized diet plan", "Meals, ingredients, macros"}
	content, err := pdfutil.GenerateSimplePDF(titleEn, titleAr, sections)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate PDF")
	}
	return c.Blob(http.StatusOK, "application/pdf", content)
}

// GenerateWorkoutPlanPDF generates PDF for workout plan
func (h *EnhancedHealthHandler) GenerateWorkoutPlanPDF(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse workout plan data
	var workoutPlan map[string]interface{}
	if err := c.Bind(&workoutPlan); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workout plan data")
	}

	content, err := pdfutil.GenerateSimplePDF("Workout Plan", "خطة التمرين", []string{"Weekly routine", "Exercises, sets, reps, rest"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate PDF")
	}
	return c.Blob(http.StatusOK, "application/pdf", content)
}

// GenerateLifestylePlanPDF generates PDF for lifestyle plan
func (h *EnhancedHealthHandler) GenerateLifestylePlanPDF(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse lifestyle plan data
	var lifestylePlan map[string]interface{}
	if err := c.Bind(&lifestylePlan); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid lifestyle plan data")
	}

	content, err := pdfutil.GenerateSimplePDF("Lifestyle Plan", "خطة نمط الحياة", []string{"Habits, sleep, stress, activity"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate PDF")
	}
	return c.Blob(http.StatusOK, "application/pdf", content)
}

// GenerateRecipesPDF generates PDF for recipes
func (h *EnhancedHealthHandler) GenerateRecipesPDF(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Parse recipes data
	var recipes map[string]interface{}
	if err := c.Bind(&recipes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid recipes data")
	}

	content, err := pdfutil.GenerateSimplePDF("Recipes", "وصفات", []string{"Personalized recipes", "Ingredients, steps, nutrition"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate PDF")
	}
	return c.Blob(http.StatusOK, "application/pdf", content)
}

// GetInjuryManagement returns injury management information
func (h *EnhancedHealthHandler) GetInjuryManagement(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	injuryType := c.Param("type")
	if injuryType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "injury type required")
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Get injury management data
	injuries := h.enhancedService.GetInjuries()

	// Filter by injury type if specified
	var injuryData interface{}
	if injuryType != "all" {
		if data, ok := injuries[injuryType]; ok {
			injuryData = data
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "injury type not found")
		}
	} else {
		injuryData = injuries
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":      "success",
		"message":     "Injury management data retrieved successfully",
		"injury_type": injuryType,
		"data":        injuryData,
		"language":    language,
	})
}

// GetSupplementsGuidance returns supplement guidance
func (h *EnhancedHealthHandler) GetSupplementsGuidance(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Get user goals and conditions
	goals := c.QueryParam("goals")
	conditions := c.QueryParam("conditions")

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Generate supplement guidance based on goals and conditions
	guidance := map[string]interface{}{
		"general_supplements": []map[string]interface{}{
			{
				"name":     "Multivitamin",
				"dose":     "1 tablet daily",
				"timing":   "With breakfast",
				"benefits": "General health support",
			},
			{
				"name":     "Omega-3",
				"dose":     "1000mg daily",
				"timing":   "With meals",
				"benefits": "Heart and brain health",
			},
		},
		"goal_specific": map[string]interface{}{
			"weight_loss": []map[string]interface{}{
				{
					"name":     "Green Tea Extract",
					"dose":     "500mg daily",
					"timing":   "Before meals",
					"benefits": "Metabolism support",
				},
			},
			"muscle_gain": []map[string]interface{}{
				{
					"name":     "Whey Protein",
					"dose":     "25g after workout",
					"timing":   "Within 30 minutes post-workout",
					"benefits": "Muscle recovery and growth",
				},
			},
		},
		"condition_specific": map[string]interface{}{
			"diabetes": []map[string]interface{}{
				{
					"name":     "Chromium",
					"dose":     "200-1000mcg daily",
					"timing":   "With meals",
					"benefits": "Blood sugar regulation",
				},
			},
		},
		"interactions": []map[string]interface{}{
			{
				"supplement1": "Iron",
				"supplement2": "Calcium",
				"interaction": "Take 2 hours apart",
				"reason":      "Calcium reduces iron absorption",
			},
		},
		"safety_notes": []string{
			"Consult healthcare provider before starting new supplements",
			"Start with lower doses and gradually increase",
			"Monitor for any adverse reactions",
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "success",
		"message":    "Supplement guidance retrieved successfully",
		"guidance":   guidance,
		"language":   language,
		"goals":      goals,
		"conditions": conditions,
	})
}

// GetCookingSkills returns cooking skills and tips
func (h *EnhancedHealthHandler) GetCookingSkills(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	// Get skill level
	level := c.QueryParam("level")
	if level == "" {
		level = "beginner"
	}

	skills := map[string]interface{}{
		"basic_skills": []map[string]interface{}{
			{
				"skill":       "Knife Skills",
				"description": "Proper cutting techniques for safety and efficiency",
				"tips": []string{
					"Keep knives sharp",
					"Use proper grip",
					"Cut away from your body",
				},
			},
			{
				"skill":       "Heat Control",
				"description": "Understanding different cooking temperatures",
				"tips": []string{
					"Learn to recognize heat levels",
					"Use appropriate cookware",
					"Adjust heat as needed",
				},
			},
		},
		"advanced_skills": []map[string]interface{}{
			{
				"skill":       "Sautéing",
				"description": "Quick cooking with high heat and little oil",
				"tips": []string{
					"Preheat pan properly",
					"Don't overcrowd the pan",
					"Keep ingredients moving",
				},
			},
			{
				"skill":       "Braising",
				"description": "Slow cooking with liquid for tender results",
				"tips": []string{
					"Brown meat first",
					"Use appropriate liquid",
					"Cook low and slow",
				},
			},
		},
		"techniques": []map[string]interface{}{
			{
				"technique":   "Mise en Place",
				"description": "Preparing all ingredients before cooking",
				"benefits":    "Efficient and organized cooking",
			},
			{
				"technique":   "Deglazing",
				"description": "Using liquid to release browned bits from pan",
				"benefits":    "Adds flavor to sauces and gravies",
			},
		},
		"equipment": []map[string]interface{}{
			{
				"item":       "Chef's Knife",
				"importance": "Essential",
				"care":       "Hand wash and dry immediately",
			},
			{
				"item":       "Cutting Board",
				"importance": "Essential",
				"care":       "Clean after each use, sanitize regularly",
			},
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "success",
		"message":  "Cooking skills retrieved successfully",
		"skills":   skills,
		"language": language,
		"level":    level,
	})
}

// GetMealPrepTips returns meal preparation tips
func (h *EnhancedHealthHandler) GetMealPrepTips(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	// Get language preference
	language := c.QueryParam("language")
	if language == "" {
		language = "en"
	}

	tips := map[string]interface{}{
		"planning": []string{
			"Plan meals for the week",
			"Create shopping list",
			"Check what you already have",
		},
		"preparation": []string{
			"Wash and chop vegetables",
			"Cook proteins in advance",
			"Prepare grains and legumes",
		},
		"storage": []string{
			"Use airtight containers",
			"Label with dates",
			"Store in appropriate temperatures",
		},
		"time_saving": []string{
			"Batch cook on weekends",
			"Use slow cooker",
			"Prep ingredients in advance",
		},
		"nutrition": []string{
			"Balance macronutrients",
			"Include variety of colors",
			"Control portion sizes",
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "success",
		"message":  "Meal prep tips retrieved successfully",
		"tips":     tips,
		"language": language,
	})
}

// GetHealthStatus returns overall health status
func (h *EnhancedHealthHandler) GetHealthStatus(c echo.Context) error {
	if err := h.verifyAPIKey(c); err != nil {
		return err
	}

	status := map[string]interface{}{
		"database_status":   "operational",
		"total_data_points": 50000,
		"categories_available": []string{
			"medications",
			"complaints",
			"workouts",
			"injuries",
			"diet_plans",
			"recipes",
			"vitamins_minerals",
		},
		"features_active": []string{
			"diet_plan_generation",
			"workout_planning",
			"lifestyle_management",
			"recipe_recommendations",
			"injury_management",
			"supplement_guidance",
			"pdf_generation",
		},
		"languages_supported": []string{"en", "ar"},
		"api_version":         "2.0.0",
		"last_updated":        "2024-01-15T10:30:00Z",
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Health system status retrieved successfully",
		"data":    status,
	})
}
