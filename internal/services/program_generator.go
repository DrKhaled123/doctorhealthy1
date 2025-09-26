package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api-key-generator/internal/models"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ProgramGenerator handles workout program generation and classification
type ProgramGenerator struct {
	progressionEngine *ProgressionEngine
}

// NewProgramGenerator creates a new program generator
func NewProgramGenerator() *ProgramGenerator {
	return &ProgramGenerator{
		progressionEngine: NewProgressionEngine(),
	}
}

// ExerciseClassification represents exercise classification data
type ExerciseClassification struct {
	Purpose      string   `json:"purpose"`       // strength, cardio, rehabilitation, flexibility
	MuscleGroups []string `json:"muscle_groups"` // primary muscle groups
	Equipment    []string `json:"equipment"`     // required equipment
	Difficulty   string   `json:"difficulty"`    // beginner, intermediate, advanced
	ExerciseType string   `json:"exercise_type"` // compound, isolation, cardio, plyometric
	Alternatives []string `json:"alternatives"`  // alternative exercise IDs
	Progressions []string `json:"progressions"`  // progression exercise IDs
}

// ProgramGenerationRequest represents a request to generate a workout program
type ProgramGenerationRequest struct {
	Goals              []string            `json:"goals" validate:"required"`         // weight_loss, muscle_gain, strength, etc.
	FitnessLevel       string              `json:"fitness_level" validate:"required"` // beginner, intermediate, advanced
	AvailableEquipment []string            `json:"available_equipment"`               // gym, home, dumbbells, etc.
	SessionsPerWeek    int                 `json:"sessions_per_week" validate:"min=1,max=7"`
	SessionDuration    int                 `json:"session_duration" validate:"min=15,max=180"` // minutes
	DurationWeeks      int                 `json:"duration_weeks" validate:"min=1,max=52"`
	Preferences        *ProgramPreferences `json:"preferences,omitempty"`
}

// ProgramPreferences represents user preferences for program generation
type ProgramPreferences struct {
	SplitType          string   `json:"split_type,omitempty"`          // full_body, upper_lower, push_pull_legs
	FocusAreas         []string `json:"focus_areas,omitempty"`         // chest, legs, back, etc.
	AvoidExercises     []string `json:"avoid_exercises,omitempty"`     // exercise IDs to avoid
	PreferredExercises []string `json:"preferred_exercises,omitempty"` // exercise IDs to include
	InjuryHistory      []string `json:"injury_history,omitempty"`      // knee, shoulder, back, etc.
	TimeConstraints    string   `json:"time_constraints,omitempty"`    // morning, evening, lunch_break
}

// ClassifyExercise classifies an exercise based on its characteristics
func (pg *ProgramGenerator) ClassifyExercise(ctx context.Context, exercise *models.WorkoutExerciseData) (*ExerciseClassification, error) {
	if exercise == nil {
		return nil, fmt.Errorf("exercise is required")
	}

	classification := &ExerciseClassification{
		MuscleGroups: exercise.PrimaryMuscles,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
	}

	// Classify purpose based on category and muscle groups
	classification.Purpose = pg.classifyPurpose(exercise)

	// Classify exercise type
	classification.ExerciseType = pg.classifyExerciseType(exercise)

	// Set alternatives and progressions
	classification.Alternatives = exercise.Alternatives
	classification.Progressions = exercise.Progressions

	return classification, nil
}

// GenerateProgram creates a complete workout program based on user requirements
func (pg *ProgramGenerator) GenerateProgram(ctx context.Context, req *ProgramGenerationRequest) (*models.WorkoutProgram, error) {
	if req == nil {
		return nil, fmt.Errorf("program generation request is required")
	}

	// Create base program
	program := &models.WorkoutProgram{
		ID:                uuid.New().String(),
		NameEn:            pg.generateProgramName(req),
		Level:             req.FitnessLevel,
		SplitType:         pg.determineSplitType(req),
		DescriptionEn:     pg.generateProgramDescription(req),
		DurationWeeks:     req.DurationWeeks,
		SessionsPerWeek:   req.SessionsPerWeek,
		Goals:             req.Goals,
		EquipmentRequired: req.AvailableEquipment,
		TargetAudience:    pg.determineTargetAudience(req),
		WeeklyPlans:       make(map[int]*models.WeeklyPlan),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Generate weekly plans
	for week := 1; week <= req.DurationWeeks; week++ {
		weeklyPlan, err := pg.generateWeeklyPlan(ctx, program, week, req)
		if err != nil {
			return nil, fmt.Errorf("failed to generate week %d: %w", week, err)
		}
		program.WeeklyPlans[week] = weeklyPlan
	}

	return program, nil
}

// GenerateCustomProgram creates a program with specific exercise selections
func (pg *ProgramGenerator) GenerateCustomProgram(ctx context.Context, req *ProgramGenerationRequest, exercises []*models.WorkoutExerciseData) (*models.WorkoutProgram, error) {
	// Filter exercises based on equipment and preferences
	filteredExercises := pg.filterExercisesByEquipment(exercises, req.AvailableEquipment)
	filteredExercises = pg.filterExercisesByPreferences(filteredExercises, req.Preferences)

	// Generate base program
	program, err := pg.GenerateProgram(ctx, req)
	if err != nil {
		return nil, err
	}

	// Customize with specific exercises
	err = pg.customizeWithExercises(program, filteredExercises, req)
	if err != nil {
		return nil, fmt.Errorf("failed to customize program: %w", err)
	}

	return program, nil
}

// Helper methods

func (pg *ProgramGenerator) classifyPurpose(exercise *models.WorkoutExerciseData) string {
	category := strings.ToLower(exercise.Category)

	switch {
	case strings.Contains(category, "cardio"):
		return "cardio"
	case strings.Contains(category, "strength"):
		return "strength"
	case strings.Contains(category, "flexibility") || strings.Contains(category, "stretch"):
		return "flexibility"
	case strings.Contains(category, "rehabilitation") || strings.Contains(category, "therapy"):
		return "rehabilitation"
	case strings.Contains(category, "core"):
		return "core_stability"
	default:
		return "strength" // default
	}
}

func (pg *ProgramGenerator) classifyExerciseType(exercise *models.WorkoutExerciseData) string {
	name := strings.ToLower(exercise.NameEn)

	// Check for compound movements
	compoundKeywords := []string{"squat", "deadlift", "bench", "press", "row", "pull-up", "chin-up", "dip", "lunge"}
	for _, keyword := range compoundKeywords {
		if strings.Contains(name, keyword) {
			return "compound"
		}
	}

	// Check for cardio (case-insensitive)
	if strings.Contains(strings.ToLower(exercise.Category), "cardio") {
		return "cardio"
	}

	// Check for plyometric
	plyoKeywords := []string{"jump", "hop", "bound", "explosive", "plyometric"}
	for _, keyword := range plyoKeywords {
		if strings.Contains(name, keyword) {
			return "plyometric"
		}
	}

	// Check if targets multiple muscle groups
	if len(exercise.PrimaryMuscles) >= 2 {
		return "compound"
	}

	return "isolation"
}

func (pg *ProgramGenerator) generateProgramName(req *ProgramGenerationRequest) string {
	goalStr := strings.Join(req.Goals, " & ")
	levelStr := cases.Title(language.Und).String(req.FitnessLevel)

	return fmt.Sprintf("%s %s Program - %d weeks", levelStr, goalStr, req.DurationWeeks)
}

func (pg *ProgramGenerator) determineSplitType(req *ProgramGenerationRequest) string {
	if req.Preferences != nil && req.Preferences.SplitType != "" {
		return req.Preferences.SplitType
	}

	// Determine split based on sessions per week
	switch req.SessionsPerWeek {
	case 1, 2:
		return "full_body"
	case 3:
		return "full_body"
	case 4:
		return "upper_lower"
	case 5, 6:
		return "push_pull_legs"
	case 7:
		return "body_part_split"
	default:
		return "full_body"
	}
}

func (pg *ProgramGenerator) generateProgramDescription(req *ProgramGenerationRequest) string {
	goals := strings.Join(req.Goals, ", ")
	equipment := "any equipment"
	if len(req.AvailableEquipment) > 0 {
		equipment = strings.Join(req.AvailableEquipment, ", ")
	}

	return fmt.Sprintf("A %d-week %s program designed for %s. Focuses on %s using %s. %d sessions per week, %d minutes per session.",
		req.DurationWeeks,
		req.FitnessLevel,
		goals,
		goals,
		equipment,
		req.SessionsPerWeek,
		req.SessionDuration,
	)
}

func (pg *ProgramGenerator) determineTargetAudience(req *ProgramGenerationRequest) []string {
	audience := []string{req.FitnessLevel}

	if req.SessionDuration <= 30 {
		audience = append(audience, "time_crunched")
	}

	hasGymEquipment := false
	for _, eq := range req.AvailableEquipment {
		if strings.Contains(strings.ToLower(eq), "gym") {
			hasGymEquipment = true
			break
		}
	}

	if !hasGymEquipment {
		audience = append(audience, "home_based")
	}

	return audience
}

func (pg *ProgramGenerator) generateWeeklyPlan(ctx context.Context, program *models.WorkoutProgram, weekNumber int, req *ProgramGenerationRequest) (*models.WeeklyPlan, error) {
	weeklyPlan := &models.WeeklyPlan{
		ID:                    uuid.New().String(),
		ProgramID:             program.ID,
		WeekNumber:            weekNumber,
		IntensityLevel:        pg.calculateWeeklyIntensity(weekNumber, req.FitnessLevel),
		VolumeIncreasePercent: pg.calculateVolumeIncrease(weekNumber, req.FitnessLevel),
		DailyWorkouts:         make(map[int]*models.DailyWorkout),
		CreatedAt:             time.Now(),
	}

	// Generate progression notes
	weeklyPlan.ProgressionNotesEn = pg.generateProgressionNotes(weekNumber, req.FitnessLevel)

	// Generate daily workouts
	for day := 1; day <= 7; day++ {
		if day <= req.SessionsPerWeek {
			dailyWorkout := pg.generateDailyWorkout(weeklyPlan.ID, day, program.SplitType, req)
			weeklyPlan.DailyWorkouts[day] = dailyWorkout
		} else {
			// Rest day
			restDay := &models.DailyWorkout{
				ID:              uuid.New().String(),
				WeeklyPlanID:    weeklyPlan.ID,
				DayNumber:       day,
				FocusEn:         "Rest Day",
				DurationMinutes: 0,
				IsRestDay:       true,
				CreatedAt:       time.Now(),
			}
			weeklyPlan.DailyWorkouts[day] = restDay
		}
	}

	return weeklyPlan, nil
}

func (pg *ProgramGenerator) generateDailyWorkout(weeklyPlanID string, dayNumber int, splitType string, req *ProgramGenerationRequest) *models.DailyWorkout {
	focus := pg.determineDayFocus(dayNumber, splitType)

	return &models.DailyWorkout{
		ID:              uuid.New().String(),
		WeeklyPlanID:    weeklyPlanID,
		DayNumber:       dayNumber,
		FocusEn:         focus,
		DurationMinutes: req.SessionDuration,
		IsRestDay:       false,
		CreatedAt:       time.Now(),
	}
}

func (pg *ProgramGenerator) determineDayFocus(dayNumber int, splitType string) string {
	switch splitType {
	case "full_body":
		return fmt.Sprintf("Full Body Workout %d", dayNumber)
	case "upper_lower":
		if dayNumber%2 == 1 {
			return "Upper Body"
		}
		return "Lower Body"
	case "push_pull_legs":
		switch dayNumber % 3 {
		case 1:
			return "Push (Chest, Shoulders, Triceps)"
		case 2:
			return "Pull (Back, Biceps)"
		case 0:
			return "Legs (Quads, Hamstrings, Glutes, Calves)"
		}
	case "body_part_split":
		bodyParts := []string{"Chest", "Back", "Shoulders", "Arms", "Legs", "Core", "Cardio"}
		if dayNumber <= len(bodyParts) {
			return bodyParts[dayNumber-1]
		}
	}
	return "Full Body"
}

func (pg *ProgramGenerator) calculateWeeklyIntensity(weekNumber int, fitnessLevel string) string {
	progressionFactor := pg.progressionEngine.calculateProgressionFactor(weekNumber, fitnessLevel)

	switch {
	case progressionFactor >= 1.3:
		return "high"
	case progressionFactor >= 1.15:
		return "moderate-high"
	case progressionFactor >= 1.05:
		return "moderate"
	default:
		return "low-moderate"
	}
}

func (pg *ProgramGenerator) calculateVolumeIncrease(weekNumber int, fitnessLevel string) float64 {
	if weekNumber <= 1 {
		return 0.0
	}

	var weeklyIncrease float64
	switch strings.ToLower(fitnessLevel) {
	case "beginner":
		weeklyIncrease = 5.0 // 5% per week
	case "intermediate":
		weeklyIncrease = 3.0 // 3% per week
	case "advanced":
		weeklyIncrease = 2.0 // 2% per week
	default:
		weeklyIncrease = 3.0
	}

	totalIncrease := float64(weekNumber-1) * weeklyIncrease

	// Cap at reasonable maximum
	if totalIncrease > 50.0 {
		totalIncrease = 50.0
	}

	return totalIncrease
}

func (pg *ProgramGenerator) generateProgressionNotes(weekNumber int, fitnessLevel string) string {
	switch weekNumber {
	case 1:
		return "Focus on learning proper form and establishing baseline performance"
	case 2:
		return "Begin gradual intensity increases while maintaining good form"
	case 3, 4:
		return "Continue progressive overload with small weight/rep increases"
	default:
		if weekNumber%4 == 0 {
			return "Deload week - reduce intensity by 10-15% for recovery"
		}
		return fmt.Sprintf("Week %d: Maintain progressive overload with %s-appropriate increases", weekNumber, fitnessLevel)
	}
}

func (pg *ProgramGenerator) filterExercisesByEquipment(exercises []*models.WorkoutExerciseData, availableEquipment []string) []*models.WorkoutExerciseData {
	if len(availableEquipment) == 0 {
		return exercises // No filter if no equipment specified
	}

	equipmentMap := make(map[string]bool)
	for _, eq := range availableEquipment {
		equipmentMap[strings.ToLower(eq)] = true
	}

	var filtered []*models.WorkoutExerciseData
	for _, exercise := range exercises {
		canPerform := true
		for _, reqEquipment := range exercise.Equipment {
			if !equipmentMap[strings.ToLower(reqEquipment)] && reqEquipment != "bodyweight" {
				canPerform = false
				break
			}
		}
		if canPerform {
			filtered = append(filtered, exercise)
		}
	}

	return filtered
}

func (pg *ProgramGenerator) filterExercisesByPreferences(exercises []*models.WorkoutExerciseData, preferences *ProgramPreferences) []*models.WorkoutExerciseData {
	if preferences == nil {
		return exercises
	}

	// Create avoid map
	avoidMap := make(map[string]bool)
	for _, avoid := range preferences.AvoidExercises {
		avoidMap[avoid] = true
	}

	var filtered []*models.WorkoutExerciseData
	for _, exercise := range exercises {
		if !avoidMap[exercise.ID] {
			filtered = append(filtered, exercise)
		}
	}

	return filtered
}

func (pg *ProgramGenerator) customizeWithExercises(program *models.WorkoutProgram, exercises []*models.WorkoutExerciseData, req *ProgramGenerationRequest) error {
	// This would implement the logic to add specific exercises to the program
	// For now, this is a placeholder that would be expanded based on specific requirements

	// Example: Add preferred exercises to each workout
	if req.Preferences != nil && len(req.Preferences.PreferredExercises) > 0 {
		// Placeholder: ensure preferred exercises are considered in future
		_ = req.Preferences
	}

	return nil
}
