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

// WorkoutService handles workout plan generation
type WorkoutService struct {
	db          *sql.DB
	userService *UserService
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(db *sql.DB, userService *UserService) *WorkoutService {
	return &WorkoutService{
		db:          db,
		userService: userService,
	}
}

// GenerateWorkoutPlan generates a personalized workout plan
func (s *WorkoutService) GenerateWorkoutPlan(ctx context.Context, req *models.GenerateWorkoutPlanRequest) (*models.WorkoutPlan, error) {
	// Get user data
	user, err := s.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate exercises based on goal and constraints
	exercises, err := s.generateExercises(ctx, user, req.Goal, req.WorkoutType, req.Injuries, req.Complaints)
	if err != nil {
		return nil, fmt.Errorf("failed to generate exercises: %w", err)
	}

	// Create workout plan
	plan := &models.WorkoutPlan{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Goal:        req.Goal,
		WorkoutType: req.WorkoutType,
		Exercises:   exercises,
		Injuries:    req.Injuries,
		Complaints:  req.Complaints,
		CreatedAt:   time.Now().UTC(),
	}

	// Save to database (optional - for tracking)
	err = s.savePlan(ctx, plan)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save workout plan: %v\n", err)
	}

	return plan, nil
}

// generateExercises generates exercises based on user requirements
func (s *WorkoutService) generateExercises(ctx context.Context, user *models.User, goal, workoutType string, injuries, complaints []string) ([]models.Exercise, error) {
	// Get exercise templates
	templates := s.getExerciseTemplates(workoutType)

	// Filter exercises based on goal
	goalFiltered := s.filterExercisesByGoal(templates, goal)

	// Filter out exercises that conflict with injuries
	injuryFiltered := s.filterExercisesByInjuries(goalFiltered, injuries)

	// Select diverse exercises
	selectedExercises := s.selectDiverseExercises(injuryFiltered, goal)

	// Add alternatives for each exercise
	for i := range selectedExercises {
		alternative := s.findAlternativeExercise(selectedExercises[i], templates)
		if alternative != nil {
			selectedExercises[i].Alternative = alternative
		}
	}

	return selectedExercises, nil
}

// getExerciseTemplates returns exercise templates based on workout type
func (s *WorkoutService) getExerciseTemplates(workoutType string) []models.Exercise {
	if workoutType == "home" {
		return s.getHomeExercises()
	}
	return s.getGymExercises()
}

// getGymExercises returns gym-based exercises
func (s *WorkoutService) getGymExercises() []models.Exercise {
	return []models.Exercise{
		{
			ID:   "gym_squat_001",
			Name: "Barbell Squat",
			Type: "compound",
			Sets: 4,
			Reps: "8-12",
			Rest: "2-3 minutes",
			Instructions: []string{
				"Stand with feet shoulder-width apart",
				"Place barbell on upper back",
				"Lower body by bending knees and hips",
				"Keep chest up and knees aligned with toes",
				"Return to starting position",
			},
			CommonMistakes: []string{
				"Knees caving inward",
				"Not going deep enough",
				"Leaning too far forward",
				"Not keeping core tight",
			},
			TargetMuscles: []string{"quadriceps", "glutes", "hamstrings", "core"},
		},
		{
			ID:   "gym_deadlift_001",
			Name: "Deadlift",
			Type: "compound",
			Sets: 3,
			Reps: "5-8",
			Rest: "3-4 minutes",
			Instructions: []string{
				"Stand with feet hip-width apart",
				"Grip barbell with hands just outside legs",
				"Keep back straight and chest up",
				"Lift by extending hips and knees",
				"Lower bar with control",
			},
			CommonMistakes: []string{
				"Rounding the back",
				"Bar drifting away from body",
				"Not engaging lats",
				"Hyperextending at the top",
			},
			TargetMuscles: []string{"hamstrings", "glutes", "erector spinae", "traps"},
		},
		{
			ID:   "gym_bench_001",
			Name: "Bench Press",
			Type: "compound",
			Sets: 4,
			Reps: "8-12",
			Rest: "2-3 minutes",
			Instructions: []string{
				"Lie on bench with feet flat on floor",
				"Grip bar slightly wider than shoulders",
				"Lower bar to chest with control",
				"Press bar up explosively",
				"Keep shoulder blades retracted",
			},
			CommonMistakes: []string{
				"Bouncing bar off chest",
				"Flaring elbows too wide",
				"Not retracting shoulder blades",
				"Lifting feet off ground",
			},
			TargetMuscles: []string{"chest", "shoulders", "triceps"},
		},
		{
			ID:   "gym_pullup_001",
			Name: "Pull-ups",
			Type: "compound",
			Sets: 3,
			Reps: "6-12",
			Rest: "2-3 minutes",
			Instructions: []string{
				"Hang from bar with palms facing away",
				"Pull body up until chin clears bar",
				"Lower with control",
				"Keep core engaged throughout",
			},
			CommonMistakes: []string{
				"Using momentum",
				"Not going full range of motion",
				"Swinging body",
				"Shrugging shoulders",
			},
			TargetMuscles: []string{"lats", "rhomboids", "biceps", "rear delts"},
		},
	}
}

// getHomeExercises returns home-based exercises
func (s *WorkoutService) getHomeExercises() []models.Exercise {
	return []models.Exercise{
		{
			ID:   "home_squat_001",
			Name: "Bodyweight Squat",
			Type: "compound",
			Sets: 3,
			Reps: "15-20",
			Rest: "1-2 minutes",
			Instructions: []string{
				"Stand with feet shoulder-width apart",
				"Lower body by bending knees and hips",
				"Keep chest up and weight on heels",
				"Return to starting position",
			},
			CommonMistakes: []string{
				"Knees caving inward",
				"Not going deep enough",
				"Leaning too far forward",
			},
			TargetMuscles: []string{"quadriceps", "glutes", "hamstrings"},
		},
		{
			ID:   "home_pushup_001",
			Name: "Push-ups",
			Type: "compound",
			Sets: 3,
			Reps: "10-15",
			Rest: "1-2 minutes",
			Instructions: []string{
				"Start in plank position",
				"Lower body until chest nearly touches floor",
				"Push back up to starting position",
				"Keep body in straight line",
			},
			CommonMistakes: []string{
				"Sagging hips",
				"Not going full range of motion",
				"Flaring elbows too wide",
			},
			TargetMuscles: []string{"chest", "shoulders", "triceps", "core"},
		},
		{
			ID:   "home_lunge_001",
			Name: "Lunges",
			Type: "compound",
			Sets: 3,
			Reps: "12-16 each leg",
			Rest: "1-2 minutes",
			Instructions: []string{
				"Step forward with one leg",
				"Lower hips until both knees are at 90 degrees",
				"Push back to starting position",
				"Alternate legs",
			},
			CommonMistakes: []string{
				"Knee extending past toes",
				"Not lowering enough",
				"Leaning forward",
			},
			TargetMuscles: []string{"quadriceps", "glutes", "hamstrings"},
		},
		{
			ID:   "home_plank_001",
			Name: "Plank",
			Type: "isometric",
			Sets: 3,
			Reps: "30-60 seconds",
			Rest: "1 minute",
			Instructions: []string{
				"Start in push-up position",
				"Hold body in straight line",
				"Keep core engaged",
				"Breathe normally",
			},
			CommonMistakes: []string{
				"Sagging hips",
				"Raising hips too high",
				"Holding breath",
			},
			TargetMuscles: []string{"core", "shoulders", "glutes"},
		},
	}
}

// filterExercisesByGoal filters exercises based on user's goal
func (s *WorkoutService) filterExercisesByGoal(exercises []models.Exercise, goal string) []models.Exercise {
	switch goal {
	case "build_muscle", "improve_strength":
		// Prefer compound movements
		var filtered []models.Exercise
		for _, ex := range exercises {
			if ex.Type == "compound" {
				filtered = append(filtered, ex)
			}
		}
		return filtered
	case "lose_weight":
		// Include all exercises but adjust reps/sets for higher volume
		for i := range exercises {
			exercises[i].Reps = "12-20"
			exercises[i].Rest = "1-2 minutes"
		}
		return exercises
	default:
		return exercises
	}
}

// filterExercisesByInjuries removes exercises that may aggravate injuries
func (s *WorkoutService) filterExercisesByInjuries(exercises []models.Exercise, injuries []string) []models.Exercise {
	if len(injuries) == 0 {
		return exercises
	}

	var filtered []models.Exercise

	for _, exercise := range exercises {
		skip := false

		for _, injury := range injuries {
			// Simple injury filtering logic
			switch strings.ToLower(injury) {
			case "knee":
				if strings.Contains(strings.ToLower(exercise.Name), "squat") ||
					strings.Contains(strings.ToLower(exercise.Name), "lunge") {
					skip = true
				}
			case "shoulder":
				if strings.Contains(strings.ToLower(exercise.Name), "press") ||
					strings.Contains(strings.ToLower(exercise.Name), "pull") {
					skip = true
				}
			case "back", "lower back":
				if strings.Contains(strings.ToLower(exercise.Name), "deadlift") {
					skip = true
				}
			}
		}

		if !skip {
			filtered = append(filtered, exercise)
		}
	}

	return filtered
}

// selectDiverseExercises selects a diverse set of exercises
func (s *WorkoutService) selectDiverseExercises(exercises []models.Exercise, goal string) []models.Exercise {
	if len(exercises) <= 6 {
		return exercises
	}

	// Select 6-8 exercises ensuring muscle group diversity
	var selected []models.Exercise
	usedMuscleGroups := make(map[string]bool)

	// First pass: select exercises targeting different muscle groups
	for _, exercise := range exercises {
		if len(selected) >= 6 {
			break
		}

		hasNewMuscleGroup := false
		for _, muscle := range exercise.TargetMuscles {
			if !usedMuscleGroups[muscle] {
				hasNewMuscleGroup = true
				break
			}
		}

		if hasNewMuscleGroup {
			selected = append(selected, exercise)
			for _, muscle := range exercise.TargetMuscles {
				usedMuscleGroups[muscle] = true
			}
		}
	}

	// Second pass: fill remaining slots
	for _, exercise := range exercises {
		if len(selected) >= 8 {
			break
		}

		// Check if already selected
		alreadySelected := false
		for _, sel := range selected {
			if sel.Name == exercise.Name {
				alreadySelected = true
				break
			}
		}

		if !alreadySelected {
			selected = append(selected, exercise)
		}
	}

	return selected
}

// findAlternativeExercise finds an alternative for the given exercise
func (s *WorkoutService) findAlternativeExercise(exercise models.Exercise, allExercises []models.Exercise) *models.Exercise {
	// Create muscle group map for faster lookup
	exerciseMuscles := make(map[string]bool)
	for _, muscle := range exercise.TargetMuscles {
		exerciseMuscles[muscle] = true
	}

	// Find exercise targeting similar muscle groups
	for _, alt := range allExercises {
		if alt.Name == exercise.Name {
			continue
		}

		// Check if targets similar muscles
		commonMuscles := 0
		for _, muscle := range alt.TargetMuscles {
			if exerciseMuscles[muscle] {
				commonMuscles++
				if commonMuscles >= 2 {
					altCopy := alt
					return &altCopy
				}
			}
		}
	}

	return nil
}

// savePlan saves the workout plan to database
func (s *WorkoutService) savePlan(ctx context.Context, plan *models.WorkoutPlan) error {
	query := `
		INSERT INTO workout_plans (
			id, user_id, goal, workout_type, created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		plan.ID,
		plan.UserID,
		plan.Goal,
		plan.WorkoutType,
		plan.CreatedAt,
	)

	return err
}

// GetAvailableGoals returns available workout goals
func (s *WorkoutService) GetAvailableGoals() []string {
	return []string{
		"build_muscle",
		"improve_strength",
		"lose_weight",
		"improve_endurance",
		"general_fitness",
		"body_recomposition",
	}
}

// GetAvailableInjuries returns common injuries to choose from
func (s *WorkoutService) GetAvailableInjuries() []string {
	return []string{
		"knee",
		"shoulder",
		"lower_back",
		"upper_back",
		"neck",
		"ankle",
		"wrist",
		"elbow",
		"hip",
	}
}

// GetAvailableComplaints returns common complaints
func (s *WorkoutService) GetAvailableComplaints() []string {
	return []string{
		"joint_pain",
		"muscle_weakness",
		"poor_posture",
		"low_energy",
		"stress",
		"sleep_issues",
		"mobility_issues",
	}
}
