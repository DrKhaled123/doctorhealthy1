package models

import (
	"time"
)

// User represents a user profile with health data
type User struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email            string    `json:"email" db:"email" validate:"required,email"`
	Age              int       `json:"age" db:"age" validate:"required,min=13,max=120"`
	Weight           float64   `json:"weight" db:"weight" validate:"required,min=30,max=300"`
	Height           float64   `json:"height" db:"height" validate:"required,min=100,max=250"`
	Gender           string    `json:"gender" db:"gender" validate:"required,oneof=male female"`
	ActivityLevel    string    `json:"activity_level" db:"activity_level" validate:"required,oneof=sedentary light moderate active very_active"`
	MetabolicRate    string    `json:"metabolic_rate" db:"metabolic_rate" validate:"required,oneof=low medium high"`
	Goal             string    `json:"goal" db:"goal" validate:"required"`
	FoodDislikes     []string  `json:"food_dislikes" db:"food_dislikes"`
	Allergies        []string  `json:"allergies" db:"allergies"`
	Diseases         []string  `json:"diseases" db:"diseases"`
	Medications      []string  `json:"medications" db:"medications"`
	PreferredCuisine string    `json:"preferred_cuisine" db:"preferred_cuisine"`
	Language         string    `json:"language" db:"language" validate:"oneof=en ar"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// NutritionPlan represents a generated nutrition plan
type NutritionPlan struct {
	ID                string       `json:"id" db:"id"`
	UserID            string       `json:"user_id" db:"user_id"`
	CaloriesPerDay    int          `json:"calories_per_day" db:"calories_per_day"`
	ProteinGrams      float64      `json:"protein_grams" db:"protein_grams"`
	CarbsGrams        float64      `json:"carbs_grams" db:"carbs_grams"`
	FatsGrams         float64      `json:"fats_grams" db:"fats_grams"`
	PlanType          string       `json:"plan_type" db:"plan_type"`
	WeeklyMeals       []WeeklyMeal `json:"weekly_meals"`
	CalculationMethod string       `json:"calculation_method" db:"calculation_method"`
	Disclaimer        string       `json:"disclaimer" db:"disclaimer"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
}

// WeeklyMeal represents a meal in the weekly plan
type WeeklyMeal struct {
	ID           string       `json:"id"`
	Day          string       `json:"day" validate:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	MealType     string       `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	Name         string       `json:"name" validate:"required"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
	Calories     int          `json:"calories"`
	Protein      float64      `json:"protein"`
	Carbs        float64      `json:"carbs"`
	Fats         float64      `json:"fats"`
	Alternative  *WeeklyMeal  `json:"alternative,omitempty"`
}

// WorkoutPlan represents a generated workout plan
type WorkoutPlan struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	Goal        string     `json:"goal" db:"goal"`
	WorkoutType string     `json:"workout_type" db:"workout_type" validate:"oneof=gym home"`
	Exercises   []Exercise `json:"exercises"`
	Injuries    []string   `json:"injuries"`
	Complaints  []string   `json:"complaints"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// Exercise represents a workout exercise
type Exercise struct {
	ID             string    `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Type           string    `json:"type" validate:"required"`
	Sets           int       `json:"sets" validate:"required,min=1,max=10"`
	Reps           string    `json:"reps" validate:"required"`
	Rest           string    `json:"rest" validate:"required"`
	Instructions   []string  `json:"instructions"`
	CommonMistakes []string  `json:"common_mistakes"`
	Alternative    *Exercise `json:"alternative,omitempty"`
	TargetMuscles  []string  `json:"target_muscles"`
}

// HealthPlan represents a health management plan
type HealthPlan struct {
	ID               string             `json:"id" db:"id"`
	UserID           string             `json:"user_id" db:"user_id"`
	Diseases         []string           `json:"diseases"`
	Medications      []string           `json:"medications"`
	Complaints       []string           `json:"complaints"`
	TreatmentPlan    []TreatmentItem    `json:"treatment_plan"`
	NutritionAdvice  []string           `json:"nutrition_advice"`
	LifestyleChanges []string           `json:"lifestyle_changes"`
	Supplements      []SupplementAdvice `json:"supplements"`
	Disclaimer       string             `json:"disclaimer"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
}

// TreatmentItem represents a treatment recommendation
type TreatmentItem struct {
	Category    string   `json:"category"`
	Advice      []string `json:"advice"`
	Precautions []string `json:"precautions"`
}

// SupplementAdvice represents supplement recommendations
type SupplementAdvice struct {
	Name         string   `json:"name"`
	Dosage       string   `json:"dosage"`
	Frequency    string   `json:"frequency"`
	Instructions string   `json:"instructions"`
	Benefits     []string `json:"benefits"`
}

// Request Models

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Name             string   `json:"name" validate:"required,min=2,max=100"`
	Email            string   `json:"email" validate:"required,email"`
	Age              int      `json:"age" validate:"required,min=13,max=120"`
	Weight           float64  `json:"weight" validate:"required,min=30,max=300"`
	Height           float64  `json:"height" validate:"required,min=100,max=250"`
	Gender           string   `json:"gender" validate:"required,oneof=male female"`
	ActivityLevel    string   `json:"activity_level" validate:"required,oneof=sedentary light moderate active very_active"`
	MetabolicRate    string   `json:"metabolic_rate" validate:"required,oneof=low medium high"`
	Goal             string   `json:"goal" validate:"required"`
	FoodDislikes     []string `json:"food_dislikes"`
	Allergies        []string `json:"allergies"`
	Diseases         []string `json:"diseases"`
	Medications      []string `json:"medications"`
	PreferredCuisine string   `json:"preferred_cuisine"`
	Language         string   `json:"language" validate:"oneof=en ar"`
}

// GenerateNutritionPlanRequest represents nutrition plan generation request
type GenerateNutritionPlanRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	PlanType string `json:"plan_type" validate:"required,oneof=low_carb keto mediterranean dash balanced high_carb vegan anti_inflammatory"`
	Duration int    `json:"duration" validate:"required,min=1,max=4"` // weeks
}

// GenerateWorkoutPlanRequest represents workout plan generation request
type GenerateWorkoutPlanRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Goal        string   `json:"goal" validate:"required"`
	WorkoutType string   `json:"workout_type" validate:"required,oneof=gym home"`
	Injuries    []string `json:"injuries"`
	Complaints  []string `json:"complaints"`
	Duration    int      `json:"duration" validate:"required,min=1,max=4"` // weeks
}

// GenerateHealthPlanRequest represents health plan generation request
type GenerateHealthPlanRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Diseases    []string `json:"diseases"`
	Medications []string `json:"medications"`
	Complaints  []string `json:"complaints"`
}

// GenerateRecipeRequest represents recipe generation request
type GenerateRecipeRequest struct {
	UserID      string `json:"user_id" validate:"required"`
	Cuisine     string `json:"cuisine" validate:"required"`
	MealType    string `json:"meal_type" validate:"oneof=breakfast lunch dinner snack"`
	Difficulty  string `json:"difficulty" validate:"oneof=easy medium hard"`
	MaxCalories int    `json:"max_calories" validate:"min=100,max=2000"`
}

// CalorieCalculation represents calorie calculation details
type CalorieCalculation struct {
	BMI                 float64 `json:"bmi"`
	BMR                 float64 `json:"bmr"`
	TDEE                float64 `json:"tdee"`
	CaloriesPerKg       int     `json:"calories_per_kg"`
	RecommendedCalories int     `json:"recommended_calories"`
	ProteinGrams        float64 `json:"protein_grams"`
	CarbsGrams          float64 `json:"carbs_grams"`
	FatsGrams           float64 `json:"fats_grams"`
	Method              string  `json:"method"`
	Explanation         string  `json:"explanation"`
}

// Available Goals
var AvailableGoals = []string{
	"lose_weight",
	"gain_weight",
	"maintain_weight",
	"build_muscle",
	"improve_strength",
	"body_recomposition",
}

// Available Plan Types
var AvailablePlanTypes = []string{
	"low_carb",
	"keto",
	"mediterranean",
	"dash",
	"balanced",
	"high_carb",
	"vegan",
	"anti_inflammatory",
}

// Available Cuisines
var AvailableCuisines = []string{
	"mediterranean",
	"middle_eastern",
	"asian",
	"indian",
	"mexican",
	"italian",
	"american",
	"french",
	"japanese",
	"thai",
	"greek",
	"turkish",
	"moroccan",
	"lebanese",
	"egyptian",
}

// Medical Disclaimer
const MedicalDisclaimer = `
**MEDICAL DISCLAIMER**: This information is for educational purposes only and is not intended as medical advice. 
Always consult with qualified healthcare professionals before making changes to your diet, exercise routine, or 
taking supplements, especially if you have medical conditions or take medications. Individual results may vary.

**إخلاء المسؤولية الطبية**: هذه المعلومات لأغراض تعليمية فقط وليست نصيحة طبية. 
استشر دائماً مختصين صحيين مؤهلين قبل إجراء تغييرات على نظامك الغذائي أو التمارين أو تناول المكملات، 
خاصة إذا كان لديك حالات طبية أو تتناول أدوية. النتائج الفردية قد تختلف.
`
