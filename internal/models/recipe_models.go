package models

import (
	"time"
)

// Recipe represents a complete recipe with all details
type Recipe struct {
	ID              string         `json:"id" db:"id"`
	Name            string         `json:"name" db:"name"`
	Cuisine         string         `json:"cuisine" db:"cuisine"`
	Category        string         `json:"category" db:"category"`
	Ingredients     []Ingredient   `json:"ingredients" db:"ingredients"`
	Instructions    []string       `json:"instructions" db:"instructions"`
	Nutrition       NutritionInfo  `json:"nutrition" db:"nutrition"`
	Macros          MacroBreakdown `json:"macros" db:"macros"`
	PrepTime        int            `json:"prep_time" db:"prep_time"`
	CookTime        int            `json:"cook_time" db:"cook_time"`
	Servings        int            `json:"servings" db:"servings"`
	Difficulty      string         `json:"difficulty" db:"difficulty"`
	Alternatives    []string       `json:"alternatives" db:"alternatives"`
	SkillsTips      []string       `json:"skills_tips" db:"skills_tips"`
	EnhancementTips []string       `json:"enhancement_tips" db:"enhancement_tips"`
	HealthBenefits  []string       `json:"health_benefits" db:"health_benefits"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	Tags            []string       `json:"tags,omitempty" db:"tags"`
	ImageURL        string         `json:"image_url,omitempty" db:"image_url"`
	VideoURL        string         `json:"video_url,omitempty" db:"video_url"`
	Rating          float64        `json:"rating,omitempty" db:"rating"`
	ReviewCount     int            `json:"review_count,omitempty" db:"review_count"`
}

// Ingredient represents a recipe ingredient
type Ingredient struct {
	Name   string  `json:"name" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Unit   string  `json:"unit" validate:"required"`
}

// NutritionInfo represents nutritional information per serving
type NutritionInfo struct {
	Calories int     `json:"calories" validate:"required,gt=0"`
	Protein  float64 `json:"protein" validate:"required,gte=0"`
	Carbs    float64 `json:"carbs" validate:"required,gte=0"`
	Fat      float64 `json:"fat" validate:"required,gte=0"`
	Fiber    float64 `json:"fiber" validate:"gte=0"`
	Sodium   float64 `json:"sodium" validate:"gte=0"`
	Sugar    float64 `json:"sugar,omitempty" validate:"gte=0"`
	Calcium  float64 `json:"calcium,omitempty" validate:"gte=0"`
	Iron     float64 `json:"iron,omitempty" validate:"gte=0"`
	VitaminC float64 `json:"vitamin_c,omitempty" validate:"gte=0"`
}

// MacroBreakdown represents macronutrient percentages
type MacroBreakdown struct {
	Protein int `json:"protein" validate:"required,gte=0,lte=100"`
	Carbs   int `json:"carbs" validate:"required,gte=0,lte=100"`
	Fat     int `json:"fat" validate:"required,gte=0,lte=100"`
}

// RecipeFilters represents filtering options for recipe queries
type RecipeFilters struct {
	Cuisine     string   `json:"cuisine,omitempty"`
	Category    string   `json:"category,omitempty"`
	Difficulty  string   `json:"difficulty,omitempty"`
	MaxCalories *int     `json:"max_calories,omitempty"`
	MinProtein  *int     `json:"min_protein,omitempty"`
	MaxPrepTime *int     `json:"max_prep_time,omitempty"`
	Allergies   []string `json:"allergies,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Limit       int      `json:"limit" validate:"required,gt=0,lte=100"`
	Offset      int      `json:"offset" validate:"gte=0"`
}

// PersonalizedRecipeRequest represents a request for personalized recipe generation
type PersonalizedRecipeRequest struct {
	UserID             string   `json:"user_id" validate:"required"`
	Cuisine            string   `json:"cuisine" validate:"required,oneof=arabian_gulf shami egyptian moroccan"`
	Category           string   `json:"category" validate:"required,oneof=appetizer main_course dessert breakfast"`
	Difficulty         string   `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard"`
	MaxCalories        int      `json:"max_calories,omitempty" validate:"omitempty,gt=0"`
	MinProtein         int      `json:"min_protein,omitempty" validate:"gte=0"`
	MaxPrepTime        int      `json:"max_prep_time,omitempty" validate:"omitempty,gt=0"`
	Allergies          []string `json:"allergies,omitempty"`
	DietaryPreferences []string `json:"dietary_preferences,omitempty"`
	DislikedFoods      []string `json:"disliked_foods,omitempty"`
	HealthConditions   []string `json:"health_conditions,omitempty"`
}

// RecipeResponse represents the response for recipe queries
type RecipeResponse struct {
	Recipes []Recipe `json:"recipes"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

// CuisinesResponse represents available cuisines
type CuisinesResponse struct {
	Cuisines []CuisineInfo `json:"cuisines"`
}

// CuisineInfo represents cuisine information
type CuisineInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RecipeCount int    `json:"recipe_count"`
}

// CategoriesResponse represents available categories
type CategoriesResponse struct {
	Categories []CategoryInfo `json:"categories"`
}

// CategoryInfo represents category information
type CategoryInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RecipeCount int    `json:"recipe_count"`
}

// NutritionalAnalysis represents detailed nutritional analysis
type NutritionalAnalysis struct {
	RecipeID         string                    `json:"recipe_id"`
	Servings         int                       `json:"servings"`
	PerServing       NutritionInfo             `json:"per_serving"`
	Total            NutritionInfo             `json:"total"`
	MacroBreakdown   MacroBreakdown            `json:"macro_breakdown"`
	HealthScore      int                       `json:"health_score"`
	DietaryFlags     []string                  `json:"dietary_flags"`
	AllergenWarnings []string                  `json:"allergen_warnings"`
	Recommendations  []NutritionRecommendation `json:"recommendations"`
}

// NutritionRecommendation represents nutritional recommendations
type NutritionRecommendation struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	Suggestion string `json:"suggestion,omitempty"`
}

// RecipeAlternativesResponse represents recipe alternatives
type RecipeAlternativesResponse struct {
	OriginalRecipe string   `json:"original_recipe"`
	Alternatives   []Recipe `json:"alternatives"`
}

// RecipeRating represents a recipe rating
type RecipeRating struct {
	ID        string    `json:"id" db:"id"`
	RecipeID  string    `json:"recipe_id" db:"recipe_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating" validate:"required,gte=1,lte=5"`
	Review    string    `json:"review,omitempty" db:"review"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RecipeCollection represents a collection of recipes
type RecipeCollection struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Description string    `json:"description,omitempty" db:"description"`
	UserID      string    `json:"user_id" db:"user_id"`
	RecipeIDs   []string  `json:"recipe_ids" db:"recipe_ids"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// MealPlan represents a meal plan
type MealPlan struct {
	ID        string              `json:"id" db:"id"`
	UserID    string              `json:"user_id" db:"user_id"`
	Name      string              `json:"name" db:"name" validate:"required"`
	StartDate time.Time           `json:"start_date" db:"start_date"`
	EndDate   time.Time           `json:"end_date" db:"end_date"`
	Meals     []PlannedMeal       `json:"meals" db:"meals"`
	Nutrition DailyNutritionGoals `json:"nutrition_goals" db:"nutrition_goals"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" db:"updated_at"`
}

// PlannedMeal represents a meal in a meal plan
type PlannedMeal struct {
	Date     time.Time `json:"date"`
	MealType string    `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	RecipeID string    `json:"recipe_id" validate:"required"`
	Servings int       `json:"servings" validate:"required,gt=0"`
}

// DailyNutritionGoals represents daily nutrition targets
type DailyNutritionGoals struct {
	Calories int     `json:"calories" validate:"required,gt=0"`
	Protein  float64 `json:"protein" validate:"required,gt=0"`
	Carbs    float64 `json:"carbs" validate:"required,gt=0"`
	Fat      float64 `json:"fat" validate:"required,gt=0"`
	Fiber    float64 `json:"fiber" validate:"gt=0"`
	Sodium   float64 `json:"sodium" validate:"gt=0"`
}

// ShoppingList represents a shopping list generated from recipes
type ShoppingList struct {
	ID        string             `json:"id" db:"id"`
	UserID    string             `json:"user_id" db:"user_id"`
	Name      string             `json:"name" db:"name" validate:"required"`
	Items     []ShoppingListItem `json:"items" db:"items"`
	RecipeIDs []string           `json:"recipe_ids" db:"recipe_ids"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
}

// ShoppingListItem represents an item in a shopping list
type ShoppingListItem struct {
	Name      string  `json:"name" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Unit      string  `json:"unit" validate:"required"`
	Category  string  `json:"category,omitempty"`
	Purchased bool    `json:"purchased"`
	Notes     string  `json:"notes,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message,omitempty"`
	Code      string            `json:"code,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	TraceID   string            `json:"trace_id,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}
