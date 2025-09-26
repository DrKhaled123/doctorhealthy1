package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api-key-generator/internal/models"
)

// WarmupService handles warmup routine generation and technique guidance
type WarmupService struct{}

// NewWarmupService creates a new warmup service
func NewWarmupService() *WarmupService {
	return &WarmupService{}
}

// GenerateWarmupRoutine creates a customized warmup routine based on workout type
func (ws *WarmupService) GenerateWarmupRoutine(ctx context.Context, req *WarmupRequest) (*models.WarmupRoutine, error) {
	if req == nil {
		return nil, fmt.Errorf("warmup request is required")
	}

	workoutType := strings.ToLower(req.WorkoutType)

	// Get base warmup template
	template := ws.getWarmupTemplate(workoutType)

	// Customize based on duration and intensity
	routine := &models.WarmupRoutine{
		ID:          generateID(),
		WorkoutType: workoutType,
		Duration:    ws.calculateWarmupDuration(req.WorkoutDuration, req.Intensity),
		Exercises:   ws.generateWarmupExercises(workoutType, req.TargetMuscles, req.Intensity),
	}

	// Set instructions based on language
	if req.Language == "ar" {
		routine.InstructionsAr = template.InstructionsAr
	} else {
		routine.InstructionsEn = template.InstructionsEn
	}

	routine.Tips = ws.getWarmupTips(workoutType, req.Language)
	routine.SafetyNotes = ws.getSafetyNotes(workoutType, req.Language)
	routine.CreatedAt = time.Now()

	return routine, nil
}

// GetExerciseTechnique provides comprehensive technique guidance for an exercise
func (ws *WarmupService) GetExerciseTechnique(ctx context.Context, exerciseID string, language string) (*ExerciseTechnique, error) {
	if exerciseID == "" {
		return nil, fmt.Errorf("exercise ID is required")
	}

	// Get technique data (in production, this would query the database)
	technique := ws.getTechniqueData(exerciseID, language)
	if technique == nil {
		return nil, fmt.Errorf("technique data not found for exercise: %s", exerciseID)
	}

	return technique, nil
}

// GenerateCooldownRoutine creates a post-workout cooldown routine
func (ws *WarmupService) GenerateCooldownRoutine(ctx context.Context, req *CooldownRequest) (*models.CooldownRoutine, error) {
	if req == nil {
		return nil, fmt.Errorf("cooldown request is required")
	}

	routine := &models.CooldownRoutine{
		Duration:     ws.calculateCooldownDuration(req.WorkoutDuration, req.Intensity),
		Stretches:    ws.generateStretches(req.TargetMuscles, req.Language),
		RecoveryTips: ws.getRecoveryTips(req.WorkoutType, req.Language),
	}

	return routine, nil
}

// Helper methods

func (ws *WarmupService) calculateWarmupDuration(workoutDuration int, intensity string) int {
	baseDuration := 10 // 10 minutes base

	// Adjust based on workout duration
	if workoutDuration > 90 {
		baseDuration = 15
	} else if workoutDuration < 45 {
		baseDuration = 8
	}

	// Adjust based on intensity
	switch strings.ToLower(intensity) {
	case "high":
		baseDuration += 3
	case "low":
		baseDuration -= 2
	}

	if baseDuration < 5 {
		baseDuration = 5
	}

	return baseDuration
}

func (ws *WarmupService) generateWarmupExercises(workoutType string, targetMuscles []string, intensity string) []*models.WarmupExercise {
	exercises := []*models.WarmupExercise{}

	// General mobility (always included)
	exercises = append(exercises, &models.WarmupExercise{
		Name:         "Arm Circles",
		Duration:     "30 seconds each direction",
		Instructions: []string{"Stand with arms extended", "Make small circles, gradually increasing size", "Reverse direction"},
		Purpose:      "Shoulder mobility",
	})

	// Workout-specific exercises
	switch workoutType {
	case "upper_body", "push", "pull":
		exercises = append(exercises, ws.getUpperBodyWarmup()...)
	case "lower_body", "legs":
		exercises = append(exercises, ws.getLowerBodyWarmup()...)
	case "full_body":
		exercises = append(exercises, ws.getFullBodyWarmup()...)
	case "cardio":
		exercises = append(exercises, ws.getCardioWarmup()...)
	}

	// Add dynamic stretches for target muscles
	for _, muscle := range targetMuscles {
		if stretch := ws.getDynamicStretch(muscle); stretch != nil {
			exercises = append(exercises, stretch)
		}
	}

	return exercises
}

func (ws *WarmupService) getUpperBodyWarmup() []*models.WarmupExercise {
	return []*models.WarmupExercise{
		{
			Name:         "Shoulder Rolls",
			Duration:     "10 reps each direction",
			Instructions: []string{"Roll shoulders backward in large circles", "Focus on full range of motion"},
			Purpose:      "Shoulder preparation",
		},
		{
			Name:         "Band Pull-Aparts",
			Duration:     "15 reps",
			Instructions: []string{"Hold resistance band at chest level", "Pull apart squeezing shoulder blades"},
			Purpose:      "Rear delt activation",
		},
	}
}

func (ws *WarmupService) getLowerBodyWarmup() []*models.WarmupExercise {
	return []*models.WarmupExercise{
		{
			Name:         "Leg Swings",
			Duration:     "10 each leg, each direction",
			Instructions: []string{"Hold wall for support", "Swing leg forward/back, then side to side"},
			Purpose:      "Hip mobility",
		},
		{
			Name:         "Bodyweight Squats",
			Duration:     "15 reps",
			Instructions: []string{"Feet shoulder-width apart", "Lower slowly, chest up", "Full range of motion"},
			Purpose:      "Movement pattern preparation",
		},
	}
}

func (ws *WarmupService) getFullBodyWarmup() []*models.WarmupExercise {
	return []*models.WarmupExercise{
		{
			Name:         "Jumping Jacks",
			Duration:     "30 seconds",
			Instructions: []string{"Start with feet together", "Jump while spreading legs and raising arms"},
			Purpose:      "Full body activation",
		},
		{
			Name:         "Inchworms",
			Duration:     "8 reps",
			Instructions: []string{"Bend forward, walk hands out to plank", "Walk hands back, stand up"},
			Purpose:      "Dynamic full body stretch",
		},
	}
}

func (ws *WarmupService) getCardioWarmup() []*models.WarmupExercise {
	return []*models.WarmupExercise{
		{
			Name:         "Light Jogging",
			Duration:     "3 minutes",
			Instructions: []string{"Start very slow", "Gradually increase pace", "Focus on breathing"},
			Purpose:      "Cardiovascular preparation",
		},
		{
			Name:         "High Knees",
			Duration:     "30 seconds",
			Instructions: []string{"Lift knees to hip level", "Stay on balls of feet", "Pump arms naturally"},
			Purpose:      "Dynamic leg activation",
		},
	}
}

func (ws *WarmupService) getDynamicStretch(muscle string) *models.WarmupExercise {
	stretches := map[string]*models.WarmupExercise{
		"chest": {
			Name:         "Arm Swings",
			Duration:     "10 reps each arm",
			Instructions: []string{"Swing arm across body", "Then swing back and open chest"},
			Purpose:      "Chest and shoulder mobility",
		},
		"hamstrings": {
			Name:         "Leg Kicks",
			Duration:     "10 reps each leg",
			Instructions: []string{"Kick leg up toward opposite hand", "Keep leg straight", "Control the movement"},
			Purpose:      "Hamstring dynamic stretch",
		},
	}

	return stretches[strings.ToLower(muscle)]
}

func (ws *WarmupService) getWarmupTemplate(workoutType string) *WarmupTemplate {
	templates := map[string]*WarmupTemplate{
		"upper_body": {
			InstructionsEn: "Focus on shoulder mobility and upper body activation. Start slowly and gradually increase intensity.",
			InstructionsAr: "ركز على حركة الكتف وتنشيط الجزء العلوي من الجسم. ابدأ ببطء وزد الشدة تدريجياً.",
		},
		"lower_body": {
			InstructionsEn: "Emphasize hip mobility and leg activation. Prepare joints for loaded movements.",
			InstructionsAr: "أكد على حركة الورك وتنشيط الساقين. حضر المفاصل للحركات المحملة.",
		},
	}

	if template, exists := templates[workoutType]; exists {
		return template
	}

	return &WarmupTemplate{
		InstructionsEn: "Complete all exercises with controlled movements. Focus on quality over speed.",
		InstructionsAr: "أكمل جميع التمارين بحركات محكومة. ركز على الجودة وليس السرعة.",
	}
}

func (ws *WarmupService) getWarmupTips(workoutType, language string) []string {
	if language == "ar" {
		return []string{
			"ابدأ ببطء وزد الشدة تدريجياً",
			"ركز على التنفس العميق",
			"توقف إذا شعرت بألم",
			"اشرب الماء حسب الحاجة",
		}
	}

	return []string{
		"Start slowly and gradually increase intensity",
		"Focus on deep breathing throughout",
		"Stop if you feel any pain",
		"Stay hydrated as needed",
	}
}

func (ws *WarmupService) getSafetyNotes(workoutType, language string) []string {
	if language == "ar" {
		return []string{
			"لا تقم بحركات قوية أو مفاجئة",
			"استمع لجسمك وتوقف عند الحاجة",
			"تأكد من وجود مساحة كافية للحركة",
		}
	}

	return []string{
		"Avoid sudden or forceful movements",
		"Listen to your body and stop when needed",
		"Ensure adequate space for movement",
	}
}

func (ws *WarmupService) calculateCooldownDuration(workoutDuration int, intensity string) int {
	baseDuration := 10

	if workoutDuration > 90 {
		baseDuration = 15
	}

	switch strings.ToLower(intensity) {
	case "high":
		baseDuration += 5
	}

	return baseDuration
}

func (ws *WarmupService) generateStretches(targetMuscles []string, language string) []*models.StretchExercise {
	stretches := []*models.StretchExercise{}

	// Always include general stretches
	stretches = append(stretches, &models.StretchExercise{
		Name:         "Child's Pose",
		Duration:     "30 seconds",
		Instructions: []string{"Kneel and sit back on heels", "Extend arms forward", "Breathe deeply"},
		TargetMuscle: "back",
	})

	// Add muscle-specific stretches
	for _, muscle := range targetMuscles {
		if stretch := ws.getStaticStretch(muscle, language); stretch != nil {
			stretches = append(stretches, stretch)
		}
	}

	return stretches
}

func (ws *WarmupService) getStaticStretch(muscle, language string) *models.StretchExercise {
	stretches := map[string]*models.StretchExercise{
		"chest": {
			Name:         "Doorway Chest Stretch",
			Duration:     "30 seconds each arm",
			Instructions: []string{"Place forearm on doorframe", "Step forward gently", "Feel stretch across chest"},
			TargetMuscle: "chest",
		},
		"hamstrings": {
			Name:         "Seated Hamstring Stretch",
			Duration:     "30 seconds each leg",
			Instructions: []string{"Sit with one leg extended", "Reach toward toes", "Keep back straight"},
			TargetMuscle: "hamstrings",
		},
	}

	return stretches[strings.ToLower(muscle)]
}

func (ws *WarmupService) getRecoveryTips(workoutType, language string) []string {
	if language == "ar" {
		return []string{
			"اشرب الماء لتعويض السوائل المفقودة",
			"تناول وجبة تحتوي على البروتين والكربوهيدرات",
			"احصل على نوم كافٍ للتعافي",
		}
	}

	return []string{
		"Hydrate to replace lost fluids",
		"Consume protein and carbs within 30 minutes",
		"Get adequate sleep for recovery",
	}
}

func (ws *WarmupService) getTechniqueData(exerciseID, language string) *ExerciseTechnique {
	// Sample technique data (in production, this would come from database)
	techniques := map[string]*ExerciseTechnique{
		"bench_press": {
			ExerciseID:   "bench_press",
			ExerciseName: "Bench Press",
			ProperForm: []string{
				"Lie flat on bench with feet on floor",
				"Grip bar slightly wider than shoulders",
				"Lower bar to chest with control",
				"Press up explosively",
			},
			CommonMistakes: []string{
				"Bouncing bar off chest",
				"Lifting feet off ground",
				"Arching back excessively",
			},
			SafetyTips: []string{
				"Always use a spotter for heavy weights",
				"Keep wrists straight and strong",
				"Don't lock elbows aggressively",
			},
			Modifications: []string{
				"Use dumbbells for better range of motion",
				"Incline bench for upper chest focus",
				"Reduce weight if form breaks down",
			},
		},
	}

	return techniques[exerciseID]
}

// Request/Response types

type WarmupRequest struct {
	WorkoutType     string   `json:"workout_type"`
	WorkoutDuration int      `json:"workout_duration_minutes"`
	Intensity       string   `json:"intensity"`
	TargetMuscles   []string `json:"target_muscles"`
	Language        string   `json:"language"`
}

type CooldownRequest struct {
	WorkoutType     string   `json:"workout_type"`
	WorkoutDuration int      `json:"workout_duration_minutes"`
	Intensity       string   `json:"intensity"`
	TargetMuscles   []string `json:"target_muscles"`
	Language        string   `json:"language"`
}

type ExerciseTechnique struct {
	ExerciseID     string   `json:"exercise_id"`
	ExerciseName   string   `json:"exercise_name"`
	ProperForm     []string `json:"proper_form"`
	CommonMistakes []string `json:"common_mistakes"`
	SafetyTips     []string `json:"safety_tips"`
	Modifications  []string `json:"modifications"`
}

type WarmupTemplate struct {
	InstructionsEn string
	InstructionsAr string
}
