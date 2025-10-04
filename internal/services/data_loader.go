package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DataLoader handles loading JSON data files
type DataLoader struct {
	dataPath string
}

// NewDataLoader creates a new data loader
func NewDataLoader(dataPath string) *DataLoader {
	return &DataLoader{
		dataPath: dataPath,
	}
}

// Disease represents disease data structure
type Disease struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Symptoms    []string `json:"symptoms"`
	RiskFactors []string `json:"risk_factors"`
	Treatment   struct {
		Lifestyle              []string `json:"lifestyle"`
		DietaryRecommendations []string `json:"dietary_recommendations"`
		FoodsToAvoid           []string `json:"foods_to_avoid"`
		RecommendedFoods       []string `json:"recommended_foods"`
		Supplements            []string `json:"supplements"`
	} `json:"treatment"`
	Complications []string `json:"complications"`
	Monitoring    []string `json:"monitoring"`
}

// TypePlan represents nutrition plan type data
type TypePlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Macros      struct {
		Carbs   int `json:"carbs"`
		Protein int `json:"protein"`
		Fat     int `json:"fat"`
	} `json:"macros"`
	Benefits     []string `json:"benefits"`
	SuitableFor  []string `json:"suitable_for"`
	FoodsAllowed []string `json:"foods_allowed"`
	FoodsToAvoid []string `json:"foods_to_avoid"`
	MealTiming   string   `json:"meal_timing"`
	Duration     string   `json:"duration"`
	Precautions  []string `json:"precautions"`
}

// Complaint represents health complaint data
type Complaint struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	Description    string   `json:"description"`
	PossibleCauses []string `json:"possible_causes"`
	Symptoms       []string `json:"symptoms"`
	Solutions      struct {
		Lifestyle   []string `json:"lifestyle"`
		Dietary     []string `json:"dietary"`
		Supplements []struct {
			Name     string `json:"name"`
			Dosage   string `json:"dosage"`
			Timing   string `json:"timing"`
			Benefits string `json:"benefits"`
		} `json:"supplements"`
	} `json:"solutions"`
	WhenToSeeDoctor []string `json:"when_to_see_doctor"`
}

// RecipeData represents recipe data structure
type RecipeData struct {
	Mediterranean []Recipe `json:"mediterranean"`
	MiddleEastern []Recipe `json:"middle_eastern"`
	Asian         []Recipe `json:"asian"`
	Indian        []Recipe `json:"indian"`
}

// Recipe represents individual recipe
type Recipe struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	MealType           string `json:"meal_type"`
	Difficulty         string `json:"difficulty"`
	PrepTime           int    `json:"prep_time"`
	CookTime           int    `json:"cook_time"`
	Servings           int    `json:"servings"`
	CaloriesPerServing int    `json:"calories_per_serving"`
	Ingredients        []struct {
		Name     string  `json:"name"`
		Amount   float64 `json:"amount"`
		Unit     string  `json:"unit"`
		Calories int     `json:"calories"`
	} `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Nutrition    struct {
		Protein int `json:"protein"`
		Carbs   int `json:"carbs"`
		Fat     int `json:"fat"`
		Fiber   int `json:"fiber"`
	} `json:"nutrition"`
	Tags      []string `json:"tags"`
	Allergens []string `json:"allergens"`
	Tips      []string `json:"tips"`
}

// VitaminMineral represents supplement data
type VitaminMineral struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Category           string            `json:"category"`
	Description        string            `json:"description"`
	Benefits           []string          `json:"benefits"`
	DeficiencySymptoms []string          `json:"deficiency_symptoms"`
	FoodSources        []string          `json:"food_sources"`
	RecommendedDosage  map[string]string `json:"recommended_dosage"`
	Timing             string            `json:"timing"`
	Interactions       []string          `json:"interactions"`
	Precautions        []string          `json:"precautions"`
	SuitableFor        []string          `json:"suitable_for"`
}

// WorkoutData represents workout data structure
type WorkoutData struct {
	Gym   []Exercise `json:"gym"`
	Home  []Exercise `json:"home"`
	Goals map[string]struct {
		Description       string   `json:"description"`
		RepRange          string   `json:"rep_range"`
		Sets              string   `json:"sets"`
		Rest              string   `json:"rest"`
		Frequency         string   `json:"frequency"`
		ExerciseSelection []string `json:"exercise_selection"`
		NutritionFocus    []string `json:"nutrition_focus"`
	} `json:"goals"`
}

// Exercise represents individual exercise
type Exercise struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Category         string   `json:"category"`
	PrimaryMuscles   []string `json:"primary_muscles"`
	SecondaryMuscles []string `json:"secondary_muscles"`
	Equipment        []string `json:"equipment"`
	Difficulty       string   `json:"difficulty"`
	Instructions     []string `json:"instructions"`
	CommonMistakes   []string `json:"common_mistakes"`
	Tips             []string `json:"tips"`
	SetsReps         map[string]struct {
		Sets int    `json:"sets"`
		Reps string `json:"reps"`
		Rest string `json:"rest"`
	} `json:"sets_reps"`
	Progressions []string `json:"progressions"`
	Alternatives []string `json:"alternatives"`
}

// Injury represents injury data
type Injury struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	CommonCauses     []string `json:"common_causes"`
	Symptoms         []string `json:"symptoms"`
	ExercisesToAvoid []struct {
		Exercise string `json:"exercise"`
		Reason   string `json:"reason"`
	} `json:"exercises_to_avoid"`
	RecommendedExercises []struct {
		Exercise     string   `json:"exercise"`
		Sets         int      `json:"sets"`
		Reps         string   `json:"reps"`
		Description  string   `json:"description"`
		Instructions []string `json:"instructions"`
	} `json:"recommended_exercises"`
	TreatmentAdvice  []string `json:"treatment_advice"`
	WhenToSeeDoctor  []string `json:"when_to_see_doctor"`
	RecoveryTimeline string   `json:"recovery_timeline"`
	Prevention       []string `json:"prevention"`
}

// LoadDiseases loads disease data from comprehensive JSON file
func (dl *DataLoader) LoadDiseases() ([]Disease, error) {
	filePath := filepath.Join(dl.dataPath, "comprehensive-diseases-database-EXPANDED.js")
	var result []Disease
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadTypePlans loads nutrition plan types from JSON file
func (dl *DataLoader) LoadTypePlans() ([]TypePlan, error) {
	filePath := filepath.Join(dl.dataPath, "type-plans.js")
	var result []TypePlan
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadComplaints loads health complaints from comprehensive JSON file
func (dl *DataLoader) LoadComplaints() ([]Complaint, error) {
	filePath := filepath.Join(dl.dataPath, "comprehensive-complaints-database-COMPLETE.js")
	var result []Complaint
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadRecipes loads recipe data from JSON file
func (dl *DataLoader) LoadRecipes() (RecipeData, error) {
	filePath := filepath.Join(dl.dataPath, "recipes.js")
	var result RecipeData
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadVitaminsAndMinerals loads supplement data from VIP drugs-nutrition file
func (dl *DataLoader) LoadVitaminsAndMinerals() ([]VitaminMineral, error) {
	filePath := filepath.Join(dl.dataPath, "vip-drugs-nutrition.js")
	var result []VitaminMineral
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadWorkouts loads workout data from VIP workouts file
func (dl *DataLoader) LoadWorkouts() (WorkoutData, error) {
	filePath := filepath.Join(dl.dataPath, "vip-workouts.js")
	var result WorkoutData
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// LoadInjuries loads injury data from VIP injuries file
func (dl *DataLoader) LoadInjuries() ([]Injury, error) {
	filePath := filepath.Join(dl.dataPath, "vip-injuries.js")
	var result []Injury
	err := dl.loadJSONFile(filePath, &result)
	return result, err
}

// loadJSONFile is a generic function to load and parse JSON files
func (dl *DataLoader) loadJSONFile(filePath string, result interface{}) error {
	// Validate file path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(dl.dataPath, cleanPath)
	}

	// Ensure the file is within the allowed data directory
	if !isWithinDirectory(cleanPath, dl.dataPath) {
		return fmt.Errorf("file path outside allowed directory: %s", filePath)
	}

	// Read the file
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", cleanPath, err)
	}

	// Convert JavaScript module.exports to JSON
	// Remove "const variableName = " and "module.exports = variableName;"
	jsonData := dl.convertJSToJSON(string(data))

	// Parse JSON
	err = json.Unmarshal([]byte(jsonData), result)
	if err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filePath, err)
	}

	return nil
}

// convertJSToJSON converts JavaScript module format to JSON
func (dl *DataLoader) convertJSToJSON(jsContent string) string {
	// This is a simple conversion - in production you might want a more robust parser
	// Remove "const variableName = " from the beginning
	start := 0
	if idx := findAssignmentStart(jsContent); idx != -1 {
		start = idx
	}

	// Remove "module.exports = variableName;" from the end
	end := len(jsContent)
	if idx := findModuleExports(jsContent); idx != -1 {
		end = idx
	}

	return jsContent[start:end]
}

// Helper functions for JS to JSON conversion
func findAssignmentStart(content string) int {
	// Find the start of the actual data (after "const name = ")
	for i := 0; i < len(content); i++ {
		if content[i] == '[' || content[i] == '{' {
			return i
		}
	}
	return -1
}

func findModuleExports(content string) int {
	// Find "module.exports" and return the position before it
	moduleExports := "module.exports"
	if idx := findLastOccurrence(content, moduleExports); idx != -1 {
		// Go backwards to find the end of the data structure
		for i := idx - 1; i >= 0; i-- {
			if content[i] == ']' || content[i] == '}' {
				return i + 1
			}
		}
	}
	return -1
}

func findLastOccurrence(content, substr string) int {
	lastIdx := -1
	for i := 0; i <= len(content)-len(substr); i++ {
		if content[i:i+len(substr)] == substr {
			lastIdx = i
		}
	}
	return lastIdx
}

// isWithinDirectory checks if a file path is within the allowed directory
func isWithinDirectory(filePath, allowedDir string) bool {
	rel, err := filepath.Rel(allowedDir, filePath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, "/")
}
