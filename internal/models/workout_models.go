package models

import (
	"time"
)

// WorkoutExerciseData represents a single exercise with comprehensive data
type WorkoutExerciseData struct {
	ID               string    `json:"id" db:"id"`
	NameEn           string    `json:"name_en" db:"name_en"`
	NameAr           string    `json:"name_ar,omitempty" db:"name_ar"`
	Category         string    `json:"category" db:"category"`
	PrimaryMuscles   []string  `json:"primary_muscles" db:"primary_muscles"`
	SecondaryMuscles []string  `json:"secondary_muscles,omitempty" db:"secondary_muscles"`
	Equipment        []string  `json:"equipment,omitempty" db:"equipment"`
	Difficulty       string    `json:"difficulty" db:"difficulty"`
	InstructionsEn   []string  `json:"instructions_en" db:"instructions_en"`
	InstructionsAr   []string  `json:"instructions_ar,omitempty" db:"instructions_ar"`
	CommonMistakesEn []string  `json:"common_mistakes_en,omitempty" db:"common_mistakes_en"`
	CommonMistakesAr []string  `json:"common_mistakes_ar,omitempty" db:"common_mistakes_ar"`
	InjuryRisksEn    []string  `json:"injury_risks_en,omitempty" db:"injury_risks_en"`
	InjuryRisksAr    []string  `json:"injury_risks_ar,omitempty" db:"injury_risks_ar"`
	TipsEn           []string  `json:"tips_en,omitempty" db:"tips_en"`
	TipsAr           []string  `json:"tips_ar,omitempty" db:"tips_ar"`
	Alternatives     []string  `json:"alternatives,omitempty" db:"alternatives"`
	Progressions     []string  `json:"progressions,omitempty" db:"progressions"`
	EvidenceLinks    []string  `json:"evidence_links,omitempty" db:"evidence_links"`
	SetsDefault      int       `json:"sets_default" db:"sets_default"`
	RepsDefault      string    `json:"reps_default" db:"reps_default"`
	RestDefault      string    `json:"rest_default" db:"rest_default"`
	IsGymExercise    bool      `json:"is_gym_exercise" db:"is_gym_exercise"`
	IsHomeExercise   bool      `json:"is_home_exercise" db:"is_home_exercise"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// WorkoutProgram represents a structured workout program with multiple purposes
type WorkoutProgram struct {
	ID                string                `json:"id" db:"id"`
	NameEn            string                `json:"name_en" db:"name_en"`
	NameAr            string                `json:"name_ar,omitempty" db:"name_ar"`
	Level             string                `json:"level" db:"level"`
	SplitType         string                `json:"split_type" db:"split_type"`
	DescriptionEn     string                `json:"description_en" db:"description_en"`
	DescriptionAr     string                `json:"description_ar,omitempty" db:"description_ar"`
	DurationWeeks     int                   `json:"duration_weeks" db:"duration_weeks"`
	SessionsPerWeek   int                   `json:"sessions_per_week" db:"sessions_per_week"`
	Goals             []string              `json:"goals" db:"goals"`
	EquipmentRequired []string              `json:"equipment_required,omitempty" db:"equipment_required"`
	TargetAudience    []string              `json:"target_audience,omitempty" db:"target_audience"`
	WeeklyPlans       map[int]*WeeklyPlan   `json:"weekly_plans,omitempty"`
	NutritionPlan     *WorkoutNutritionPlan `json:"nutrition_plan,omitempty"`
	SupplementPlan    *SupplementProtocol   `json:"supplement_plan,omitempty"`
	CreatedAt         time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at" db:"updated_at"`
}

// WeeklyPlan represents a week-specific workout plan with progression tracking
type WeeklyPlan struct {
	ID                    string                `json:"id" db:"id"`
	ProgramID             string                `json:"program_id" db:"program_id"`
	WeekNumber            int                   `json:"week_number" db:"week_number"`
	ProgressionNotesEn    string                `json:"progression_notes_en,omitempty" db:"progression_notes_en"`
	ProgressionNotesAr    string                `json:"progression_notes_ar,omitempty" db:"progression_notes_ar"`
	IntensityLevel        string                `json:"intensity_level" db:"intensity_level"`
	VolumeIncreasePercent float64               `json:"volume_increase_percent" db:"volume_increase_percent"`
	DailyWorkouts         map[int]*DailyWorkout `json:"daily_workouts,omitempty"`
	CreatedAt             time.Time             `json:"created_at" db:"created_at"`
}

// DailyWorkout represents a single day's workout within a weekly plan
type DailyWorkout struct {
	ID               string              `json:"id" db:"id"`
	WeeklyPlanID     string              `json:"weekly_plan_id" db:"weekly_plan_id"`
	DayNumber        int                 `json:"day_number" db:"day_number"`
	FocusEn          string              `json:"focus_en" db:"focus_en"`
	FocusAr          string              `json:"focus_ar,omitempty" db:"focus_ar"`
	DurationMinutes  int                 `json:"duration_minutes" db:"duration_minutes"`
	IsRestDay        bool                `json:"is_rest_day" db:"is_rest_day"`
	WarmupRoutine    *WarmupRoutine      `json:"warmup_routine,omitempty"`
	WorkoutExercises []*WorkoutExercise  `json:"workout_exercises,omitempty"`
	CooldownRoutine  *CooldownRoutine    `json:"cooldown_routine,omitempty"`
	PreWorkoutMeal   *NutritionTiming    `json:"pre_workout_meal,omitempty"`
	PostWorkoutMeal  *NutritionTiming    `json:"post_workout_meal,omitempty"`
	Supplements      []*SupplementTiming `json:"supplements,omitempty"`
	CreatedAt        time.Time           `json:"created_at" db:"created_at"`
}

// WorkoutExercise represents an exercise within a daily workout
type WorkoutExercise struct {
	ID             string               `json:"id" db:"id"`
	DailyWorkoutID string               `json:"daily_workout_id" db:"daily_workout_id"`
	ExerciseID     string               `json:"exercise_id" db:"exercise_id"`
	ExerciseOrder  int                  `json:"exercise_order" db:"exercise_order"`
	Sets           int                  `json:"sets" db:"sets"`
	Reps           string               `json:"reps" db:"reps"`
	RestSeconds    int                  `json:"rest_seconds" db:"rest_seconds"`
	WeightKg       *float64             `json:"weight_kg,omitempty" db:"weight_kg"`
	IntensityNotes string               `json:"intensity_notes,omitempty" db:"intensity_notes"`
	Exercise       *WorkoutExerciseData `json:"exercise,omitempty"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
}

// WarmupRoutine represents a comprehensive warmup routine
type WarmupRoutine struct {
	ID             string            `json:"id" db:"id"`
	ProgramID      string            `json:"program_id" db:"program_id"`
	WorkoutType    string            `json:"workout_type" db:"workout_type"`
	Duration       int               `json:"duration_minutes" db:"duration_minutes"`
	Exercises      []*WarmupExercise `json:"exercises"`
	InstructionsEn string            `json:"instructions_en" db:"instructions_en"`
	InstructionsAr string            `json:"instructions_ar,omitempty" db:"instructions_ar"`
	Tips           []string          `json:"tips,omitempty"`
	SafetyNotes    []string          `json:"safety_notes,omitempty"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}

// WarmupExercise represents a single warmup exercise
type WarmupExercise struct {
	Name         string   `json:"name"`
	Duration     string   `json:"duration"`
	Instructions []string `json:"instructions"`
	Purpose      string   `json:"purpose"`
}

// CooldownRoutine represents post-workout cooldown and stretching
type CooldownRoutine struct {
	Duration     int                `json:"duration_minutes"`
	Stretches    []*StretchExercise `json:"stretches"`
	RecoveryTips []string           `json:"recovery_tips,omitempty"`
}

// StretchExercise represents a stretching exercise
type StretchExercise struct {
	Name         string   `json:"name"`
	Duration     string   `json:"duration"`
	Instructions []string `json:"instructions"`
	TargetMuscle string   `json:"target_muscle"`
}

// WorkoutNutritionPlan represents comprehensive nutrition guidance with timing
type WorkoutNutritionPlan struct {
	PreWorkout    *NutritionTiming   `json:"pre_workout,omitempty"`
	PostWorkout   *NutritionTiming   `json:"post_workout,omitempty"`
	HydrationPlan *HydrationGuidance `json:"hydration_plan,omitempty"`
	GeneralTips   []string           `json:"general_tips,omitempty"`
}

// NutritionTiming represents meal timing recommendations
type NutritionTiming struct {
	ID                string                 `json:"id,omitempty" db:"id"`
	ProgramID         string                 `json:"program_id,omitempty" db:"program_id"`
	MealTiming        string                 `json:"meal_timing" db:"meal_timing"`
	TimingMinutes     int                    `json:"timing_minutes" db:"timing_minutes"`
	Recommendations   string                 `json:"recommendations_en" db:"recommendations_en"`
	RecommendationsAr string                 `json:"recommendations_ar,omitempty" db:"recommendations_ar"`
	FoodSuggestions   []string               `json:"food_suggestions,omitempty" db:"food_suggestions"`
	Portions          map[string]string      `json:"portions,omitempty"`
	Macros            *WorkoutMacroBreakdown `json:"macros,omitempty"`
	Tips              []string               `json:"tips,omitempty"`
	CreatedAt         time.Time              `json:"created_at,omitempty" db:"created_at"`
}

// WorkoutMacroBreakdown represents macronutrient breakdown
type WorkoutMacroBreakdown struct {
	Protein  string `json:"protein"`
	Carbs    string `json:"carbs"`
	Fats     string `json:"fats"`
	Calories string `json:"calories"`
	Fiber    string `json:"fiber,omitempty"`
}

// HydrationGuidance represents hydration recommendations
type HydrationGuidance struct {
	PreWorkout  string   `json:"pre_workout"`
	During      string   `json:"during"`
	PostWorkout string   `json:"post_workout"`
	Daily       string   `json:"daily"`
	Tips        []string `json:"tips,omitempty"`
}

// WorkoutProgress represents user progress tracking
type WorkoutProgress struct {
	ID                   string                       `json:"id" db:"id"`
	UserID               string                       `json:"user_id" db:"user_id"`
	ProgramID            string                       `json:"program_id" db:"program_id"`
	CurrentWeek          int                          `json:"current_week" db:"current_week"`
	CompletedSessions    int                          `json:"completed_sessions" db:"completed_sessions"`
	TotalSessions        int                          `json:"total_sessions" db:"total_sessions"`
	StartDate            time.Time                    `json:"start_date" db:"start_date"`
	LastSessionDate      *time.Time                   `json:"last_session_date,omitempty" db:"last_session_date"`
	CompletionPercentage float64                      `json:"completion_percentage" db:"completion_percentage"`
	Status               string                       `json:"status" db:"status"`
	Notes                string                       `json:"notes,omitempty" db:"notes"`
	ExerciseProgress     map[string]*ExerciseProgress `json:"exercise_progress,omitempty"`
	Metrics              *ProgressMetrics             `json:"metrics,omitempty"`
	CreatedAt            time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at" db:"updated_at"`
}

// ExerciseProgress represents progress for individual exercises
type ExerciseProgress struct {
	ExerciseID     string          `json:"exercise_id"`
	Weight         float64         `json:"weight"`
	Reps           int             `json:"reps"`
	Sets           int             `json:"sets"`
	RPE            int             `json:"rpe"` // Rate of Perceived Exertion (1-10)
	Notes          string          `json:"notes,omitempty"`
	LastPerformed  time.Time       `json:"last_performed"`
	PersonalRecord *PersonalRecord `json:"personal_record,omitempty"`
}

// PersonalRecord represents a user's personal record for an exercise
type PersonalRecord struct {
	Weight    float64   `json:"weight"`
	Reps      int       `json:"reps"`
	OneRepMax float64   `json:"one_rep_max"`
	Date      time.Time `json:"date"`
}

// ProgressMetrics represents overall progress metrics
type ProgressMetrics struct {
	StrengthGains     map[string]float64 `json:"strength_gains,omitempty"`
	VolumeProgression float64            `json:"volume_progression"`
	ConsistencyScore  float64            `json:"consistency_score"`
	AverageRPE        float64            `json:"average_rpe"`
	TotalWorkoutTime  int                `json:"total_workout_time_minutes"`
	WeightProgression map[string]float64 `json:"weight_progression,omitempty"`
}

// SessionLog represents a completed workout session
type SessionLog struct {
	ID              string         `json:"id" db:"id"`
	UserID          string         `json:"user_id" db:"user_id"`
	DailyWorkoutID  string         `json:"daily_workout_id" db:"daily_workout_id"`
	SessionDate     time.Time      `json:"session_date" db:"session_date"`
	DurationMinutes int            `json:"duration_minutes,omitempty" db:"duration_minutes"`
	RPEScore        int            `json:"rpe_score,omitempty" db:"rpe_score"`
	Notes           string         `json:"notes,omitempty" db:"notes"`
	Completed       bool           `json:"completed" db:"completed"`
	ExerciseLogs    []*ExerciseLog `json:"exercise_logs,omitempty"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
}

// ExerciseLog represents performance for individual exercises in a session
type ExerciseLog struct {
	ID            string    `json:"id" db:"id"`
	SessionLogID  string    `json:"session_log_id" db:"session_log_id"`
	ExerciseID    string    `json:"exercise_id" db:"exercise_id"`
	SetNumber     int       `json:"set_number" db:"set_number"`
	RepsCompleted int       `json:"reps_completed,omitempty" db:"reps_completed"`
	WeightKg      *float64  `json:"weight_kg,omitempty" db:"weight_kg"`
	RestSeconds   int       `json:"rest_seconds,omitempty" db:"rest_seconds"`
	RPEScore      int       `json:"rpe_score,omitempty" db:"rpe_score"`
	Notes         string    `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ProgramRequest represents a request to generate a custom workout program
type ProgramRequest struct {
	Goals              []string            `json:"goals" validate:"required,min=1"`
	FitnessLevel       string              `json:"fitness_level" validate:"required,oneof=beginner intermediate advanced"`
	AvailableEquipment []string            `json:"available_equipment"`
	SessionsPerWeek    int                 `json:"sessions_per_week" validate:"required,min=1,max=7"`
	SessionDuration    int                 `json:"session_duration_minutes" validate:"required,min=15,max=180"`
	DurationWeeks      int                 `json:"duration_weeks" validate:"required,min=1,max=52"`
	Preferences        *ProgramPreferences `json:"preferences,omitempty"`
}

// ProgramPreferences represents user preferences for program generation
type ProgramPreferences struct {
	SplitType          string   `json:"split_type,omitempty"`
	FocusAreas         []string `json:"focus_areas,omitempty"`
	AvoidExercises     []string `json:"avoid_exercises,omitempty"`
	PreferredExercises []string `json:"preferred_exercises,omitempty"`
	InjuryHistory      []string `json:"injury_history,omitempty"`
	TimeConstraints    string   `json:"time_constraints,omitempty"`
}

// ExerciseSearchRequest represents a request to search exercises
type ExerciseSearchRequest struct {
	Query          string   `json:"query,omitempty"`
	Category       string   `json:"category,omitempty"`
	MuscleGroups   []string `json:"muscle_groups,omitempty"`
	Equipment      []string `json:"equipment,omitempty"`
	Difficulty     string   `json:"difficulty,omitempty"`
	IsGymExercise  *bool    `json:"is_gym_exercise,omitempty"`
	IsHomeExercise *bool    `json:"is_home_exercise,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Offset         int      `json:"offset,omitempty"`
}

// ExerciseSearchResponse represents the response for exercise search
type ExerciseSearchResponse struct {
	Exercises []*Exercise `json:"exercises"`
	Total     int         `json:"total"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
}

// ProgramSearchRequest represents a request to search workout programs
type ProgramSearchRequest struct {
	Query           string   `json:"query,omitempty"`
	Level           string   `json:"level,omitempty"`
	Goals           []string `json:"goals,omitempty"`
	SplitType       string   `json:"split_type,omitempty"`
	DurationWeeks   int      `json:"duration_weeks,omitempty"`
	SessionsPerWeek int      `json:"sessions_per_week,omitempty"`
	Equipment       []string `json:"equipment,omitempty"`
	Limit           int      `json:"limit,omitempty"`
	Offset          int      `json:"offset,omitempty"`
}

// ProgramSearchResponse represents the response for program search
type ProgramSearchResponse struct {
	Programs []*WorkoutProgram `json:"programs"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
