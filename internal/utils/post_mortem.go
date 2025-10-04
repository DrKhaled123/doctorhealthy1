package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api-key-generator/internal/models"
)

// Incident represents a production incident for post-mortem analysis
type Incident struct {
	ID              string                         `json:"id"`
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	Severity        string                         `json:"severity"`
	Status          string                         `json:"status"`
	StartTime       time.Time                      `json:"start_time"`
	EndTime         *time.Time                     `json:"end_time,omitempty"`
	Duration        string                         `json:"duration"`
	RootCause       string                         `json:"root_cause"`
	Impact          string                         `json:"impact"`
	AffectedSystems []string                       `json:"affected_systems"`
	RelatedErrors   []models.EnhancedErrorResponse `json:"related_errors"`
	Resolution      string                         `json:"resolution"`
	Prevention      string                         `json:"prevention"`
	LessonsLearned  []string                       `json:"lessons_learned"`
	Tags            []string                       `json:"tags"`
	CreatedBy       string                         `json:"created_by"`
	UpdatedAt       time.Time                      `json:"updated_at"`
}

// PostMortemManager manages incident documentation and analysis
type PostMortemManager struct {
	IncidentsPath string
	Incidents     []Incident
	AutoCreate    bool
}

// NewPostMortemManager creates a new post-mortem manager
func NewPostMortemManager(incidentsPath string) *PostMortemManager {
	return &PostMortemManager{
		IncidentsPath: incidentsPath,
		Incidents:     make([]Incident, 0),
		AutoCreate:    true,
	}
}

// CreateIncident creates a new incident from an error response
func (pmm *PostMortemManager) CreateIncident(errorResp *models.EnhancedErrorResponse, title, description string) (*Incident, error) {
	incident := &Incident{
		ID:            generateIncidentID(),
		Title:         title,
		Description:   description,
		Severity:      errorResp.Severity,
		Status:        "investigating",
		StartTime:     time.Now(),
		RelatedErrors: []models.EnhancedErrorResponse{*errorResp},
		Tags:          []string{"automated", "error_response"},
		CreatedBy:     "system",
		UpdatedAt:     time.Now(),
	}

	// Determine affected systems based on error category
	incident.AffectedSystems = pmm.determineAffectedSystems(errorResp.Category)

	// Load existing incidents
	if err := pmm.loadIncidents(); err != nil {
		log.Printf("Warning: Could not load existing incidents: %v", err)
	}

	// Add to incidents list
	pmm.Incidents = append(pmm.Incidents, *incident)

	// Save incidents
	if err := pmm.saveIncidents(); err != nil {
		return nil, fmt.Errorf("failed to save incident: %w", err)
	}

	log.Printf("Created incident: %s - %s", incident.ID, incident.Title)
	return incident, nil
}

// UpdateIncident updates an existing incident
func (pmm *PostMortemManager) UpdateIncident(id string, updates map[string]interface{}) error {
	// Load incidents
	if err := pmm.loadIncidents(); err != nil {
		return fmt.Errorf("failed to load incidents: %w", err)
	}

	// Find and update incident
	for i := range pmm.Incidents {
		if pmm.Incidents[i].ID == id {
			// Apply updates
			if rootCause, ok := updates["root_cause"].(string); ok {
				pmm.Incidents[i].RootCause = rootCause
			}
			if resolution, ok := updates["resolution"].(string); ok {
				pmm.Incidents[i].Resolution = resolution
			}
			if prevention, ok := updates["prevention"].(string); ok {
				pmm.Incidents[i].Prevention = prevention
			}
			if status, ok := updates["status"].(string); ok {
				pmm.Incidents[i].Status = status
				if status == "resolved" && pmm.Incidents[i].EndTime == nil {
					endTime := time.Now()
					pmm.Incidents[i].EndTime = &endTime
					duration := endTime.Sub(pmm.Incidents[i].StartTime)
					pmm.Incidents[i].Duration = duration.String()
				}
			}
			if lessons, ok := updates["lessons_learned"].([]string); ok {
				pmm.Incidents[i].LessonsLearned = lessons
			}

			pmm.Incidents[i].UpdatedAt = time.Now()

			// Save updated incidents
			return pmm.saveIncidents()
		}
	}

	return fmt.Errorf("incident with ID %s not found", id)
}

// GetIncident retrieves an incident by ID
func (pmm *PostMortemManager) GetIncident(id string) (*Incident, error) {
	if err := pmm.loadIncidents(); err != nil {
		return nil, err
	}

	for _, incident := range pmm.Incidents {
		if incident.ID == id {
			return &incident, nil
		}
	}

	return nil, fmt.Errorf("incident with ID %s not found", id)
}

// GetAllIncidents returns all incidents
func (pmm *PostMortemManager) GetAllIncidents() []Incident {
	pmm.loadIncidents()
	return pmm.Incidents
}

// GetIncidentsByStatus returns incidents filtered by status
func (pmm *PostMortemManager) GetIncidentsByStatus(status string) []Incident {
	allIncidents := pmm.GetAllIncidents()
	var filtered []Incident

	for _, incident := range allIncidents {
		if incident.Status == status {
			filtered = append(filtered, incident)
		}
	}

	return filtered
}

// GetIncidentsBySeverity returns incidents filtered by severity
func (pmm *PostMortemManager) GetIncidentsBySeverity(severity string) []Incident {
	allIncidents := pmm.GetAllIncidents()
	var filtered []Incident

	for _, incident := range allIncidents {
		if incident.Severity == severity {
			filtered = append(filtered, incident)
		}
	}

	return filtered
}

// GeneratePostMortemReport generates a comprehensive post-mortem report
func (pmm *PostMortemManager) GeneratePostMortemReport(id string) (string, error) {
	incident, err := pmm.GetIncident(id)
	if err != nil {
		return "", err
	}

	report := fmt.Sprintf(`# Post-Mortem Report: %s

## Incident Summary
- **ID**: %s
- **Title**: %s
- **Severity**: %s
- **Status**: %s
- **Start Time**: %s
- **Duration**: %s

## Description
%s

## Root Cause
%s

## Impact
%s

## Affected Systems
%s

## Resolution
%s

## Prevention Measures
%s

## Lessons Learned
%s

## Timeline
- **Started**: %s
- **Detected**: %s
- **Resolved**: %s

## Related Errors
%s

---
*Generated on: %s*
`,
		incident.Title,
		incident.ID,
		incident.Title,
		incident.Severity,
		incident.Status,
		incident.StartTime.Format(time.RFC3339),
		incident.Duration,
		incident.Description,
		incident.RootCause,
		incident.Impact,
		fmt.Sprintf("- %s", strings.Join(incident.AffectedSystems, "\n- ")),
		incident.Resolution,
		incident.Prevention,
		fmt.Sprintf("- %s", strings.Join(incident.LessonsLearned, "\n- ")),
		incident.StartTime.Format(time.RFC3339),
		incident.StartTime.Format(time.RFC3339),
		func() string {
			if incident.EndTime != nil {
				return incident.EndTime.Format(time.RFC3339)
			}
			return "Not yet resolved"
		}(),
		pmm.formatRelatedErrors(incident.RelatedErrors),
		time.Now().Format(time.RFC3339),
	)

	return report, nil
}

// determineAffectedSystems determines which systems are affected based on error category
func (pmm *PostMortemManager) determineAffectedSystems(category string) []string {
	switch category {
	case "database":
		return []string{"Database Layer", "Data Persistence", "Query Engine"}
	case "network":
		return []string{"Network Layer", "API Gateway", "Load Balancer"}
	case "wasm":
		return []string{"WebAssembly Runtime", "Browser Compatibility", "Performance Layer"}
	case "security":
		return []string{"Security Layer", "Authentication", "Authorization"}
	case "performance":
		return []string{"Performance Layer", "Resource Management", "Scaling"}
	default:
		return []string{"General Application", "Unknown System"}
	}
}

// formatRelatedErrors formats related errors for the report
func (pmm *PostMortemManager) formatRelatedErrors(errors []models.EnhancedErrorResponse) string {
	if len(errors) == 0 {
		return "No related errors recorded"
	}

	result := ""
	for _, err := range errors {
		result += fmt.Sprintf(`
### Error: %s
- **Message**: %s
- **Category**: %s
- **Severity**: %s
- **Timestamp**: %s
- **Suggestions**: %s
`,
			err.Error,
			err.Message,
			err.Category,
			err.Severity,
			err.Timestamp.Format(time.RFC3339),
			strings.Join(err.Suggestions, ", "),
		)
	}

	return result
}

// loadIncidents loads incidents from storage
func (pmm *PostMortemManager) loadIncidents() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(pmm.IncidentsPath), 0755); err != nil {
		return fmt.Errorf("failed to create incidents directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(pmm.IncidentsPath); os.IsNotExist(err) {
		// Create empty incidents file
		pmm.Incidents = make([]Incident, 0)
		return pmm.saveIncidents()
	}

	// Read incidents file
	data, err := os.ReadFile(pmm.IncidentsPath)
	if err != nil {
		return fmt.Errorf("failed to read incidents file: %w", err)
	}

	if len(data) == 0 {
		pmm.Incidents = make([]Incident, 0)
		return nil
	}

	// Parse JSON
	if err := json.Unmarshal(data, &pmm.Incidents); err != nil {
		return fmt.Errorf("failed to parse incidents JSON: %w", err)
	}

	return nil
}

// saveIncidents saves incidents to storage
func (pmm *PostMortemManager) saveIncidents() error {
	data, err := json.MarshalIndent(pmm.Incidents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal incidents: %w", err)
	}

	if err := os.WriteFile(pmm.IncidentsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write incidents file: %w", err)
	}

	return nil
}

// generateIncidentID generates a unique incident ID
func generateIncidentID() string {
	return fmt.Sprintf("INC-%d", time.Now().Unix())
}

// AutoCreateIncidentFromError automatically creates an incident from a critical error
func (pmm *PostMortemManager) AutoCreateIncidentFromError(errorResp *models.EnhancedErrorResponse) {
	if !pmm.AutoCreate || errorResp.Severity != "critical" {
		return
	}

	title := fmt.Sprintf("Critical %s Error", errorResp.Category)
	description := fmt.Sprintf("Automated incident created from critical error: %s", errorResp.Message)

	incident, err := pmm.CreateIncident(errorResp, title, description)
	if err != nil {
		log.Printf("Failed to auto-create incident: %v", err)
		return
	}

	log.Printf("Auto-created incident from critical error: %s", incident.ID)
}

// GetIncidentStatistics returns statistics about incidents
func (pmm *PostMortemManager) GetIncidentStatistics() map[string]interface{} {
	incidents := pmm.GetAllIncidents()

	stats := map[string]interface{}{
		"total_incidents":     len(incidents),
		"by_status":           make(map[string]int),
		"by_severity":         make(map[string]int),
		"by_category":         make(map[string]int),
		"resolved_count":      0,
		"avg_resolution_time": "0s",
	}

	var totalResolutionTime time.Duration

	for _, incident := range incidents {
		// Count by status
		stats["by_status"].(map[string]int)[incident.Status]++

		// Count by severity
		stats["by_severity"].(map[string]int)[incident.Severity]++

		// Count by category (from related errors)
		if len(incident.RelatedErrors) > 0 {
			category := incident.RelatedErrors[0].Category
			stats["by_category"].(map[string]int)[category]++
		}

		// Track resolved incidents
		if incident.Status == "resolved" && incident.EndTime != nil {
			stats["resolved_count"] = stats["resolved_count"].(int) + 1
			totalResolutionTime += incident.EndTime.Sub(incident.StartTime)
		}
	}

	// Calculate average resolution time
	if resolvedCount := stats["resolved_count"].(int); resolvedCount > 0 {
		avgTime := totalResolutionTime / time.Duration(resolvedCount)
		stats["avg_resolution_time"] = avgTime.String()
	}

	return stats
}

// ExportIncidents exports incidents to various formats
func (pmm *PostMortemManager) ExportIncidents(format string) (string, error) {
	incidents := pmm.GetAllIncidents()

	switch format {
	case "json":
		data, err := json.MarshalIndent(incidents, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal incidents: %w", err)
		}
		return string(data), nil

	case "summary":
		return pmm.generateSummaryReport(incidents), nil

	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// generateSummaryReport generates a summary report of all incidents
func (pmm *PostMortemManager) generateSummaryReport(incidents []Incident) string {
	stats := pmm.GetIncidentStatistics()

	report := "# Incident Summary Report\n\n"
	report += "## Overview\n"
	report += fmt.Sprintf("- **Total Incidents**: %v\n", stats["total_incidents"])
	report += fmt.Sprintf("- **Resolved**: %v\n", stats["resolved_count"])
	report += fmt.Sprintf("- **Average Resolution Time**: %v\n", stats["avg_resolution_time"])

	report += "\n## By Status\n"
	for status, count := range stats["by_status"].(map[string]int) {
		report += fmt.Sprintf("- **%s**: %d\n", status, count)
	}

	report += "\n## By Severity\n"
	for severity, count := range stats["by_severity"].(map[string]int) {
		report += fmt.Sprintf("- **%s**: %d\n", severity, count)
	}

	report += "\n## Recent Incidents\n"
	recentCount := 0
	for _, incident := range incidents {
		if recentCount >= 10 { // Show last 10 incidents
			break
		}

		report += fmt.Sprintf("### %s (%s)\n", incident.Title, incident.ID)
		report += fmt.Sprintf("- **Status**: %s\n", incident.Status)
		report += fmt.Sprintf("- **Severity**: %s\n", incident.Severity)
		report += fmt.Sprintf("- **Duration**: %s\n", incident.Duration)
		report += "\n"
		recentCount++
	}

	return report
}

// Global post-mortem manager instance
var GlobalPostMortemManager = NewPostMortemManager("./incidents.json")

// Convenience functions for post-mortem management

// CreateIncident creates a new incident
func CreateIncident(errorResp *models.EnhancedErrorResponse, title, description string) (*Incident, error) {
	return GlobalPostMortemManager.CreateIncident(errorResp, title, description)
}

// UpdateIncident updates an existing incident
func UpdateIncident(id string, updates map[string]interface{}) error {
	return GlobalPostMortemManager.UpdateIncident(id, updates)
}

// GetIncident retrieves an incident by ID
func GetIncident(id string) (*Incident, error) {
	return GlobalPostMortemManager.GetIncident(id)
}

// GetAllIncidents returns all incidents
func GetAllIncidents() []Incident {
	return GlobalPostMortemManager.GetAllIncidents()
}

// GeneratePostMortemReport generates a post-mortem report
func GeneratePostMortemReport(id string) (string, error) {
	return GlobalPostMortemManager.GeneratePostMortemReport(id)
}

// GetIncidentStatistics returns incident statistics
func GetIncidentStatistics() map[string]interface{} {
	return GlobalPostMortemManager.GetIncidentStatistics()
}

// AutoCreateIncidentFromError automatically creates an incident from a critical error
func AutoCreateIncidentFromError(errorResp *models.EnhancedErrorResponse) {
	GlobalPostMortemManager.AutoCreateIncidentFromError(errorResp)
}
