package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api-key-generator/internal/models"
)

// ProgressionEngine handles sets, reps, rest period management and progression calculations
type ProgressionEngine struct {
	// Dependencies would be injected here
}

// NewProgressionEngine creates a new progression engine instance
func NewProgressionEngine() *ProgressionEngine {
	return &ProgressionEngine{}
}

// SetsRepsRecommendation represents dynamic sets/reps recommendations
type SetsRepsRecommendation struct {
	Sets        int     `json:"sets"`
	RepsMin     int     `json:"reps_min"`
	RepsMax     int     `json:"reps_max"`
	RestSeconds int     `json:"rest_seconds"`
	Intensity   float64 `json:"intensity_percent"` // Percentage of 1RM
	Notes       string  `json:"notes"`
}

// GoalBasedRecommendations defines sets/reps based on training goals
var GoalBasedRecommendations = map[string]SetsRepsRecommendation{
	"strength": {
		Sets:        3,
		RepsMin:     1,
		RepsMax:     5,
		RestSeconds: 180, // 3 minutes
		Intensity:   85.0,
		Notes:       "Focus on heavy weight, low reps for maximum strength gains",
	},
	"hypertrophy": {
		Sets:        3,
		RepsMin:     6,
		RepsMax:     12,
		RestSeconds: 90, // 1.5 minutes
		Intensity:   70.0,
		Notes:       "Moderate weight, moderate reps for muscle growth",
	},
	"endurance": {
		Sets:        2,
		RepsMin:     15,
		RepsMax:     25,
		RestSeconds: 60, // 1 minute
		Intensity:   50.0,
		Notes:       "Light weight, high reps for muscular endurance",
	},
	"power": {
		Sets:        4,
		RepsMin:     3,
		RepsMax:     6,
		RestSeconds: 120, // 2 minutes
		Intensity:   75.0,
		Notes:       "Explosive movements with adequate rest for power development",
	},
	"weight_loss": {
		Sets:        3,
		RepsMin:     12,
		RepsMax:     20,
		RestSeconds: 45, // 45 seconds
		Intensity:   60.0,
		Notes:       "Higher volume with shorter rest for calorie burn",
	},
}

// ExerciseTypeRestPeriods defines rest periods based on exercise type
var ExerciseTypeRestPeriods = map[string]int{
	"compound":   180, // 3 minutes for squats, deadlifts, bench press
	"isolation":  90,  // 1.5 minutes for bicep curls, tricep extensions
	"cardio":     30,  // 30 seconds for cardio exercises
	"core":       60,  // 1 minute for core exercises
	"plyometric": 120, // 2 minutes for explosive movements
	"stretching": 15,  // 15 seconds between stretches
}

// GenerateSetsRepsRecommendation creates dynamic recommendations based on goals and exercise type
func (pe *ProgressionEngine) GenerateSetsRepsRecommendation(ctx context.Context, goals []string, exerciseType string, fitnessLevel string, weekNumber int) (*SetsRepsRecommendation, error) {
	// Determine primary goal (first goal takes priority)
	primaryGoal := "hypertrophy" // default
	if len(goals) > 0 {
		primaryGoal = strings.ToLower(goals[0])
	}

	// Get base recommendation for the goal
	baseRec, exists := GoalBasedRecommendations[primaryGoal]
	if !exists {
		baseRec = GoalBasedRecommendations["hypertrophy"]
	}

	// Create a copy to modify
	recommendation := SetsRepsRecommendation{
		Sets:        baseRec.Sets,
		RepsMin:     baseRec.RepsMin,
		RepsMax:     baseRec.RepsMax,
		RestSeconds: baseRec.RestSeconds,
		Intensity:   baseRec.Intensity,
		Notes:       baseRec.Notes,
	}

	// Adjust based on fitness level
	switch strings.ToLower(fitnessLevel) {
	case "beginner":
		recommendation.Sets = max(1, recommendation.Sets-1)
		recommendation.Intensity *= 0.8  // Reduce intensity by 20%
		recommendation.RestSeconds += 30 // Add 30 seconds rest
	case "advanced":
		recommendation.Sets += 1
		recommendation.Intensity *= 1.1 // Increase intensity by 10%
		if recommendation.Intensity > 95.0 {
			recommendation.Intensity = 95.0
		}
	}

	// Adjust rest periods based on exercise type
	if restPeriod, exists := ExerciseTypeRestPeriods[strings.ToLower(exerciseType)]; exists {
		recommendation.RestSeconds = restPeriod
	}

	// Apply weekly progression (gradual increase over time)
	progressionFactor := pe.calculateProgressionFactor(weekNumber, fitnessLevel)
	recommendation.Intensity *= progressionFactor
	if recommendation.Intensity > 95.0 {
		recommendation.Intensity = 95.0
	}

	return &recommendation, nil
}

// calculateProgressionFactor determines how to progress intensity over weeks
func (pe *ProgressionEngine) calculateProgressionFactor(weekNumber int, fitnessLevel string) float64 {
	baseProgression := 1.0

	// Different progression rates based on fitness level
	var weeklyIncrease float64
	switch strings.ToLower(fitnessLevel) {
	case "beginner":
		weeklyIncrease = 0.025 // 2.5% per week
	case "intermediate":
		weeklyIncrease = 0.015 // 1.5% per week
	case "advanced":
		weeklyIncrease = 0.01 // 1% per week
	default:
		weeklyIncrease = 0.015
	}

	// Calculate progression with diminishing returns
	progression := baseProgression + (float64(weekNumber-1) * weeklyIncrease)

	// Cap progression to prevent unrealistic increases
	maxProgression := 1.5 // Maximum 50% increase
	if progression > maxProgression {
		progression = maxProgression
	}

	return progression
}

// CalculateRestPeriod determines optimal rest period based on exercise and goals
func (pe *ProgressionEngine) CalculateRestPeriod(exerciseType string, intensity float64, goals []string) int {
	baseRest := ExerciseTypeRestPeriods["isolation"] // default 90 seconds

	// Get base rest from exercise type
	if rest, exists := ExerciseTypeRestPeriods[strings.ToLower(exerciseType)]; exists {
		baseRest = rest
	}

	// Adjust based on intensity
	if intensity >= 85.0 {
		baseRest += 60 // Add 1 minute for high intensity
	} else if intensity <= 60.0 {
		baseRest -= 30 // Reduce 30 seconds for low intensity
	}

	// Adjust based on goals
	if len(goals) > 0 {
		primaryGoal := strings.ToLower(goals[0])
		if goalRec, exists := GoalBasedRecommendations[primaryGoal]; exists {
			// Blend with goal-specific rest period
			baseRest = (baseRest + goalRec.RestSeconds) / 2
		}
	}

	// Ensure minimum rest period
	if baseRest < 15 {
		baseRest = 15
	}

	return baseRest
}

// TrackExerciseProgress records and analyzes exercise performance
func (pe *ProgressionEngine) TrackExerciseProgress(ctx context.Context, userID string, exerciseID string, sessionData *models.ExerciseLog) (*models.ExerciseProgress, error) {
	// Create or update exercise progress
	progress := &models.ExerciseProgress{
		ExerciseID:    exerciseID,
		Weight:        0,
		Reps:          sessionData.RepsCompleted,
		Sets:          sessionData.SetNumber,
		RPE:           sessionData.RPEScore,
		Notes:         sessionData.Notes,
		LastPerformed: time.Now(),
	}

	if sessionData.WeightKg != nil {
		progress.Weight = *sessionData.WeightKg
	}

	// Calculate one-rep max if weight is provided
	if progress.Weight > 0 && progress.Reps > 0 {
		oneRepMax := pe.calculateOneRepMax(progress.Weight, progress.Reps)

		// Update personal record if this is a new PR
		if progress.PersonalRecord == nil || oneRepMax > progress.PersonalRecord.OneRepMax {
			progress.PersonalRecord = &models.PersonalRecord{
				Weight:    progress.Weight,
				Reps:      progress.Reps,
				OneRepMax: oneRepMax,
				Date:      time.Now(),
			}
		}
	}

	return progress, nil
}

// calculateOneRepMax estimates 1RM using Epley formula
func (pe *ProgressionEngine) calculateOneRepMax(weight float64, reps int) float64 {
	if reps == 1 {
		return weight
	}
	// Epley formula: 1RM = weight × (1 + reps/30)
	return weight * (1 + float64(reps)/30.0)
}

// GenerateProgressAnalytics creates comprehensive progress analytics
func (pe *ProgressionEngine) GenerateProgressAnalytics(ctx context.Context, userID string, programID string, exerciseProgress map[string]*models.ExerciseProgress) (*models.ProgressMetrics, error) {
	metrics := &models.ProgressMetrics{
		StrengthGains:     make(map[string]float64),
		WeightProgression: make(map[string]float64),
	}

	totalRPE := 0.0
	rpeCount := 0
	totalVolume := 0.0

	// Analyze each exercise's progress
	for exerciseID, progress := range exerciseProgress {
		if progress == nil {
			continue
		}

		// Calculate strength gains (based on 1RM improvements)
		if progress.PersonalRecord != nil {
			metrics.StrengthGains[exerciseID] = progress.PersonalRecord.OneRepMax
		}

		// Calculate weight progression
		if progress.Weight > 0 {
			metrics.WeightProgression[exerciseID] = progress.Weight
		}

		// Calculate volume (sets × reps × weight)
		volume := float64(progress.Sets) * float64(progress.Reps) * progress.Weight
		totalVolume += volume

		// Track RPE for average calculation
		if progress.RPE > 0 {
			totalRPE += float64(progress.RPE)
			rpeCount++
		}
	}

	// Calculate overall metrics
	metrics.VolumeProgression = totalVolume

	if rpeCount > 0 {
		metrics.AverageRPE = totalRPE / float64(rpeCount)
	}

	// Calculate consistency score (placeholder - would need session history)
	metrics.ConsistencyScore = pe.calculateConsistencyScore(exerciseProgress)

	return metrics, nil
}

// calculateConsistencyScore determines workout consistency
func (pe *ProgressionEngine) calculateConsistencyScore(exerciseProgress map[string]*models.ExerciseProgress) float64 {
	// Simplified consistency calculation
	// In a real implementation, this would analyze session frequency over time

	completedExercises := 0
	for _, progress := range exerciseProgress {
		if progress != nil && progress.LastPerformed.After(time.Now().AddDate(0, 0, -7)) {
			completedExercises++
		}
	}

	totalExercises := len(exerciseProgress)
	if totalExercises == 0 {
		return 0.0
	}

	return (float64(completedExercises) / float64(totalExercises)) * 100.0
}

// SuggestProgressionAdjustment recommends changes based on performance
func (pe *ProgressionEngine) SuggestProgressionAdjustment(ctx context.Context, progress *models.ExerciseProgress, currentRecommendation *SetsRepsRecommendation) (*SetsRepsRecommendation, error) {
	if progress == nil || currentRecommendation == nil {
		return currentRecommendation, nil
	}

	// Create adjusted recommendation
	adjusted := *currentRecommendation

	// Adjust based on RPE (Rate of Perceived Exertion)
	if progress.RPE > 0 {
		switch {
		case progress.RPE <= 6: // Too easy
			adjusted.Intensity *= 1.05 // Increase by 5%
			adjusted.Notes = "Increase weight - previous session was too easy"
		case progress.RPE >= 9: // Too hard
			adjusted.Intensity *= 0.95 // Decrease by 5%
			adjusted.Notes = "Reduce weight - previous session was too challenging"
		case progress.RPE >= 7 && progress.RPE <= 8: // Perfect range
			adjusted.Notes = "Maintain current intensity - good progression"
		}
	}

	// Adjust based on rep completion
	if progress.Reps > 0 {
		targetRepsMax := currentRecommendation.RepsMax
		if progress.Reps > targetRepsMax {
			// Completed more reps than target - increase weight
			adjusted.Intensity *= 1.025 // Increase by 2.5%
			adjusted.Notes = "Increase weight - exceeded target reps"
		}
	}

	// Ensure intensity stays within reasonable bounds
	if adjusted.Intensity > 95.0 {
		adjusted.Intensity = 95.0
	}
	if adjusted.Intensity < 40.0 {
		adjusted.Intensity = 40.0
	}

	return &adjusted, nil
}

// GenerateWeeklyProgressionPlan creates a progression plan for the week
func (pe *ProgressionEngine) GenerateWeeklyProgressionPlan(ctx context.Context, baseProgram *models.WorkoutProgram, weekNumber int, userProgress map[string]*models.ExerciseProgress) (*models.WeeklyPlan, error) {
	if baseProgram == nil {
		return nil, fmt.Errorf("base program is required")
	}

	weeklyPlan := &models.WeeklyPlan{
		ProgramID:             baseProgram.ID,
		WeekNumber:            weekNumber,
		IntensityLevel:        pe.calculateWeeklyIntensity(weekNumber, baseProgram.Level),
		VolumeIncreasePercent: pe.calculateVolumeIncrease(weekNumber, baseProgram.Level),
		DailyWorkouts:         make(map[int]*models.DailyWorkout),
	}

	// Generate progression notes
	weeklyPlan.ProgressionNotesEn = pe.generateProgressionNotes(weekNumber, baseProgram.Level)

	return weeklyPlan, nil
}

// calculateWeeklyIntensity determines intensity level for the week
func (pe *ProgressionEngine) calculateWeeklyIntensity(weekNumber int, fitnessLevel string) string {
	progressionFactor := pe.calculateProgressionFactor(weekNumber, fitnessLevel)

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

// calculateVolumeIncrease determines volume increase percentage
func (pe *ProgressionEngine) calculateVolumeIncrease(weekNumber int, fitnessLevel string) float64 {
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

// generateProgressionNotes creates helpful progression notes
func (pe *ProgressionEngine) generateProgressionNotes(weekNumber int, fitnessLevel string) string {
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

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
