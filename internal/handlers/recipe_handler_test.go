package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-key-generator/internal/models"
	"api-key-generator/internal/services"

	"github.com/go-playground/validator/v10"
)

// MockRecipeService for testing
type MockRecipeService struct {
	mock.Mock
}

func (m *MockRecipeService) GetRecipes(filters models.RecipeFilters) ([]models.Recipe, int, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Recipe), args.Int(1), args.Error(2)
}

func (m *MockRecipeService) GetRecipeByID(id string) (*models.Recipe, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *MockRecipeService) GeneratePersonalizedRecipe(req models.PersonalizedRecipeRequest) (*models.Recipe, error) {
	args := m.Called(req)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *MockRecipeService) SearchRecipes(query string, limit, offset int) ([]models.Recipe, int, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]models.Recipe), args.Int(1), args.Error(2)
}

func (m *MockRecipeService) GetAvailableCuisines() []models.CuisineInfo {
	args := m.Called()
	return args.Get(0).([]models.CuisineInfo)
}

func (m *MockRecipeService) GetAvailableCategories() []models.CategoryInfo {
	args := m.Called()
	return args.Get(0).([]models.CategoryInfo)
}

func (m *MockRecipeService) GetNutritionalAnalysis(recipeID string, servings int) (*models.NutritionalAnalysis, error) {
	args := m.Called(recipeID, servings)
	return args.Get(0).(*models.NutritionalAnalysis), args.Error(1)
}

func (m *MockRecipeService) GetRecipeAlternatives(recipeID string, allergies, preferences []string) ([]models.Recipe, error) {
	args := m.Called(recipeID, allergies, preferences)
	return args.Get(0).([]models.Recipe), args.Error(1)
}

func TestGetRecipes(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	// Mock data
	mockRecipes := []models.Recipe{
		{
			ID:       "test_001",
			Name:     "Test Recipe",
			Cuisine:  "arabian_gulf",
			Category: "main_course",
		},
	}

	mockService.On("GetRecipes", mock.AnythingOfType("models.RecipeFilters")).Return(mockRecipes, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes?cuisine=arabian_gulf", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRecipes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.RecipeResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Recipes))
	assert.Equal(t, "test_001", response.Recipes[0].ID)

	mockService.AssertExpectations(t)
}

func TestGetRecipeByID(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	mockRecipe := &models.Recipe{
		ID:       "test_001",
		Name:     "Test Recipe",
		Cuisine:  "arabian_gulf",
		Category: "main_course",
	}

	mockService.On("GetRecipeByID", "test_001").Return(mockRecipe, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes/test_001", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test_001")

	err := handler.GetRecipeByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Recipe
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test_001", response.ID)

	mockService.AssertExpectations(t)
}

func TestGetRecipeByID_NotFound(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	mockService.On("GetRecipeByID", "nonexistent").Return((*models.Recipe)(nil), services.ErrRecipeNotFound)

	req := httptest.NewRequest(http.MethodGet, "/recipes/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := handler.GetRecipeByID(c)

	// Expect an HTTPError with 404 now that handlers return errors
	assert.Error(t, err)
	if httpErr, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	} else {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	mockService.AssertExpectations(t)
}

func TestGeneratePersonalizedRecipe(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService, validator: validator.New()}

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

	mockService.On("GeneratePersonalizedRecipe", mock.AnythingOfType("models.PersonalizedRecipeRequest")).Return(mockRecipe, nil)

	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/recipes/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GeneratePersonalizedRecipe(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Recipe
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "generated_001", response.ID)

	mockService.AssertExpectations(t)
}

func TestSearchRecipes(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	mockRecipes := []models.Recipe{
		{
			ID:   "search_001",
			Name: "Chicken Kabsa",
		},
	}

	mockService.On("SearchRecipes", "chicken", 20, 0).Return(mockRecipes, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?q=chicken", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.SearchRecipes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.RecipeResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Recipes))

	mockService.AssertExpectations(t)
}

func TestGetCuisines(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	mockCuisines := []models.CuisineInfo{
		{ID: "arabian_gulf", Name: "Arabian Gulf"},
		{ID: "shami", Name: "Shami"},
	}

	mockService.On("GetAvailableCuisines").Return(mockCuisines)

	req := httptest.NewRequest(http.MethodGet, "/recipes/cuisines", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCuisines(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.CuisinesResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Cuisines))

	mockService.AssertExpectations(t)
}

func TestGetNutritionalAnalysis(t *testing.T) {
	e := echo.New()
	mockService := new(MockRecipeService)
	handler := &RecipeHandler{recipeService: mockService}

	mockAnalysis := &models.NutritionalAnalysis{
		RecipeID: "test_001",
		Servings: 2,
		PerServing: models.NutritionInfo{
			Calories: 300,
			Protein:  25,
		},
	}

	mockService.On("GetNutritionalAnalysis", "test_001", 2).Return(mockAnalysis, nil)

	req := httptest.NewRequest(http.MethodGet, "/recipes/test_001/nutrition?servings=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test_001")

	err := handler.GetNutritionalAnalysis(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.NutritionalAnalysis
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test_001", response.RecipeID)
	assert.Equal(t, 2, response.Servings)

	mockService.AssertExpectations(t)
}
