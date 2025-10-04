package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ComprehensiveDataLoader handles loading comprehensive VIP JSON data files
type ComprehensiveDataLoader struct {
	dataPath string
}

// NewComprehensiveDataLoader creates a new comprehensive data loader
func NewComprehensiveDataLoader(dataPath string) *ComprehensiveDataLoader {
	return &ComprehensiveDataLoader{
		dataPath: dataPath,
	}
}

// ComprehensiveHealthDatabase represents the complete health database structure
type ComprehensiveHealthDatabase struct {
	Medications struct {
		WeightLossDrugs      []Medication `json:"weight_loss_drugs"`
		DiabetesMedications  []Medication `json:"diabetes_medications"`
		AppetiteSuppressants []Medication `json:"appetite_suppressants"`
		InvestigationalDrugs []Medication `json:"investigational_drugs"`
		HormonalMedications  []Medication `json:"hormonal_medications"`
		PerformanceEnhancers []Medication `json:"performance_enhancers"`
	} `json:"medications"`

	Diseases struct {
		Type2Diabetes struct {
			BasicInfo                      DiseaseBasicInfo                 `json:"basic_info"`
			CausesAndDifferentialDiagnosis []CauseWithDifferentialDiagnosis `json:"causes_and_differential_diagnosis"`
		} `json:"type_2_diabetes"`
		CardiovascularDiseases map[string]CardiovascularDisease `json:"cardiovascular_diseases"`
		EndocrineDiseases      map[string]EndocrineDisease      `json:"endocrine_diseases"`
		MetabolicDiseases      map[string]MetabolicDisease      `json:"metabolic_diseases"`
	} `json:"diseases"`

	HalalHaramFoods struct {
		HaramIngredients    []HaramIngredient    `json:"haram_ingredients"`
		MashboohIngredients []MashboohIngredient `json:"mashbooh_ingredients"`
	} `json:"halal_haram_foods"`

	Vitamins struct {
		WaterSoluble []Vitamin `json:"water_soluble"`
		FatSoluble   []Vitamin `json:"fat_soluble"`
	} `json:"vitamins"`

	Minerals []Mineral `json:"minerals"`

	WorkoutTechniques   []WorkoutTechnique  `json:"workout_techniques"`
	WeightLossDrugs     []WeightLossDrug    `json:"weight_loss_drugs"`
	CookingSkills       []CookingSkill      `json:"cooking_skills"`
	SportsNutrition     SportsNutrition     `json:"sports_nutrition"`
	InjuryManagement    InjuryManagement    `json:"injury_management"`
	DietPlans           DietPlans           `json:"diet_plans"`
	NutritionGuidelines NutritionGuidelines `json:"nutrition_guidelines"`
}

// ComprehensiveComplaintsDatabase represents the complete complaints database
type ComprehensiveComplaintsDatabase struct {
	Cases                     []ComplaintCase `json:"cases"`
	CommonSupplementProtocols struct {
		WeightManagement  []SupplementProtocol `json:"weight_management"`
		EnergyEnhancement []SupplementProtocol `json:"energy_enhancement"`
		MetabolicSupport  []SupplementProtocol `json:"metabolic_support"`
	} `json:"common_supplement_protocols"`
}

// ComprehensiveDiseasesDatabase represents the complete diseases database
type ComprehensiveDiseasesDatabase struct {
	Diseases map[string]ComprehensiveDisease `json:"diseases"`
}

// ComprehensiveDisease represents a comprehensive disease entry
type ComprehensiveDisease struct {
	BasicInfo     DiseaseBasicInfo  `json:"basic_info"`
	Symptoms      []string          `json:"symptoms,omitempty"`
	RiskFactors   []string          `json:"risk_factors,omitempty"`
	Complications []string          `json:"complications,omitempty"`
	Treatment     DiseasesTreatment `json:"treatment"`
}

// DiseasesTreatment represents treatment options for diseases
type DiseasesTreatment struct {
	Lifestyle    []string `json:"lifestyle,omitempty"`
	Medications  []string `json:"medications,omitempty"`
	Procedures   []string `json:"procedures,omitempty"`
	Supportive   []string `json:"supportive,omitempty"`
	Acute        []string `json:"acute,omitempty"`
	Preventive   []string `json:"preventive,omitempty"`
	Prevention   []string `json:"prevention,omitempty"`
	Topical      []string `json:"topical,omitempty"`
	Systemic     []string `json:"systemic,omitempty"`
	Surgical     []string `json:"surgical,omitempty"`
	Conservative []string `json:"conservative,omitempty"`
	Supplements  []string `json:"supplements,omitempty"`
	Therapy      []string `json:"therapy,omitempty"`
	Monitoring   []string `json:"monitoring,omitempty"`
	Duration     []string `json:"duration,omitempty"`
	Dietary      []string `json:"dietary,omitempty"`
}

// Supporting data structures
type Medication struct {
	Generic           string   `json:"generic"`
	Brand             []string `json:"brand"`
	Category          string   `json:"category"`
	Dose              string   `json:"dose"`
	Administration    string   `json:"administration"`
	Mechanism         string   `json:"mechanism"`
	Contraindications []string `json:"contraindications,omitempty"`
	SideEffects       []string `json:"side_effects,omitempty"`
	Interactions      []string `json:"interactions,omitempty"`
}

type DiseaseBasicInfo struct {
	Name        string `json:"name"`
	ICD10       string `json:"icd10"`
	Category    string `json:"category"`
	Prevalence  string `json:"prevalence"`
	Description string `json:"description"`
}

type CauseWithDifferentialDiagnosis struct {
	Cause                 string                  `json:"cause"`
	Symptoms              []string                `json:"symptoms"`
	DifferentialDiagnosis []DifferentialDiagnosis `json:"differential_diagnosis"`
}

type DifferentialDiagnosis struct {
	Condition string   `json:"condition"`
	Symptoms  []string `json:"symptoms"`
}

type CardiovascularDisease struct {
	BasicInfo     DiseaseBasicInfo `json:"basic_info"`
	Causes        []string         `json:"causes,omitempty"`
	Symptoms      []string         `json:"symptoms"`
	Complications []string         `json:"complications,omitempty"`
	Treatment     Treatment        `json:"treatment"`
}

type EndocrineDisease struct {
	BasicInfo DiseaseBasicInfo `json:"basic_info"`
	Symptoms  []string         `json:"symptoms"`
	Treatment Treatment        `json:"treatment"`
}

type MetabolicDisease struct {
	BasicInfo       DiseaseBasicInfo `json:"basic_info"`
	Criteria        []string         `json:"criteria,omitempty"`
	Classifications []string         `json:"classifications,omitempty"`
	Treatment       Treatment        `json:"treatment"`
}

type Treatment struct {
	Lifestyle   []string `json:"lifestyle"`
	Medications []string `json:"medications"`
	Other       []string `json:"other,omitempty"`
	Monitoring  []string `json:"monitoring,omitempty"`
}

type HaramIngredient struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	CommonIn     []string `json:"common_in"`
	Alternatives []string `json:"alternatives"`
}

type MashboohIngredient struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Concern      string   `json:"concern"`
	Verification []string `json:"verification"`
}

type Vitamin struct {
	Name       BilingualName     `json:"name"`
	Function   string            `json:"function"`
	Deficiency string            `json:"deficiency"`
	Dosage     map[string]string `json:"dosage"`
}

type Mineral struct {
	Name       BilingualName     `json:"name"`
	Function   string            `json:"function"`
	Deficiency string            `json:"deficiency"`
	Dosage     map[string]string `json:"dosage"`
}

type BilingualName struct {
	En string `json:"en"`
	Ar string `json:"ar"`
}

type ScientificReference struct {
	Study   string `json:"study"`
	Finding string `json:"finding"`
	Link    string `json:"link"`
}

type PopulationDosing struct {
	Population string `json:"population"`
	Dosage     string `json:"dosage"`
	Notes      string `json:"notes"`
}

type WorkoutTechnique struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Benefits     []string         `json:"benefits"`
	HowToPerform []string         `json:"how_to_perform"`
	Examples     []WorkoutExample `json:"examples"`
}

type WorkoutExample struct {
	Exercise string `json:"exercise"`
	Sets     string `json:"sets"`
	Reps     string `json:"reps"`
	Rest     string `json:"rest"`
}

type WeightLossDrug struct {
	Name              string   `json:"name"`
	Mechanism         string   `json:"mechanism"`
	Dosage            string   `json:"dosage"`
	SideEffects       []string `json:"side_effects"`
	Contraindications []string `json:"contraindications"`
}

type CookingSkill struct {
	Skill       string   `json:"skill"`
	Description string   `json:"description"`
	Techniques  []string `json:"techniques"`
	Tips        []string `json:"tips"`
}

type SportsNutrition struct {
	Sports map[string]SportNutrition `json:"sports"`
}

type SportNutrition struct {
	PreWorkout    NutritionTiming `json:"pre_workout"`
	DuringWorkout NutritionTiming `json:"during_workout"`
	PostWorkout   NutritionTiming `json:"post_workout"`
}

type NutritionTiming struct {
	Timing      string   `json:"timing"`
	Foods       []string `json:"foods"`
	Supplements []string `json:"supplements"`
	Hydration   string   `json:"hydration"`
}

type InjuryManagement struct {
	GymSprains InjuryTreatment `json:"gym_sprains"`
}

type InjuryTreatment struct {
	Description   string                     `json:"description"`
	ImmediateCare []string                   `json:"immediate_care"`
	Supplements   []InjurySupplementProtocol `json:"supplements"`
	Medications   []InjuryMedication         `json:"medications"`
	Recovery      []string                   `json:"recovery"`
}

type InjurySupplementProtocol struct {
	Supplement string `json:"supplement"`
	Dosage     string `json:"dosage"`
	Duration   string `json:"duration"`
	Purpose    string `json:"purpose"`
}

type InjuryMedication struct {
	Medication string `json:"medication"`
	Dosage     string `json:"dosage"`
	Duration   string `json:"duration"`
	Purpose    string `json:"purpose"`
}

type DietPlans struct {
	Keto KetoDiet `json:"keto"`
}

type KetoDiet struct {
	Description    string   `json:"description"`
	MacroRatio     string   `json:"macro_ratio"`
	AllowedFoods   []string `json:"allowed_foods"`
	ForbiddenFoods []string `json:"forbidden_foods"`
	Benefits       []string `json:"benefits"`
	Considerations []string `json:"considerations"`
}

type NutritionGuidelines struct {
	CalorieRecommendations CalorieRecommendations `json:"calorie_recommendations"`
	MealTiming             MealTiming             `json:"meal_timing"`
}

type CalorieRecommendations struct {
	SedentaryMale   string `json:"sedentary_male"`
	ActiveMale      string `json:"active_male"`
	SedentaryFemale string `json:"sedentary_female"`
	ActiveFemale    string `json:"active_female"`
}

type MealTiming struct {
	Breakfast string `json:"breakfast"`
	Lunch     string `json:"lunch"`
	Dinner    string `json:"dinner"`
	Snacks    string `json:"snacks"`
}

type ComplaintCase struct {
	ID                      int                      `json:"id"`
	ConditionEn             string                   `json:"condition_en"`
	ConditionAr             string                   `json:"condition_ar"`
	Recommendations         CaseRecommendations      `json:"recommendations"`
	EnhancedRecommendations *EnhancedRecommendations `json:"enhanced_recommendations,omitempty"`
}

type CaseRecommendations struct {
	Nutrition           BilingualRecommendation  `json:"nutrition"`
	SpecificFoods       BilingualRecommendation  `json:"specific_foods"`
	VitaminsSupplements BilingualRecommendation  `json:"vitamins_supplements"`
	Exercise            *BilingualRecommendation `json:"exercise,omitempty"`
	Medications         BilingualRecommendation  `json:"medications"`
}

type BilingualRecommendation struct {
	En string `json:"en"`
	Ar string `json:"ar"`
}

// LoadComprehensiveHealthDatabase loads the complete health database
func (cdl *ComprehensiveDataLoader) LoadComprehensiveHealthDatabase() (*ComprehensiveHealthDatabase, error) {
	// Try different file paths for VIP health database
	filePath := filepath.Join(cdl.dataPath, "vip-workouts.js")
	var result ComprehensiveHealthDatabase
	err := cdl.loadJSONFile(filePath, &result)
	return &result, err
}

// LoadComprehensiveComplaintsDatabase loads the complete complaints database
func (cdl *ComprehensiveDataLoader) LoadComprehensiveComplaintsDatabase() (*ComprehensiveComplaintsDatabase, error) {
	filePath := filepath.Join(cdl.dataPath, "comprehensive-complaints-database-COMPLETE.js")
	var result ComprehensiveComplaintsDatabase
	err := cdl.loadJSONFile(filePath, &result)
	return &result, err
}

// LoadComprehensiveDiseasesDatabase loads the complete diseases database
func (cdl *ComprehensiveDataLoader) LoadComprehensiveDiseasesDatabase() (*ComprehensiveDiseasesDatabase, error) {
	filePath := filepath.Join(cdl.dataPath, "comprehensive-diseases-database-EXPANDED.js")
	var result ComprehensiveDiseasesDatabase
	err := cdl.loadJSONFile(filePath, &result)
	return &result, err
}

// loadJSONFile is a generic function to load and parse comprehensive JSON files
func (cdl *ComprehensiveDataLoader) loadJSONFile(filePath string, result interface{}) error {
	// Validate file path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(cdl.dataPath, cleanPath)
	}

	// Ensure the file is within the allowed data directory
	if !isWithinVIPDirectory(cleanPath, cdl.dataPath) {
		return fmt.Errorf("file path outside allowed directory: %s", filePath)
	}

	// Read the file
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read comprehensive file %s: %w", cleanPath, err)
	}

	// Convert JavaScript module.exports to JSON
	jsonData := cdl.convertJSToJSON(string(data))

	// Parse JSON
	err = json.Unmarshal([]byte(jsonData), result)
	if err != nil {
		return fmt.Errorf("failed to parse comprehensive JSON from %s: %w", filePath, err)
	}

	return nil
}

// convertJSToJSON converts JavaScript module format to JSON for comprehensive files
func (cdl *ComprehensiveDataLoader) convertJSToJSON(jsContent string) string {
	// Handle comprehensive database format
	// Find the start of the actual data structure
	start := 0
	if idx := strings.Index(jsContent, "const comprehensive"); idx != -1 {
		// Find the opening brace after the const declaration
		for i := idx; i < len(jsContent); i++ {
			if jsContent[i] == '{' {
				start = i
				break
			}
		}
	} else if idx := strings.Index(jsContent, "const complete"); idx != -1 {
		// Find the opening brace after the const declaration
		for i := idx; i < len(jsContent); i++ {
			if jsContent[i] == '{' {
				start = i
				break
			}
		}
	}

	// Find the end of the data structure (before module.exports)
	end := len(jsContent)
	if idx := strings.Index(jsContent, "module.exports"); idx != -1 {
		end = idx
		// Trim any whitespace before module.exports
		for end > 0 && (jsContent[end-1] == ' ' || jsContent[end-1] == '\n' || jsContent[end-1] == '\r' || jsContent[end-1] == '\t') {
			end--
		}
	}

	jsonContent := strings.TrimSpace(jsContent[start:end])

	// Remove JavaScript comments
	jsonContent = cdl.removeJSComments(jsonContent)

	return jsonContent
}

// removeJSComments removes JavaScript-style comments from JSON content
func (cdl *ComprehensiveDataLoader) removeJSComments(content string) string {
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Remove single-line comments (//)
		if idx := strings.Index(line, "//"); idx != -1 {
			// Make sure it's not inside a string
			inString := false
			escaped := false
			for i, char := range line {
				if escaped {
					escaped = false
					continue
				}
				if char == '\\' {
					escaped = true
					continue
				}
				if char == '"' {
					inString = !inString
				}
				if i == idx && !inString {
					line = line[:idx]
					break
				}
			}
		}

		// Keep non-empty lines
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// GetAllMedications extracts all medications from the comprehensive database
func (cdl *ComprehensiveDataLoader) GetAllMedications() ([]Medication, error) {
	db, err := cdl.LoadComprehensiveHealthDatabase()
	if err != nil {
		return nil, err
	}

	var allMedications []Medication
	allMedications = append(allMedications, db.Medications.WeightLossDrugs...)
	allMedications = append(allMedications, db.Medications.DiabetesMedications...)
	allMedications = append(allMedications, db.Medications.AppetiteSuppressants...)
	allMedications = append(allMedications, db.Medications.InvestigationalDrugs...)
	allMedications = append(allMedications, db.Medications.HormonalMedications...)
	allMedications = append(allMedications, db.Medications.PerformanceEnhancers...)

	return allMedications, nil
}

// GetAllVitaminsAndMinerals extracts all vitamins and minerals from the comprehensive database
func (cdl *ComprehensiveDataLoader) GetAllVitaminsAndMinerals() ([]interface{}, error) {
	db, err := cdl.LoadComprehensiveHealthDatabase()
	if err != nil {
		return nil, err
	}

	var allNutrients []interface{}
	for _, vitamin := range db.Vitamins.WaterSoluble {
		allNutrients = append(allNutrients, vitamin)
	}
	for _, vitamin := range db.Vitamins.FatSoluble {
		allNutrients = append(allNutrients, vitamin)
	}
	for _, mineral := range db.Minerals {
		allNutrients = append(allNutrients, mineral)
	}

	return allNutrients, nil
}

// GetAllComplaintCases extracts all complaint cases from the comprehensive database
func (cdl *ComprehensiveDataLoader) GetAllComplaintCases() ([]ComplaintCase, error) {
	db, err := cdl.LoadComprehensiveComplaintsDatabase()
	if err != nil {
		return nil, err
	}

	return db.Cases, nil
}

// isWithinVIPDirectory checks if a file path is within the allowed VIP directory
func isWithinVIPDirectory(filePath, allowedDir string) bool {
	rel, err := filepath.Rel(allowedDir, filePath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, "/")
}
