package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"api-key-generator/internal/models"

	"github.com/google/uuid"
)

// NutritionService handles nutrition plan generation
type NutritionService struct {
	db          *sql.DB
	userService *UserService
}

// NewNutritionService creates a new nutrition service
func NewNutritionService(db *sql.DB, userService *UserService) *NutritionService {
	return &NutritionService{
		db:          db,
		userService: userService,
	}
}

// GenerateNutritionPlan generates a personalized nutrition plan
func (s *NutritionService) GenerateNutritionPlan(ctx context.Context, req *models.GenerateNutritionPlanRequest) (*models.NutritionPlan, error) {
	// Get user data
	user, err := s.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Calculate calories and macros
	calculation, err := s.userService.CalculateCalories(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate calories: %w", err)
	}

	// Adjust macros based on plan type
	s.adjustMacrosForPlanType(calculation, req.PlanType)

	// Generate weekly meals
	weeklyMeals, err := s.generateWeeklyMeals(ctx, user, calculation, req.PlanType, req.Duration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate weekly meals: %w", err)
	}

	// Create nutrition plan
	plan := &models.NutritionPlan{
		ID:                uuid.New().String(),
		UserID:            req.UserID,
		CaloriesPerDay:    calculation.RecommendedCalories,
		ProteinGrams:      calculation.ProteinGrams,
		CarbsGrams:        calculation.CarbsGrams,
		FatsGrams:         calculation.FatsGrams,
		PlanType:          req.PlanType,
		WeeklyMeals:       weeklyMeals,
		CalculationMethod: calculation.Method + " - " + calculation.Explanation,
		Disclaimer:        s.getDisclaimer(user.Language),
		CreatedAt:         time.Now().UTC(),
	}

	// Save to database (optional - for tracking)
	err = s.savePlan(ctx, plan)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save nutrition plan: %v\n", err)
	}

	return plan, nil
}

// generateWeeklyMeals generates meals for the specified duration
func (s *NutritionService) generateWeeklyMeals(ctx context.Context, user *models.User, calc *models.CalorieCalculation, planType string, duration int) ([]models.WeeklyMeal, error) {
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	mealTypes := s.getMealDistribution(user.Goal)

	// Pre-allocate slice for better performance
	allMeals := make([]models.WeeklyMeal, 0, duration*len(days)*len(mealTypes))

	for week := 0; week < duration; week++ {
		for _, day := range days {
			for _, mealType := range mealTypes {
				meal, err := s.generateMeal(ctx, user, calc, planType, day, mealType)
				if err != nil {
					return nil, fmt.Errorf("failed to generate meal for %s %s: %w", day, mealType, err)
				}

				// Generate alternative
				alternative, err := s.generateMeal(ctx, user, calc, planType, day, mealType)
				if err == nil {
					meal.Alternative = alternative
				}

				allMeals = append(allMeals, *meal)
			}
		}
	}

	return allMeals, nil
}

// generateMeal generates a single meal based on user preferences and plan type
func (s *NutritionService) generateMeal(ctx context.Context, user *models.User, calc *models.CalorieCalculation, planType, day, mealType string) (*models.WeeklyMeal, error) {
	// Calculate calories for this meal
	mealCalories := s.calculateMealCalories(calc.RecommendedCalories, mealType)

	// Get meal template based on plan type and cuisine
	mealTemplate := s.getMealTemplate(planType, user.PreferredCuisine, mealType)

	// Filter out disliked foods and allergens
	filteredIngredients := s.filterIngredients(mealTemplate.Ingredients, user.FoodDislikes, user.Allergies)

	// Adjust portions to match calorie target
	adjustedIngredients := s.adjustPortions(filteredIngredients, mealCalories)

	// Calculate actual macros
	protein, carbs, fats := s.calculateMacros(adjustedIngredients)

	meal := &models.WeeklyMeal{
		ID:           uuid.New().String(),
		Day:          day,
		MealType:     mealType,
		Name:         mealTemplate.Name,
		Ingredients:  adjustedIngredients,
		Instructions: mealTemplate.Instructions,
		Calories:     mealCalories,
		Protein:      protein,
		Carbs:        carbs,
		Fats:         fats,
	}

	return meal, nil
}

// getMealDistribution returns meal types based on goal
func (s *NutritionService) getMealDistribution(goal string) []string {
	switch goal {
	case "gain_weight", "build_muscle":
		// 4 meals or 3 meals + 2 snacks
		return []string{"breakfast", "lunch", "dinner", "snack"}
	case "lose_weight":
		// 3 meals + optional snack
		return []string{"breakfast", "lunch", "dinner"}
	default:
		// Balanced approach
		return []string{"breakfast", "lunch", "dinner", "snack"}
	}
}

// calculateMealCalories calculates calories for each meal type
func (s *NutritionService) calculateMealCalories(totalCalories int, mealType string) int {
	switch mealType {
	case "breakfast":
		return int(float64(totalCalories) * 0.25) // 25%
	case "lunch":
		return int(float64(totalCalories) * 0.35) // 35%
	case "dinner":
		return int(float64(totalCalories) * 0.30) // 30%
	case "snack":
		return int(float64(totalCalories) * 0.10) // 10%
	default:
		return int(float64(totalCalories) * 0.25)
	}
}

// getMealTemplate returns a meal template based on plan type and cuisine
func (s *NutritionService) getMealTemplate(planType, cuisine, mealType string) *models.WeeklyMeal {
	// This would normally come from a database or external service
	// For now, we'll use hardcoded templates

	templates := s.getMealTemplates()

	// Find matching template
	for _, template := range templates {
		if (template.MealType == mealType) &&
			(cuisine == "" || strings.Contains(template.Name, cuisine)) {
			// Create local copy to avoid loop pointer issue
			matchedTemplate := template
			return &matchedTemplate
		}
	}

	// Return default template if no match
	return &templates[0]
}

// getMealTemplates returns hardcoded meal templates
func (s *NutritionService) getMealTemplates() []models.WeeklyMeal {
	return []models.WeeklyMeal{
		{
			Name:     "Mediterranean Breakfast Bowl",
			MealType: "breakfast",
			Ingredients: []models.Ingredient{
				{Name: "Greek yogurt", Amount: 200, Unit: "g"},
				{Name: "Oats", Amount: 50, Unit: "g"},
				{Name: "Honey", Amount: 15, Unit: "g"},
				{Name: "Almonds", Amount: 20, Unit: "g"},
				{Name: "Berries", Amount: 100, Unit: "g"},
			},
			Instructions: []string{
				"Mix oats with Greek yogurt",
				"Add honey and mix well",
				"Top with almonds and berries",
				"Serve immediately",
			},
		},
		{
			Name:     "Grilled Chicken Salad",
			MealType: "lunch",
			Ingredients: []models.Ingredient{
				{Name: "Chicken breast", Amount: 150, Unit: "g"},
				{Name: "Mixed greens", Amount: 100, Unit: "g"},
				{Name: "Olive oil", Amount: 15, Unit: "ml"},
				{Name: "Tomatoes", Amount: 100, Unit: "g"},
				{Name: "Cucumber", Amount: 50, Unit: "g"},
			},
			Instructions: []string{
				"Season and grill chicken breast",
				"Prepare salad with mixed greens",
				"Add tomatoes and cucumber",
				"Drizzle with olive oil",
				"Slice chicken and serve on top",
			},
		},
		{
			Name:     "Baked Salmon with Vegetables",
			MealType: "dinner",
			Ingredients: []models.Ingredient{
				{Name: "Salmon fillet", Amount: 150, Unit: "g"},
				{Name: "Broccoli", Amount: 150, Unit: "g"},
				{Name: "Sweet potato", Amount: 100, Unit: "g"},
				{Name: "Olive oil", Amount: 10, Unit: "ml"},
			},
			Instructions: []string{
				"Preheat oven to 200°C",
				"Season salmon with herbs",
				"Cut vegetables and drizzle with olive oil",
				"Bake for 20-25 minutes",
				"Serve hot",
			},
		},
		{
			Name:     "Protein Smoothie",
			MealType: "snack",
			Ingredients: []models.Ingredient{
				{Name: "Protein powder", Amount: 30, Unit: "g"},
				{Name: "Banana", Amount: 100, Unit: "g"},
				{Name: "Almond milk", Amount: 250, Unit: "ml"},
				{Name: "Peanut butter", Amount: 15, Unit: "g"},
			},
			Instructions: []string{
				"Add all ingredients to blender",
				"Blend until smooth",
				"Serve immediately",
			},
		},
	}
}

// filterIngredients removes disliked foods and allergens
func (s *NutritionService) filterIngredients(ingredients []models.Ingredient, dislikes, allergies []string) []models.Ingredient {
	var filtered []models.Ingredient

	for _, ingredient := range ingredients {
		skip := false

		// Check against dislikes
		for _, dislike := range dislikes {
			if strings.Contains(strings.ToLower(ingredient.Name), strings.ToLower(dislike)) {
				skip = true
				break
			}
		}

		// Check against allergies
		if !skip {
			for _, allergy := range allergies {
				if strings.Contains(strings.ToLower(ingredient.Name), strings.ToLower(allergy)) {
					skip = true
					break
				}
			}
		}

		if !skip {
			filtered = append(filtered, ingredient)
		}
	}

	return filtered
}

// adjustPortions adjusts ingredient portions to match target calories
func (s *NutritionService) adjustPortions(ingredients []models.Ingredient, targetCalories int) []models.Ingredient {
	if len(ingredients) == 0 {
		return ingredients
	}

	// Simple adjustment - increase all portions proportionally
	ratio := 1.2 // Default 20% increase
	if targetCalories > 500 {
		ratio = 1.5
	} else if targetCalories < 300 {
		ratio = 0.8
	}

	// Adjust each ingredient
	var adjusted []models.Ingredient
	for _, ingredient := range ingredients {
		newAmount := ingredient.Amount * ratio

		adjusted = append(adjusted, models.Ingredient{
			Name:   ingredient.Name,
			Amount: newAmount,
			Unit:   ingredient.Unit,
		})
	}

	return adjusted
}

// calculateMacros calculates protein, carbs, and fats from ingredients
func (s *NutritionService) calculateMacros(ingredients []models.Ingredient) (protein, carbs, fats float64) {
	// Simplified calculation based on ingredient count and types
	estimatedCalories := len(ingredients) * 100 // Rough estimate

	// Rough estimates based on ingredient types
	protein = float64(estimatedCalories) * 0.25 / 4 // 25% protein
	carbs = float64(estimatedCalories) * 0.45 / 4   // 45% carbs
	fats = float64(estimatedCalories) * 0.30 / 9    // 30% fats

	return protein, carbs, fats
}

// adjustMacrosForPlanType adjusts macro distribution based on plan type
func (s *NutritionService) adjustMacrosForPlanType(calc *models.CalorieCalculation, planType string) {
	totalCalories := float64(calc.RecommendedCalories)

	switch planType {
	case "keto":
		// 70% fat, 25% protein, 5% carbs
		calc.FatsGrams = (totalCalories * 0.70) / 9
		calc.ProteinGrams = (totalCalories * 0.25) / 4
		calc.CarbsGrams = (totalCalories * 0.05) / 4
	case "low_carb":
		// 40% fat, 35% protein, 25% carbs
		calc.FatsGrams = (totalCalories * 0.40) / 9
		calc.ProteinGrams = (totalCalories * 0.35) / 4
		calc.CarbsGrams = (totalCalories * 0.25) / 4
	case "high_carb":
		// 20% fat, 20% protein, 60% carbs
		calc.FatsGrams = (totalCalories * 0.20) / 9
		calc.ProteinGrams = (totalCalories * 0.20) / 4
		calc.CarbsGrams = (totalCalories * 0.60) / 4
	case "mediterranean", "dash", "anti_inflammatory":
		// 30% fat, 25% protein, 45% carbs
		calc.FatsGrams = (totalCalories * 0.30) / 9
		calc.ProteinGrams = (totalCalories * 0.25) / 4
		calc.CarbsGrams = (totalCalories * 0.45) / 4
	default: // balanced
		// Keep original calculation
	}
}

// savePlan saves the nutrition plan to database
func (s *NutritionService) savePlan(ctx context.Context, plan *models.NutritionPlan) error {
	query := `
		INSERT INTO nutrition_plans (
			id, user_id, calories_per_day, protein_grams, carbs_grams, 
			fats_grams, plan_type, calculation_method, disclaimer, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		plan.ID,
		plan.UserID,
		plan.CaloriesPerDay,
		plan.ProteinGrams,
		plan.CarbsGrams,
		plan.FatsGrams,
		plan.PlanType,
		plan.CalculationMethod,
		plan.Disclaimer,
		plan.CreatedAt,
	)

	return err
}

// getDisclaimer returns appropriate disclaimer based on language
func (s *NutritionService) getDisclaimer(language string) string {
	if language == "ar" {
		return `إخلاء المسؤولية الطبية: هذه المعلومات لأغراض تعليمية فقط وليست نصيحة طبية. 
استشر دائماً مختصين صحيين مؤهلين قبل إجراء تغييرات على نظامك الغذائي، خاصة إذا كان لديك حالات طبية أو تتناول أدوية.`
	}

	return `MEDICAL DISCLAIMER: This information is for educational purposes only and is not intended as medical advice. 
Always consult with qualified healthcare professionals before making changes to your diet, especially if you have medical conditions or take medications.`
}
