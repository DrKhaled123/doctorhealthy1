package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"api-key-generator/internal/models"
)

// RawRecipe represents the JSON structure from the recipe files
type RawRecipe struct {
	ID              interface{}       `json:"id"` // Can be int or string
	Name            string            `json:"name"`
	Cuisine         string            `json:"cuisine"`
	MealType        string            `json:"meal_type,omitempty"` // Old format
	Category        string            `json:"category,omitempty"`  // New format
	Difficulty      string            `json:"difficulty"`
	PrepTime        int               `json:"prep_time"`
	CookTime        int               `json:"cook_time"`
	Servings        int               `json:"servings"`
	Calories        int               `json:"calories,omitempty"`
	Ingredients     []interface{}     `json:"ingredients"`            // Can be old or new format
	Preparation     []string          `json:"preparation,omitempty"`  // Old format
	Instructions    []string          `json:"instructions,omitempty"` // New format
	SkillsTips      []string          `json:"skills_tips"`
	EnhancementTips []string          `json:"enhancement_tips"`
	NutritionalInfo RawNutritionInfo  `json:"nutritional_info,omitempty"` // Old format
	Nutrition       RawNutritionInfo  `json:"nutrition,omitempty"`        // New format
	Macros          RawMacroBreakdown `json:"macros,omitempty"`           // New format
	HealthBenefits  interface{}       `json:"health_benefits"`            // Can be string or []string
}

type RawIngredient struct {
	Item   string  `json:"item"`
	Amount string  `json:"amount"`
	Grams  float64 `json:"grams,omitempty"`
	Ml     float64 `json:"ml,omitempty"`
}

type RawNutritionInfo struct {
	Protein float64 `json:"protein"`
	Carbs   float64 `json:"carbs"`
	Fat     float64 `json:"fat"`
	Fiber   float64 `json:"fiber"`
	Sodium  float64 `json:"sodium"`
}

type RawMacroBreakdown struct {
	Protein int `json:"protein"`
	Carbs   int `json:"carbs"`
	Fat     int `json:"fat"`
}

type RecipeLoader struct {
	db *sql.DB
}

func NewRecipeLoader(db *sql.DB) *RecipeLoader {
	return &RecipeLoader{db: db}
}

func (rl *RecipeLoader) LoadAllRecipes(recipesDir string) error {
	// Get all recipe JSON files
	files, err := filepath.Glob(filepath.Join(recipesDir, "recipes_*.json"))
	if err != nil {
		return fmt.Errorf("failed to find recipe files: %w", err)
	}

	totalLoaded := 0
	for _, file := range files {
		count, err := rl.loadRecipeFile(file)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", file, err)
		}
		totalLoaded += count
		fmt.Printf("Loaded %d recipes from %s\n", count, filepath.Base(file))
	}

	fmt.Printf("Total recipes loaded: %d\n", totalLoaded)
	return nil
}

func (rl *RecipeLoader) loadRecipeFile(filename string) (int, error) {
	data, err := os.ReadFile(filename) // nosec G304 - filename comes from controlled glob pattern
	if err != nil {
		return 0, err
	}

	var rawRecipes []RawRecipe
	if err := json.Unmarshal(data, &rawRecipes); err != nil {
		return 0, err
	}

	count := 0
	for _, rawRecipe := range rawRecipes {
		recipe := rl.convertRawRecipe(rawRecipe)
		if err := rl.insertRecipe(recipe); err != nil {
			fmt.Printf("Warning: failed to insert recipe %s: %v\n", recipe.ID, err)
			continue
		}
		count++
	}

	return count, nil
}

func (rl *RecipeLoader) convertRawRecipe(raw RawRecipe) models.Recipe {
	// Convert ID
	var id string
	switch v := raw.ID.(type) {
	case int:
		id = strconv.Itoa(v)
	case string:
		id = v
	default:
		id = "unknown"
	}

	// Determine category
	category := raw.Category
	if category == "" {
		category = raw.MealType
	}

	// Convert ingredients
	ingredients := make([]models.Ingredient, len(raw.Ingredients))
	for i, ingInterface := range raw.Ingredients {
		switch ing := ingInterface.(type) {
		case map[string]interface{}:
			// New format: {"name": "rice", "amount": 200, "unit": "g"}
			name, _ := ing["name"].(string)
			amount, _ := ing["amount"].(float64)
			unit, _ := ing["unit"].(string)
			ingredients[i] = models.Ingredient{
				Name:   name,
				Amount: amount,
				Unit:   unit,
			}
		case RawIngredient:
			// Old format
			unit := "pieces"
			if ing.Grams > 0 {
				unit = "g"
			} else if ing.Ml > 0 {
				unit = "ml"
			}
			ingredients[i] = models.Ingredient{
				Name:   ing.Item,
				Amount: ing.Grams + ing.Ml,
				Unit:   unit,
			}
		}
	}

	// Determine instructions
	instructions := raw.Instructions
	if len(instructions) == 0 {
		instructions = raw.Preparation
	}

	// Convert nutrition info
	var nutrition models.NutritionInfo
	if raw.Nutrition.Protein > 0 || raw.Nutrition.Carbs > 0 {
		// New format
		nutrition = models.NutritionInfo{
			Calories: raw.Calories,
			Protein:  raw.Nutrition.Protein,
			Carbs:    raw.Nutrition.Carbs,
			Fat:      raw.Nutrition.Fat,
			Fiber:    raw.Nutrition.Fiber,
			Sodium:   raw.Nutrition.Sodium,
		}
	} else {
		// Old format
		nutrition = models.NutritionInfo{
			Calories: raw.Calories,
			Protein:  raw.NutritionalInfo.Protein,
			Carbs:    raw.NutritionalInfo.Carbs,
			Fat:      raw.NutritionalInfo.Fat,
			Fiber:    raw.NutritionalInfo.Fiber,
			Sodium:   raw.NutritionalInfo.Sodium,
		}
	}

	// Convert macros
	var macros models.MacroBreakdown
	if raw.Macros.Protein > 0 || raw.Macros.Carbs > 0 || raw.Macros.Fat > 0 {
		// New format
		macros = models.MacroBreakdown{
			Protein: raw.Macros.Protein,
			Carbs:   raw.Macros.Carbs,
			Fat:     raw.Macros.Fat,
		}
	} else {
		// Calculate from nutrition (old format)
		totalMacros := raw.NutritionalInfo.Protein + raw.NutritionalInfo.Carbs + raw.NutritionalInfo.Fat
		if totalMacros > 0 {
			macros.Protein = int((raw.NutritionalInfo.Protein / totalMacros) * 100)
			macros.Carbs = int((raw.NutritionalInfo.Carbs / totalMacros) * 100)
			macros.Fat = int((raw.NutritionalInfo.Fat / totalMacros) * 100)
		}
	}

	// Convert health benefits
	var healthBenefits []string
	switch v := raw.HealthBenefits.(type) {
	case string:
		healthBenefits = []string{v}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				healthBenefits = append(healthBenefits, str)
			}
		}
	case []string:
		healthBenefits = v
	}

	return models.Recipe{
		ID:              id,
		Name:            raw.Name,
		Cuisine:         raw.Cuisine,
		Category:        category,
		Ingredients:     ingredients,
		Instructions:    instructions,
		Nutrition:       nutrition,
		Macros:          macros,
		PrepTime:        raw.PrepTime,
		CookTime:        raw.CookTime,
		Servings:        raw.Servings,
		Difficulty:      raw.Difficulty,
		SkillsTips:      raw.SkillsTips,
		EnhancementTips: raw.EnhancementTips,
		HealthBenefits:  healthBenefits,
	}
}

func (rl *RecipeLoader) insertRecipe(recipe models.Recipe) error {
	// Set timestamps
	now := time.Now()
	recipe.CreatedAt = now
	recipe.UpdatedAt = now

	// Convert arrays/objects to JSON
	ingredientsJSON, _ := json.Marshal(recipe.Ingredients)
	instructionsJSON, _ := json.Marshal(recipe.Instructions)
	nutritionJSON, _ := json.Marshal(recipe.Nutrition)
	macrosJSON, _ := json.Marshal(recipe.Macros)
	alternativesJSON, _ := json.Marshal(recipe.Alternatives)
	skillsTipsJSON, _ := json.Marshal(recipe.SkillsTips)
	enhancementTipsJSON, _ := json.Marshal(recipe.EnhancementTips)
	healthBenefitsJSON, _ := json.Marshal(recipe.HealthBenefits)

	query := `INSERT OR REPLACE INTO recipes (
		id, name, cuisine, category, ingredients, instructions, nutrition, macros,
		prep_time, cook_time, servings, difficulty, alternatives, skills_tips,
		enhancement_tips, health_benefits, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := rl.db.Exec(query,
		recipe.ID, recipe.Name, recipe.Cuisine, recipe.Category,
		string(ingredientsJSON), string(instructionsJSON), string(nutritionJSON), string(macrosJSON),
		recipe.PrepTime, recipe.CookTime, recipe.Servings, recipe.Difficulty,
		string(alternativesJSON), string(skillsTipsJSON), string(enhancementTipsJSON), string(healthBenefitsJSON),
		recipe.CreatedAt, recipe.UpdatedAt,
	)

	return err
}

func (rl *RecipeLoader) InitializeDatabase() error {
	// Read and execute schema
	schemaPath := filepath.Join("internal", "database", "recipe_schema.sql") // nosec G304 - hardcoded path to internal schema file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split schema into individual statements
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := rl.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}

func (rl *RecipeLoader) GetRecipeStats() (map[string]int, error) {
	stats := make(map[string]int)

	// Total recipes
	var total int
	err := rl.db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// By cuisine
	cuisines := []string{"arabian_gulf", "shami", "egyptian", "moroccan"}
	for _, cuisine := range cuisines {
		var count int
		err := rl.db.QueryRow("SELECT COUNT(*) FROM recipes WHERE cuisine = ?", cuisine).Scan(&count)
		if err != nil {
			return nil, err
		}
		stats[cuisine] = count
	}

	// By category
	categories := []string{"appetizer", "main_course", "dessert", "breakfast"}
	for _, category := range categories {
		var count int
		err := rl.db.QueryRow("SELECT COUNT(*) FROM recipes WHERE category = ?", category).Scan(&count)
		if err != nil {
			return nil, err
		}
		stats[category] = count
	}

	return stats, nil
}
