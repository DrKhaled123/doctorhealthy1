package handlers

import (
	"net/http"

	"api-key-generator/internal/models"
	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

// Services container for dependency injection
type Services struct {
	User      *services.UserService
	Nutrition *services.NutritionService
	Workout   *services.WorkoutService
	Health    *services.HealthService
	Recipe    *services.RecipeService
}

// HealthHandler handles health management HTTP requests
type HealthHandler struct {
	services *Services
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(services *Services) *HealthHandler {
	return &HealthHandler{
		services: services,
	}
}

// User Management Endpoints

// CreateUser creates a new user profile
func (h *HealthHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.services.User.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User profile created successfully",
		"user":    user,
	})
}

// GetUser retrieves a user profile
func (h *HealthHandler) GetUser(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	user, err := h.services.User.GetUser(c.Request().Context(), userID)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, user)
}

// CalculateCalories calculates daily calorie needs for a user
func (h *HealthHandler) CalculateCalories(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	user, err := h.services.User.GetUser(c.Request().Context(), userID)
	if err != nil {
		return handleServiceError(err)
	}

	calculation, err := h.services.User.CalculateCalories(c.Request().Context(), user)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"calculation": calculation,
		"user_data": map[string]interface{}{
			"name":   user.Name,
			"age":    user.Age,
			"weight": user.Weight,
			"height": user.Height,
			"goal":   user.Goal,
		},
	})
}

// Nutrition Plan Endpoints

// GenerateNutritionPlan generates a personalized nutrition plan
func (h *HealthHandler) GenerateNutritionPlan(c echo.Context) error {
	var req models.GenerateNutritionPlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	plan, err := h.services.Nutrition.GenerateNutritionPlan(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Nutrition plan generated successfully",
		"nutrition_plan": plan,
	})
}

// GetAvailablePlanTypes returns available nutrition plan types
func (h *HealthHandler) GetAvailablePlanTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"plan_types": models.AvailablePlanTypes,
		"descriptions": map[string]string{
			"low_carb":          "Low carbohydrate diet (25% carbs, 35% protein, 40% fat)",
			"keto":              "Ketogenic diet (5% carbs, 25% protein, 70% fat)",
			"mediterranean":     "Mediterranean diet rich in healthy fats and whole foods",
			"dash":              "DASH diet for blood pressure management",
			"balanced":          "Balanced macronutrient distribution",
			"high_carb":         "High carbohydrate diet for athletes (60% carbs, 20% protein, 20% fat)",
			"vegan":             "Plant-based diet excluding all animal products",
			"anti_inflammatory": "Anti-inflammatory foods to reduce inflammation",
		},
	})
}

// Workout Plan Endpoints

// GenerateWorkoutPlan generates a personalized workout plan
func (h *HealthHandler) GenerateWorkoutPlan(c echo.Context) error {
	var req models.GenerateWorkoutPlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	plan, err := h.services.Workout.GenerateWorkoutPlan(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Workout plan generated successfully",
		"workout_plan": plan,
	})
}

// GetWorkoutGoals returns available workout goals
func (h *HealthHandler) GetWorkoutGoals(c echo.Context) error {
	goals := h.services.Workout.GetAvailableGoals()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"goals": goals,
		"descriptions": map[string]string{
			"build_muscle":       "Focus on muscle hypertrophy and strength",
			"improve_strength":   "Increase overall strength and power",
			"lose_weight":        "Weight loss through cardio and resistance training",
			"improve_endurance":  "Cardiovascular and muscular endurance",
			"general_fitness":    "Overall health and fitness improvement",
			"body_recomposition": "Simultaneous muscle gain and fat loss",
		},
	})
}

// GetAvailableInjuries returns common injuries for selection
func (h *HealthHandler) GetAvailableInjuries(c echo.Context) error {
	injuries := h.services.Workout.GetAvailableInjuries()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"injuries": injuries,
	})
}

// GetWorkoutComplaints returns available workout-related complaints
func (h *HealthHandler) GetWorkoutComplaints(c echo.Context) error {
	complaints := h.services.Workout.GetAvailableComplaints()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"complaints": complaints,
	})
}

// Health Plan Endpoints

// GenerateHealthPlan generates a personalized health management plan
func (h *HealthHandler) GenerateHealthPlan(c echo.Context) error {
	var req models.GenerateHealthPlanRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	plan, err := h.services.Health.GenerateHealthPlan(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Health plan generated successfully",
		"health_plan": plan,
	})
}

// GetAvailableDiseases returns common diseases for selection
func (h *HealthHandler) GetAvailableDiseases(c echo.Context) error {
	diseases := h.services.Health.GetAvailableDiseases()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"diseases": diseases,
	})
}

// GetHealthComplaints returns available health complaints
func (h *HealthHandler) GetHealthComplaints(c echo.Context) error {
	complaints := h.services.Health.GetAvailableComplaints()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"complaints": complaints,
	})
}

// Recipe Endpoints

// GenerateRecipe generates a recipe based on user preferences
func (h *HealthHandler) GenerateRecipe(c echo.Context) error {
	var req models.GenerateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	recipe, err := h.services.Recipe.GenerateRecipe(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Recipe generated successfully",
		"recipe":  recipe,
	})
}

// GetAvailableCuisines returns available cuisines
func (h *HealthHandler) GetAvailableCuisines(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"cuisines": models.AvailableCuisines,
		"descriptions": map[string]string{
			"mediterranean":  "Fresh vegetables, olive oil, fish, and whole grains",
			"middle_eastern": "Spices, legumes, grains, and lean meats",
			"asian":          "Rice, vegetables, lean proteins, and minimal oil",
			"indian":         "Spices, legumes, vegetables, and yogurt",
			"mexican":        "Beans, vegetables, lean meats, and whole grains",
			"italian":        "Tomatoes, olive oil, herbs, and lean proteins",
			"american":       "Balanced portions of protein, vegetables, and grains",
			"french":         "Fresh ingredients, herbs, and moderate portions",
			"japanese":       "Fish, rice, vegetables, and fermented foods",
			"thai":           "Fresh herbs, vegetables, lean proteins, and spices",
			"greek":          "Olive oil, fish, vegetables, and yogurt",
			"turkish":        "Vegetables, legumes, grains, and lean meats",
			"moroccan":       "Spices, vegetables, legumes, and lean meats",
			"lebanese":       "Fresh vegetables, olive oil, legumes, and herbs",
			"egyptian":       "Legumes, vegetables, whole grains, and lean proteins",
		},
	})
}

// GetUserRecipes gets recipes generated for a specific user
func (h *HealthHandler) GetUserRecipes(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	limit := 10
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		// Parse limit parameter (simplified)
		limit = 10 // Default for now
	}

	recipes, err := h.services.Recipe.GetRecipesByUser(c.Request().Context(), userID, limit)
	if err != nil {
		return handleServiceError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recipes": recipes,
		"count":   len(recipes),
	})
}

// Utility Endpoints

// GetAvailableGoals returns all available goals
func (h *HealthHandler) GetAvailableGoals(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"goals": models.AvailableGoals,
		"descriptions": map[string]string{
			"lose_weight":        "Reduce body weight through calorie deficit",
			"gain_weight":        "Increase body weight through calorie surplus",
			"maintain_weight":    "Maintain current weight while improving health",
			"build_muscle":       "Increase muscle mass and strength",
			"improve_strength":   "Enhance overall physical strength",
			"body_recomposition": "Simultaneously build muscle and lose fat",
		},
	})
}

// GetMedicalDisclaimer returns the medical disclaimer
func (h *HealthHandler) GetMedicalDisclaimer(c echo.Context) error {
	language := c.QueryParam("lang")
	if language == "" {
		language = "en"
	}

	var disclaimer string
	if language == "ar" {
		disclaimer = `إخلاء المسؤولية الطبية: هذه المعلومات لأغراض تعليمية فقط وليست نصيحة طبية. 
استشر دائماً مختصين صحيين مؤهلين قبل إجراء تغييرات على نظامك الغذائي أو التمارين أو تناول المكملات، 
خاصة إذا كان لديك حالات طبية أو تتناول أدوية. النتائج الفردية قد تختلف.`
	} else {
		disclaimer = `MEDICAL DISCLAIMER: This information is for educational purposes only and is not intended as medical advice. 
Always consult with qualified healthcare professionals before making changes to your diet, exercise routine, or 
taking supplements, especially if you have medical conditions or take medications. Individual results may vary.`
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"disclaimer": disclaimer,
		"language":   language,
	})
}

// RegisterHealthRoutes registers all health management routes
func RegisterHealthRoutes(g *echo.Group,
	userService *services.UserService,
	nutritionService *services.NutritionService,
	workoutService *services.WorkoutService,
	healthService *services.HealthService,
	recipeService *services.RecipeService,
) {
	servicesContainer := &Services{
		User:      userService,
		Nutrition: nutritionService,
		Workout:   workoutService,
		Health:    healthService,
		Recipe:    recipeService,
	}
	handler := NewHealthHandler(servicesContainer)

	// User management routes
	g.POST("/users", handler.CreateUser)
	g.GET("/users/:id", handler.GetUser)
	g.GET("/users/:id/calories", handler.CalculateCalories)

	// Nutrition routes
	g.POST("/nutrition/generate-plan", handler.GenerateNutritionPlan)
	g.GET("/nutrition/plan-types", handler.GetAvailablePlanTypes)

	// Workout routes
	g.POST("/workouts/generate-plan", handler.GenerateWorkoutPlan)
	g.GET("/workouts/goals", handler.GetWorkoutGoals)
	g.GET("/workouts/injuries", handler.GetAvailableInjuries)
	g.GET("/workouts/complaints", handler.GetWorkoutComplaints)

	// Health management routes
	g.POST("/health/generate-plan", handler.GenerateHealthPlan)
	g.GET("/health/diseases", handler.GetAvailableDiseases)
	g.GET("/health/complaints", handler.GetHealthComplaints)

	// Recipe routes
	g.POST("/recipes/generate", handler.GenerateRecipe)
	g.GET("/recipes/cuisines", handler.GetAvailableCuisines)
	g.GET("/users/:id/recipes", handler.GetUserRecipes)

	// Utility routes
	g.GET("/goals", handler.GetAvailableGoals)
	g.GET("/disclaimer", handler.GetMedicalDisclaimer)
}
