package services

import (
	"context"
	"fmt"

	"api-key-generator/internal/models"
)

// NutritionTimingService handles workout nutrition and meal timing
type NutritionTimingService struct{}

// NewNutritionTimingService creates a new nutrition timing service
func NewNutritionTimingService() *NutritionTimingService {
	return &NutritionTimingService{}
}

// GenerateWorkoutNutrition creates comprehensive nutrition plan for workouts
func (nts *NutritionTimingService) GenerateWorkoutNutrition(ctx context.Context, req *WorkoutNutritionRequest) (*models.WorkoutNutritionPlan, error) {
	if req == nil {
		return nil, fmt.Errorf("nutrition request is required")
	}

	plan := &models.WorkoutNutritionPlan{
		PreWorkout:    nts.generatePreWorkoutNutrition(req),
		PostWorkout:   nts.generatePostWorkoutNutrition(req),
		HydrationPlan: nts.generateHydrationPlan(req),
		GeneralTips:   nts.getGeneralNutritionTips(req.Goal, req.Language),
	}

	return plan, nil
}

func (nts *NutritionTimingService) generatePreWorkoutNutrition(req *WorkoutNutritionRequest) *models.NutritionTiming {
	preWorkoutMeals := map[string]map[string]*models.NutritionTiming{
		"fat_loss": {
			"hiit": {
				MealTiming:      "pre_workout",
				TimingMinutes:   30,
				Recommendations: "Light carbs for quick energy without digestive stress",
				FoodSuggestions: []string{"Rice cake", "Black coffee", "Small banana"},
				Portions:        map[string]string{"carbs": "15-20g", "caffeine": "100-200mg"},
				Macros:          &models.WorkoutMacroBreakdown{Protein: "0g", Carbs: "15-20g", Fats: "0g", Calories: "60-80"},
				Tips:            []string{"Avoid fats and fiber", "Stay hydrated", "Time caffeine 30min before"},
			},
			"strength": {
				MealTiming:      "pre_workout",
				TimingMinutes:   60,
				Recommendations: "Balanced meal with moderate carbs and protein",
				FoodSuggestions: []string{"Greek yogurt", "Apple", "Almonds"},
				Portions:        map[string]string{"protein": "15-20g", "carbs": "20-25g", "fats": "5-10g"},
				Macros:          &models.WorkoutMacroBreakdown{Protein: "15-20g", Carbs: "20-25g", Fats: "5-10g", Calories: "180-220"},
				Tips:            []string{"Allow 60min digestion", "Include some protein", "Moderate portion size"},
			},
		},
		"muscle_gain": {
			"strength": {
				MealTiming:      "pre_workout",
				TimingMinutes:   90,
				Recommendations: "Substantial meal with complex carbs and quality protein",
				FoodSuggestions: []string{"Oatmeal", "Whey protein", "Banana", "Peanut butter"},
				Portions:        map[string]string{"protein": "25-30g", "carbs": "40-50g", "fats": "10-15g"},
				Macros:          &models.WorkoutMacroBreakdown{Protein: "25-30g", Carbs: "40-50g", Fats: "10-15g", Calories: "350-420"},
				Tips:            []string{"Allow 90min digestion", "Focus on complex carbs", "Include healthy fats"},
			},
		},
	}

	if goalMeals, exists := preWorkoutMeals[req.Goal]; exists {
		if meal, exists := goalMeals[req.WorkoutType]; exists {
			return meal
		}
	}

	// Default pre-workout nutrition
	return &models.NutritionTiming{
		MealTiming:      "pre_workout",
		TimingMinutes:   45,
		Recommendations: "Balanced pre-workout meal for sustained energy",
		FoodSuggestions: []string{"Banana", "Almond butter", "Green tea"},
		Portions:        map[string]string{"carbs": "25-30g", "fats": "8-12g"},
		Macros:          &models.WorkoutMacroBreakdown{Protein: "5g", Carbs: "25-30g", Fats: "8-12g", Calories: "180-220"},
		Tips:            []string{"Easy to digest foods", "Avoid high fiber", "Stay hydrated"},
	}
}

func (nts *NutritionTimingService) generatePostWorkoutNutrition(req *WorkoutNutritionRequest) *models.NutritionTiming {
	postWorkoutMeals := map[string]*models.NutritionTiming{
		"fat_loss": {
			MealTiming:      "post_workout",
			TimingMinutes:   30,
			Recommendations: "High protein with moderate carbs for recovery without excess calories",
			FoodSuggestions: []string{"Whey protein", "Berries", "Spinach", "Egg whites"},
			Portions:        map[string]string{"protein": "25-30g", "carbs": "15-20g", "fats": "0-5g"},
			Macros:          &models.WorkoutMacroBreakdown{Protein: "25-30g", Carbs: "15-20g", Fats: "0-5g", Calories: "160-220"},
			Tips:            []string{"Prioritize protein", "Keep carbs moderate", "Minimize fats"},
		},
		"muscle_gain": {
			MealTiming:      "post_workout",
			TimingMinutes:   30,
			Recommendations: "High protein and carbs for maximum muscle protein synthesis",
			FoodSuggestions: []string{"Chicken breast", "White rice", "Sweet potato", "Olive oil"},
			Portions:        map[string]string{"protein": "35-40g", "carbs": "50-60g", "fats": "10-15g"},
			Macros:          &models.WorkoutMacroBreakdown{Protein: "35-40g", Carbs: "50-60g", Fats: "10-15g", Calories: "420-500"},
			Tips:            []string{"Fast-digesting protein", "High glycemic carbs", "Don't skip this meal"},
		},
		"endurance": {
			MealTiming:      "post_workout",
			TimingMinutes:   30,
			Recommendations: "Carb-focused recovery with adequate protein for glycogen replenishment",
			FoodSuggestions: []string{"Chocolate milk", "Banana", "Turkey sandwich", "Honey"},
			Portions:        map[string]string{"protein": "20-25g", "carbs": "40-50g", "fats": "5-10g"},
			Macros:          &models.WorkoutMacroBreakdown{Protein: "20-25g", Carbs: "40-50g", Fats: "5-10g", Calories: "280-360"},
			Tips:            []string{"3:1 or 4:1 carb to protein ratio", "Include some sodium", "Rehydrate adequately"},
		},
	}

	if meal, exists := postWorkoutMeals[req.Goal]; exists {
		return meal
	}

	return postWorkoutMeals["muscle_gain"] // Default
}

func (nts *NutritionTimingService) generateHydrationPlan(req *WorkoutNutritionRequest) *models.HydrationGuidance {
	intensityMultipliers := map[string]float64{
		"low":      1.0,
		"moderate": 1.2,
		"high":     1.5,
		"hiit":     1.8,
	}

	multiplier := intensityMultipliers["moderate"]
	if mult, exists := intensityMultipliers[req.WorkoutType]; exists {
		multiplier = mult
	}

	baseFluid := 500 // ml
	duringFluid := int(float64(baseFluid) * multiplier)
	postFluid := int(float64(baseFluid) * multiplier * 1.5)

	return &models.HydrationGuidance{
		PreWorkout:  fmt.Sprintf("%d-500ml water 2-3 hours before", baseFluid),
		During:      fmt.Sprintf("%d-250ml every 15-20min during workout", duringFluid/4),
		PostWorkout: fmt.Sprintf("%dml within 2 hours post-workout", postFluid),
		Daily:       "35-40ml per kg bodyweight daily",
		Tips: []string{
			"Monitor urine color (pale yellow is ideal)",
			"Add electrolytes for sessions >60min",
			"Weigh yourself before/after to assess fluid loss",
			"Don't wait until thirsty to drink",
		},
	}
}

func (nts *NutritionTimingService) getGeneralNutritionTips(goal, language string) []string {
	tips := map[string]map[string][]string{
		"en": {
			"fat_loss": {
				"Create moderate calorie deficit (300-500 kcal/day)",
				"Prioritize protein (1.8-2.2g/kg bodyweight)",
				"Time carbs around workouts",
				"Include fiber-rich vegetables with every meal",
				"Avoid liquid calories except post-workout",
			},
			"muscle_gain": {
				"Maintain calorie surplus (300-500 kcal/day)",
				"Eat protein every 3-4 hours",
				"Include complex carbs with most meals",
				"Don't fear healthy fats (nuts, avocado, olive oil)",
				"Eat before bed (casein protein or cottage cheese)",
			},
		},
		"ar": {
			"fat_loss": {
				"أنشئ عجزاً معتدلاً في السعرات (300-500 سعرة/يوم)",
				"أعط الأولوية للبروتين (1.8-2.2جم/كجم من وزن الجسم)",
				"وقت الكربوهيدرات حول التمارين",
				"أدرج الخضروات الغنية بالألياف مع كل وجبة",
				"تجنب السعرات السائلة عدا ما بعد التمرين",
			},
			"muscle_gain": {
				"حافظ على فائض السعرات (300-500 سعرة/يوم)",
				"تناول البروتين كل 3-4 ساعات",
				"أدرج الكربوهيدرات المعقدة مع معظم الوجبات",
				"لا تخف من الدهون الصحية (المكسرات، الأفوكادو، زيت الزيتون)",
				"تناول الطعام قبل النوم (بروتين الكازين أو الجبن القريش)",
			},
		},
	}

	lang := "en"
	if language == "ar" {
		lang = "ar"
	}

	if langTips, exists := tips[lang]; exists {
		if goalTips, exists := langTips[goal]; exists {
			return goalTips
		}
	}

	return tips["en"]["muscle_gain"] // Default
}

// GetNutritionRecommendations provides detailed nutrition guidelines
func (nts *NutritionTimingService) GetNutritionRecommendations(ctx context.Context, goal, fitnessLevel string) (*NutritionRecommendations, error) {
	recommendations := map[string]*NutritionRecommendations{
		"fat_loss": {
			CalorieAdjustment: "Deficit of 300-500 kcal/day",
			MacroRatio:        "40% protein, 30% carbs, 30% fat",
			ProteinIntake:     "1.8-2.2g/kg bodyweight",
			CarbTiming:        "Primarily around workouts",
			FatSources:        "Nuts, avocado, olive oil, fatty fish",
			MealFrequency:     "4-5 meals/day",
			Hydration:         "3-4L water/day",
			SpecialNotes:      "Prioritize whole foods, limit processed sugars",
		},
		"muscle_gain": {
			CalorieAdjustment: "Surplus of 300-500 kcal/day",
			MacroRatio:        "30% protein, 50% carbs, 20% fat",
			ProteinIntake:     "2.0-2.5g/kg bodyweight",
			CarbTiming:        "Throughout the day, emphasize post-workout",
			FatSources:        "Whole milk, red meat, nuts, oils",
			MealFrequency:     "5-6 meals/day including before bed",
			Hydration:         "4-5L water/day",
			SpecialNotes:      "Include calorie-dense foods, don't skip meals",
		},
		"endurance": {
			CalorieAdjustment: "Match energy expenditure",
			MacroRatio:        "15% protein, 60% carbs, 25% fat",
			ProteinIntake:     "1.4-1.8g/kg bodyweight",
			CarbTiming:        "Before, during, and after long sessions",
			FatSources:        "Healthy oils, nuts, seeds for sustained energy",
			MealFrequency:     "5-6 smaller meals/day",
			Hydration:         "5-10ml/kg bodyweight/hour during exercise",
			SpecialNotes:      "Focus on glycogen replenishment, include electrolytes",
		},
	}

	if rec, exists := recommendations[goal]; exists {
		return rec, nil
	}

	return recommendations["muscle_gain"], nil // Default
}

// Response types

type NutritionRecommendations struct {
	CalorieAdjustment string `json:"calorie_adjustment"`
	MacroRatio        string `json:"macro_ratio"`
	ProteinIntake     string `json:"protein_intake"`
	CarbTiming        string `json:"carb_timing"`
	FatSources        string `json:"fat_sources"`
	MealFrequency     string `json:"meal_frequency"`
	Hydration         string `json:"hydration"`
	SpecialNotes      string `json:"special_notes"`
}
