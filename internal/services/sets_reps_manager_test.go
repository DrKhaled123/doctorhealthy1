package services

import (
	"context"
	"testing"
	"time"

	"api-key-generator/internal/models"
)

func TestSetsRepsManager_GenerateExerciseRecommendations(t *testing.T) {
	srm := NewSetsRepsManager()
	ctx := context.Background()

	// Create sample exercises
	exercises := []*models.WorkoutExerciseData{
		{
			ID:             "bench_press",
			NameEn:         "Bench Press",
			Category:       "Strength",
			PrimaryMuscles: []string{"chest", "triceps"},
		},
		{
			ID:             "bicep_curl",
			NameEn:         "Bicep Curl",
			Category:       "Isolation",
			PrimaryMuscles: []string{"biceps"},
		},
	}

	req := &ExerciseRecommendationRequest{
		Exercises:    exercises,
		Goals:        []string{"hypertrophy"},
		FitnessLevel: "intermediate",
		WeekNumber:   1,
	}

	response, err := srm.GenerateExerciseRecommendations(ctx, req)

	if err != nil {
		t.Errorf("GenerateExerciseRecommendations() error = %v", err)
		return
	}

	if len(response.Recommendations) != 2 {
		t.Errorf("GenerateExerciseRecommendations() recommendations count = %v, want 2", len(response.Recommendations))
	}

	// Check first recommendation (bench press - should be compound)
	benchRec := response.Recommendations[0]
	if benchRec.ExerciseID != "bench_press" {
		t.Errorf("First recommendation exerciseID = %v, want bench_press", benchRec.ExerciseID)
	}

	if benchRec.ExerciseType != "compound" {
		t.Errorf("Bench press should be classified as compound, got %v", benchRec.ExerciseType)
	}

	// Check that sets and reps are reasonable for hypertrophy
	if benchRec.Sets < 2 || benchRec.Sets > 5 {
		t.Errorf("Bench press sets = %v, should be between 2-5", benchRec.Sets)
	}

	if benchRec.RepsMin < 6 || benchRec.RepsMax > 12 {
		t.Errorf("Bench press reps = %v-%v, should be in hypertrophy range (6-12)", benchRec.RepsMin, benchRec.RepsMax)
	}

	// Check second recommendation (bicep curl - should be isolation)
	bicepRec := response.Recommendations[1]
	if bicepRec.ExerciseType != "isolation" {
		t.Errorf("Bicep curl should be classified as isolation, got %v", bicepRec.ExerciseType)
	}
}

func TestSetsRepsManager_GenerateExerciseRecommendationsWithProgress(t *testing.T) {
	srm := NewSetsRepsManager()
	ctx := context.Background()

	exercises := []*models.WorkoutExerciseData{
		{
			ID:             "squat",
			NameEn:         "Squat",
			Category:       "Strength",
			PrimaryMuscles: []string{"quadriceps", "glutes"},
		},
	}

	// Create user progress data
	userProgress := map[string]*models.ExerciseProgress{
		"squat": {
			ExerciseID:    "squat",
			Weight:        100.0,
			Reps:          8,
			RPE:           6, // Too easy - should suggest increase
			LastPerformed: time.Now().AddDate(0, 0, -2),
		},
	}

	req := &ExerciseRecommendationRequest{
		Exercises:    exercises,
		Goals:        []string{"strength"},
		FitnessLevel: "intermediate",
		WeekNumber:   2,
		UserProgress: userProgress,
	}

	response, err := srm.GenerateExerciseRecommendations(ctx, req)

	if err != nil {
		t.Errorf("GenerateExerciseRecommendations() error = %v", err)
		return
	}

	rec := response.Recommendations[0]

	// Should have previous performance data
	if rec.PreviousWeight != 100.0 {
		t.Errorf("Previous weight = %v, want 100.0", rec.PreviousWeight)
	}

	if rec.PreviousReps != 8 {
		t.Errorf("Previous reps = %v, want 8", rec.PreviousReps)
	}

	if rec.LastPerformed == nil {
		t.Error("Last performed should be set")
	}

	// Should have progression notes suggesting increase (RPE was too low)
	if rec.ProgressionNotes == "" {
		t.Error("Should have progression notes when RPE is too low")
	}
}

func TestSetsRepsManager_TrackWorkoutSession(t *testing.T) {
	srm := NewSetsRepsManager()
	ctx := context.Background()

	req := &WorkoutSessionRequest{
		UserID:          "user123",
		ProgramID:       "program456",
		DailyWorkoutID:  "workout789",
		DurationMinutes: 60,
		OverallRPE:      7,
		Notes:           "Good workout",
		Completed:       true,
		ExercisePerformance: []*ExercisePerformance{
			{
				ExerciseID: "bench_press",
				Sets: []*SetPerformance{
					{
						RepsCompleted: 8,
						WeightKg:      &[]float64{80.0}[0],
						RestSeconds:   120,
						RPE:           7,
						Notes:         "Good form",
					},
					{
						RepsCompleted: 7,
						WeightKg:      &[]float64{80.0}[0],
						RestSeconds:   120,
						RPE:           8,
					},
					{
						RepsCompleted: 6,
						WeightKg:      &[]float64{80.0}[0],
						RestSeconds:   120,
						RPE:           9,
					},
				},
			},
		},
	}

	response, err := srm.TrackWorkoutSession(ctx, req)

	if err != nil {
		t.Errorf("TrackWorkoutSession() error = %v", err)
		return
	}

	// Check session log
	if response.SessionLog.UserID != "user123" {
		t.Errorf("Session log userID = %v, want user123", response.SessionLog.UserID)
	}

	if response.SessionLog.DurationMinutes != 60 {
		t.Errorf("Session log duration = %v, want 60", response.SessionLog.DurationMinutes)
	}

	if response.SessionLog.RPEScore != 7 {
		t.Errorf("Session log RPE = %v, want 7", response.SessionLog.RPEScore)
	}

	if !response.SessionLog.Completed {
		t.Error("Session should be marked as completed")
	}

	// Check exercise logs (should have 3 sets)
	if len(response.SessionLog.ExerciseLogs) != 3 {
		t.Errorf("Exercise logs count = %v, want 3", len(response.SessionLog.ExerciseLogs))
	}

	// Check exercise progress
	if len(response.ExerciseProgress) != 1 {
		t.Errorf("Exercise progress count = %v, want 1", len(response.ExerciseProgress))
	}

	benchProgress, exists := response.ExerciseProgress["bench_press"]
	if !exists {
		t.Error("Should have bench press progress")
	} else {
		// Should use the best set (first set with 8 reps)
		if benchProgress.Reps != 8 {
			t.Errorf("Progress reps = %v, want 8 (best set)", benchProgress.Reps)
		}
		if benchProgress.Weight != 80.0 {
			t.Errorf("Progress weight = %v, want 80.0", benchProgress.Weight)
		}
	}

	// Check analytics
	if response.Analytics == nil {
		t.Error("Should have analytics")
	}

	// Check recommendations
	if len(response.NextRecommendations) == 0 {
		t.Error("Should have next session recommendations")
	}
}

func TestSetsRepsManager_determineExerciseType(t *testing.T) {
	srm := NewSetsRepsManager()

	tests := []struct {
		name     string
		exercise *models.WorkoutExerciseData
		want     string
	}{
		{
			name: "Compound exercise - squat",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Squat",
				Category:       "Strength",
				PrimaryMuscles: []string{"quadriceps", "glutes"},
			},
			want: "compound",
		},
		{
			name: "Isolation exercise - bicep curl",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Bicep Curl",
				Category:       "Isolation",
				PrimaryMuscles: []string{"biceps"},
			},
			want: "isolation",
		},
		{
			name: "Cardio exercise",
			exercise: &models.WorkoutExerciseData{
				NameEn:   "Running",
				Category: "Cardio",
			},
			want: "cardio",
		},
		{
			name: "Core exercise",
			exercise: &models.WorkoutExerciseData{
				NameEn:   "Plank",
				Category: "Core",
			},
			want: "core",
		},
		{
			name: "Plyometric exercise",
			exercise: &models.WorkoutExerciseData{
				NameEn:   "Box Jump",
				Category: "Plyometric",
			},
			want: "plyometric",
		},
		{
			name: "Compound by name - bench press",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Bench Press",
				Category:       "Strength",
				PrimaryMuscles: []string{"chest"},
			},
			want: "compound",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srm.determineExerciseType(tt.exercise)
			if got != tt.want {
				t.Errorf("determineExerciseType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetsRepsManager_isCompoundExercise(t *testing.T) {
	srm := NewSetsRepsManager()

	tests := []struct {
		name     string
		exercise *models.WorkoutExerciseData
		want     bool
	}{
		{
			name: "Squat - compound by name",
			exercise: &models.WorkoutExerciseData{
				NameEn: "Back Squat",
			},
			want: true,
		},
		{
			name: "Deadlift - compound by name",
			exercise: &models.WorkoutExerciseData{
				NameEn: "Romanian Deadlift",
			},
			want: true,
		},
		{
			name: "Multiple muscle groups - compound",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Pull-up",
				PrimaryMuscles: []string{"lats", "biceps"},
			},
			want: true,
		},
		{
			name: "Single muscle group - isolation",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Bicep Curl",
				PrimaryMuscles: []string{"biceps"},
			},
			want: false,
		},
		{
			name: "Bench press - compound by name",
			exercise: &models.WorkoutExerciseData{
				NameEn: "Incline Bench Press",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srm.isCompoundExercise(tt.exercise)
			if got != tt.want {
				t.Errorf("isCompoundExercise() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetsRepsManager_formatRepsDisplay(t *testing.T) {
	srm := NewSetsRepsManager()

	tests := []struct {
		name string
		min  int
		max  int
		want string
	}{
		{
			name: "Same min and max",
			min:  8,
			max:  8,
			want: "8",
		},
		{
			name: "Different min and max",
			min:  6,
			max:  12,
			want: "6-12",
		},
		{
			name: "Single rep",
			min:  1,
			max:  1,
			want: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srm.formatRepsDisplay(tt.min, tt.max)
			if got != tt.want {
				t.Errorf("formatRepsDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetsRepsManager_formatRestDisplay(t *testing.T) {
	srm := NewSetsRepsManager()

	tests := []struct {
		name    string
		seconds int
		want    string
	}{
		{
			name:    "Less than a minute",
			seconds: 45,
			want:    "45s",
		},
		{
			name:    "Exactly one minute",
			seconds: 60,
			want:    "1m",
		},
		{
			name:    "One and half minutes",
			seconds: 90,
			want:    "1m 30s",
		},
		{
			name:    "Three minutes",
			seconds: 180,
			want:    "3m",
		},
		{
			name:    "Two minutes thirty seconds",
			seconds: 150,
			want:    "2m 30s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srm.formatRestDisplay(tt.seconds)
			if got != tt.want {
				t.Errorf("formatRestDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetsRepsManager_findBestSet(t *testing.T) {
	srm := NewSetsRepsManager()

	sets := []*SetPerformance{
		{
			RepsCompleted: 8,
			WeightKg:      &[]float64{80.0}[0], // Score: 640
		},
		{
			RepsCompleted: 10,
			WeightKg:      &[]float64{70.0}[0], // Score: 700 (best)
		},
		{
			RepsCompleted: 6,
			WeightKg:      &[]float64{90.0}[0], // Score: 540
		},
	}

	bestSet := srm.findBestSet(sets)

	if bestSet == nil {
		t.Error("findBestSet() should return a set")
		return
	}

	if bestSet.RepsCompleted != 10 || *bestSet.WeightKg != 70.0 {
		t.Errorf("findBestSet() = reps:%v weight:%v, want reps:10 weight:70.0", bestSet.RepsCompleted, *bestSet.WeightKg)
	}
}

func TestSetsRepsManager_calculateSetScore(t *testing.T) {
	srm := NewSetsRepsManager()

	tests := []struct {
		name string
		set  *SetPerformance
		want float64
	}{
		{
			name: "With weight",
			set: &SetPerformance{
				RepsCompleted: 8,
				WeightKg:      &[]float64{80.0}[0],
			},
			want: 640.0, // 8 * 80
		},
		{
			name: "Without weight",
			set: &SetPerformance{
				RepsCompleted: 10,
				WeightKg:      nil,
			},
			want: 0.0, // 10 * 0
		},
		{
			name: "Zero reps",
			set: &SetPerformance{
				RepsCompleted: 0,
				WeightKg:      &[]float64{100.0}[0],
			},
			want: 0.0, // 0 * 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srm.calculateSetScore(tt.set)
			if got != tt.want {
				t.Errorf("calculateSetScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetsRepsManager_GenerateProgressReport(t *testing.T) {
	srm := NewSetsRepsManager()
	ctx := context.Background()

	timeRange := TimeRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}

	report, err := srm.GenerateProgressReport(ctx, "user123", "program456", timeRange)

	if err != nil {
		t.Errorf("GenerateProgressReport() error = %v", err)
		return
	}

	if report.UserID != "user123" {
		t.Errorf("Report userID = %v, want user123", report.UserID)
	}

	if report.ProgramID != "program456" {
		t.Errorf("Report programID = %v, want program456", report.ProgramID)
	}

	// Check that report has required sections
	if report.StrengthImprovements == nil {
		t.Error("Report should have strength improvements")
	}

	if report.VolumeProgression == nil {
		t.Error("Report should have volume progression")
	}

	if report.PerformanceMetrics == nil {
		t.Error("Report should have performance metrics")
	}

	// Check sample data structure
	if len(report.StrengthImprovements) == 0 {
		t.Error("Report should have sample strength improvements")
	}
}
