package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// VIPWorkoutExtractor extracts and processes VIP workout data
type VIPWorkoutExtractor struct {
	nutritionService *NutritionTimingService
}

// NewVIPWorkoutExtractor creates a new VIP workout extractor
func NewVIPWorkoutExtractor() *VIPWorkoutExtractor {
	return &VIPWorkoutExtractor{
		nutritionService: NewNutritionTimingService(),
	}
}

// VIPWorkoutData represents the structure from workouts.js
type VIPWorkoutData struct {
	WorkoutPlans []struct {
		Level       string `json:"level"`
		Split       string `json:"split"`
		Description struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"description"`
		Weeks []struct {
			WeekNumber  int    `json:"week_number"`
			Progression string `json:"progression"`
			Days        []struct {
				Day   int `json:"day"`
				Focus struct {
					En string `json:"en"`
					Ar string `json:"ar"`
				} `json:"focus"`
				Exercises []struct {
					Name struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"name"`
					Sets         int    `json:"sets"`
					Reps         string `json:"reps"`
					Rest         string `json:"rest"`
					Instructions struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"instructions"`
					CommonMistakes struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"common_mistakes"`
					EvidenceLink struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"evidence_link"`
				} `json:"exercises"`
			} `json:"days"`
		} `json:"weeks"`
	} `json:"workout_plans"`
}

// ExtractVIPWorkouts processes VIP workout data and creates comprehensive programs
func (vwe *VIPWorkoutExtractor) ExtractVIPWorkouts(ctx context.Context, dataPath string) ([]*ComprehensiveWorkoutProgram, error) {
	// Read VIP workouts.js file
	filePath := filepath.Join(dataPath, "vip json", " workouts.js") // nosec G304 - controlled path within application data directory
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workouts file: %w", err)
	}

	var vipData VIPWorkoutData
	if err := json.Unmarshal(data, &vipData); err != nil {
		return nil, fmt.Errorf("failed to parse workout data: %w", err)
	}

	programs := make([]*ComprehensiveWorkoutProgram, 0)

	// Process each workout plan
	for _, plan := range vipData.WorkoutPlans {
		program := vwe.createComprehensiveProgram(plan)
		programs = append(programs, program)
	}

	// Add additional specialized programs
	programs = append(programs, vwe.generateSpecializedPrograms()...)

	return programs, nil
}

func (vwe *VIPWorkoutExtractor) createComprehensiveProgram(plan interface{}) *ComprehensiveWorkoutProgram {
	// Type assertion to access plan data
	planData, ok := plan.(struct {
		Level       string `json:"level"`
		Split       string `json:"split"`
		Description struct {
			En string `json:"en"`
			Ar string `json:"ar"`
		} `json:"description"`
		Weeks []struct {
			WeekNumber  int    `json:"week_number"`
			Progression string `json:"progression"`
			Days        []struct {
				Day   int `json:"day"`
				Focus struct {
					En string `json:"en"`
					Ar string `json:"ar"`
				} `json:"focus"`
				Exercises []struct {
					Name struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"name"`
					Sets         int    `json:"sets"`
					Reps         string `json:"reps"`
					Rest         string `json:"rest"`
					Instructions struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"instructions"`
					CommonMistakes struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"common_mistakes"`
					EvidenceLink struct {
						En string `json:"en"`
						Ar string `json:"ar"`
					} `json:"evidence_link"`
				} `json:"exercises"`
			} `json:"days"`
		} `json:"weeks"`
	})

	if !ok {
		// Return default program if type assertion fails
		return vwe.getDefaultProgram()
	}

	program := &ComprehensiveWorkoutProgram{
		ID:          generateID(),
		NameEnglish: fmt.Sprintf("%s %s Program", cases.Title(language.Und).String(planData.Level), cases.Title(language.Und).String(planData.Split)),
		NameArabic:  fmt.Sprintf("برنامج %s %s", planData.Level, planData.Split),
		Purpose:     planData.Description.En,
		Level:       planData.Level,
		SplitType:   planData.Split,
		CreatedAt:   time.Now(),
	}

	// Generate comprehensive details
	program.WorkoutDescription = vwe.generateWorkoutDescription(planData)
	program.TechniqueGuidance = vwe.generateTechniqueGuidance(planData)
	program.SupplementProtocol = vwe.generateSupplementProtocol(planData.Level)
	program.NutritionTiming = vwe.generateNutritionTiming(planData.Level)
	program.WarmupRecovery = vwe.generateWarmupRecovery(planData.Split)
	program.NutritionGuidelines = vwe.generateNutritionGuidelines(planData.Level)
	program.ReferenceLinks = vwe.generateReferenceLinks(planData)

	return program
}

func (vwe *VIPWorkoutExtractor) generateWorkoutDescription(planData interface{}) string {
	// Generate based on level and split
	return "Comprehensive workout program with progressive overload, proper form emphasis, and evidence-based exercise selection. 3-5 sessions per week with structured progression."
}

func (vwe *VIPWorkoutExtractor) generateTechniqueGuidance(planData interface{}) string {
	return "Prioritize form over weight. Control eccentric phase (2-3 sec). Maintain neutral spine. Engage core throughout. Use full range of motion. Progress gradually to prevent injury."
}

func (vwe *VIPWorkoutExtractor) generateSupplementProtocol(level string) map[string]string {
	protocols := map[string]map[string]string{
		"beginner": {
			"protein":   "20-25g post-workout",
			"creatine":  "3g/day",
			"vitamin_d": "1000 IU/day",
			"omega_3":   "1g/day",
		},
		"intermediate": {
			"protein":      "25-30g post-workout",
			"creatine":     "5g/day",
			"beta_alanine": "3g/day",
			"vitamin_d":    "2000 IU/day",
			"magnesium":    "400mg/day",
		},
		"advanced": {
			"protein":      "30-40g post-workout",
			"creatine":     "5-10g/day",
			"beta_alanine": "4-6g/day",
			"hmb":          "3g/day",
			"zma":          "30mg zinc + 450mg magnesium",
		},
	}

	if protocol, exists := protocols[level]; exists {
		return protocol
	}
	return protocols["intermediate"]
}

func (vwe *VIPWorkoutExtractor) generateNutritionTiming(level string) *WorkoutNutritionTiming {
	timings := map[string]*WorkoutNutritionTiming{
		"beginner": {
			PreWorkout:  "Banana + almond butter (25g carbs + 8g fat) 45 min before",
			PostWorkout: "Whey protein (25g) + apple (20g carbs) within 30 min",
		},
		"intermediate": {
			PreWorkout:  "Oats (40g carbs) + whey protein (25g) 60 min before",
			PostWorkout: "Chicken breast (30g protein) + rice (40g carbs) within 45 min",
		},
		"advanced": {
			PreWorkout:  "White rice (50g carbs) + whey protein (30g) 90 min before",
			PostWorkout: "Lean beef (40g protein) + potatoes (60g carbs) within 30 min",
		},
	}

	if timing, exists := timings[level]; exists {
		return timing
	}
	return timings["intermediate"]
}

func (vwe *VIPWorkoutExtractor) generateWarmupRecovery(splitType string) string {
	warmups := map[string]string{
		"full_body":      "10 min dynamic warm-up: arm circles, leg swings, torso twists. Post-workout: 15 min full-body stretching + foam rolling.",
		"upper_lower":    "8 min specific warm-up: band pull-aparts, shoulder circles, hip circles. Post-workout: 10 min targeted stretching + cold therapy.",
		"push_pull_legs": "12 min activation warm-up: movement-specific patterns. Post-workout: 20 min recovery protocol + contrast showers.",
	}

	if warmup, exists := warmups[splitType]; exists {
		return warmup
	}
	return warmups["full_body"]
}

func (vwe *VIPWorkoutExtractor) generateNutritionGuidelines(level string) string {
	guidelines := map[string]string{
		"beginner":     "Balanced diet with slight calorie adjustment. 1.6g protein/kg bodyweight. Focus on whole foods. 3-4 meals/day.",
		"intermediate": "Structured nutrition with macro tracking. 1.8-2.2g protein/kg. Nutrient timing around workouts. 4-5 meals/day.",
		"advanced":     "Precision nutrition with periodization. 2.2-2.5g protein/kg. Advanced nutrient timing. 5-6 meals/day including pre-bed.",
	}

	if guideline, exists := guidelines[level]; exists {
		return guideline
	}
	return guidelines["intermediate"]
}

func (vwe *VIPWorkoutExtractor) generateReferenceLinks(planData interface{}) []string {
	return []string{
		"https://www.ncbi.nlm.nih.gov/pmc/articles/PMC6303127/",
		"https://journals.lww.com/nsca-jscr/",
		"https://pubmed.ncbi.nlm.nih.gov/",
		"https://www.acsm.org/",
	}
}

func (vwe *VIPWorkoutExtractor) generateSpecializedPrograms() []*ComprehensiveWorkoutProgram {
	return []*ComprehensiveWorkoutProgram{
		{
			ID:                  "hiit_fat_loss",
			NameEnglish:         "HIIT Fat Loss Accelerator",
			NameArabic:          "مسرع حرق الدهون HIIT",
			Purpose:             "Maximum fat loss in minimum time through high-intensity intervals",
			WorkoutDescription:  "20-25 min HIIT sessions: 30 sec all-out + 90 sec recovery. 3x/week.",
			TechniqueGuidance:   "Maintain form during high intensity. Land softly. Keep core engaged.",
			SupplementProtocol:  map[string]string{"caffeine": "200mg pre", "l_carnitine": "2g/day"},
			NutritionTiming:     &WorkoutNutritionTiming{PreWorkout: "Coffee + rice cake 30min before", PostWorkout: "Protein shake immediately after"},
			WarmupRecovery:      "5 min dynamic warm-up. Post: contrast showers + stretching.",
			NutritionGuidelines: "Calorie deficit 500 kcal/day. High protein. Time carbs around workouts.",
			ReferenceLinks:      []string{"https://journals.lww.com/acsm-msse/"},
			CreatedAt:           time.Now(),
		},
	}
}

func (vwe *VIPWorkoutExtractor) getDefaultProgram() *ComprehensiveWorkoutProgram {
	return &ComprehensiveWorkoutProgram{
		ID:                  "default_program",
		NameEnglish:         "General Fitness Program",
		NameArabic:          "برنامج اللياقة العامة",
		Purpose:             "Balanced fitness development",
		WorkoutDescription:  "3x/week full-body workouts",
		TechniqueGuidance:   "Focus on proper form",
		SupplementProtocol:  map[string]string{"protein": "25g post-workout"},
		NutritionTiming:     &WorkoutNutritionTiming{PreWorkout: "Light snack", PostWorkout: "Protein meal"},
		WarmupRecovery:      "Standard warm-up and cool-down",
		NutritionGuidelines: "Balanced nutrition",
		ReferenceLinks:      []string{"https://www.acsm.org/"},
		CreatedAt:           time.Now(),
	}
}

// ComprehensiveWorkoutProgram represents a complete workout program
type ComprehensiveWorkoutProgram struct {
	ID                  string                  `json:"id"`
	NameEnglish         string                  `json:"name_english"`
	NameArabic          string                  `json:"name_arabic"`
	Purpose             string                  `json:"purpose"`
	Level               string                  `json:"level"`
	SplitType           string                  `json:"split_type"`
	WorkoutDescription  string                  `json:"workout_description"`
	TechniqueGuidance   string                  `json:"technique_guidance"`
	SupplementProtocol  map[string]string       `json:"supplement_protocol"`
	NutritionTiming     *WorkoutNutritionTiming `json:"nutrition_timing"`
	WarmupRecovery      string                  `json:"warmup_recovery"`
	NutritionGuidelines string                  `json:"nutrition_guidelines"`
	ReferenceLinks      []string                `json:"reference_links"`
	CreatedAt           time.Time               `json:"created_at"`
}
