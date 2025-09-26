package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-key-generator/internal/models"
)

// MockAPIKeyService for testing
type MockAPIKeyService struct {
	mock.Mock
}

func (m *MockAPIKeyService) GetAPIKeyByKey(ctx context.Context, key string) (*models.APIKey, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*models.APIKey), args.Error(1)
}

func (m *MockAPIKeyService) RecordUsage(ctx context.Context, apiKeyID, endpoint, method string, status int, ip, userAgent string) error {
	args := m.Called(ctx, apiKeyID, endpoint, method, status, ip, userAgent)
	return args.Error(0)
}

func TestRecipeHandler_GetRecipes_WithValidAPIKey(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	// Mock valid API key
	rl := 100
	validAPIKey := &models.APIKey{
		ID:            "test-key-id",
		Key:           "ak_test_key_12345",
		Permissions:   []string{"recipes:read"},
		IsActive:      true,
		ExpiresAt:     &time.Time{},
		RateLimit:     &rl,
		RateLimitUsed: 10,
	}

	// Mock recipe data
	mockRecipes := []models.Recipe{
		{
			ID:       "recipe_001",
			Name:     "Test Recipe",
			Cuisine:  "arabian_gulf",
			Category: "main_course",
		},
	}

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "ak_test_key_12345").Return(validAPIKey, nil)
	mockRecipeService.On("GetRecipes", mock.AnythingOfType("models.RecipeFilters")).Return(mockRecipes, 1, nil)
	mockAPIKeyService.On("RecordUsage", mock.Anything, "test-key-id", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes?cuisine=arabian_gulf", nil)
	req.Header.Set("X-API-Key", "ak_test_key_12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.RecipeResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Recipes))
	assert.Equal(t, "recipe_001", response.Recipes[0].ID)

	mockAPIKeyService.AssertExpectations(t)
	mockRecipeService.AssertExpectations(t)
}

func TestRecipeHandler_GetRecipes_WithInvalidAPIKey(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "invalid_key").Return((*models.APIKey)(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	req.Header.Set("X-API-Key", "invalid_key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)

	mockAPIKeyService.AssertExpectations(t)
}

func TestRecipeHandler_GetRecipes_WithInsufficientPermissions(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	// Mock API key without recipes:read permission
	limitedAPIKey := &models.APIKey{
		ID:          "limited-key-id",
		Key:         "ak_limited_key_12345",
		Permissions: []string{"users:read"}, // No recipes:read permission
		IsActive:    true,
		ExpiresAt:   &time.Time{},
	}

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "ak_limited_key_12345").Return(limitedAPIKey, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	req.Header.Set("X-API-Key", "ak_limited_key_12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)

	mockAPIKeyService.AssertExpectations(t)
}

func TestRecipeHandler_GeneratePersonalizedRecipe_WithValidRequest(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	// Mock valid API key with write permissions
	validAPIKey := &models.APIKey{
		ID:          "test-key-id",
		Key:         "ak_test_key_12345",
		Permissions: []string{"recipes:write"},
		IsActive:    true,
		ExpiresAt:   &time.Time{},
	}

	mockRecipe := &models.Recipe{
		ID:       "generated_001",
		Name:     "Generated Recipe",
		Cuisine:  "arabian_gulf",
		Category: "main_course",
	}

	request := models.PersonalizedRecipeRequest{
		UserID:   "user_123",
		Cuisine:  "arabian_gulf",
		Category: "main_course",
	}

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "ak_test_key_12345").Return(validAPIKey, nil)
	mockRecipeService.On("GeneratePersonalizedRecipe", mock.AnythingOfType("models.PersonalizedRecipeRequest")).Return(mockRecipe, nil)
	mockAPIKeyService.On("RecordUsage", mock.Anything, "test-key-id", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "ak_test_key_12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GeneratePersonalizedRecipe(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Recipe
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "generated_001", response.ID)

	mockAPIKeyService.AssertExpectations(t)
	mockRecipeService.AssertExpectations(t)
}

func TestRecipeHandler_InputSanitization(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	validAPIKey := &models.APIKey{
		ID:          "test-key-id",
		Key:         "ak_test_key_12345",
		Permissions: []string{"recipes:read"},
		IsActive:    true,
		ExpiresAt:   &time.Time{},
	}

	mockRecipes := []models.Recipe{}

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "ak_test_key_12345").Return(validAPIKey, nil)
	mockRecipeService.On("GetRecipes", mock.MatchedBy(func(filters models.RecipeFilters) bool {
		// Verify that malicious input was sanitized
		return filters.Cuisine == "arabian_gulf" // Should be sanitized from malicious input
	})).Return(mockRecipes, 0, nil)
	mockAPIKeyService.On("RecordUsage", mock.Anything, "test-key-id", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Test with potentially malicious input
	req := httptest.NewRequest(http.MethodGet, "/recipes?cuisine=arabian_gulf<script>alert('xss')</script>", nil)
	req.Header.Set("X-API-Key", "ak_test_key_12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockAPIKeyService.AssertExpectations(t)
	mockRecipeService.AssertExpectations(t)
}

func TestRecipeHandler_RateLimitHandling(t *testing.T) {
	e := echo.New()
	mockRecipeService := new(MockRecipeService)
	mockAPIKeyService := new(MockAPIKeyService)
	handler := NewRecipeHandler(mockRecipeService, mockAPIKeyService)

	// Mock API key that has exceeded rate limit
	rl2 := 100
	rateLimitedAPIKey := &models.APIKey{
		ID:            "rate-limited-key",
		Key:           "ak_rate_limited_12345",
		Permissions:   []string{"recipes:read"},
		IsActive:      true,
		ExpiresAt:     &time.Time{},
		RateLimit:     &rl2,
		RateLimitUsed: 100, // At limit
	}

	mockAPIKeyService.On("GetAPIKeyByKey", mock.Anything, "ak_rate_limited_12345").Return(rateLimitedAPIKey, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	req.Header.Set("X-API-Key", "ak_rate_limited_12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusTooManyRequests, httpErr.Code)

	mockAPIKeyService.AssertExpectations(t)
}
