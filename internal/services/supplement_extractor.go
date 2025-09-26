package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"api-key-generator/internal/models"
)

// SupplementExtractor extracts supplement data from VIP JSON files
type SupplementExtractor struct {
	dataPath string
}

// NewSupplementExtractor creates a new supplement extractor
func NewSupplementExtractor(dataPath string) *SupplementExtractor {
	return &SupplementExtractor{
		dataPath: dataPath,
	}
}

// VIPSupplementData represents the structure of VIP supplement data
type VIPSupplementData struct {
	WeightLossDrugs []VIPDrug       `json:"weight_loss_drugs"`
	Supplements     []VIPSupplement `json:"supplements"`
}

// VIPDrug represents drug information from VIP data
type VIPDrug struct {
	DrugName struct {
		Generic string   `json:"generic"`
		Brand   []string `json:"brand"`
	} `json:"drug_name"`
	Doses struct {
		TypicalDose  string `json:"typical_dose"`
		StartingDose string `json:"starting_dose"`
		MaximumDose  string `json:"maximum_dose"`
	} `json:"doses"`
	Mechanism   string `json:"mechanism_of_action"`
	SideEffects struct {
		Common  []string `json:"common"`
		Serious []string `json:"serious"`
	} `json:"side_effects"`
	Interactions      []string `json:"interactions"`
	Contraindications []string `json:"contraindications"`
}

// VIPSupplement represents supplement information
type VIPSupplement struct {
	Name              string   `json:"name"`
	Category          string   `json:"category"`
	RecommendedDose   string   `json:"recommended_dose"`
	MaxDose           string   `json:"max_dose"`
	Benefits          []string `json:"benefits"`
	SideEffects       []string `json:"side_effects"`
	Interactions      []string `json:"interactions"`
	Contraindications []string `json:"contraindications"`
	Timing            string   `json:"timing"`
}

// ExtractSupplementData extracts supplement data from VIP JSON files
func (e *SupplementExtractor) ExtractSupplementData() ([]models.SupplementTiming, error) {
	// Read the VIP JSON file
	filePath := fmt.Sprintf("%s/drugs and nutrition.js", e.dataPath)
	data, err := e.readVIPFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read VIP file: %w", err)
	}

	// Parse the JSON data
	var vipData VIPSupplementData
	if err := json.Unmarshal(data, &vipData); err != nil {
		return nil, fmt.Errorf("failed to parse VIP JSON: %w", err)
	}

	// Convert to our supplement timing format
	supplements := []models.SupplementTiming{}

	// Extract from weight loss drugs that can be used as supplements
	for _, drug := range vipData.WeightLossDrugs {
		if e.isSupplement(drug.DrugName.Generic) {
			supplement := e.convertDrugToSupplement(drug)
			supplements = append(supplements, supplement)
		}
	}

	// Extract from supplements section
	for _, supp := range vipData.Supplements {
		supplement := e.convertVIPSupplement(supp)
		supplements = append(supplements, supplement)
	}

	return supplements, nil
}

// readVIPFile reads and processes VIP JavaScript file
func (e *SupplementExtractor) readVIPFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath) // nosec G304 - controlled path within application data directory
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Convert JavaScript to JSON
	jsonContent := e.convertJSToJSON(string(content))
	return []byte(jsonContent), nil
}

// convertJSToJSON converts JavaScript module format to JSON
func (e *SupplementExtractor) convertJSToJSON(jsContent string) string {
	// Find the start of the data structure
	start := strings.Index(jsContent, "{")
	if start == -1 {
		return "{}"
	}

	// Find the end before module.exports
	end := strings.Index(jsContent, "module.exports")
	if end == -1 {
		end = len(jsContent)
	}

	// Extract JSON content
	jsonContent := strings.TrimSpace(jsContent[start:end])

	// Remove JavaScript comments
	jsonContent = e.removeJSComments(jsonContent)

	return jsonContent
}

// removeJSComments removes JavaScript-style comments
func (e *SupplementExtractor) removeJSComments(content string) string {
	// Remove single-line comments
	re := regexp.MustCompile(`//.*`)
	content = re.ReplaceAllString(content, "")

	// Remove multi-line comments
	re = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	content = re.ReplaceAllString(content, "")

	return content
}

// isSupplement checks if a drug can be considered a supplement
func (e *SupplementExtractor) isSupplement(drugName string) bool {
	supplements := []string{"metformin", "orlistat", "caffeine", "green tea extract"}
	drugLower := strings.ToLower(drugName)

	for _, supp := range supplements {
		if strings.Contains(drugLower, supp) {
			return true
		}
	}
	return false
}

// convertDrugToSupplement converts VIP drug data to supplement timing
func (e *SupplementExtractor) convertDrugToSupplement(drug VIPDrug) models.SupplementTiming {
	return models.SupplementTiming{
		SupplementName: drug.DrugName.Generic,
		Dosage:         e.extractDosage(drug.Doses.TypicalDose),
		Unit:           e.extractUnit(drug.Doses.TypicalDose),
		TimingMinutes:  0,
		Instructions:   fmt.Sprintf("Take as directed: %s", drug.Doses.TypicalDose),
		Benefits:       e.extractBenefits(drug.Mechanism),
		SideEffects:    drug.SideEffects.Common,
		MaxDailyDose:   drug.Doses.MaximumDose,
	}
}

// convertVIPSupplement converts VIP supplement data to our format
func (e *SupplementExtractor) convertVIPSupplement(supp VIPSupplement) models.SupplementTiming {
	return models.SupplementTiming{
		SupplementName: supp.Name,
		Dosage:         e.extractDosage(supp.RecommendedDose),
		Unit:           e.extractUnit(supp.RecommendedDose),
		TimingMinutes:  e.parseTimingMinutes(supp.Timing),
		Instructions:   fmt.Sprintf("Take %s", supp.Timing),
		Benefits:       supp.Benefits,
		SideEffects:    supp.SideEffects,
		MaxDailyDose:   supp.MaxDose,
	}
}

// extractDosage extracts numeric dosage from dose string
func (e *SupplementExtractor) extractDosage(doseStr string) string {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(doseStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return "1"
}

// extractUnit extracts unit from dose string
func (e *SupplementExtractor) extractUnit(doseStr string) string {
	units := []string{"mg", "g", "ml", "mcg", "iu", "tablet", "capsule", "scoop"}
	doseStrLower := strings.ToLower(doseStr)

	for _, unit := range units {
		if strings.Contains(doseStrLower, unit) {
			return unit
		}
	}
	return "mg"
}

// extractBenefits extracts benefits from mechanism description
func (e *SupplementExtractor) extractBenefits(mechanism string) []string {
	benefits := []string{}

	if strings.Contains(strings.ToLower(mechanism), "weight") {
		benefits = append(benefits, "Weight management")
	}
	if strings.Contains(strings.ToLower(mechanism), "energy") {
		benefits = append(benefits, "Energy boost")
	}
	if strings.Contains(strings.ToLower(mechanism), "muscle") {
		benefits = append(benefits, "Muscle support")
	}

	if len(benefits) == 0 {
		benefits = append(benefits, "Health support")
	}

	return benefits
}

// parseTimingMinutes converts timing description to minutes
func (e *SupplementExtractor) parseTimingMinutes(timing string) int {
	timingLower := strings.ToLower(timing)

	if strings.Contains(timingLower, "before") {
		if strings.Contains(timingLower, "30") {
			return -30
		}
		return -15
	}
	if strings.Contains(timingLower, "after") {
		if strings.Contains(timingLower, "30") {
			return 30
		}
		return 15
	}

	return 0 // During or with meals
}
