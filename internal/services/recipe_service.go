package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"api-key-generator/internal/models"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
)

type RecipeService struct {
	db *sql.DB
}

func NewRecipeService(db *sql.DB) *RecipeService {
	return &RecipeService{db: db}
}

// GenerateRecipe provides a minimal wrapper to maintain backward compatibility
func (s *RecipeService) GenerateRecipe(ctx context.Context, req *models.GenerateRecipeRequest) (*models.Recipe, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	pr := models.PersonalizedRecipeRequest{
		Cuisine:     req.Cuisine,
		Category:    req.MealType,
		Difficulty:  req.Difficulty,
		MaxCalories: req.MaxCalories,
		MinProtein:  0,
		Allergies:   []string{},
	}
	return s.GeneratePersonalizedRecipe(pr)
}

func (s *RecipeService) GetRecipes(filters models.RecipeFilters) ([]models.Recipe, int, error) {
	query := `SELECT id, name, cuisine, category, ingredients, instructions, nutrition, macros, 
			  prep_time, cook_time, servings, difficulty, alternatives, skills_tips, 
			  enhancement_tips, health_benefits, created_at, updated_at 
			  FROM recipes WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filters.Cuisine != "" {
		query += fmt.Sprintf(" AND cuisine = $%d", argIndex)
		args = append(args, filters.Cuisine)
		argIndex++
	}

	if filters.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Difficulty != "" {
		query += fmt.Sprintf(" AND difficulty = $%d", argIndex)
		args = append(args, filters.Difficulty)
		argIndex++
	}

	if filters.MaxCalories != nil {
		query += fmt.Sprintf(" AND JSON_EXTRACT(nutrition, '$.calories') <= $%d", argIndex)
		args = append(args, *filters.MaxCalories)
		argIndex++
	}

	if filters.MinProtein != nil {
		query += fmt.Sprintf(" AND JSON_EXTRACT(nutrition, '$.protein') >= $%d", argIndex)
		args = append(args, *filters.MinProtein)
		argIndex++
	}

	// Count total
	countQuery := strings.Replace(query, "SELECT id, name, cuisine, category, ingredients, instructions, nutrition, macros, prep_time, cook_time, servings, difficulty, alternatives, skills_tips, enhancement_tips, health_benefits, created_at, updated_at FROM recipes", "SELECT COUNT(*) FROM recipes", 1)

	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var ingredientsJSON, instructionsJSON, nutritionJSON, macrosJSON string
		var alternativesJSON, skillsTipsJSON, enhancementTipsJSON, healthBenefitsJSON string

		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Cuisine, &recipe.Category,
			&ingredientsJSON, &instructionsJSON, &nutritionJSON, &macrosJSON,
			&recipe.PrepTime, &recipe.CookTime, &recipe.Servings, &recipe.Difficulty,
			&alternativesJSON, &skillsTipsJSON, &enhancementTipsJSON, &healthBenefitsJSON,
			&recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse JSON fields
		_ = json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients)
		_ = json.Unmarshal([]byte(instructionsJSON), &recipe.Instructions)
		_ = json.Unmarshal([]byte(nutritionJSON), &recipe.Nutrition)
		_ = json.Unmarshal([]byte(macrosJSON), &recipe.Macros)
		_ = json.Unmarshal([]byte(alternativesJSON), &recipe.Alternatives)
		_ = json.Unmarshal([]byte(skillsTipsJSON), &recipe.SkillsTips)
		_ = json.Unmarshal([]byte(enhancementTipsJSON), &recipe.EnhancementTips)
		_ = json.Unmarshal([]byte(healthBenefitsJSON), &recipe.HealthBenefits)

		recipes = append(recipes, recipe)
	}

	return recipes, total, nil
}

func (s *RecipeService) GetRecipeByID(id string) (*models.Recipe, error) {
	query := `SELECT id, name, cuisine, category, ingredients, instructions, nutrition, macros, 
			  prep_time, cook_time, servings, difficulty, alternatives, skills_tips, 
			  enhancement_tips, health_benefits, created_at, updated_at 
			  FROM recipes WHERE id = $1`

	var recipe models.Recipe
	var ingredientsJSON, instructionsJSON, nutritionJSON, macrosJSON string
	var alternativesJSON, skillsTipsJSON, enhancementTipsJSON, healthBenefitsJSON string

	err := s.db.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.Name, &recipe.Cuisine, &recipe.Category,
		&ingredientsJSON, &instructionsJSON, &nutritionJSON, &macrosJSON,
		&recipe.PrepTime, &recipe.CookTime, &recipe.Servings, &recipe.Difficulty,
		&alternativesJSON, &skillsTipsJSON, &enhancementTipsJSON, &healthBenefitsJSON,
		&recipe.CreatedAt, &recipe.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}

	// Parse JSON fields
	_ = json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients)
	_ = json.Unmarshal([]byte(instructionsJSON), &recipe.Instructions)
	_ = json.Unmarshal([]byte(nutritionJSON), &recipe.Nutrition)
	_ = json.Unmarshal([]byte(macrosJSON), &recipe.Macros)
	_ = json.Unmarshal([]byte(alternativesJSON), &recipe.Alternatives)
	_ = json.Unmarshal([]byte(skillsTipsJSON), &recipe.SkillsTips)
	_ = json.Unmarshal([]byte(enhancementTipsJSON), &recipe.EnhancementTips)
	_ = json.Unmarshal([]byte(healthBenefitsJSON), &recipe.HealthBenefits)

	return &recipe, nil
}

// GetRecipesByUser returns recent recipes for a user (placeholder using filters only)
func (s *RecipeService) GetRecipesByUser(ctx context.Context, userID string, limit int) ([]models.Recipe, error) {
	filters := models.RecipeFilters{Limit: limit, Offset: 0}
	recipes, _, err := s.GetRecipes(filters)
	return recipes, err
}

func (s *RecipeService) GeneratePersonalizedRecipe(req models.PersonalizedRecipeRequest) (*models.Recipe, error) {
	// Get user preferences and restrictions
	filters := models.RecipeFilters{
		Cuisine:     req.Cuisine,
		Category:    req.Category,
		Difficulty:  req.Difficulty,
		MaxCalories: &req.MaxCalories,
		MinProtein:  &req.MinProtein,
		Allergies:   req.Allergies,
		Limit:       10,
		Offset:      0,
	}

	recipes, _, err := s.GetRecipes(filters)
	if err != nil {
		return nil, err
	}

	if len(recipes) == 0 {
		return nil, errors.New("no recipes found matching criteria")
	}

	// Return the first matching recipe (could be enhanced with ML scoring)
	return &recipes[0], nil
}

func (s *RecipeService) SearchRecipes(query string, limit, offset int) ([]models.Recipe, int, error) {
	searchQuery := `SELECT id, name, cuisine, category, ingredients, instructions, nutrition, macros, 
					prep_time, cook_time, servings, difficulty, alternatives, skills_tips, 
					enhancement_tips, health_benefits, created_at, updated_at 
					FROM recipes 
					WHERE name LIKE $1 OR ingredients LIKE $1 OR health_benefits LIKE $1
					ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"

	// Count total
	countQuery := `SELECT COUNT(*) FROM recipes 
				   WHERE name LIKE $1 OR ingredients LIKE $1 OR health_benefits LIKE $1`

	var total int
	err := s.db.QueryRow(countQuery, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var ingredientsJSON, instructionsJSON, nutritionJSON, macrosJSON string
		var alternativesJSON, skillsTipsJSON, enhancementTipsJSON, healthBenefitsJSON string

		err := rows.Scan(
			&recipe.ID, &recipe.Name, &recipe.Cuisine, &recipe.Category,
			&ingredientsJSON, &instructionsJSON, &nutritionJSON, &macrosJSON,
			&recipe.PrepTime, &recipe.CookTime, &recipe.Servings, &recipe.Difficulty,
			&alternativesJSON, &skillsTipsJSON, &enhancementTipsJSON, &healthBenefitsJSON,
			&recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse JSON fields
		_ = json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients)
		_ = json.Unmarshal([]byte(instructionsJSON), &recipe.Instructions)
		_ = json.Unmarshal([]byte(nutritionJSON), &recipe.Nutrition)
		_ = json.Unmarshal([]byte(macrosJSON), &recipe.Macros)
		_ = json.Unmarshal([]byte(alternativesJSON), &recipe.Alternatives)
		_ = json.Unmarshal([]byte(skillsTipsJSON), &recipe.SkillsTips)
		_ = json.Unmarshal([]byte(enhancementTipsJSON), &recipe.EnhancementTips)
		_ = json.Unmarshal([]byte(healthBenefitsJSON), &recipe.HealthBenefits)

		recipes = append(recipes, recipe)
	}

	return recipes, total, nil
}

func (s *RecipeService) GetAvailableCuisines() []models.CuisineInfo {
	return []models.CuisineInfo{
		{ID: "arabian_gulf", Name: "Arabian Gulf", Description: "Traditional Gulf cuisine"},
		{ID: "shami", Name: "Shami", Description: "Levantine cuisine"},
		{ID: "egyptian", Name: "Egyptian", Description: "Traditional Egyptian dishes"},
		{ID: "moroccan", Name: "Moroccan", Description: "North African Moroccan cuisine"},
	}
}

func (s *RecipeService) GetAvailableCategories() []models.CategoryInfo {
	return []models.CategoryInfo{
		{ID: "appetizer", Name: "Appetizer", Description: "Starters and small plates"},
		{ID: "main_course", Name: "Main Course", Description: "Main dishes"},
		{ID: "dessert", Name: "Dessert", Description: "Sweet treats"},
		{ID: "breakfast", Name: "Breakfast", Description: "Morning meals"},
	}
}

func (s *RecipeService) GetNutritionalAnalysis(recipeID string, servings int) (*models.NutritionalAnalysis, error) {
	recipe, err := s.GetRecipeByID(recipeID)
	if err != nil {
		return nil, err
	}

	analysis := &models.NutritionalAnalysis{
		RecipeID:   recipeID,
		Servings:   servings,
		PerServing: recipe.Nutrition,
		Total: models.NutritionInfo{
			Calories: recipe.Nutrition.Calories * servings,
			Protein:  recipe.Nutrition.Protein * float64(servings),
			Carbs:    recipe.Nutrition.Carbs * float64(servings),
			Fat:      recipe.Nutrition.Fat * float64(servings),
			Fiber:    recipe.Nutrition.Fiber * float64(servings),
			Sodium:   recipe.Nutrition.Sodium * float64(servings),
		},
		MacroBreakdown: recipe.Macros,
		HealthScore:    calculateHealthScore(recipe.Nutrition),
		DietaryFlags:   getDietaryFlags(recipe),
	}

	return analysis, nil
}

func (s *RecipeService) GetRecipeAlternatives(recipeID string, allergies, preferences []string) ([]models.Recipe, error) {
	recipe, err := s.GetRecipeByID(recipeID)
	if err != nil {
		return nil, err
	}

	filters := models.RecipeFilters{
		Cuisine:   recipe.Cuisine,
		Category:  recipe.Category,
		Allergies: allergies,
		Limit:     5,
		Offset:    0,
	}

	alternatives, _, err := s.GetRecipes(filters)
	if err != nil {
		return nil, err
	}

	// Remove the original recipe from alternatives
	var filtered []models.Recipe
	for _, alt := range alternatives {
		if alt.ID != recipeID {
			filtered = append(filtered, alt)
		}
	}

	return filtered, nil
}

func calculateHealthScore(nutrition models.NutritionInfo) int {
	score := 50 // Base score

	// Adjust based on nutritional content
	if nutrition.Fiber > 5 {
		score += 10
	}
	if nutrition.Protein > 20 {
		score += 10
	}
	if nutrition.Sodium < 400 {
		score += 10
	}
	if nutrition.Calories < 400 {
		score += 10
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

func getDietaryFlags(recipe *models.Recipe) []string {
	flags := []string{}

	if recipe.Nutrition.Calories < 300 {
		flags = append(flags, "low_calorie")
	}
	if recipe.Nutrition.Protein > 25 {
		flags = append(flags, "high_protein")
	}
	if recipe.Nutrition.Fiber > 8 {
		flags = append(flags, "high_fiber")
	}
	if recipe.Nutrition.Sodium < 300 {
		flags = append(flags, "low_sodium")
	}

	return flags
}
