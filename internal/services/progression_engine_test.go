package services

import (
	"context"
	"testing"
	"time"

	"api-key-generator/internal/models"
)

func TestProgressionEngine_GenerateSetsRepsRecommendation(t *testing.T) {
	pe := NewProgressionEngine()
	ctx := context.Background()

	tests := []struct {
		name         string
		goals        []string
		exerciseType string
		fitnessLevel string
		weekNumber   int
		wantSets     int
		wantRepsMin  int
		wantRepsMax  int
		wantErr      bool
	}{
		{
			name:         "Strength goal for beginner",
			goals:        []string{"strength"},
			exerciseType: "compound",
			fitnessLevel: "beginner",
			weekNumber:   1,
			wantSets:     2, // reduced from 3 for beginner
			wantRepsMin:  1,
			wantRepsMax:  5,
			wantErr:      false,
		},
		{
			name:         "Hypertrophy goal for intermediate",
			goals:        []string{"hypertrophy"},
			exerciseType: "isolation",
			fitnessLevel: "intermediate",
			weekNumber:   1,
			wantSets:     3,
			wantRepsMin:  6,
			wantRepsMax:  12,
			wantErr:      false,
		},
		{
			name:         "Weight loss goal for advanced",
			goals:        []string{"weight_loss"},
			exerciseType: "cardio",
			fitnessLevel: "advanced",
			weekNumber:   1,
			wantSets:     4, // increased from 3 for advanced
			wantRepsMin:  12,
			wantRepsMax:  20,
			wantErr:      false,
		},
		{
			name:         "Endurance goal with progression",
			goals:        []string{"endurance"},
			exerciseType: "isolation",
			fitnessLevel: "intermediate",
			weekNumber:   4,
			wantSets:     2,
			wantRepsMin:  15,
			wantRepsMax:  25,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, err := pe.GenerateSetsRepsRecommendation(ctx, tt.goals, tt.exerciseType, tt.fitnessLevel, tt.weekNumber)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSetsRepsRecommendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if rec.Sets != tt.wantSets {
					t.Errorf("GenerateSetsRepsRecommendation() sets = %v, want %v", rec.Sets, tt.wantSets)
				}
				if rec.RepsMin != tt.wantRepsMin {
					t.Errorf("GenerateSetsRepsRecommendation() repsMin = %v, want %v", rec.RepsMin, tt.wantRepsMin)
				}
				if rec.RepsMax != tt.wantRepsMax {
					t.Errorf("GenerateSetsRepsRecommendation() repsMax = %v, want %v", rec.RepsMax, tt.wantRepsMax)
				}
			}
		})
	}
}

func TestProgressionEngine_CalculateRestPeriod(t *testing.T) {
	pe := NewProgressionEngine()

	tests := []struct {
		name         string
		exerciseType string
		intensity    float64
		goals        []string
		wantMin      int
		wantMax      int
	}{
		{
			name:         "Compound exercise high intensity",
			exerciseType: "compound",
			intensity:    90.0,
			goals:        []string{"strength"},
			wantMin:      180, // Should be at least 3 minutes
			wantMax:      300, // Should not exceed 5 minutes
		},
		{
			name:         "Isolation exercise moderate intensity",
			exerciseType: "isolation",
			intensity:    70.0,
			goals:        []string{"hypertrophy"},
			wantMin:      60,  // Should be at least 1 minute
			wantMax:      120, // Should not exceed 2 minutes
		},
		{
			name:         "Cardio exercise low intensity",
			exerciseType: "cardio",
			intensity:    50.0,
			goals:        []string{"weight_loss"},
			wantMin:      15, // Should be at least 15 seconds
			wantMax:      60, // Should not exceed 1 minute
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := pe.CalculateRestPeriod(tt.exerciseType, tt.intensity, tt.goals)

			if rest < tt.wantMin {
				t.Errorf("CalculateRestPeriod() = %v, want at least %v", rest, tt.wantMin)
			}
			if rest > tt.wantMax {
				t.Errorf("CalculateRestPeriod() = %v, want at most %v", rest, tt.wantMax)
			}
		})
	}
}

func TestProgressionEngine_TrackExerciseProgress(t *testing.T) {
	pe := NewProgressionEngine()
	ctx := context.Background()

	sessionData := &models.ExerciseLog{
		ExerciseID:    "bench_press",
		RepsCompleted: 8,
		WeightKg:      &[]float64{80.0}[0],
		RPEScore:      7,
		Notes:         "Good form",
	}

	progress, err := pe.TrackExerciseProgress(ctx, "user123", "bench_press", sessionData)

	if err != nil {
		t.Errorf("TrackExerciseProgress() error = %v", err)
		return
	}

	if progress.ExerciseID != "bench_press" {
		t.Errorf("TrackExerciseProgress() exerciseID = %v, want bench_press", progress.ExerciseID)
	}

	if progress.Weight != 80.0 {
		t.Errorf("TrackExerciseProgress() weight = %v, want 80.0", progress.Weight)
	}

	if progress.Reps != 8 {
		t.Errorf("TrackExerciseProgress() reps = %v, want 8", progress.Reps)
	}

	if progress.RPE != 7 {
		t.Errorf("TrackExerciseProgress() RPE = %v, want 7", progress.RPE)
	}

	// Check that one-rep max was calculated
	if progress.PersonalRecord == nil {
		t.Error("TrackExerciseProgress() should have calculated personal record")
	} else {
		expectedOneRM := 80.0 * (1 + 8.0/30.0) // Epley formula
		if progress.PersonalRecord.OneRepMax != expectedOneRM {
			t.Errorf("TrackExerciseProgress() oneRepMax = %v, want %v", progress.PersonalRecord.OneRepMax, expectedOneRM)
		}
	}
}

func TestProgressionEngine_calculateOneRepMax(t *testing.T) {
	pe := NewProgressionEngine()

	tests := []struct {
		name   string
		weight float64
		reps   int
		want   float64
	}{
		{
			name:   "One rep max for 1 rep",
			weight: 100.0,
			reps:   1,
			want:   100.0,
		},
		{
			name:   "One rep max for 5 reps",
			weight: 80.0,
			reps:   5,
			want:   93.33333333333334, // 80.0 * (1 + 5.0/30.0)
		},
		{
			name:   "One rep max for 10 reps",
			weight: 60.0,
			reps:   10,
			want:   60.0 * (1 + 10.0/30.0), // 80.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pe.calculateOneRepMax(tt.weight, tt.reps)
			if got != tt.want {
				t.Errorf("calculateOneRepMax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressionEngine_GenerateProgressAnalytics(t *testing.T) {
	pe := NewProgressionEngine()
	ctx := context.Background()

	// Create sample exercise progress data
	exerciseProgress := map[string]*models.ExerciseProgress{
		"bench_press": {
			ExerciseID: "bench_press",
			Weight:     80.0,
			Reps:       8,
			Sets:       3,
			RPE:        7,
			PersonalRecord: &models.PersonalRecord{
				Weight:    80.0,
				Reps:      8,
				OneRepMax: 106.67,
				Date:      time.Now(),
			},
		},
		"squat": {
			ExerciseID: "squat",
			Weight:     100.0,
			Reps:       5,
			Sets:       3,
			RPE:        8,
			PersonalRecord: &models.PersonalRecord{
				Weight:    100.0,
				Reps:      5,
				OneRepMax: 116.67,
				Date:      time.Now(),
			},
		},
	}

	analytics, err := pe.GenerateProgressAnalytics(ctx, "user123", "program456", exerciseProgress)

	if err != nil {
		t.Errorf("GenerateProgressAnalytics() error = %v", err)
		return
	}

	// Check that strength gains were recorded
	if len(analytics.StrengthGains) != 2 {
		t.Errorf("GenerateProgressAnalytics() strengthGains count = %v, want 2", len(analytics.StrengthGains))
	}

	// Check that weight progression was recorded
	if len(analytics.WeightProgression) != 2 {
		t.Errorf("GenerateProgressAnalytics() weightProgression count = %v, want 2", len(analytics.WeightProgression))
	}

	// Check average RPE calculation
	expectedAvgRPE := (7.0 + 8.0) / 2.0
	if analytics.AverageRPE != expectedAvgRPE {
		t.Errorf("GenerateProgressAnalytics() averageRPE = %v, want %v", analytics.AverageRPE, expectedAvgRPE)
	}

	// Check volume progression calculation
	expectedVolume := (3.0 * 8.0 * 80.0) + (3.0 * 5.0 * 100.0) // 1920 + 1500 = 3420
	if analytics.VolumeProgression != expectedVolume {
		t.Errorf("GenerateProgressAnalytics() volumeProgression = %v, want %v", analytics.VolumeProgression, expectedVolume)
	}
}

func TestProgressionEngine_SuggestProgressionAdjustment(t *testing.T) {
	pe := NewProgressionEngine()
	ctx := context.Background()

	currentRec := &SetsRepsRecommendation{
		Sets:        3,
		RepsMin:     6,
		RepsMax:     12,
		RestSeconds: 90,
		Intensity:   70.0,
		Notes:       "Base recommendation",
	}

	tests := []struct {
		name         string
		progress     *models.ExerciseProgress
		wantIncrease bool
		wantDecrease bool
	}{
		{
			name: "RPE too low - should increase intensity",
			progress: &models.ExerciseProgress{
				RPE: 5, // Too easy
			},
			wantIncrease: true,
			wantDecrease: false,
		},
		{
			name: "RPE too high - should decrease intensity",
			progress: &models.ExerciseProgress{
				RPE: 9, // Too hard
			},
			wantIncrease: false,
			wantDecrease: true,
		},
		{
			name: "RPE perfect - should maintain",
			progress: &models.ExerciseProgress{
				RPE: 7, // Perfect range
			},
			wantIncrease: false,
			wantDecrease: false,
		},
		{
			name: "Exceeded target reps - should increase intensity",
			progress: &models.ExerciseProgress{
				Reps: 15, // More than max target of 12
				RPE:  7,
			},
			wantIncrease: true,
			wantDecrease: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adjusted, err := pe.SuggestProgressionAdjustment(ctx, tt.progress, currentRec)

			if err != nil {
				t.Errorf("SuggestProgressionAdjustment() error = %v", err)
				return
			}

			if tt.wantIncrease && adjusted.Intensity <= currentRec.Intensity {
				t.Errorf("SuggestProgressionAdjustment() should increase intensity, got %v, original %v", adjusted.Intensity, currentRec.Intensity)
			}

			if tt.wantDecrease && adjusted.Intensity >= currentRec.Intensity {
				t.Errorf("SuggestProgressionAdjustment() should decrease intensity, got %v, original %v", adjusted.Intensity, currentRec.Intensity)
			}

			if !tt.wantIncrease && !tt.wantDecrease && adjusted.Intensity != currentRec.Intensity {
				// Allow small variations for perfect RPE case
				diff := adjusted.Intensity - currentRec.Intensity
				if diff > 1.0 || diff < -1.0 {
					t.Errorf("SuggestProgressionAdjustment() should maintain intensity, got %v, original %v", adjusted.Intensity, currentRec.Intensity)
				}
			}
		})
	}
}

func TestProgressionEngine_calculateProgressionFactor(t *testing.T) {
	pe := NewProgressionEngine()

	tests := []struct {
		name         string
		weekNumber   int
		fitnessLevel string
		wantMin      float64
		wantMax      float64
	}{
		{
			name:         "Beginner week 1",
			weekNumber:   1,
			fitnessLevel: "beginner",
			wantMin:      1.0,
			wantMax:      1.1,
		},
		{
			name:         "Beginner week 4",
			weekNumber:   4,
			fitnessLevel: "beginner",
			wantMin:      1.05,
			wantMax:      1.2,
		},
		{
			name:         "Advanced week 8",
			weekNumber:   8,
			fitnessLevel: "advanced",
			wantMin:      1.05,
			wantMax:      1.15,
		},
		{
			name:         "Very long program should cap progression",
			weekNumber:   50,
			fitnessLevel: "beginner",
			wantMin:      1.4,
			wantMax:      1.5, // Should be capped at 1.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor := pe.calculateProgressionFactor(tt.weekNumber, tt.fitnessLevel)

			if factor < tt.wantMin {
				t.Errorf("calculateProgressionFactor() = %v, want at least %v", factor, tt.wantMin)
			}
			if factor > tt.wantMax {
				t.Errorf("calculateProgressionFactor() = %v, want at most %v", factor, tt.wantMax)
			}
		})
	}
}

func TestProgressionEngine_GenerateWeeklyProgressionPlan(t *testing.T) {
	pe := NewProgressionEngine()
	ctx := context.Background()

	baseProgram := &models.WorkoutProgram{
		ID:    "program123",
		Level: "intermediate",
	}

	weeklyPlan, err := pe.GenerateWeeklyProgressionPlan(ctx, baseProgram, 3, nil)

	if err != nil {
		t.Errorf("GenerateWeeklyProgressionPlan() error = %v", err)
		return
	}

	if weeklyPlan.ProgramID != "program123" {
		t.Errorf("GenerateWeeklyProgressionPlan() programID = %v, want program123", weeklyPlan.ProgramID)
	}

	if weeklyPlan.WeekNumber != 3 {
		t.Errorf("GenerateWeeklyProgressionPlan() weekNumber = %v, want 3", weeklyPlan.WeekNumber)
	}

	if weeklyPlan.IntensityLevel == "" {
		t.Error("GenerateWeeklyProgressionPlan() should set intensity level")
	}

	if weeklyPlan.VolumeIncreasePercent < 0 {
		t.Errorf("GenerateWeeklyProgressionPlan() volumeIncrease = %v, should be non-negative", weeklyPlan.VolumeIncreasePercent)
	}

	if weeklyPlan.ProgressionNotesEn == "" {
		t.Error("GenerateWeeklyProgressionPlan() should set progression notes")
	}
}
