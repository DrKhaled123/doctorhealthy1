package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"api-key-generator/internal/models"

	"github.com/google/uuid"
)

// HealthService handles health management and disease-related advice
type HealthService struct {
	db          *sql.DB
	userService *UserService
	dataLoader  *DataLoader
}

// NewHealthService creates a new health service
func NewHealthService(db *sql.DB, userService *UserService, dataLoader *DataLoader) *HealthService {
	return &HealthService{
		db:          db,
		userService: userService,
		dataLoader:  dataLoader,
	}
}

// GenerateHealthPlan generates a personalized health management plan
func (s *HealthService) GenerateHealthPlan(ctx context.Context, req *models.GenerateHealthPlanRequest) (*models.HealthPlan, error) {
	// Get user data
	user, err := s.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate treatment plan based on diseases and medications
	treatmentPlan := s.generateTreatmentPlan(req.Diseases, req.Medications)

	// Generate nutrition advice based on health conditions
	nutritionAdvice := s.generateNutritionAdvice(req.Diseases, req.Medications)

	// Generate lifestyle recommendations
	lifestyleChanges := s.generateLifestyleRecommendations(req.Diseases, req.Complaints)

	// Generate supplement recommendations
	supplements := s.generateSupplementRecommendations(req.Diseases, req.Complaints, user.Weight)

	// Create health plan
	plan := &models.HealthPlan{
		ID:               uuid.New().String(),
		UserID:           req.UserID,
		Diseases:         req.Diseases,
		Medications:      req.Medications,
		Complaints:       req.Complaints,
		TreatmentPlan:    treatmentPlan,
		NutritionAdvice:  nutritionAdvice,
		LifestyleChanges: lifestyleChanges,
		Supplements:      supplements,
		Disclaimer:       s.getHealthDisclaimer(user.Language),
		CreatedAt:        time.Now().UTC(),
	}

	// Save to database (optional - for tracking)
	err = s.savePlan(ctx, plan)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save health plan: %v\n", err)
	}

	return plan, nil
}

// generateTreatmentPlan generates treatment recommendations based on diseases
func (s *HealthService) generateTreatmentPlan(diseases, medications []string) []models.TreatmentItem {
	var treatmentPlan []models.TreatmentItem

	// Get disease information from knowledge base
	diseaseInfo := s.getDiseaseInformation()

	for _, disease := range diseases {
		if info, exists := diseaseInfo[strings.ToLower(disease)]; exists {
			treatmentPlan = append(treatmentPlan, info)
		}
	}

	// Add medication-specific advice
	if len(medications) > 0 {
		medicationAdvice := models.TreatmentItem{
			Category: "Medication Management",
			Advice: []string{
				"Take medications as prescribed by your healthcare provider",
				"Do not stop or change medications without consulting your doctor",
				"Keep a medication list and share it with all healthcare providers",
				"Be aware of potential drug interactions",
				"Monitor for side effects and report them to your doctor",
			},
			Precautions: []string{
				"Never share medications with others",
				"Store medications properly according to instructions",
				"Check expiration dates regularly",
				"Inform healthcare providers of all medications and supplements",
			},
		}
		treatmentPlan = append(treatmentPlan, medicationAdvice)
	}

	return treatmentPlan
}

// generateNutritionAdvice generates nutrition recommendations based on health conditions
func (s *HealthService) generateNutritionAdvice(diseases, medications []string) []string {
	var advice []string

	// General healthy eating advice
	advice = append(advice,
		"Follow a balanced diet rich in fruits, vegetables, whole grains, and lean proteins",
		"Stay hydrated by drinking adequate water throughout the day",
		"Limit processed foods, added sugars, and excessive sodium",
		"Practice portion control and mindful eating",
	)

	// Disease-specific nutrition advice
	for _, disease := range diseases {
		switch strings.ToLower(disease) {
		case "diabetes", "type 2 diabetes":
			advice = append(advice,
				"Monitor carbohydrate intake and choose complex carbohydrates",
				"Include fiber-rich foods to help control blood sugar",
				"Eat regular meals to maintain stable blood glucose levels",
				"Limit sugary drinks and refined carbohydrates",
			)
		case "hypertension", "high blood pressure":
			advice = append(advice,
				"Follow a low-sodium diet (less than 2300mg per day)",
				"Increase potassium-rich foods like bananas, spinach, and beans",
				"Limit alcohol consumption",
				"Choose lean proteins and reduce saturated fats",
			)
		case "heart disease", "cardiovascular disease":
			advice = append(advice,
				"Follow a heart-healthy diet low in saturated and trans fats",
				"Include omega-3 fatty acids from fish, nuts, and seeds",
				"Increase soluble fiber intake to help lower cholesterol",
				"Limit sodium and processed foods",
			)
		case "obesity", "overweight":
			advice = append(advice,
				"Create a moderate calorie deficit for gradual weight loss",
				"Focus on nutrient-dense, low-calorie foods",
				"Increase protein intake to preserve muscle mass",
				"Practice portion control and mindful eating habits",
			)
		}
	}

	return advice
}

// generateLifestyleRecommendations generates lifestyle change recommendations
func (s *HealthService) generateLifestyleRecommendations(diseases, complaints []string) []string {
	var recommendations []string

	// General lifestyle recommendations
	recommendations = append(recommendations,
		"Maintain a regular sleep schedule (7-9 hours per night)",
		"Engage in regular physical activity as approved by your healthcare provider",
		"Practice stress management techniques like meditation or deep breathing",
		"Avoid smoking and limit alcohol consumption",
		"Maintain social connections and seek support when needed",
	)

	// Disease-specific recommendations
	for _, disease := range diseases {
		switch strings.ToLower(disease) {
		case "diabetes":
			recommendations = append(recommendations,
				"Monitor blood glucose levels as recommended",
				"Maintain a consistent meal and exercise schedule",
				"Learn to recognize signs of high and low blood sugar",
				"Keep emergency glucose supplies available",
			)
		case "hypertension":
			recommendations = append(recommendations,
				"Monitor blood pressure regularly at home",
				"Practice relaxation techniques to manage stress",
				"Maintain a healthy weight",
				"Limit caffeine intake if sensitive",
			)
		case "arthritis":
			recommendations = append(recommendations,
				"Engage in low-impact exercises like swimming or walking",
				"Apply heat or cold therapy as needed for pain relief",
				"Maintain joint flexibility through gentle stretching",
				"Use ergonomic tools to reduce joint stress",
			)
		}
	}

	// Complaint-specific recommendations
	for _, complaint := range complaints {
		switch strings.ToLower(complaint) {
		case "fatigue", "low energy":
			recommendations = append(recommendations,
				"Establish a consistent sleep routine",
				"Take short breaks during the day",
				"Consider iron levels if experiencing persistent fatigue",
				"Limit caffeine late in the day",
			)
		case "stress", "anxiety":
			recommendations = append(recommendations,
				"Practice daily stress-reduction techniques",
				"Consider counseling or therapy if needed",
				"Limit exposure to stressful situations when possible",
				"Engage in regular physical activity to reduce stress",
			)
		}
	}

	return recommendations
}

// generateSupplementRecommendations generates supplement advice based on conditions and weight
func (s *HealthService) generateSupplementRecommendations(diseases, complaints []string, weight float64) []models.SupplementAdvice {
	var supplements []models.SupplementAdvice

	// General supplements for most adults
	supplements = append(supplements, models.SupplementAdvice{
		Name:         "Vitamin D3",
		Dosage:       "1000-2000 IU",
		Frequency:    "Daily",
		Instructions: "Take with a meal containing fat for better absorption",
		Benefits:     []string{"Bone health", "Immune function", "Mood support"},
	})

	supplements = append(supplements, models.SupplementAdvice{
		Name:         "Omega-3 Fatty Acids",
		Dosage:       "1000-2000 mg EPA/DHA",
		Frequency:    "Daily",
		Instructions: "Take with meals to reduce fishy aftertaste",
		Benefits:     []string{"Heart health", "Brain function", "Anti-inflammatory"},
	})

	// Disease-specific supplements
	for _, disease := range diseases {
		switch strings.ToLower(disease) {
		case "diabetes":
			supplements = append(supplements, models.SupplementAdvice{
				Name:         "Chromium",
				Dosage:       "200-400 mcg",
				Frequency:    "Daily",
				Instructions: "Take with meals",
				Benefits:     []string{"Blood sugar control", "Insulin sensitivity"},
			})
		case "hypertension":
			supplements = append(supplements, models.SupplementAdvice{
				Name:         "Magnesium",
				Dosage:       fmt.Sprintf("%.0f mg", weight*4), // 4mg per kg body weight
				Frequency:    "Daily",
				Instructions: "Take with food to avoid stomach upset",
				Benefits:     []string{"Blood pressure support", "Muscle function", "Heart health"},
			})
		case "arthritis":
			supplements = append(supplements, models.SupplementAdvice{
				Name:         "Glucosamine & Chondroitin",
				Dosage:       "1500mg Glucosamine + 1200mg Chondroitin",
				Frequency:    "Daily",
				Instructions: "Take with meals, effects may take 2-3 months",
				Benefits:     []string{"Joint health", "Cartilage support", "Pain reduction"},
			})
		}
	}

	// Complaint-specific supplements
	for _, complaint := range complaints {
		switch strings.ToLower(complaint) {
		case "fatigue", "low energy":
			supplements = append(supplements, models.SupplementAdvice{
				Name:         "B-Complex",
				Dosage:       "50-100 mg",
				Frequency:    "Daily",
				Instructions: "Take in the morning with breakfast",
				Benefits:     []string{"Energy metabolism", "Nervous system support", "Mood support"},
			})
		case "stress", "anxiety":
			supplements = append(supplements, models.SupplementAdvice{
				Name:         "Magnesium Glycinate",
				Dosage:       "200-400 mg",
				Frequency:    "Evening",
				Instructions: "Take 1-2 hours before bedtime",
				Benefits:     []string{"Stress reduction", "Better sleep", "Muscle relaxation"},
			})
		}
	}

	return supplements
}

// getDiseaseInformation returns treatment information for common diseases
func (s *HealthService) getDiseaseInformation() map[string]models.TreatmentItem {
	return map[string]models.TreatmentItem{
		"diabetes": {
			Category: "Diabetes Management",
			Advice: []string{
				"Monitor blood glucose levels regularly",
				"Follow a consistent meal plan",
				"Take medications as prescribed",
				"Stay physically active",
				"Maintain a healthy weight",
				"Get regular check-ups with your healthcare team",
			},
			Precautions: []string{
				"Watch for signs of high or low blood sugar",
				"Keep emergency glucose supplies available",
				"Inform all healthcare providers about your diabetes",
				"Check feet daily for cuts or sores",
				"Get regular eye and kidney function tests",
			},
		},
		"hypertension": {
			Category: "Blood Pressure Management",
			Advice: []string{
				"Take blood pressure medications as prescribed",
				"Monitor blood pressure at home",
				"Follow a low-sodium diet",
				"Maintain a healthy weight",
				"Exercise regularly",
				"Limit alcohol consumption",
				"Manage stress effectively",
			},
			Precautions: []string{
				"Don't stop medications without consulting your doctor",
				"Be aware of medication side effects",
				"Avoid sudden position changes to prevent dizziness",
				"Limit caffeine if sensitive",
			},
		},
		"heart disease": {
			Category: "Cardiovascular Health",
			Advice: []string{
				"Take prescribed medications consistently",
				"Follow a heart-healthy diet",
				"Exercise as recommended by your doctor",
				"Quit smoking if applicable",
				"Manage stress levels",
				"Get adequate sleep",
				"Attend regular cardiology appointments",
			},
			Precautions: []string{
				"Know the signs of heart attack and stroke",
				"Carry emergency medications if prescribed",
				"Avoid excessive physical exertion without clearance",
				"Monitor for medication side effects",
			},
		},
		"arthritis": {
			Category: "Joint Health Management",
			Advice: []string{
				"Stay physically active with low-impact exercises",
				"Maintain a healthy weight to reduce joint stress",
				"Use heat and cold therapy for pain relief",
				"Take anti-inflammatory medications as prescribed",
				"Practice joint protection techniques",
				"Consider physical therapy",
			},
			Precautions: []string{
				"Avoid activities that cause excessive joint pain",
				"Use assistive devices when needed",
				"Be aware of medication side effects",
				"Don't ignore persistent joint swelling",
			},
		},
	}
}

// savePlan saves the health plan to database
func (s *HealthService) savePlan(ctx context.Context, plan *models.HealthPlan) error {
	query := `
		INSERT INTO health_plans (
			id, user_id, disclaimer, created_at
		) VALUES (?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		plan.ID,
		plan.UserID,
		plan.Disclaimer,
		plan.CreatedAt,
	)

	return err
}

// getHealthDisclaimer returns appropriate health disclaimer based on language
func (s *HealthService) getHealthDisclaimer(language string) string {
	if language == "ar" {
		return `إخلاء المسؤولية الطبية: هذه المعلومات لأغراض تعليمية فقط وليست نصيحة طبية. 
لا تحل محل استشارة طبيب مؤهل. استشر دائماً مقدم الرعاية الصحية قبل إجراء تغييرات على العلاج أو الأدوية أو نمط الحياة. 
في حالات الطوارئ الطبية، اتصل بخدمات الطوارئ فوراً.`
	}

	return `MEDICAL DISCLAIMER: This information is for educational purposes only and is not intended as medical advice. 
It does not replace consultation with a qualified healthcare provider. Always consult your healthcare provider before making 
changes to treatment, medications, or lifestyle. In case of medical emergencies, contact emergency services immediately.`
}

// GetAvailableDiseases returns common diseases for selection
func (s *HealthService) GetAvailableDiseases() []string {
	return []string{
		"diabetes",
		"type_2_diabetes",
		"hypertension",
		"high_blood_pressure",
		"heart_disease",
		"cardiovascular_disease",
		"arthritis",
		"osteoarthritis",
		"rheumatoid_arthritis",
		"obesity",
		"overweight",
		"high_cholesterol",
		"thyroid_disorders",
		"asthma",
		"copd",
		"depression",
		"anxiety",
		"osteoporosis",
		"kidney_disease",
		"liver_disease",
	}
}

// GetAvailableComplaints returns common health complaints
func (s *HealthService) GetAvailableComplaints() []string {
	return []string{
		"fatigue",
		"low_energy",
		"joint_pain",
		"muscle_pain",
		"headaches",
		"sleep_problems",
		"stress",
		"anxiety",
		"digestive_issues",
		"weight_gain",
		"weight_loss",
		"mood_changes",
		"memory_problems",
		"concentration_issues",
		"skin_problems",
		"hair_loss",
		"frequent_infections",
	}
}
