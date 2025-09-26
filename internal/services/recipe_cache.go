package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"api-key-generator/internal/models"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// RecipeCache provides in-memory caching for recipe data
type RecipeCache struct {
	cache map[string]*CacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewRecipeCache creates a new recipe cache
func NewRecipeCache(ttl time.Duration) *RecipeCache {
	cache := &RecipeCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves an item from cache
func (c *RecipeCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(c.cache, key)
		return nil, false
	}

	return entry.Data, true
}

// Set stores an item in cache
func (c *RecipeCache) Set(key string, data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete removes an item from cache
func (c *RecipeCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}

// Clear removes all items from cache
func (c *RecipeCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*CacheEntry)
}

// Size returns the number of items in cache
func (c *RecipeCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.cache)
}

// cleanup removes expired entries periodically
func (c *RecipeCache) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.After(entry.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// CachedRecipeService wraps RecipeService with caching
type CachedRecipeService struct {
	service *RecipeService
	cache   *RecipeCache
}

// NewCachedRecipeService creates a new cached recipe service
func NewCachedRecipeService(service *RecipeService, cacheTTL time.Duration) *CachedRecipeService {
	return &CachedRecipeService{
		service: service,
		cache:   NewRecipeCache(cacheTTL),
	}
}

// GetRecipes with caching
func (s *CachedRecipeService) GetRecipes(filters models.RecipeFilters) ([]models.Recipe, int, error) {
	// Create cache key from filters
	key := s.createCacheKey("recipes", filters)

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if result, ok := cached.(models.RecipeResponse); ok {
			return result.Recipes, result.Total, nil
		}
	}

	// Cache miss - get from service
	recipes, total, err := s.service.GetRecipes(filters)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result
	result := models.RecipeResponse{
		Recipes: recipes,
		Total:   total,
		Limit:   filters.Limit,
		Offset:  filters.Offset,
	}
	s.cache.Set(key, result)

	return recipes, total, nil
}

// GetRecipeByID with caching
func (s *CachedRecipeService) GetRecipeByID(id string) (*models.Recipe, error) {
	key := fmt.Sprintf("recipe:%s", id)

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if recipe, ok := cached.(*models.Recipe); ok {
			return recipe, nil
		}
	}

	// Cache miss - get from service
	recipe, err := s.service.GetRecipeByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.Set(key, recipe)

	return recipe, nil
}

// SearchRecipes with caching
func (s *CachedRecipeService) SearchRecipes(query string, limit, offset int) ([]models.Recipe, int, error) {
	key := fmt.Sprintf("search:%s:%d:%d", query, limit, offset)

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if result, ok := cached.(models.RecipeResponse); ok {
			return result.Recipes, result.Total, nil
		}
	}

	// Cache miss - get from service
	recipes, total, err := s.service.SearchRecipes(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Cache the result
	result := models.RecipeResponse{
		Recipes: recipes,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}
	s.cache.Set(key, result)

	return recipes, total, nil
}

// GetAvailableCuisines with caching
func (s *CachedRecipeService) GetAvailableCuisines() []models.CuisineInfo {
	key := "cuisines"

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if cuisines, ok := cached.([]models.CuisineInfo); ok {
			return cuisines
		}
	}

	// Cache miss - get from service
	cuisines := s.service.GetAvailableCuisines()

	// Cache the result
	s.cache.Set(key, cuisines)

	return cuisines
}

// GetAvailableCategories with caching
func (s *CachedRecipeService) GetAvailableCategories() []models.CategoryInfo {
	key := "categories"

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if categories, ok := cached.([]models.CategoryInfo); ok {
			return categories
		}
	}

	// Cache miss - get from service
	categories := s.service.GetAvailableCategories()

	// Cache the result
	s.cache.Set(key, categories)

	return categories
}

// GeneratePersonalizedRecipe - not cached as it's personalized
func (s *CachedRecipeService) GeneratePersonalizedRecipe(req models.PersonalizedRecipeRequest) (*models.Recipe, error) {
	return s.service.GeneratePersonalizedRecipe(req)
}

// GetNutritionalAnalysis with caching
func (s *CachedRecipeService) GetNutritionalAnalysis(recipeID string, servings int) (*models.NutritionalAnalysis, error) {
	key := fmt.Sprintf("nutrition:%s:%d", recipeID, servings)

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if analysis, ok := cached.(*models.NutritionalAnalysis); ok {
			return analysis, nil
		}
	}

	// Cache miss - get from service
	analysis, err := s.service.GetNutritionalAnalysis(recipeID, servings)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.Set(key, analysis)

	return analysis, nil
}

// GetRecipeAlternatives with caching
func (s *CachedRecipeService) GetRecipeAlternatives(recipeID string, allergies, preferences []string) ([]models.Recipe, error) {
	key := s.createCacheKey("alternatives", map[string]interface{}{
		"recipe_id":   recipeID,
		"allergies":   allergies,
		"preferences": preferences,
	})

	// Try cache first
	if cached, found := s.cache.Get(key); found {
		if alternatives, ok := cached.([]models.Recipe); ok {
			return alternatives, nil
		}
	}

	// Cache miss - get from service
	alternatives, err := s.service.GetRecipeAlternatives(recipeID, allergies, preferences)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.Set(key, alternatives)

	return alternatives, nil
}

// InvalidateRecipe removes recipe-related cache entries
func (s *CachedRecipeService) InvalidateRecipe(recipeID string) {
	s.cache.Delete(fmt.Sprintf("recipe:%s", recipeID))
	// Could also implement pattern-based invalidation for related entries
}

// GetCacheStats returns cache statistics
func (s *CachedRecipeService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"size": s.cache.Size(),
		"ttl":  s.cache.ttl.String(),
	}
}

// createCacheKey creates a consistent cache key from data
func (s *CachedRecipeService) createCacheKey(prefix string, data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return fmt.Sprintf("%s:%x", prefix, jsonData)
}
