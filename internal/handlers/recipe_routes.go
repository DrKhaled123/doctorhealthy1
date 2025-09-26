package handlers

import (
	"api-key-generator/internal/middleware"
	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

func SetupRecipeRoutes(e *echo.Echo, recipeService *services.RecipeService, apiKeyService *services.APIKeyService) {
	// Create recipe handler with security integration
	recipeHandler := NewRecipeHandler(recipeService, apiKeyService)

	// Recipe API group with comprehensive middleware stack
	api := e.Group("/api/v1")
	api.Use(middleware.RecipeSecurity())
	api.Use(middleware.RecipeLogging())
	api.Use(middleware.RecipeRateLimit(apiKeyService))
	api.Use(middleware.RecipeValidation())
	api.Use(middleware.RecipeMetrics())

	// Recipe endpoints - all protected by middleware stack
	// Read operations require recipes:read scope
	api.GET("/recipes", recipeHandler.GetRecipes, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/:id", recipeHandler.GetRecipeByID, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/search", recipeHandler.SearchRecipes, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/cuisines", recipeHandler.GetCuisines, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/categories", recipeHandler.GetCategories, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/:id/nutrition", recipeHandler.GetNutritionalAnalysis, middleware.ScopesAny(apiKeyService, "recipes:read"))
	api.GET("/recipes/:id/alternatives", recipeHandler.GetRecipeAlternatives, middleware.ScopesAny(apiKeyService, "recipes:read"))

	// Write/generation requires recipes:write
	api.POST("/recipes/generate", recipeHandler.GeneratePersonalizedRecipe, middleware.ScopesAny(apiKeyService, "recipes:write"))
}
