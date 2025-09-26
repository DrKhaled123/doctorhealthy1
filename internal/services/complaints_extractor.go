package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// ComplaintData represents the raw complaint data from JSON
type ComplaintData struct {
	ID                      int                     `json:"id"`
	ConditionEN             string                  `json:"condition_en"`
	ConditionAR             string                  `json:"condition_ar"`
	Recommendations         Recommendations         `json:"recommendations"`
	EnhancedRecommendations EnhancedRecommendations `json:"enhanced_recommendations"`
}

// Recommendations represents basic treatment recommendations
type Recommendations struct {
	Nutrition           BilingualContent `json:"nutrition"`
	SpecificFoods       BilingualContent `json:"specific_foods"`
	VitaminsSupplements BilingualContent `json:"vitamins_supplements"`
	Exercise            BilingualContent `json:"exercise"`
	Medications         BilingualContent `json:"medications"`
}

// EnhancedRecommendations represents advanced treatment protocols
type EnhancedRecommendations struct {
	AdvancedNutrition      BilingualContentWithRefs `json:"advanced_nutrition"`
	AdvancedWorkout        BilingualContentWithRefs `json:"advanced_workout"`
	LifestyleModifications BilingualContentWithRefs `json:"lifestyle_modifications"`
	AdditionalSupplements  BilingualContentWithRefs `json:"additional_supplements"`
	ClinicalInsights       BilingualContentWithRefs `json:"clinical_insights"`
}

// BilingualContent represents content in both English and Arabic
type BilingualContent struct {
	EN string `json:"en"`
	AR string `json:"ar"`
}

// BilingualContentWithRefs represents bilingual content with reference links
type BilingualContentWithRefs struct {
	EN         string   `json:"en"`
	AR         string   `json:"ar"`
	References []string `json:"references"`
}

// MedicationProtocol represents parsed medication information
type MedicationProtocol struct {
	Name                string   `json:"name"`
	Dosage              string   `json:"dosage"`
	Frequency           string   `json:"frequency"`
	Duration            string   `json:"duration"`
	SideEffects         []string `json:"side_effects"`
	Contraindications   []string `json:"contraindications"`
	Interactions        []string `json:"interactions"`
	SupervisionRequired bool     `json:"supervision_required"`
	References          []string `json:"references"`
}

// SupplementProtocol represents parsed supplement information
type SupplementProtocol struct {
	Name         string   `json:"name"`
	Dosage       string   `json:"dosage"`
	Timing       string   `json:"timing"`
	Purpose      string   `json:"purpose"`
	SafetyNotes  []string `json:"safety_notes"`
	Interactions []string `json:"interactions"`
	References   []string `json:"references"`
}

// ExtractionStats provides statistics about the extraction process
type ExtractionStats struct {
	TotalCases           int           `json:"total_cases"`
	ProcessedCases       int           `json:"processed_cases"`
	FailedCases          int           `json:"failed_cases"`
	ExtractedMedications int           `json:"extracted_medications"`
	ExtractedSupplements int           `json:"extracted_supplements"`
	ExtractedReferences  int           `json:"extracted_references"`
	ProcessingTime       time.Duration `json:"processing_time"`
	Errors               []string      `json:"errors"`
}

// DataExtractor interface for complaint data extraction
type DataExtractor interface {
	ExtractComplaints(source string) ([]ComplaintData, error)
	ValidateSource(source string) error
	GetExtractionStats() ExtractionStats
	ParseMedicationProtocols(content BilingualContent) ([]MedicationProtocol, error)
	ParseSupplementProtocols(content BilingualContent) ([]SupplementProtocol, error)
	ExtractReferences(content BilingualContentWithRefs) []string
}

// ComplaintsExtractor implements the DataExtractor interface
type ComplaintsExtractor struct {
	stats ExtractionStats
}

// NewComplaintsExtractor creates a new complaints extractor
func NewComplaintsExtractor() *ComplaintsExtractor {
	return &ComplaintsExtractor{
		stats: ExtractionStats{
			Errors: make([]string, 0),
		},
	}
}

// ValidateSource validates the JSON source file
func (e *ComplaintsExtractor) ValidateSource(source string) error {
	file, err := os.Open(source) // nosec G304 - controlled path within application data directory
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Check if file is readable and has content
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	if stat.Size() == 0 {
		return fmt.Errorf("source file is empty")
	}

	// Try to read first few bytes to ensure it's valid JSON
	buffer := make([]byte, 1024)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Basic JSON validation - should start with { or [
	content := strings.TrimSpace(string(buffer))
	if !strings.HasPrefix(content, "{") && !strings.HasPrefix(content, "[") {
		return fmt.Errorf("source file does not appear to be valid JSON")
	}

	return nil
}

// ExtractComplaints extracts all complaint data from the JSON source
func (e *ComplaintsExtractor) ExtractComplaints(source string) ([]ComplaintData, error) {
	startTime := time.Now()
	e.stats = ExtractionStats{
		Errors: make([]string, 0),
	}

	// Validate source first
	if err := e.ValidateSource(source); err != nil {
		e.stats.Errors = append(e.stats.Errors, fmt.Sprintf("Source validation failed: %v", err))
		return nil, fmt.Errorf("source validation failed: %w", err)
	}

	// Read the JSON file
	file, err := os.Open(source) // nosec G304 - controlled path within application data directory
	if err != nil {
		e.stats.Errors = append(e.stats.Errors, fmt.Sprintf("Failed to open file: %v", err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		e.stats.Errors = append(e.stats.Errors, fmt.Sprintf("Failed to read file: %v", err))
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON structure
	var jsonData struct {
		Cases []ComplaintData `json:"cases"`
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		e.stats.Errors = append(e.stats.Errors, fmt.Sprintf("Failed to parse JSON: %v", err))
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	e.stats.TotalCases = len(jsonData.Cases)
	complaints := make([]ComplaintData, 0, len(jsonData.Cases))

	// Process each complaint case
	for i, complaint := range jsonData.Cases {
		// Create a copy to avoid memory aliasing issues
		complaintCopy := complaint
		if err := e.validateComplaintData(&complaintCopy); err != nil {
			e.stats.FailedCases++
			e.stats.Errors = append(e.stats.Errors, fmt.Sprintf("Case %d validation failed: %v", complaintCopy.ID, err))
			continue
		}

		// Extract and parse references from enhanced recommendations
		complaint.EnhancedRecommendations.AdvancedNutrition.References = e.ExtractReferences(complaint.EnhancedRecommendations.AdvancedNutrition)
		complaint.EnhancedRecommendations.AdvancedWorkout.References = e.ExtractReferences(complaint.EnhancedRecommendations.AdvancedWorkout)
		complaint.EnhancedRecommendations.LifestyleModifications.References = e.ExtractReferences(complaint.EnhancedRecommendations.LifestyleModifications)
		complaint.EnhancedRecommendations.AdditionalSupplements.References = e.ExtractReferences(complaint.EnhancedRecommendations.AdditionalSupplements)
		complaint.EnhancedRecommendations.ClinicalInsights.References = e.ExtractReferences(complaint.EnhancedRecommendations.ClinicalInsights)

		// Count extracted references
		e.stats.ExtractedReferences += len(complaint.EnhancedRecommendations.AdvancedNutrition.References)
		e.stats.ExtractedReferences += len(complaint.EnhancedRecommendations.AdvancedWorkout.References)
		e.stats.ExtractedReferences += len(complaint.EnhancedRecommendations.LifestyleModifications.References)
		e.stats.ExtractedReferences += len(complaint.EnhancedRecommendations.AdditionalSupplements.References)
		e.stats.ExtractedReferences += len(complaint.EnhancedRecommendations.ClinicalInsights.References)

		complaints = append(complaints, complaint)
		e.stats.ProcessedCases++

		// Log progress for large datasets
		if (i+1)%100 == 0 {
			fmt.Printf("Processed %d/%d complaints...\n", i+1, len(jsonData.Cases))
		}
	}

	e.stats.ProcessingTime = time.Since(startTime)

	fmt.Printf("Extraction completed: %d/%d cases processed successfully\n",
		e.stats.ProcessedCases, e.stats.TotalCases)

	return complaints, nil
}

// validateComplaintData validates individual complaint data
func (e *ComplaintsExtractor) validateComplaintData(complaint *ComplaintData) error {
	if complaint.ID <= 0 {
		return fmt.Errorf("invalid complaint ID: %d", complaint.ID)
	}

	if strings.TrimSpace(complaint.ConditionEN) == "" {
		return fmt.Errorf("missing English condition description")
	}

	if strings.TrimSpace(complaint.ConditionAR) == "" {
		return fmt.Errorf("missing Arabic condition description")
	}

	// Validate that essential recommendations exist
	if strings.TrimSpace(complaint.Recommendations.Nutrition.EN) == "" {
		return fmt.Errorf("missing nutrition recommendations")
	}

	return nil
}

// ParseMedicationProtocols extracts medication information from bilingual content
func (e *ComplaintsExtractor) ParseMedicationProtocols(content BilingualContent) ([]MedicationProtocol, error) {
	protocols := make([]MedicationProtocol, 0)

	// Parse English content for medication details
	medications := e.extractMedicationsFromText(content.EN)

	for _, med := range medications {
		protocol := MedicationProtocol{
			Name:                med.Name,
			Dosage:              med.Dosage,
			Frequency:           med.Frequency,
			Duration:            med.Duration,
			SupervisionRequired: strings.Contains(strings.ToLower(content.EN), "medical supervision required"),
			SideEffects:         make([]string, 0),
			Contraindications:   make([]string, 0),
			Interactions:        make([]string, 0),
		}

		protocols = append(protocols, protocol)
		e.stats.ExtractedMedications++
	}

	return protocols, nil
}

// ParseSupplementProtocols extracts supplement information from bilingual content
func (e *ComplaintsExtractor) ParseSupplementProtocols(content BilingualContent) ([]SupplementProtocol, error) {
	protocols := make([]SupplementProtocol, 0)

	// Parse English content for supplement details
	supplements := e.extractSupplementsFromText(content.EN)

	for _, supp := range supplements {
		protocol := SupplementProtocol{
			Name:         supp.Name,
			Dosage:       supp.Dosage,
			Timing:       supp.Timing,
			Purpose:      supp.Purpose,
			SafetyNotes:  make([]string, 0),
			Interactions: make([]string, 0),
		}

		protocols = append(protocols, protocol)
		e.stats.ExtractedSupplements++
	}

	return protocols, nil
}

// ExtractReferences extracts reference URLs from bilingual content
func (e *ComplaintsExtractor) ExtractReferences(content BilingualContentWithRefs) []string {
	references := make([]string, 0)

	// Regex pattern to match URLs in square brackets
	urlPattern := regexp.MustCompile(`\[Ref:\s*(https?://[^\]]+)\]`)

	// Extract from English content
	matches := urlPattern.FindAllStringSubmatch(content.EN, -1)
	for _, match := range matches {
		if len(match) > 1 {
			references = append(references, strings.TrimSpace(match[1]))
		}
	}

	// Extract from Arabic content (in case there are different references)
	matches = urlPattern.FindAllStringSubmatch(content.AR, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			// Avoid duplicates
			found := false
			for _, existing := range references {
				if existing == url {
					found = true
					break
				}
			}
			if !found {
				references = append(references, url)
			}
		}
	}

	return references
}

// GetExtractionStats returns the current extraction statistics
func (e *ComplaintsExtractor) GetExtractionStats() ExtractionStats {
	return e.stats
}

// Helper structures for parsing
type medicationInfo struct {
	Name      string
	Dosage    string
	Frequency string
	Duration  string
}

type supplementInfo struct {
	Name    string
	Dosage  string
	Timing  string
	Purpose string
}

// extractMedicationsFromText parses medication information from text
func (e *ComplaintsExtractor) extractMedicationsFromText(text string) []medicationInfo {
	medications := make([]medicationInfo, 0)

	// Pattern to match medication with dosage: "MedicationName: dosage info"
	medPattern := regexp.MustCompile(`([A-Za-z\-\s]+):\s*([^.]+)`)

	matches := medPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 2 {
			name := strings.TrimSpace(match[1])
			dosageInfo := strings.TrimSpace(match[2])

			// Skip if it doesn't look like a medication
			if strings.Contains(strings.ToLower(name), "start") ||
				strings.Contains(strings.ToLower(name), "increase") {
				continue
			}

			med := medicationInfo{
				Name:   name,
				Dosage: dosageInfo,
			}

			// Extract frequency and duration if present
			if strings.Contains(dosageInfo, "/day") {
				med.Frequency = "daily"
			}
			if strings.Contains(dosageInfo, "/week") {
				med.Frequency = "weekly"
			}

			medications = append(medications, med)
		}
	}

	return medications
}

// extractSupplementsFromText parses supplement information from text
func (e *ComplaintsExtractor) extractSupplementsFromText(text string) []supplementInfo {
	supplements := make([]supplementInfo, 0)

	// Pattern to match supplements with dosage: "SupplementName: dosage"
	suppPattern := regexp.MustCompile(`([A-Za-z\-\s]+):\s*([^.]+)`)

	matches := suppPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 2 {
			name := strings.TrimSpace(match[1])
			dosageInfo := strings.TrimSpace(match[2])

			supp := supplementInfo{
				Name:   name,
				Dosage: dosageInfo,
			}

			// Determine timing based on context
			if strings.Contains(strings.ToLower(text), "pre-workout") {
				supp.Timing = "pre-workout"
			} else if strings.Contains(strings.ToLower(text), "post-workout") {
				supp.Timing = "post-workout"
			} else {
				supp.Timing = "daily"
			}

			supplements = append(supplements, supp)
		}
	}

	return supplements
}
