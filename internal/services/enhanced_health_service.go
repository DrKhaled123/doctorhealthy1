package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// EnhancedHealthService handles all comprehensive health functionality
type EnhancedHealthService struct {
	dataPath string
	data     map[string]interface{}
}

// NewEnhancedHealthService creates a new enhanced health service
func NewEnhancedHealthService(dataPath string) *EnhancedHealthService {
	service := &EnhancedHealthService{
		dataPath: dataPath,
		data:     make(map[string]interface{}),
	}

	// Load enhanced database
	if err := service.loadEnhancedDatabase(); err != nil {
		log.Printf("Warning: Failed to load enhanced database: %v", err)
	}

	return service
}

// loadEnhancedDatabase loads the enhanced health database
func (s *EnhancedHealthService) loadEnhancedDatabase() error {
	// Try to load the enhanced database file
	enhancedPath := filepath.Join(s.dataPath, "ULTIMATE-ENHANCED-HEALTH-DATABASE.js")

	if _, err := os.Stat(enhancedPath); os.IsNotExist(err) {
		// Fallback to existing database
		return s.loadFallbackDatabase()
	}

	// Load enhanced database (simplified for Go - in real implementation, you'd parse the JS file)
	// For now, we'll use the existing structure
	return s.loadFallbackDatabase()
}

// loadFallbackDatabase loads fallback database
func (s *EnhancedHealthService) loadFallbackDatabase() error {
	// Load existing comprehensive database
	comprehensivePath := filepath.Join(s.dataPath, "ULTIMATE-COMPREHENSIVE-HEALTH-DATABASE-COMPLETE.js")

	if _, err := os.Stat(comprehensivePath); os.IsNotExist(err) {
		return fmt.Errorf("no database file found")
	}

	// In a real implementation, you would parse the JavaScript file
	// For now, we'll create a mock structure
	s.data = map[string]interface{}{
		"medications": map[string]interface{}{
			"weight_loss_drugs": []interface{}{
				map[string]interface{}{
					"generic": "Orlistat",
					"brand":   []string{"Xenical", "Alli"},
					"dose":    "120 mg three times daily",
				},
			},
		},
		"complaints": map[string]interface{}{
			"cases": []interface{}{
				map[string]interface{}{
					"id":           1,
					"condition_en": "Cases of Rapid Weight Gain After Eating",
					"condition_ar": "حالات زيادة الوزن بسرعة بعد الأكل",
				},
			},
		},
		"workouts": map[string]interface{}{
			"workout_plans": []interface{}{
				map[string]interface{}{
					"level": "beginner",
					"split": "3_day_full_body",
				},
			},
		},
		"injuries": map[string]interface{}{
			"gym_sprain": map[string]interface{}{
				"title": map[string]string{
					"english": "GymSprain Fix",
					"arabic":  "علاج سريع لالتواء الجيم",
				},
			},
		},
		"diet_plans": map[string]interface{}{
			"diets": []interface{}{
				map[string]interface{}{
					"name": map[string]string{
						"en": "Keto (Ketogenic) Diet",
						"ar": "نظام كيتو (الكيتوجينيك)",
					},
				},
			},
		},
		"recipes": map[string]interface{}{
			"mediterranean": []interface{}{
				map[string]interface{}{
					"name": map[string]string{
						"en": "Mediterranean Quinoa Bowl",
						"ar": "وعاء الكينوا المتوسطي",
					},
				},
			},
		},
		"vitamins_minerals": map[string]interface{}{
			"water_soluble_vitamins": []interface{}{
				map[string]interface{}{
					"name": "Vitamin B1 (Thiamine)",
					"dose": "1.1-1.2 mg daily",
				},
			},
		},
	}

	return nil
}

// GetAllData returns all enhanced health data
func (s *EnhancedHealthService) GetAllData() map[string]interface{} {
	return s.data
}

// GetMedications returns all medications
func (s *EnhancedHealthService) GetMedications() map[string]interface{} {
	if medications, ok := s.data["medications"].(map[string]interface{}); ok {
		return medications
	}
	return make(map[string]interface{})
}

// GetComplaints returns all complaint cases
func (s *EnhancedHealthService) GetComplaints() map[string]interface{} {
	if complaints, ok := s.data["complaints"].(map[string]interface{}); ok {
		return complaints
	}
	return make(map[string]interface{})
}

// GetWorkouts returns all workout plans
func (s *EnhancedHealthService) GetWorkouts() map[string]interface{} {
	if workouts, ok := s.data["workouts"].(map[string]interface{}); ok {
		return workouts
	}
	return make(map[string]interface{})
}

// GetInjuries returns all injury management data
func (s *EnhancedHealthService) GetInjuries() map[string]interface{} {
	if injuries, ok := s.data["injuries"].(map[string]interface{}); ok {
		return injuries
	}
	return make(map[string]interface{})
}

// GetDietPlans returns all diet plans
func (s *EnhancedHealthService) GetDietPlans() map[string]interface{} {
	if dietPlans, ok := s.data["diet_plans"].(map[string]interface{}); ok {
		return dietPlans
	}
	return make(map[string]interface{})
}

// GetRecipes returns all recipes
func (s *EnhancedHealthService) GetRecipes() map[string]interface{} {
	if recipes, ok := s.data["recipes"].(map[string]interface{}); ok {
		return recipes
	}
	return make(map[string]interface{})
}

// GetVitaminsMinerals returns all vitamins and minerals
func (s *EnhancedHealthService) GetVitaminsMinerals() map[string]interface{} {
	if vitamins, ok := s.data["vitamins_minerals"].(map[string]interface{}); ok {
		return vitamins
	}
	return make(map[string]interface{})
}

// GenerateDietPlan generates a personalized diet plan
func (s *EnhancedHealthService) GenerateDietPlan(userData map[string]interface{}, language string) map[string]interface{} {
	// Extract user preferences
	age := 30
	if a, ok := userData["age"].(int); ok {
		age = a
	}

	weight := 70.0
	if w, ok := userData["weight"].(float64); ok {
		weight = w
	}

	height := 170.0
	if h, ok := userData["height"].(float64); ok {
		height = h
	}

	goal := "weight_loss"
	if g, ok := userData["goal"].(string); ok {
		goal = g
	}

	preferredCuisine := "mediterranean"
	if c, ok := userData["preferred_cuisine"].(string); ok {
		preferredCuisine = c
	}

	// Calculate daily calories (simplified)
	bmr := 10*weight + 6.25*height - 5*float64(age) + 5
	activityMultiplier := 1.5
	tdee := bmr * activityMultiplier

	// Generate diet plan based on goal
	var dietPlan map[string]interface{}

	switch goal {
	case "weight_loss":
		calories := tdee * 0.8 // 20% deficit
		dietPlan = s.generateWeightLossPlan(calories, preferredCuisine, language)
	case "muscle_gain":
		calories := tdee * 1.2 // 20% surplus
		dietPlan = s.generateMuscleGainPlan(calories, preferredCuisine, language)
	default:
		dietPlan = s.generateMaintenancePlan(tdee, preferredCuisine, language)
	}

	return dietPlan
}

// generateWeightLossPlan creates a weight loss diet plan
func (s *EnhancedHealthService) generateWeightLossPlan(calories float64, cuisine string, language string) map[string]interface{} {
	protein := calories * 0.3 / 4 // 30% protein
	carbs := calories * 0.4 / 4   // 40% carbs
	fats := calories * 0.3 / 9    // 30% fats

	plan := map[string]interface{}{
		"plan_type":      "weight_loss",
		"daily_calories": int(calories),
		"macros": map[string]interface{}{
			"protein_grams": int(protein),
			"carbs_grams":   int(carbs),
			"fats_grams":    int(fats),
		},
		"meals":           s.generateMeals(calories, cuisine, language),
		"recommendations": s.getWeightLossRecommendations(language),
		"supplements":     s.getWeightLossSupplements(language),
	}

	return plan
}

// generateMuscleGainPlan creates a muscle gain diet plan
func (s *EnhancedHealthService) generateMuscleGainPlan(calories float64, cuisine string, language string) map[string]interface{} {
	protein := calories * 0.35 / 4 // 35% protein
	carbs := calories * 0.45 / 4   // 45% carbs
	fats := calories * 0.2 / 9     // 20% fats

	plan := map[string]interface{}{
		"plan_type":      "muscle_gain",
		"daily_calories": int(calories),
		"macros": map[string]interface{}{
			"protein_grams": int(protein),
			"carbs_grams":   int(carbs),
			"fats_grams":    int(fats),
		},
		"meals":           s.generateMeals(calories, cuisine, language),
		"recommendations": s.getMuscleGainRecommendations(language),
		"supplements":     s.getMuscleGainSupplements(language),
	}

	return plan
}

// generateMaintenancePlan creates a maintenance diet plan
func (s *EnhancedHealthService) generateMaintenancePlan(calories float64, cuisine string, language string) map[string]interface{} {
	protein := calories * 0.25 / 4 // 25% protein
	carbs := calories * 0.5 / 4    // 50% carbs
	fats := calories * 0.25 / 9    // 25% fats

	plan := map[string]interface{}{
		"plan_type":      "maintenance",
		"daily_calories": int(calories),
		"macros": map[string]interface{}{
			"protein_grams": int(protein),
			"carbs_grams":   int(carbs),
			"fats_grams":    int(fats),
		},
		"meals":           s.generateMeals(calories, cuisine, language),
		"recommendations": s.getMaintenanceRecommendations(language),
		"supplements":     s.getMaintenanceSupplements(language),
	}

	return plan
}

// generateMeals creates meal plans
func (s *EnhancedHealthService) generateMeals(calories float64, cuisine string, language string) map[string]interface{} {
	meals := map[string]interface{}{
		"breakfast": s.generateBreakfast(calories*0.25, cuisine, language),
		"lunch":     s.generateLunch(calories*0.35, cuisine, language),
		"dinner":    s.generateDinner(calories*0.3, cuisine, language),
		"snacks":    s.generateSnacks(calories*0.1, cuisine, language),
	}

	return meals
}

// generateBreakfast creates breakfast meal
func (s *EnhancedHealthService) generateBreakfast(calories float64, cuisine string, language string) map[string]interface{} {
	meal := map[string]interface{}{
		"name":     s.getMealName("breakfast", cuisine, language),
		"calories": int(calories),
		"ingredients": []map[string]interface{}{
			{"name": "Oatmeal", "amount": "1 cup", "calories": 150},
			{"name": "Berries", "amount": "1/2 cup", "calories": 40},
			{"name": "Almonds", "amount": "1 tbsp", "calories": 35},
		},
		"preparation":  s.getPreparationMethod("breakfast", language),
		"alternatives": s.getMealAlternatives("breakfast", cuisine, language),
	}

	return meal
}

// generateLunch creates lunch meal
func (s *EnhancedHealthService) generateLunch(calories float64, cuisine string, language string) map[string]interface{} {
	meal := map[string]interface{}{
		"name":     s.getMealName("lunch", cuisine, language),
		"calories": int(calories),
		"ingredients": []map[string]interface{}{
			{"name": "Grilled Chicken", "amount": "150g", "calories": 200},
			{"name": "Quinoa", "amount": "1/2 cup", "calories": 110},
			{"name": "Mixed Vegetables", "amount": "1 cup", "calories": 50},
		},
		"preparation":  s.getPreparationMethod("lunch", language),
		"alternatives": s.getMealAlternatives("lunch", cuisine, language),
	}

	return meal
}

// generateDinner creates dinner meal
func (s *EnhancedHealthService) generateDinner(calories float64, cuisine string, language string) map[string]interface{} {
	meal := map[string]interface{}{
		"name":     s.getMealName("dinner", cuisine, language),
		"calories": int(calories),
		"ingredients": []map[string]interface{}{
			{"name": "Salmon", "amount": "150g", "calories": 250},
			{"name": "Sweet Potato", "amount": "1 medium", "calories": 100},
			{"name": "Broccoli", "amount": "1 cup", "calories": 30},
		},
		"preparation":  s.getPreparationMethod("dinner", language),
		"alternatives": s.getMealAlternatives("dinner", cuisine, language),
	}

	return meal
}

// generateSnacks creates snack options
func (s *EnhancedHealthService) generateSnacks(calories float64, cuisine string, language string) []map[string]interface{} {
	snacks := []map[string]interface{}{
		{
			"name":     s.getSnackName(1, language),
			"calories": int(calories / 2),
			"ingredients": []map[string]interface{}{
				{"name": "Greek Yogurt", "amount": "1 cup", "calories": 100},
				{"name": "Honey", "amount": "1 tsp", "calories": 20},
			},
		},
		{
			"name":     s.getSnackName(2, language),
			"calories": int(calories / 2),
			"ingredients": []map[string]interface{}{
				{"name": "Apple", "amount": "1 medium", "calories": 80},
				{"name": "Almonds", "amount": "10 pieces", "calories": 70},
			},
		},
	}

	return snacks
}

// GenerateWorkoutPlan generates a personalized workout plan
func (s *EnhancedHealthService) GenerateWorkoutPlan(userData map[string]interface{}, language string) map[string]interface{} {
	goal := "general_fitness"
	if g, ok := userData["goal"].(string); ok {
		goal = g
	}

	level := "beginner"
	if l, ok := userData["level"].(string); ok {
		level = l
	}

	injuries := []string{}
	if i, ok := userData["injuries"].([]string); ok {
		injuries = i
	}

	workoutPlan := map[string]interface{}{
		"goal":               goal,
		"level":              level,
		"duration_weeks":     4,
		"frequency_per_week": 3,
		"workouts":           s.generateWorkoutSessions(goal, level, injuries, language),
		"warmup":             s.generateWarmupRoutine(language),
		"cooldown":           s.generateCooldownRoutine(language),
		"nutrition_timing":   s.generateNutritionTiming(language),
		"supplements":        s.generateWorkoutSupplements(language),
		"injury_prevention":  s.generateInjuryPreventionTips(injuries, language),
	}

	return workoutPlan
}

// generateWorkoutSessions creates workout sessions
func (s *EnhancedHealthService) generateWorkoutSessions(goal, level string, injuries []string, language string) []map[string]interface{} {
	sessions := []map[string]interface{}{}

	// Generate 3 sessions per week for 4 weeks
	for week := 1; week <= 4; week++ {
		for day := 1; day <= 3; day++ {
			session := map[string]interface{}{
				"week":            week,
				"day":             day,
				"focus":           s.getWorkoutFocus(day, language),
				"exercises":       s.generateExercises(goal, level, injuries, language),
				"sets_reps":       s.generateSetsReps(level, language),
				"rest_periods":    s.generateRestPeriods(level, language),
				"common_mistakes": s.getCommonMistakes(goal, language),
				"progression":     s.getProgressionTips(week, language),
			}
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// GenerateLifestylePlan generates a lifestyle management plan
func (s *EnhancedHealthService) GenerateLifestylePlan(userData map[string]interface{}, language string) map[string]interface{} {
	diseases := []string{}
	if d, ok := userData["diseases"].([]string); ok {
		diseases = d
	}

	complaints := []string{}
	if c, ok := userData["complaints"].([]string); ok {
		complaints = c
	}

	medications := []string{}
	if m, ok := userData["medications"].([]string); ok {
		medications = m
	}

	lifestylePlan := map[string]interface{}{
		"diseases_management":       s.generateDiseaseManagement(diseases, language),
		"complaints_guidance":       s.generateComplaintsGuidance(complaints, language),
		"medication_guidance":       s.generateMedicationGuidance(medications, language),
		"nutrition_recommendations": s.generateDiseaseNutrition(diseases, language),
		"lifestyle_modifications":   s.generateLifestyleModifications(diseases, language),
		"supplements":               s.generateDiseaseSupplements(diseases, language),
		"monitoring":                s.generateMonitoringPlan(diseases, language),
		"emergency_contacts":        s.getEmergencyContacts(language),
	}

	return lifestylePlan
}

// GenerateRecipes generates personalized recipes
func (s *EnhancedHealthService) GenerateRecipes(userData map[string]interface{}, language string) map[string]interface{} {
	cuisine := "mediterranean"
	if c, ok := userData["preferred_cuisine"].(string); ok {
		cuisine = c
	}

	dislikes := []string{}
	if d, ok := userData["food_dislikes"].([]string); ok {
		dislikes = d
	}

	allergies := []string{}
	if a, ok := userData["allergies"].([]string); ok {
		allergies = a
	}

	recipes := map[string]interface{}{
		"cuisine":                  cuisine,
		"recipes":                  s.generateRecipeList(cuisine, dislikes, allergies, language),
		"cooking_skills":           s.getCookingSkills(language),
		"cooking_tips":             s.getCookingTips(language),
		"ingredient_substitutions": s.getIngredientSubstitutions(language),
		"meal_prep_tips":           s.getMealPrepTips(language),
		"storage_guidelines":       s.getStorageGuidelines(language),
	}

	return recipes
}

// Helper methods for generating content
func (s *EnhancedHealthService) getMealName(mealType, cuisine, language string) string {
	if language == "ar" {
		switch mealType {
		case "breakfast":
			return "وجبة الإفطار"
		case "lunch":
			return "وجبة الغداء"
		case "dinner":
			return "وجبة العشاء"
		}
	}

	switch mealType {
	case "breakfast":
		return "Healthy Breakfast"
	case "lunch":
		return "Balanced Lunch"
	case "dinner":
		return "Nutritious Dinner"
	}
	return "Meal"
}

func (s *EnhancedHealthService) getPreparationMethod(mealType, language string) string {
	if language == "ar" {
		return "طريقة التحضير: اتبع التعليمات خطوة بخطوة للحصول على أفضل النتائج."
	}
	return "Preparation method: Follow step-by-step instructions for best results."
}

func (s *EnhancedHealthService) getMealAlternatives(mealType, cuisine, language string) []map[string]interface{} {
	alternatives := []map[string]interface{}{
		{
			"name":        "Alternative 1",
			"calories":    300,
			"ingredients": []string{"Ingredient 1", "Ingredient 2"},
		},
		{
			"name":        "Alternative 2",
			"calories":    350,
			"ingredients": []string{"Ingredient 3", "Ingredient 4"},
		},
	}
	return alternatives
}

func (s *EnhancedHealthService) getSnackName(index int, language string) string {
	if language == "ar" {
		return fmt.Sprintf("وجبة خفيفة %d", index)
	}
	return fmt.Sprintf("Snack %d", index)
}

// Additional helper methods would be implemented here...
func (s *EnhancedHealthService) getWeightLossRecommendations(language string) []string {
	if language == "ar" {
		return []string{
			"اشرب الماء بانتظام",
			"تناول وجبات صغيرة متكررة",
			"تجنب الأطعمة المصنعة",
		}
	}
	return []string{
		"Drink water regularly",
		"Eat small frequent meals",
		"Avoid processed foods",
	}
}

func (s *EnhancedHealthService) getWeightLossSupplements(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":   "Multivitamin",
			"dose":   "1 tablet daily",
			"timing": "With breakfast",
		},
	}
}

func (s *EnhancedHealthService) getMuscleGainRecommendations(language string) []string {
	if language == "ar" {
		return []string{
			"تناول البروتين بعد التمرين",
			"احصل على قسط كاف من النوم",
			"تدرب بانتظام",
		}
	}
	return []string{
		"Eat protein after workout",
		"Get adequate sleep",
		"Train regularly",
	}
}

func (s *EnhancedHealthService) getMuscleGainSupplements(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":   "Whey Protein",
			"dose":   "25g after workout",
			"timing": "Within 30 minutes post-workout",
		},
	}
}

func (s *EnhancedHealthService) getMaintenanceRecommendations(language string) []string {
	if language == "ar" {
		return []string{
			"حافظ على نظام غذائي متوازن",
			"مارس الرياضة بانتظام",
			"راقب صحتك",
		}
	}
	return []string{
		"Maintain balanced diet",
		"Exercise regularly",
		"Monitor your health",
	}
}

func (s *EnhancedHealthService) getMaintenanceSupplements(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":   "Omega-3",
			"dose":   "1000mg daily",
			"timing": "With meals",
		},
	}
}

// Additional helper methods for workouts, lifestyle, and recipes would be implemented here...
func (s *EnhancedHealthService) getWorkoutFocus(day int, language string) string {
	if language == "ar" {
		switch day {
		case 1:
			return "الجزء العلوي من الجسم"
		case 2:
			return "الجزء السفلي من الجسم"
		case 3:
			return "الجسم بالكامل"
		}
	}

	switch day {
	case 1:
		return "Upper Body"
	case 2:
		return "Lower Body"
	case 3:
		return "Full Body"
	}
	return "General"
}

func (s *EnhancedHealthService) generateExercises(goal, level string, injuries []string, language string) []map[string]interface{} {
	exercises := []map[string]interface{}{
		{
			"name": "Push-ups",
			"sets": 3,
			"reps": "8-12",
			"rest": "60 seconds",
		},
		{
			"name": "Squats",
			"sets": 3,
			"reps": "10-15",
			"rest": "60 seconds",
		},
	}
	return exercises
}

func (s *EnhancedHealthService) generateSetsReps(level, language string) map[string]interface{} {
	return map[string]interface{}{
		"beginner":     "2-3 sets, 8-12 reps",
		"intermediate": "3-4 sets, 6-10 reps",
		"advanced":     "4-5 sets, 4-8 reps",
	}
}

func (s *EnhancedHealthService) generateRestPeriods(level, language string) map[string]interface{} {
	return map[string]interface{}{
		"beginner":     "60-90 seconds",
		"intermediate": "90-120 seconds",
		"advanced":     "120-180 seconds",
	}
}

func (s *EnhancedHealthService) getCommonMistakes(goal, language string) []string {
	if language == "ar" {
		return []string{
			"عدم الإحماء قبل التمرين",
			"استخدام أوزان ثقيلة جداً",
			"عدم الحفاظ على الشكل الصحيح",
		}
	}
	return []string{
		"Not warming up before exercise",
		"Using weights that are too heavy",
		"Not maintaining proper form",
	}
}

func (s *EnhancedHealthService) getProgressionTips(week int, language string) string {
	if language == "ar" {
		return fmt.Sprintf("نصائح التقدم للأسبوع %d: زيادة الوزن تدريجياً", week)
	}
	return fmt.Sprintf("Progression tips for week %d: Gradually increase weight", week)
}

// Additional helper methods for lifestyle and recipes...
func (s *EnhancedHealthService) generateDiseaseManagement(diseases []string, language string) map[string]interface{} {
	management := make(map[string]interface{})
	for _, disease := range diseases {
		management[disease] = map[string]interface{}{
			"treatment_plan":    "Consult your healthcare provider",
			"lifestyle_changes": "Follow medical recommendations",
			"monitoring":        "Regular check-ups required",
		}
	}
	return management
}

func (s *EnhancedHealthService) generateComplaintsGuidance(complaints []string, language string) map[string]interface{} {
	guidance := make(map[string]interface{})
	for _, complaint := range complaints {
		guidance[complaint] = map[string]interface{}{
			"recommendations":    "Seek medical advice",
			"self_care":          "Rest and hydration",
			"when_to_see_doctor": "If symptoms persist",
		}
	}
	return guidance
}

func (s *EnhancedHealthService) generateMedicationGuidance(medications []string, language string) map[string]interface{} {
	guidance := make(map[string]interface{})
	for _, medication := range medications {
		guidance[medication] = map[string]interface{}{
			"dosage":       "As prescribed by doctor",
			"timing":       "Follow prescription instructions",
			"side_effects": "Monitor for adverse reactions",
		}
	}
	return guidance
}

func (s *EnhancedHealthService) generateDiseaseNutrition(diseases []string, language string) map[string]interface{} {
	nutrition := make(map[string]interface{})
	for _, disease := range diseases {
		nutrition[disease] = map[string]interface{}{
			"recommended_foods": []string{"Vegetables", "Lean proteins", "Whole grains"},
			"foods_to_avoid":    []string{"Processed foods", "Excess sugar", "Trans fats"},
			"meal_timing":       "Regular meal times",
		}
	}
	return nutrition
}

func (s *EnhancedHealthService) generateLifestyleModifications(diseases []string, language string) map[string]interface{} {
	modifications := make(map[string]interface{})
	for _, disease := range diseases {
		modifications[disease] = map[string]interface{}{
			"exercise":          "Regular physical activity as tolerated",
			"sleep":             "7-9 hours per night",
			"stress_management": "Meditation and relaxation techniques",
		}
	}
	return modifications
}

func (s *EnhancedHealthService) generateDiseaseSupplements(diseases []string, language string) map[string]interface{} {
	supplements := make(map[string]interface{})
	for _, disease := range diseases {
		supplements[disease] = []map[string]interface{}{
			{
				"name":    "Multivitamin",
				"dose":    "As recommended by healthcare provider",
				"purpose": "General health support",
			},
		}
	}
	return supplements
}

func (s *EnhancedHealthService) generateMonitoringPlan(diseases []string, language string) map[string]interface{} {
	monitoring := make(map[string]interface{})
	for _, disease := range diseases {
		monitoring[disease] = map[string]interface{}{
			"frequency":       "As recommended by healthcare provider",
			"parameters":      []string{"Symptoms", "Medication adherence", "Lifestyle factors"},
			"emergency_signs": "Seek immediate medical attention if severe symptoms occur",
		}
	}
	return monitoring
}

func (s *EnhancedHealthService) getEmergencyContacts(language string) []map[string]interface{} {
	if language == "ar" {
		return []map[string]interface{}{
			{
				"type":        "طوارئ طبية",
				"number":      "911",
				"description": "للحالات الطبية الطارئة",
			},
		}
	}
	return []map[string]interface{}{
		{
			"type":        "Medical Emergency",
			"number":      "911",
			"description": "For medical emergencies",
		},
	}
}

func (s *EnhancedHealthService) generateRecipeList(cuisine string, dislikes, allergies []string, language string) []map[string]interface{} {
	recipes := []map[string]interface{}{
		{
			"name":         "Recipe 1",
			"cuisine":      cuisine,
			"prep_time":    30,
			"cook_time":    45,
			"servings":     4,
			"calories":     350,
			"difficulty":   "easy",
			"ingredients":  []string{"Ingredient 1", "Ingredient 2"},
			"instructions": "Step-by-step cooking instructions",
		},
	}
	return recipes
}

func (s *EnhancedHealthService) getCookingSkills(language string) []string {
	if language == "ar" {
		return []string{
			"مهارات السكين الأساسية",
			"التحكم في الحرارة",
			"تقنيات الطهي المختلفة",
		}
	}
	return []string{
		"Basic knife skills",
		"Heat control",
		"Different cooking techniques",
	}
}

func (s *EnhancedHealthService) getCookingTips(language string) []string {
	if language == "ar" {
		return []string{
			"استخدم مكونات طازجة",
			"تذوق الطعام أثناء الطهي",
			"نظف المطبخ بعد كل استخدام",
		}
	}
	return []string{
		"Use fresh ingredients",
		"Taste food while cooking",
		"Clean kitchen after each use",
	}
}

func (s *EnhancedHealthService) getIngredientSubstitutions(language string) map[string]interface{} {
	if language == "ar" {
		return map[string]interface{}{
			"الزبدة":        "زيت الزيتون",
			"السكر":         "العسل",
			"الدقيق الأبيض": "دقيق القمح الكامل",
		}
	}
	return map[string]interface{}{
		"Butter":      "Olive oil",
		"Sugar":       "Honey",
		"White flour": "Whole wheat flour",
	}
}

func (s *EnhancedHealthService) getMealPrepTips(language string) []string {
	if language == "ar" {
		return []string{
			"خطط للوجبات مسبقاً",
			"احضر المكونات في عطلة نهاية الأسبوع",
			"استخدم حاويات مناسبة للتخزين",
		}
	}
	return []string{
		"Plan meals in advance",
		"Prep ingredients on weekends",
		"Use proper storage containers",
	}
}

func (s *EnhancedHealthService) getStorageGuidelines(language string) map[string]interface{} {
	if language == "ar" {
		return map[string]interface{}{
			"الثلاجة": "3-5 أيام",
			"الفريزر": "2-3 أشهر",
			"الخزانة": "6-12 شهر",
		}
	}
	return map[string]interface{}{
		"Refrigerator": "3-5 days",
		"Freezer":      "2-3 months",
		"Pantry":       "6-12 months",
	}
}

func (s *EnhancedHealthService) generateWarmupRoutine(language string) map[string]interface{} {
	if language == "ar" {
		return map[string]interface{}{
			"duration":  "5-10 دقائق",
			"exercises": []string{"المشي في المكان", "تمارين الإطالة", "تحريك المفاصل"},
		}
	}
	return map[string]interface{}{
		"duration":  "5-10 minutes",
		"exercises": []string{"Walking in place", "Stretching", "Joint mobility"},
	}
}

func (s *EnhancedHealthService) generateCooldownRoutine(language string) map[string]interface{} {
	if language == "ar" {
		return map[string]interface{}{
			"duration":  "5-10 دقائق",
			"exercises": []string{"تمارين الإطالة", "التنفس العميق", "الاسترخاء"},
		}
	}
	return map[string]interface{}{
		"duration":  "5-10 minutes",
		"exercises": []string{"Stretching", "Deep breathing", "Relaxation"},
	}
}

func (s *EnhancedHealthService) generateNutritionTiming(language string) map[string]interface{} {
	if language == "ar" {
		return map[string]interface{}{
			"قبل_التمرين": "تناول وجبة خفيفة قبل 30-60 دقيقة",
			"بعد_التمرين": "تناول البروتين والكربوهيدرات خلال 30 دقيقة",
		}
	}
	return map[string]interface{}{
		"pre_workout":  "Eat light snack 30-60 minutes before",
		"post_workout": "Eat protein and carbs within 30 minutes",
	}
}

func (s *EnhancedHealthService) generateWorkoutSupplements(language string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":   "Creatine",
			"dose":   "3-5g daily",
			"timing": "Any time of day",
		},
	}
}

func (s *EnhancedHealthService) generateInjuryPreventionTips(injuries []string, language string) []string {
	if language == "ar" {
		return []string{
			"قم بالإحماء قبل التمرين",
			"استخدم الشكل الصحيح",
			"لا تفرط في التدريب",
		}
	}
	return []string{
		"Warm up before exercise",
		"Use proper form",
		"Don't overtrain",
	}
}
