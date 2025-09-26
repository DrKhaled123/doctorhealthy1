package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-key-generator/internal/models"
)

// SetsRepsManager handles dynamic sets, reps, and rest period management
type SetsRepsManager struct {
	progressionEngine *ProgressionEngine
	// In a real implementation, you'd inject database repositories here
}

// NewSetsRepsManager creates a new sets/reps manager
func NewSetsRepsManager() *SetsRepsManager {
	return &SetsRepsManager{
		progressionEngine: NewProgressionEngine(),
	}
}

// GenerateExerciseRecommendations creates dynamic sets/reps recommendations for exercises
func (srm *SetsRepsManager) GenerateExerciseRecommendations(ctx context.Context, req *ExerciseRecommendationRequest) (*ExerciseRecommendationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	recommendations := make([]*ExerciseRecommendation, 0, len(req.Exercises))

	for _, exercise := range req.Exercises {
		// Determine exercise type for rest period calculation
		exerciseType := srm.determineExerciseType(exercise)

		// Generate sets/reps recommendation
		setsRepsRec, err := srm.progressionEngine.GenerateSetsRepsRecommendation(
			ctx,
			req.Goals,
			exerciseType,
			req.FitnessLevel,
			req.WeekNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recommendation for exercise %s: %w", exercise.ID, err)
		}

		// Calculate rest period
		restPeriod := srm.progressionEngine.CalculateRestPeriod(
			exerciseType,
			setsRepsRec.Intensity,
			req.Goals,
		)

		// Create exercise recommendation
		recommendation := &ExerciseRecommendation{
			ExerciseID:       exercise.ID,
			ExerciseName:     exercise.NameEn,
			Sets:             setsRepsRec.Sets,
			RepsMin:          setsRepsRec.RepsMin,
			RepsMax:          setsRepsRec.RepsMax,
			RepsDisplay:      srm.formatRepsDisplay(setsRepsRec.RepsMin, setsRepsRec.RepsMax),
			RestSeconds:      restPeriod,
			RestDisplay:      srm.formatRestDisplay(restPeriod),
			IntensityPercent: setsRepsRec.Intensity,
			ExerciseType:     exerciseType,
			Notes:            setsRepsRec.Notes,
		}

		// Add progression tracking if user has previous data
		if req.UserProgress != nil {
			if progress, exists := req.UserProgress[exercise.ID]; exists {
				recommendation.PreviousWeight = progress.Weight
				recommendation.PreviousReps = progress.Reps
				recommendation.LastPerformed = &progress.LastPerformed

				// Suggest progression adjustments
				adjusted, err := srm.progressionEngine.SuggestProgressionAdjustment(ctx, progress, setsRepsRec)
				if err == nil && adjusted != nil {
					recommendation.IntensityPercent = adjusted.Intensity
					recommendation.ProgressionNotes = adjusted.Notes
				}
			}
		}

		recommendations = append(recommendations, recommendation)
	}

	return &ExerciseRecommendationResponse{
		Recommendations: recommendations,
		WeekNumber:      req.WeekNumber,
		GeneratedAt:     time.Now(),
	}, nil
}

// TrackWorkoutSession records a completed workout session with performance data
func (srm *SetsRepsManager) TrackWorkoutSession(ctx context.Context, req *WorkoutSessionRequest) (*WorkoutSessionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("session request is required")
	}

	sessionLog := &models.SessionLog{
		ID:              generateID(),
		UserID:          req.UserID,
		DailyWorkoutID:  req.DailyWorkoutID,
		SessionDate:     time.Now(),
		DurationMinutes: req.DurationMinutes,
		RPEScore:        req.OverallRPE,
		Notes:           req.Notes,
		Completed:       req.Completed,
		ExerciseLogs:    make([]*models.ExerciseLog, 0, len(req.ExercisePerformance)),
	}

	exerciseProgress := make(map[string]*models.ExerciseProgress)

	// Process each exercise performance
	for _, performance := range req.ExercisePerformance {
		// Create exercise logs for each set
		for setNum, setData := range performance.Sets {
			exerciseLog := &models.ExerciseLog{
				ID:            generateID(),
				SessionLogID:  sessionLog.ID,
				ExerciseID:    performance.ExerciseID,
				SetNumber:     setNum + 1,
				RepsCompleted: setData.RepsCompleted,
				WeightKg:      setData.WeightKg,
				RestSeconds:   setData.RestSeconds,
				RPEScore:      setData.RPE,
				Notes:         setData.Notes,
				CreatedAt:     time.Now(),
			}
			sessionLog.ExerciseLogs = append(sessionLog.ExerciseLogs, exerciseLog)
		}

		// Track overall exercise progress (use best set for progress tracking)
		if len(performance.Sets) > 0 {
			bestSet := srm.findBestSet(performance.Sets)
			progress, err := srm.progressionEngine.TrackExerciseProgress(ctx, req.UserID, performance.ExerciseID, &models.ExerciseLog{
				ExerciseID:    performance.ExerciseID,
				RepsCompleted: bestSet.RepsCompleted,
				WeightKg:      bestSet.WeightKg,
				RPEScore:      bestSet.RPE,
				Notes:         bestSet.Notes,
			})
			if err == nil {
				exerciseProgress[performance.ExerciseID] = progress
			}
		}
	}

	// Generate progress analytics
	analytics, err := srm.progressionEngine.GenerateProgressAnalytics(ctx, req.UserID, req.ProgramID, exerciseProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate progress analytics: %w", err)
	}

	return &WorkoutSessionResponse{
		SessionLog:          sessionLog,
		ExerciseProgress:    exerciseProgress,
		Analytics:           analytics,
		NextRecommendations: srm.generateNextSessionRecommendations(exerciseProgress),
	}, nil
}

// GenerateProgressReport creates comprehensive progress analytics
func (srm *SetsRepsManager) GenerateProgressReport(ctx context.Context, userID string, programID string, timeRange TimeRange) (*ProgressReport, error) {
	// In a real implementation, this would query the database for historical data
	// For now, we'll create a sample report structure

	report := &ProgressReport{
		UserID:      userID,
		ProgramID:   programID,
		TimeRange:   timeRange,
		GeneratedAt: time.Now(),
	}

	// Calculate strength improvements (placeholder data)
	report.StrengthImprovements = map[string]*StrengthImprovement{
		"bench_press": {
			ExerciseName:       "Bench Press",
			StartingWeight:     60.0,
			CurrentWeight:      75.0,
			ImprovementKg:      15.0,
			ImprovementPercent: 25.0,
			PersonalRecord: &models.PersonalRecord{
				Weight:    75.0,
				Reps:      5,
				OneRepMax: 84.4,
				Date:      time.Now().AddDate(0, 0, -2),
			},
		},
		"squat": {
			ExerciseName:       "Squat",
			StartingWeight:     80.0,
			CurrentWeight:      100.0,
			ImprovementKg:      20.0,
			ImprovementPercent: 25.0,
			PersonalRecord: &models.PersonalRecord{
				Weight:    100.0,
				Reps:      8,
				OneRepMax: 126.7,
				Date:      time.Now().AddDate(0, 0, -1),
			},
		},
	}

	// Calculate volume progression
	report.VolumeProgression = &VolumeProgression{
		StartingVolume:     5000.0,
		CurrentVolume:      7500.0,
		ImprovementPercent: 50.0,
		WeeklyTrend:        []float64{5000, 5250, 5500, 6000, 6500, 7000, 7500},
	}

	// Performance metrics
	report.PerformanceMetrics = &PerformanceMetrics{
		AverageRPE:             7.2,
		ConsistencyScore:       85.0,
		WorkoutFrequency:       4.2,
		AverageSessionDuration: 65,
		TotalWorkouts:          24,
		CompletionRate:         92.0,
	}

	return report, nil
}

// Helper methods

func (srm *SetsRepsManager) determineExerciseType(exercise *models.WorkoutExerciseData) string {
	// Analyze exercise characteristics to determine type
	category := strings.ToLower(exercise.Category)

	// Check for compound movements
	if srm.isCompoundExercise(exercise) {
		return "compound"
	}

	// Check for cardio exercises
	if strings.Contains(category, "cardio") || strings.Contains(category, "conditioning") {
		return "cardio"
	}

	// Check for core exercises
	if strings.Contains(category, "core") || strings.Contains(category, "abs") {
		return "core"
	}

	// Check for plyometric exercises
	if strings.Contains(category, "plyometric") || strings.Contains(category, "explosive") {
		return "plyometric"
	}

	// Default to isolation
	return "isolation"
}

func (srm *SetsRepsManager) isCompoundExercise(exercise *models.WorkoutExerciseData) bool {
	compoundKeywords := []string{
		"squat", "deadlift", "bench", "press", "row", "pull-up", "chin-up",
		"dip", "lunge", "clean", "snatch", "thruster",
	}

	exerciseName := strings.ToLower(exercise.NameEn)
	for _, keyword := range compoundKeywords {
		if strings.Contains(exerciseName, keyword) {
			return true
		}
	}

	// Also check if it targets multiple muscle groups
	return len(exercise.PrimaryMuscles) >= 2
}

func (srm *SetsRepsManager) formatRepsDisplay(min, max int) string {
	if min == max {
		return strconv.Itoa(min)
	}
	return fmt.Sprintf("%d-%d", min, max)
}

func (srm *SetsRepsManager) formatRestDisplay(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60

	if remainingSeconds == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

func (srm *SetsRepsManager) findBestSet(sets []*SetPerformance) *SetPerformance {
	if len(sets) == 0 {
		return nil
	}

	bestSet := sets[0]
	bestScore := srm.calculateSetScore(bestSet)

	for _, set := range sets[1:] {
		score := srm.calculateSetScore(set)
		if score > bestScore {
			bestSet = set
			bestScore = score
		}
	}

	return bestSet
}

func (srm *SetsRepsManager) calculateSetScore(set *SetPerformance) float64 {
	// Simple scoring: weight * reps (volume)
	weight := 0.0
	if set.WeightKg != nil {
		weight = *set.WeightKg
	}
	return weight * float64(set.RepsCompleted)
}

func (srm *SetsRepsManager) generateNextSessionRecommendations(progress map[string]*models.ExerciseProgress) []string {
	recommendations := []string{}

	for exerciseID, prog := range progress {
		if prog.RPE <= 6 {
			recommendations = append(recommendations, fmt.Sprintf("Increase weight for %s - previous session was too easy", exerciseID))
		} else if prog.RPE >= 9 {
			recommendations = append(recommendations, fmt.Sprintf("Reduce weight for %s - previous session was too challenging", exerciseID))
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Great work! Continue with current progression plan")
	}

	return recommendations
}

// generateID creates a simple ID (in production, use proper UUID)
func generateID() string {
	return fmt.Sprintf("id_%d", time.Now().UnixNano())
}

// Request/Response types for the service

type ExerciseRecommendationRequest struct {
	Exercises    []*models.WorkoutExerciseData       `json:"exercises"`
	Goals        []string                            `json:"goals"`
	FitnessLevel string                              `json:"fitness_level"`
	WeekNumber   int                                 `json:"week_number"`
	UserProgress map[string]*models.ExerciseProgress `json:"user_progress,omitempty"`
}

type ExerciseRecommendationResponse struct {
	Recommendations []*ExerciseRecommendation `json:"recommendations"`
	WeekNumber      int                       `json:"week_number"`
	GeneratedAt     time.Time                 `json:"generated_at"`
}

type ExerciseRecommendation struct {
	ExerciseID       string     `json:"exercise_id"`
	ExerciseName     string     `json:"exercise_name"`
	Sets             int        `json:"sets"`
	RepsMin          int        `json:"reps_min"`
	RepsMax          int        `json:"reps_max"`
	RepsDisplay      string     `json:"reps_display"`
	RestSeconds      int        `json:"rest_seconds"`
	RestDisplay      string     `json:"rest_display"`
	IntensityPercent float64    `json:"intensity_percent"`
	ExerciseType     string     `json:"exercise_type"`
	Notes            string     `json:"notes"`
	PreviousWeight   float64    `json:"previous_weight,omitempty"`
	PreviousReps     int        `json:"previous_reps,omitempty"`
	LastPerformed    *time.Time `json:"last_performed,omitempty"`
	ProgressionNotes string     `json:"progression_notes,omitempty"`
}

type WorkoutSessionRequest struct {
	UserID              string                 `json:"user_id"`
	ProgramID           string                 `json:"program_id"`
	DailyWorkoutID      string                 `json:"daily_workout_id"`
	DurationMinutes     int                    `json:"duration_minutes"`
	OverallRPE          int                    `json:"overall_rpe"`
	Notes               string                 `json:"notes"`
	Completed           bool                   `json:"completed"`
	ExercisePerformance []*ExercisePerformance `json:"exercise_performance"`
}

type ExercisePerformance struct {
	ExerciseID string            `json:"exercise_id"`
	Sets       []*SetPerformance `json:"sets"`
}

type SetPerformance struct {
	RepsCompleted int      `json:"reps_completed"`
	WeightKg      *float64 `json:"weight_kg,omitempty"`
	RestSeconds   int      `json:"rest_seconds"`
	RPE           int      `json:"rpe"`
	Notes         string   `json:"notes,omitempty"`
}

type WorkoutSessionResponse struct {
	SessionLog          *models.SessionLog                  `json:"session_log"`
	ExerciseProgress    map[string]*models.ExerciseProgress `json:"exercise_progress"`
	Analytics           *models.ProgressMetrics             `json:"analytics"`
	NextRecommendations []string                            `json:"next_recommendations"`
}

type TimeRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type ProgressReport struct {
	UserID               string                          `json:"user_id"`
	ProgramID            string                          `json:"program_id"`
	TimeRange            TimeRange                       `json:"time_range"`
	GeneratedAt          time.Time                       `json:"generated_at"`
	StrengthImprovements map[string]*StrengthImprovement `json:"strength_improvements"`
	VolumeProgression    *VolumeProgression              `json:"volume_progression"`
	PerformanceMetrics   *PerformanceMetrics             `json:"performance_metrics"`
}

type StrengthImprovement struct {
	ExerciseName       string                 `json:"exercise_name"`
	StartingWeight     float64                `json:"starting_weight"`
	CurrentWeight      float64                `json:"current_weight"`
	ImprovementKg      float64                `json:"improvement_kg"`
	ImprovementPercent float64                `json:"improvement_percent"`
	PersonalRecord     *models.PersonalRecord `json:"personal_record"`
}

type VolumeProgression struct {
	StartingVolume     float64   `json:"starting_volume"`
	CurrentVolume      float64   `json:"current_volume"`
	ImprovementPercent float64   `json:"improvement_percent"`
	WeeklyTrend        []float64 `json:"weekly_trend"`
}

type PerformanceMetrics struct {
	AverageRPE             float64 `json:"average_rpe"`
	ConsistencyScore       float64 `json:"consistency_score"`
	WorkoutFrequency       float64 `json:"workout_frequency"`
	AverageSessionDuration int     `json:"average_session_duration"`
	TotalWorkouts          int     `json:"total_workouts"`
	CompletionRate         float64 `json:"completion_rate"`
}
