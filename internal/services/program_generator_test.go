package services

import (
	"context"
	"testing"

	"api-key-generator/internal/models"
)

func TestProgramGenerator_ClassifyExercise(t *testing.T) {
	pg := NewProgramGenerator()
	ctx := context.Background()

	tests := []struct {
		name     string
		exercise *models.WorkoutExerciseData
		want     *ExerciseClassification
		wantErr  bool
	}{
		{
			name: "Compound strength exercise",
			exercise: &models.WorkoutExerciseData{
				ID:             "squat_001",
				NameEn:         "Barbell Squat",
				Category:       "Strength",
				PrimaryMuscles: []string{"quadriceps", "glutes"},
				Equipment:      []string{"barbell"},
				Difficulty:     "intermediate",
			},
			want: &ExerciseClassification{
				Purpose:      "strength",
				MuscleGroups: []string{"quadriceps", "glutes"},
				Equipment:    []string{"barbell"},
				Difficulty:   "intermediate",
				ExerciseType: "compound",
			},
			wantErr: false,
		},
		{
			name: "Cardio exercise",
			exercise: &models.WorkoutExerciseData{
				ID:             "running_001",
				NameEn:         "Treadmill Running",
				Category:       "Cardio",
				PrimaryMuscles: []string{"legs"},
				Equipment:      []string{"treadmill"},
				Difficulty:     "beginner",
			},
			want: &ExerciseClassification{
				Purpose:      "cardio",
				MuscleGroups: []string{"legs"},
				Equipment:    []string{"treadmill"},
				Difficulty:   "beginner",
				ExerciseType: "cardio",
			},
			wantErr: false,
		},
		{
			name: "Isolation exercise",
			exercise: &models.WorkoutExerciseData{
				ID:             "bicep_curl_001",
				NameEn:         "Bicep Curl",
				Category:       "Strength",
				PrimaryMuscles: []string{"biceps"},
				Equipment:      []string{"dumbbells"},
				Difficulty:     "beginner",
			},
			want: &ExerciseClassification{
				Purpose:      "strength",
				MuscleGroups: []string{"biceps"},
				Equipment:    []string{"dumbbells"},
				Difficulty:   "beginner",
				ExerciseType: "isolation",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pg.ClassifyExercise(ctx, tt.exercise)

			if (err != nil) != tt.wantErr {
				t.Errorf("ClassifyExercise() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.Purpose != tt.want.Purpose {
					t.Errorf("ClassifyExercise() purpose = %v, want %v", got.Purpose, tt.want.Purpose)
				}
				if got.ExerciseType != tt.want.ExerciseType {
					t.Errorf("ClassifyExercise() exerciseType = %v, want %v", got.ExerciseType, tt.want.ExerciseType)
				}
				if got.Difficulty != tt.want.Difficulty {
					t.Errorf("ClassifyExercise() difficulty = %v, want %v", got.Difficulty, tt.want.Difficulty)
				}
			}
		})
	}
}

func TestProgramGenerator_GenerateProgram(t *testing.T) {
	pg := NewProgramGenerator()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *ProgramGenerationRequest
		wantErr bool
	}{
		{
			name: "Basic program generation",
			req: &ProgramGenerationRequest{
				Goals:              []string{"muscle_gain"},
				FitnessLevel:       "intermediate",
				AvailableEquipment: []string{"gym"},
				SessionsPerWeek:    3,
				SessionDuration:    60,
				DurationWeeks:      8,
			},
			wantErr: false,
		},
		{
			name: "Home workout program",
			req: &ProgramGenerationRequest{
				Goals:              []string{"weight_loss", "general_fitness"},
				FitnessLevel:       "beginner",
				AvailableEquipment: []string{"bodyweight", "dumbbells"},
				SessionsPerWeek:    4,
				SessionDuration:    45,
				DurationWeeks:      12,
				Preferences: &ProgramPreferences{
					SplitType:  "upper_lower",
					FocusAreas: []string{"core", "legs"},
				},
			},
			wantErr: false,
		},
		{
			name: "Advanced strength program",
			req: &ProgramGenerationRequest{
				Goals:              []string{"strength"},
				FitnessLevel:       "advanced",
				AvailableEquipment: []string{"gym", "barbell", "dumbbells"},
				SessionsPerWeek:    5,
				SessionDuration:    90,
				DurationWeeks:      16,
				Preferences: &ProgramPreferences{
					SplitType: "push_pull_legs",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pg.GenerateProgram(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Verify basic program structure
				if got.Level != tt.req.FitnessLevel {
					t.Errorf("GenerateProgram() level = %v, want %v", got.Level, tt.req.FitnessLevel)
				}

				if got.SessionsPerWeek != tt.req.SessionsPerWeek {
					t.Errorf("GenerateProgram() sessionsPerWeek = %v, want %v", got.SessionsPerWeek, tt.req.SessionsPerWeek)
				}

				if got.DurationWeeks != tt.req.DurationWeeks {
					t.Errorf("GenerateProgram() durationWeeks = %v, want %v", got.DurationWeeks, tt.req.DurationWeeks)
				}

				// Verify weekly plans were generated
				if len(got.WeeklyPlans) != tt.req.DurationWeeks {
					t.Errorf("GenerateProgram() weekly plans count = %v, want %v", len(got.WeeklyPlans), tt.req.DurationWeeks)
				}

				// Verify each week has correct number of workouts
				for week, plan := range got.WeeklyPlans {
					if plan.WeekNumber != week {
						t.Errorf("Week %d plan has incorrect week number %d", week, plan.WeekNumber)
					}

					workoutCount := 0
					for _, workout := range plan.DailyWorkouts {
						if !workout.IsRestDay {
							workoutCount++
						}
					}

					if workoutCount != tt.req.SessionsPerWeek {
						t.Errorf("Week %d has %d workouts, want %d", week, workoutCount, tt.req.SessionsPerWeek)
					}
				}
			}
		})
	}
}

func TestProgramGenerator_determineSplitType(t *testing.T) {
	pg := NewProgramGenerator()

	tests := []struct {
		name            string
		sessionsPerWeek int
		preferences     *ProgramPreferences
		want            string
	}{
		{
			name:            "2 sessions - full body",
			sessionsPerWeek: 2,
			want:            "full_body",
		},
		{
			name:            "3 sessions - full body",
			sessionsPerWeek: 3,
			want:            "full_body",
		},
		{
			name:            "4 sessions - upper lower",
			sessionsPerWeek: 4,
			want:            "upper_lower",
		},
		{
			name:            "5 sessions - push pull legs",
			sessionsPerWeek: 5,
			want:            "push_pull_legs",
		},
		{
			name:            "6 sessions - push pull legs",
			sessionsPerWeek: 6,
			want:            "push_pull_legs",
		},
		{
			name:            "7 sessions - body part split",
			sessionsPerWeek: 7,
			want:            "body_part_split",
		},
		{
			name:            "Preference override",
			sessionsPerWeek: 4,
			preferences: &ProgramPreferences{
				SplitType: "full_body",
			},
			want: "full_body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ProgramGenerationRequest{
				SessionsPerWeek: tt.sessionsPerWeek,
				Preferences:     tt.preferences,
			}

			got := pg.determineSplitType(req)
			if got != tt.want {
				t.Errorf("determineSplitType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramGenerator_classifyPurpose(t *testing.T) {
	pg := NewProgramGenerator()

	tests := []struct {
		name     string
		exercise *models.WorkoutExerciseData
		want     string
	}{
		{
			name: "Cardio exercise",
			exercise: &models.WorkoutExerciseData{
				Category: "Cardio",
			},
			want: "cardio",
		},
		{
			name: "Strength exercise",
			exercise: &models.WorkoutExerciseData{
				Category: "Strength",
			},
			want: "strength",
		},
		{
			name: "Flexibility exercise",
			exercise: &models.WorkoutExerciseData{
				Category: "Flexibility",
			},
			want: "flexibility",
		},
		{
			name: "Core exercise",
			exercise: &models.WorkoutExerciseData{
				Category: "Core",
			},
			want: "core_stability",
		},
		{
			name: "Rehabilitation exercise",
			exercise: &models.WorkoutExerciseData{
				Category: "Rehabilitation",
			},
			want: "rehabilitation",
		},
		{
			name: "Unknown category defaults to strength",
			exercise: &models.WorkoutExerciseData{
				Category: "Unknown",
			},
			want: "strength",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pg.classifyPurpose(tt.exercise)
			if got != tt.want {
				t.Errorf("classifyPurpose() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramGenerator_classifyExerciseType(t *testing.T) {
	pg := NewProgramGenerator()

	tests := []struct {
		name     string
		exercise *models.WorkoutExerciseData
		want     string
	}{
		{
			name: "Compound by name - squat",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Barbell Squat",
				Category:       "Strength",
				PrimaryMuscles: []string{"quadriceps"},
			},
			want: "compound",
		},
		{
			name: "Compound by name - deadlift",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Romanian Deadlift",
				Category:       "Strength",
				PrimaryMuscles: []string{"hamstrings"},
			},
			want: "compound",
		},
		{
			name: "Compound by muscle groups",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Pull-up",
				Category:       "Strength",
				PrimaryMuscles: []string{"lats", "biceps"},
			},
			want: "compound",
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
			name: "Plyometric exercise",
			exercise: &models.WorkoutExerciseData{
				NameEn:   "Box Jump",
				Category: "Plyometric",
			},
			want: "plyometric",
		},
		{
			name: "Isolation exercise",
			exercise: &models.WorkoutExerciseData{
				NameEn:         "Bicep Curl",
				Category:       "Strength",
				PrimaryMuscles: []string{"biceps"},
			},
			want: "isolation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pg.classifyExerciseType(tt.exercise)
			if got != tt.want {
				t.Errorf("classifyExerciseType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramGenerator_filterExercisesByEquipment(t *testing.T) {
	pg := NewProgramGenerator()

	exercises := []*models.WorkoutExerciseData{
		{
			ID:        "ex1",
			NameEn:    "Barbell Squat",
			Equipment: []string{"barbell"},
		},
		{
			ID:        "ex2",
			NameEn:    "Push-up",
			Equipment: []string{"bodyweight"},
		},
		{
			ID:        "ex3",
			NameEn:    "Dumbbell Curl",
			Equipment: []string{"dumbbells"},
		},
		{
			ID:        "ex4",
			NameEn:    "Treadmill Run",
			Equipment: []string{"treadmill"},
		},
	}

	tests := []struct {
		name               string
		availableEquipment []string
		wantCount          int
		wantIDs            []string
	}{
		{
			name:               "No equipment filter",
			availableEquipment: []string{},
			wantCount:          4,
			wantIDs:            []string{"ex1", "ex2", "ex3", "ex4"},
		},
		{
			name:               "Bodyweight only",
			availableEquipment: []string{"bodyweight"},
			wantCount:          1,
			wantIDs:            []string{"ex2"},
		},
		{
			name:               "Gym equipment",
			availableEquipment: []string{"barbell", "dumbbells", "treadmill"},
			wantCount:          4, // All exercises (bodyweight is always available)
			wantIDs:            []string{"ex1", "ex2", "ex3", "ex4"},
		},
		{
			name:               "Home equipment",
			availableEquipment: []string{"dumbbells"},
			wantCount:          2,
			wantIDs:            []string{"ex2", "ex3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pg.filterExercisesByEquipment(exercises, tt.availableEquipment)

			if len(got) != tt.wantCount {
				t.Errorf("filterExercisesByEquipment() count = %v, want %v", len(got), tt.wantCount)
			}

			// Check that expected exercises are included
			gotIDs := make(map[string]bool)
			for _, ex := range got {
				gotIDs[ex.ID] = true
			}

			for _, wantID := range tt.wantIDs {
				if !gotIDs[wantID] {
					t.Errorf("filterExercisesByEquipment() missing expected exercise %v", wantID)
				}
			}
		})
	}
}

func TestProgramGenerator_determineDayFocus(t *testing.T) {
	pg := NewProgramGenerator()

	tests := []struct {
		name      string
		dayNumber int
		splitType string
		want      string
	}{
		{
			name:      "Full body day 1",
			dayNumber: 1,
			splitType: "full_body",
			want:      "Full Body Workout 1",
		},
		{
			name:      "Upper lower - day 1 (upper)",
			dayNumber: 1,
			splitType: "upper_lower",
			want:      "Upper Body",
		},
		{
			name:      "Upper lower - day 2 (lower)",
			dayNumber: 2,
			splitType: "upper_lower",
			want:      "Lower Body",
		},
		{
			name:      "Push pull legs - day 1 (push)",
			dayNumber: 1,
			splitType: "push_pull_legs",
			want:      "Push (Chest, Shoulders, Triceps)",
		},
		{
			name:      "Push pull legs - day 2 (pull)",
			dayNumber: 2,
			splitType: "push_pull_legs",
			want:      "Pull (Back, Biceps)",
		},
		{
			name:      "Push pull legs - day 3 (legs)",
			dayNumber: 3,
			splitType: "push_pull_legs",
			want:      "Legs (Quads, Hamstrings, Glutes, Calves)",
		},
		{
			name:      "Body part split - chest",
			dayNumber: 1,
			splitType: "body_part_split",
			want:      "Chest",
		},
		{
			name:      "Body part split - back",
			dayNumber: 2,
			splitType: "body_part_split",
			want:      "Back",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pg.determineDayFocus(tt.dayNumber, tt.splitType)
			if got != tt.want {
				t.Errorf("determineDayFocus() = %v, want %v", got, tt.want)
			}
		})
	}
}
