package services

import (
	"context"
	"fmt"
	"time"

	"api-key-generator/internal/models"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ComprehensiveWorkoutService handles advanced workout program generation
type ComprehensiveWorkoutService struct {
	nutritionService *NutritionTimingService
}

// NewComprehensiveWorkoutService creates a new comprehensive workout service
func NewComprehensiveWorkoutService() *ComprehensiveWorkoutService {
	return &ComprehensiveWorkoutService{
		nutritionService: NewNutritionTimingService(),
	}
}

// WorkoutProgramData represents comprehensive workout program information
type WorkoutProgramData struct {
	ID                  string                  `json:"id"`
	NameEnglish         string                  `json:"name_english"`
	NameArabic          string                  `json:"name_arabic"`
	Purpose             string                  `json:"purpose"`
	WorkoutDescription  string                  `json:"workouts"`
	TechniqueGuidance   string                  `json:"mistakes_to_avoid_and_technique"`
	SupplementProtocol  map[string]string       `json:"supplements_and_vitamins_doses"`
	NutritionTiming     *WorkoutNutritionTiming `json:"pre_and_post_workout_meals"`
	WarmupRecovery      string                  `json:"warm_up_and_recovery_tips"`
	NutritionGuidelines string                  `json:"nutrition_recommendations"`
	CreatedAt           time.Time               `json:"created_at"`
}

type WorkoutNutritionTiming struct {
	PreWorkout  string `json:"pre"`
	PostWorkout string `json:"post"`
}

// GetComprehensivePrograms returns all available comprehensive workout programs
func (cws *ComprehensiveWorkoutService) GetComprehensivePrograms(ctx context.Context) ([]*WorkoutProgramData, error) {
	programs := []*WorkoutProgramData{
		{
			ID:                  "full_body_circuit",
			NameEnglish:         "Full-Body Circuit Training",
			NameArabic:          "تمارين الدائرة للجسم بالكامل",
			Purpose:             "Simultaneous fat loss and muscle gain for body reshaping",
			WorkoutDescription:  "3-4 circuits of 6-8 exercises (squats, push-ups, rows, lunges, planks, shoulder press) with 12-15 reps each, 30-45 sec rest between exercises, 2-3 min rest between circuits. 3 sessions/week.",
			TechniqueGuidance:   "Avoid rushing through movements; prioritize form over speed. Keep core engaged during all exercises. Maintain neutral spine in squats and lunges. Control eccentric (lowering) phase.",
			SupplementProtocol:  map[string]string{"protein": "20-30g post-workout", "creatine": "3-5g/day", "vitamin_d": "1000-2000 IU/day", "omega_3": "1-2g/day"},
			NutritionTiming:     &WorkoutNutritionTiming{PreWorkout: "Banana with almond butter (30g carbs + 10g fat) 45 min before", PostWorkout: "Grilled chicken (30g protein) + sweet potato (40g carbs) + vegetables within 45 min"},
			WarmupRecovery:      "5-10 min dynamic warm-up (arm circles, leg swings, cat-cow). Post-workout: 10 min foam rolling + static stretching (hold 30 sec per muscle). Active recovery days: light cycling or swimming.",
			NutritionGuidelines: "Calorie deficit of 300-500 kcal/day. Macro ratio: 40% protein, 30% carbs, 30% fat. Eat 1.6-2g protein/kg bodyweight daily. Prioritize whole foods: lean meats, fish, eggs, legumes, vegetables, whole grains.",
			CreatedAt:           time.Now(),
		},
		{
			ID:                  "hiit_training",
			NameEnglish:         "High-Intensity Interval Training (HIIT)",
			NameArabic:          "تمرينات الفواصل عالية الكثافة",
			Purpose:             "Maximize calorie burn and fat loss in minimal time",
			WorkoutDescription:  "4-6 cycles of 30 sec all-out effort (sprints, burpees, battle ropes) followed by 90 sec active rest (walking/jogging). Total time: 20-25 min. 2-3 sessions/week.",
			TechniqueGuidance:   "Avoid sacrificing form for speed. Maintain proper posture during sprints (lean forward slightly, drive knees up). For burpees: land softly, keep core tight to protect lower back.",
			SupplementProtocol:  map[string]string{"caffeine": "200mg pre-workout", "beta_alanine": "2-5g/day", "electrolytes": "500mg sodium + 200mg potassium per liter water", "vitamin_c": "500mg/day"},
			NutritionTiming:     &WorkoutNutritionTiming{PreWorkout: "Black coffee + rice cake (20g carbs) 30 min before", PostWorkout: "Whey protein shake (25g) + berries (15g carbs) immediately after"},
			WarmupRecovery:      "Dynamic warm-up: high knees, butt kicks, jumping jacks (5 min). Post-workout: contrast showers (30 sec cold/30 sec hot) + 10 min light cycling. Prioritize sleep: 7-9 hours/night.",
			NutritionGuidelines: "Moderate calorie deficit (500 kcal/day). Time carbs around workouts. Focus on fiber-rich carbs (oats, quinoa) and lean protein. Limit processed sugars. Hydrate: 3-4L water/day.",
			CreatedAt:           time.Now(),
		},
		{
			ID:                  "metabolic_resistance",
			NameEnglish:         "Metabolic Resistance Training (MRT)",
			NameArabic:          "تدريب المقاومة الأيضي",
			Purpose:             "Burn fat while preserving muscle during cutting phases",
			WorkoutDescription:  "Supersets (2 exercises back-to-back) with 8-12 reps each, 45 sec rest between supersets. Example: Dumbbell rows + push-ups, lunges + shoulder press. 4-5 supersets/session. 3 sessions/week.",
			TechniqueGuidance:   "Avoid resting between exercises in a superset. Maintain controlled tempo (2 sec up, 3 sec down). Keep joints aligned (knees over toes in lunges, elbows at 45° in push-ups).",
			SupplementProtocol:  map[string]string{"bcaas": "5-10g during workout", "l_carnitine": "1-2g/day", "zinc_magnesium": "30mg zinc + 450mg magnesium before bed", "multivitamin": "1 serving/day"},
			NutritionTiming:     &WorkoutNutritionTiming{PreWorkout: "Greek yogurt (20g protein) + apple (20g carbs) 60 min before", PostWorkout: "Salmon (30g protein) + asparagus + brown rice (30g carbs) within 30 min"},
			WarmupRecovery:      "Activation warm-up: band pull-aparts, glute bridges, planks (8 min). Post-workout: 15 min deep tissue massage + cold bath (10-15°C for 10 min). Use compression garments for 2 hours post-workout.",
			NutritionGuidelines: "High protein diet (2.2g/kg bodyweight). Moderate carbs (30-40% calories), healthy fats (25-30%). Cycle calories: higher on training days, lower on rest days. Include green vegetables 3x/day.",
			CreatedAt:           time.Now(),
		},
	}

	return programs, nil
}

// GenerateCustomProgram creates a personalized workout program
func (cws *ComprehensiveWorkoutService) GenerateCustomProgram(ctx context.Context, req *CustomProgramRequest) (*models.WorkoutProgram, error) {
	if req == nil {
		return nil, fmt.Errorf("program request is required")
	}

	// Generate base program structure
	program := &models.WorkoutProgram{
		ID:                generateID(),
		NameEn:            fmt.Sprintf("Custom %s Program", cases.Title(language.Und).String(req.Goal)),
		Level:             req.FitnessLevel,
		SplitType:         cws.determineSplitType(req),
		DescriptionEn:     cws.generateDescription(req),
		DurationWeeks:     req.DurationWeeks,
		SessionsPerWeek:   req.SessionsPerWeek,
		Goals:             []string{req.Goal},
		EquipmentRequired: req.AvailableEquipment,
		WeeklyPlans:       make(map[int]*models.WeeklyPlan),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Generate nutrition and supplement plans
	nutritionPlan, err := cws.nutritionService.GenerateWorkoutNutrition(ctx, &WorkoutNutritionRequest{
		Goal:         req.Goal,
		FitnessLevel: req.FitnessLevel,
		WorkoutType:  cws.determineSplitType(req),
		Language:     req.Language,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate nutrition plan: %w", err)
	}

	program.NutritionPlan = nutritionPlan
	program.SupplementPlan = cws.generateSupplementPlan(req)

	// Generate first week plan
	weeklyPlan, err := cws.generateWeeklyPlan(ctx, program, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate weekly plan: %w", err)
	}

	program.WeeklyPlans[1] = weeklyPlan

	return program, nil
}

// Helper methods

func (cws *ComprehensiveWorkoutService) determineSplitType(req *CustomProgramRequest) string {
	if req.SessionsPerWeek <= 3 {
		return "full_body"
	} else if req.SessionsPerWeek == 4 {
		return "upper_lower"
	}
	return "push_pull_legs"
}

func (cws *ComprehensiveWorkoutService) generateDescription(req *CustomProgramRequest) string {
	return fmt.Sprintf("Customized %s program for %s level, %d sessions per week focusing on %s",
		cws.determineSplitType(req), req.FitnessLevel, req.SessionsPerWeek, req.Goal)
}

func (cws *ComprehensiveWorkoutService) generateSupplementPlan(req *CustomProgramRequest) *models.SupplementProtocol {

	// Base supplements for all goals
	supplementSlice := make([]models.SupplementTiming, 0)
	supplementSlice = append(supplementSlice, models.SupplementTiming{
		SupplementName: "Protein Powder",
		Dosage:         "25-30g",
		Unit:           "grams",
		TimingMinutes:  0,
		Instructions:   "Take immediately post-workout",
		Benefits:       []string{"Muscle recovery and growth"},
	})

	supplementSlice = append(supplementSlice, models.SupplementTiming{
		SupplementName: "Creatine",
		Dosage:         "3-5g",
		Unit:           "grams",
		TimingMinutes:  0,
		Instructions:   "Take daily with water",
		Benefits:       []string{"Strength and power enhancement"},
	})

	// Goal-specific supplements
	switch req.Goal {
	case "fat_loss":
		supplementSlice = append(supplementSlice, models.SupplementTiming{
			SupplementName: "L-Carnitine",
			Dosage:         "1-2g",
			Unit:           "grams",
			TimingMinutes:  -30,
			Instructions:   "Take 30 minutes before workout",
			Benefits:       []string{"Enhanced fat oxidation"},
		})
	case "muscle_gain":
		supplementSlice = append(supplementSlice, models.SupplementTiming{
			SupplementName: "Mass Gainer",
			Dosage:         "1 serving",
			Unit:           "scoop",
			TimingMinutes:  0,
			Instructions:   "Take immediately post-workout",
			Benefits:       []string{"Calorie and protein boost"},
		})
	}

	return &models.SupplementProtocol{
		PostWorkout: supplementSlice,
		CreatedAt:   time.Now(),
	}
}

func (cws *ComprehensiveWorkoutService) generateWeeklyPlan(ctx context.Context, program *models.WorkoutProgram, weekNumber int) (*models.WeeklyPlan, error) {
	weeklyPlan := &models.WeeklyPlan{
		ID:                    generateID(),
		ProgramID:             program.ID,
		WeekNumber:            weekNumber,
		IntensityLevel:        cws.calculateIntensity(weekNumber, program.Level),
		VolumeIncreasePercent: cws.calculateVolumeIncrease(weekNumber),
		DailyWorkouts:         make(map[int]*models.DailyWorkout),
		CreatedAt:             time.Now(),
	}

	// Generate daily workouts based on split type
	for day := 1; day <= program.SessionsPerWeek; day++ {
		dailyWorkout := &models.DailyWorkout{
			ID:              generateID(),
			WeeklyPlanID:    weeklyPlan.ID,
			DayNumber:       day,
			FocusEn:         cws.getDayFocus(program.SplitType, day),
			DurationMinutes: cws.calculateDuration(program.Level),
			IsRestDay:       false,
			CreatedAt:       time.Now(),
		}

		// Add nutrition timing
		if program.NutritionPlan != nil {
			dailyWorkout.PreWorkoutMeal = program.NutritionPlan.PreWorkout
			dailyWorkout.PostWorkoutMeal = program.NutritionPlan.PostWorkout
		}

		weeklyPlan.DailyWorkouts[day] = dailyWorkout
	}

	return weeklyPlan, nil
}

func (cws *ComprehensiveWorkoutService) calculateIntensity(weekNumber int, level string) string {
	baseIntensity := map[string]int{"beginner": 60, "intermediate": 70, "advanced": 80}
	base := baseIntensity[level]
	adjusted := base + (weekNumber-1)*2

	if adjusted >= 85 {
		return "high"
	} else if adjusted >= 75 {
		return "moderate-high"
	} else if adjusted >= 65 {
		return "moderate"
	}
	return "low-moderate"
}

func (cws *ComprehensiveWorkoutService) calculateVolumeIncrease(weekNumber int) float64 {
	if weekNumber <= 1 {
		return 0.0
	}
	return float64(weekNumber-1) * 2.5 // 2.5% increase per week
}

func (cws *ComprehensiveWorkoutService) getDayFocus(splitType string, day int) string {
	focuses := map[string][]string{
		"full_body":      {"Full Body A", "Full Body B", "Full Body C"},
		"upper_lower":    {"Upper Body", "Lower Body", "Upper Body", "Lower Body"},
		"push_pull_legs": {"Push", "Pull", "Legs", "Push", "Pull", "Legs"},
	}

	if dayFocuses, exists := focuses[splitType]; exists && day <= len(dayFocuses) {
		return dayFocuses[day-1]
	}
	return "Full Body"
}

func (cws *ComprehensiveWorkoutService) calculateDuration(level string) int {
	durations := map[string]int{"beginner": 45, "intermediate": 60, "advanced": 75}
	if duration, exists := durations[level]; exists {
		return duration
	}
	return 60
}

// Request types

type CustomProgramRequest struct {
	Goal               string   `json:"goal" validate:"required"`
	FitnessLevel       string   `json:"fitness_level" validate:"required"`
	SessionsPerWeek    int      `json:"sessions_per_week" validate:"required,min=2,max=6"`
	DurationWeeks      int      `json:"duration_weeks" validate:"required,min=4,max=16"`
	AvailableEquipment []string `json:"available_equipment"`
	Language           string   `json:"language"`
}

type WorkoutNutritionRequest struct {
	Goal         string `json:"goal"`
	FitnessLevel string `json:"fitness_level"`
	WorkoutType  string `json:"workout_type"`
	Language     string `json:"language"`
}
