package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"api-key-generator/internal/models"
	"api-key-generator/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// RecipeServiceContract defines the service methods used by the handler
type RecipeServiceContract interface {
	GetRecipes(filters models.RecipeFilters) ([]models.Recipe, int, error)
	GetRecipeByID(id string) (*models.Recipe, error)
	GeneratePersonalizedRecipe(req models.PersonalizedRecipeRequest) (*models.Recipe, error)
	SearchRecipes(query string, limit, offset int) ([]models.Recipe, int, error)
	GetAvailableCuisines() []models.CuisineInfo
	GetAvailableCategories() []models.CategoryInfo
	GetNutritionalAnalysis(recipeID string, servings int) (*models.NutritionalAnalysis, error)
	GetRecipeAlternatives(recipeID string, allergies, preferences []string) ([]models.Recipe, error)
}

// APIKeyProvider is the minimal API key interface used here
type APIKeyProvider interface {
	GetAPIKeyByKey(ctx context.Context, key string) (*models.APIKey, error)
	RecordUsage(ctx context.Context, apiKeyID, endpoint, method string, status int, ipAddress, userAgent string) error
}

type RecipeHandler struct {
	recipeService RecipeServiceContract
	apiKeyService APIKeyProvider
	validator     *validator.Validate
}

func NewRecipeHandler(recipeService RecipeServiceContract, apiKeyService APIKeyProvider) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
		apiKeyService: apiKeyService,
		validator:     validator.New(),
	}
}

// hasPermission checks if granted permissions contain the required permission
func hasPermission(granted []string, required string) bool {
	if required == "" {
		return true
	}
	reqLower := strings.ToLower(required)
	for _, p := range granted {
		pLower := strings.ToLower(p)
		if pLower == reqLower || pLower == "admin" || pLower == "admin:all" {
			return true
		}
	}
	return false
}

// GetRecipes godoc
// @Summary Get recipes with filtering
// @Description Get recipes with optional filtering by cuisine, category, difficulty, etc.
// @Tags recipes
// @Accept json
// @Produce json
// @Param cuisine query string false "Cuisine type" Enums(arabian_gulf,shami,egyptian,moroccan)
// @Param category query string false "Recipe category" Enums(appetizer,main_course,dessert,breakfast)
// @Param difficulty query string false "Difficulty level" Enums(easy,medium,hard)
// @Param max_calories query int false "Maximum calories"
// @Param min_protein query int false "Minimum protein (g)"
// @Param allergies query string false "Comma-separated allergies to exclude"
// @Param limit query int false "Number of recipes to return" default(20)
// @Param offset query int false "Number of recipes to skip" default(0)
// @Success 200 {object} models.RecipeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes [get]
func (h *RecipeHandler) GetRecipes(c echo.Context) error {
	// Authorize API key with required permissions (skip if no service, e.g., unit tests)
	var apiKeyID string
	if h.apiKeyService != nil {
		rawKey := c.Request().Header.Get("X-API-Key")
		if rawKey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
		}
		apiKey, err := h.apiKeyService.GetAPIKeyByKey(c.Request().Context(), rawKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
		}
		// Rate limit check (short-circuit before doing any work)
		if apiKey.RateLimit != nil && *apiKey.RateLimit > 0 && apiKey.RateLimitUsed >= *apiKey.RateLimit {
			return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
		}
		// Permission check
		if !hasPermission(apiKey.Permissions, "recipes:read") {
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
		apiKeyID = apiKey.ID
	}

	// Sanitize and validate input
	filters := models.RecipeFilters{
		Cuisine:    utils.SanitizeInput(c.QueryParam("cuisine")),
		Category:   utils.SanitizeInput(c.QueryParam("category")),
		Difficulty: utils.SanitizeInput(c.QueryParam("difficulty")),
		Allergies:  h.sanitizeStringSlice(strings.Split(c.QueryParam("allergies"), ",")),
		Limit:      20,
		Offset:     0,
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	if offset := c.QueryParam("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filters.Offset = o
		}
	}

	if maxCal := c.QueryParam("max_calories"); maxCal != "" {
		if cal, err := strconv.Atoi(maxCal); err == nil && cal > 0 {
			filters.MaxCalories = &cal
		}
	}

	if minProt := c.QueryParam("min_protein"); minProt != "" {
		if prot, err := strconv.Atoi(minProt); err == nil && prot > 0 {
			filters.MinProtein = &prot
		}
	}

	recipes, total, err := h.recipeService.GetRecipes(filters)
	if err != nil {
		// Log error with sanitization
		c.Logger().Error("Recipe fetch failed:", utils.SanitizeForLog(err.Error()))
		return h.handleServiceError(err)
	}

	// Record API usage
	if apiKeyID != "" {
		h.recordUsage(c, apiKeyID, http.StatusOK)
	}

	return c.JSON(http.StatusOK, models.RecipeResponse{
		Recipes: recipes,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	})
}

// GetRecipeByID godoc
// @Summary Get recipe by ID
// @Description Get a specific recipe by its ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/{id} [get]
func (h *RecipeHandler) GetRecipeByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "recipe ID is required")
	}

	recipe, err := h.recipeService.GetRecipeByID(id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, "recipe not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch recipe")
	}

	return c.JSON(http.StatusOK, recipe)
}

// GeneratePersonalizedRecipe godoc
// @Summary Generate personalized recipe
// @Description Generate a recipe based on user preferences and dietary requirements
// @Tags recipes
// @Accept json
// @Produce json
// @Param request body models.PersonalizedRecipeRequest true "Recipe generation request"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/generate [post]
func (h *RecipeHandler) GeneratePersonalizedRecipe(c echo.Context) error {
	// Authorize write permission (skip in unit tests when apiKeyService is nil)
	var apiKeyID string
	if h.apiKeyService != nil {
		rawKey := c.Request().Header.Get("X-API-Key")
		if rawKey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key required")
		}
		apiKey, err := h.apiKeyService.GetAPIKeyByKey(c.Request().Context(), rawKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid API key")
		}
		if !hasPermission(apiKey.Permissions, "recipes:write") {
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
		apiKeyID = apiKey.ID
	}

	var req models.PersonalizedRecipeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := h.validator.Struct(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed")
	}

	recipe, err := h.recipeService.GeneratePersonalizedRecipe(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate recipe")
	}

	// Record usage
	if apiKeyID != "" {
		h.recordUsage(c, apiKeyID, http.StatusOK)
	}

	return c.JSON(http.StatusOK, recipe)
}

// GetCuisines godoc
// @Summary Get available cuisines
// @Description Get list of all available cuisines
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {object} models.CuisinesResponse
// @Security ApiKeyAuth
// @Router /recipes/cuisines [get]
func (h *RecipeHandler) GetCuisines(c echo.Context) error {
	cuisines := h.recipeService.GetAvailableCuisines()
	return c.JSON(http.StatusOK, models.CuisinesResponse{
		Cuisines: cuisines,
	})
}

// GetCategories godoc
// @Summary Get recipe categories
// @Description Get list of all recipe categories
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {object} models.CategoriesResponse
// @Security ApiKeyAuth
// @Router /recipes/categories [get]
func (h *RecipeHandler) GetCategories(c echo.Context) error {
	categories := h.recipeService.GetAvailableCategories()
	return c.JSON(http.StatusOK, models.CategoriesResponse{
		Categories: categories,
	})
}

// SearchRecipes godoc
// @Summary Search recipes
// @Description Search recipes by name, ingredients, or description
// @Tags recipes
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of results" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} models.RecipeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/search [get]
func (h *RecipeHandler) SearchRecipes(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query is required")
	}

	limit := 20
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	recipes, total, err := h.recipeService.SearchRecipes(query, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "search failed")
	}

	return c.JSON(http.StatusOK, models.RecipeResponse{
		Recipes: recipes,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetNutritionalAnalysis godoc
// @Summary Get nutritional analysis
// @Description Get detailed nutritional analysis for a recipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param servings query int false "Number of servings" default(1)
// @Success 200 {object} models.NutritionalAnalysis
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/{id}/nutrition [get]
func (h *RecipeHandler) GetNutritionalAnalysis(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "recipe ID is required")
	}

	servings := 1
	if s := c.QueryParam("servings"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 {
			servings = parsed
		}
	}

	analysis, err := h.recipeService.GetNutritionalAnalysis(id, servings)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, "recipe not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to analyze nutrition")
	}

	return c.JSON(http.StatusOK, analysis)
}

// GetRecipeAlternatives godoc
// @Summary Get recipe alternatives
// @Description Get alternative recipes based on dietary restrictions or preferences
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param allergies query string false "Comma-separated allergies"
// @Param dietary_preferences query string false "Comma-separated dietary preferences"
// @Success 200 {object} models.RecipeAlternativesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /recipes/{id}/alternatives [get]
func (h *RecipeHandler) GetRecipeAlternatives(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "recipe ID is required")
	}

	allergies := strings.Split(c.QueryParam("allergies"), ",")
	preferences := strings.Split(c.QueryParam("dietary_preferences"), ",")

	alternatives, err := h.recipeService.GetRecipeAlternatives(id, allergies, preferences)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, "recipe not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get alternatives")
	}

	return c.JSON(http.StatusOK, models.RecipeAlternativesResponse{
		OriginalRecipe: id,
		Alternatives:   alternatives,
	})
}
