package services

import (
	"os"
	"testing"
)

func TestComplaintsExtractor_ValidateSource(t *testing.T) {
	extractor := NewComplaintsExtractor()

	tests := []struct {
		name        string
		source      string
		expectError bool
	}{
		// This subtest depends on local file; skip if missing
		// { name: "Valid JSON source", source: "vip json/complaints.js", expectError: false },
		{
			name:        "Non-existent file",
			source:      "non-existent-file.json",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.source == "vip json/complaints.js" {
				if _, err := os.Stat(tt.source); os.IsNotExist(err) {
					t.Skip("Source file not present; skipping")
				}
			}
			err := extractor.ValidateSource(tt.source)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateSource() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestComplaintsExtractor_ExtractComplaints(t *testing.T) {
	extractor := NewComplaintsExtractor()

	// Check if the source file exists
	if _, err := os.Stat("vip json/complaints.js"); os.IsNotExist(err) {
		t.Skip("Source file 'vip json/complaints.js' not found, skipping test")
	}

	complaints, err := extractor.ExtractComplaints("vip json/complaints.js")
	if err != nil {
		t.Fatalf("ExtractComplaints() error = %v", err)
	}

	if len(complaints) == 0 {
		t.Error("Expected to extract at least one complaint, got 0")
	}

	// Validate first complaint structure
	if len(complaints) > 0 {
		complaint := complaints[0]

		if complaint.ID <= 0 {
			t.Error("Expected valid complaint ID")
		}

		if complaint.ConditionEN == "" {
			t.Error("Expected non-empty English condition")
		}

		if complaint.ConditionAR == "" {
			t.Error("Expected non-empty Arabic condition")
		}

		if complaint.Recommendations.Nutrition.EN == "" {
			t.Error("Expected non-empty nutrition recommendations")
		}
	}

	// Check extraction stats
	stats := extractor.GetExtractionStats()
	if stats.ProcessedCases == 0 {
		t.Error("Expected processed cases > 0")
	}

	if stats.ProcessedCases != len(complaints) {
		t.Errorf("Expected processed cases (%d) to match complaints length (%d)",
			stats.ProcessedCases, len(complaints))
	}

	t.Logf("Successfully extracted %d complaints", len(complaints))
	t.Logf("Extraction stats: %+v", stats)
}

func TestComplaintsExtractor_ExtractReferences(t *testing.T) {
	extractor := NewComplaintsExtractor()

	testContent := BilingualContentWithRefs{
		EN: "Some text with reference [Ref: https://www.ncbi.nlm.nih.gov/pmc/articles/PMC6353095/] and more text",
		AR: "نص عربي مع مرجع [Ref: https://pubmed.ncbi.nlm.nih.gov/29172425/] ونص إضافي",
	}

	references := extractor.ExtractReferences(testContent)

	if len(references) != 2 {
		t.Errorf("Expected 2 references, got %d", len(references))
	}

	expectedRefs := []string{
		"https://www.ncbi.nlm.nih.gov/pmc/articles/PMC6353095/",
		"https://pubmed.ncbi.nlm.nih.gov/29172425/",
	}

	for i, expected := range expectedRefs {
		if i >= len(references) || references[i] != expected {
			t.Errorf("Expected reference %d to be %s, got %s", i, expected, references[i])
		}
	}
}

func TestComplaintsExtractor_ParseMedicationProtocols(t *testing.T) {
	extractor := NewComplaintsExtractor()

	testContent := BilingualContent{
		EN: "Metformin: Start 500 mg/day, increase to 1500-2000 mg/day (divided doses) [Medical supervision required]",
		AR: "ميتفورمين: ابدأ 500 ملغ/يوم، زيادة إلى 1500-2000 ملغ/يوم (جرعات مقسمة) [يتطلب إشرافًا طبيًا]",
	}

	protocols, err := extractor.ParseMedicationProtocols(testContent)
	if err != nil {
		t.Fatalf("ParseMedicationProtocols() error = %v", err)
	}

	if len(protocols) == 0 {
		t.Error("Expected at least one medication protocol")
	}

	if len(protocols) > 0 {
		protocol := protocols[0]
		if protocol.Name == "" {
			t.Error("Expected non-empty medication name")
		}
		if !protocol.SupervisionRequired {
			t.Error("Expected supervision required to be true")
		}
	}
}

func TestComplaintsExtractor_ParseSupplementProtocols(t *testing.T) {
	extractor := NewComplaintsExtractor()

	testContent := BilingualContent{
		EN: "Chromium Picolinate: 200-1000 mcg/day. Berberine: 500-1500 mg/day. Probiotics (specific strains)",
		AR: "كروم بيكولينات: 200-1000 ميكروغرام/يوم. بربرين: 500-1500 ملغ/يوم. البروبيوتيك (سلالات محددة)",
	}

	protocols, err := extractor.ParseSupplementProtocols(testContent)
	if err != nil {
		t.Fatalf("ParseSupplementProtocols() error = %v", err)
	}

	if len(protocols) == 0 {
		t.Error("Expected at least one supplement protocol")
	}

	// Check that we extracted some supplements
	stats := extractor.GetExtractionStats()
	if stats.ExtractedSupplements == 0 {
		t.Error("Expected extracted supplements count > 0")
	}
}
