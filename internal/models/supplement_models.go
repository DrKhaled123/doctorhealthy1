package models

import (
	"time"
)

// SupplementProtocol represents a comprehensive supplement recommendation
type SupplementProtocol struct {
	ID                string             `json:"id" db:"id"`
	UserID            string             `json:"user_id" db:"user_id"`
	WorkoutProgramID  string             `json:"workout_program_id,omitempty" db:"workout_program_id"`
	Name              string             `json:"name" db:"name"`
	Category          string             `json:"category" db:"category"`
	Purpose           string             `json:"purpose" db:"purpose"`
	PreWorkout        []SupplementTiming `json:"pre_workout,omitempty"`
	DuringWorkout     []SupplementTiming `json:"during_workout,omitempty"`
	PostWorkout       []SupplementTiming `json:"post_workout,omitempty"`
	Daily             []SupplementTiming `json:"daily,omitempty"`
	SafetyNotes       []string           `json:"safety_notes"`
	Contraindications []string           `json:"contraindications"`
	Interactions      []DrugInteraction  `json:"interactions"`
	References        []string           `json:"references"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`
}

// SupplementTiming represents specific supplement timing and dosage
type SupplementTiming struct {
	SupplementName string   `json:"supplement_name"`
	Dosage         string   `json:"dosage"`
	Unit           string   `json:"unit"`
	TimingMinutes  int      `json:"timing_minutes"`
	Instructions   string   `json:"instructions"`
	Benefits       []string `json:"benefits"`
	SideEffects    []string `json:"side_effects,omitempty"`
	MaxDailyDose   string   `json:"max_daily_dose,omitempty"`
}

// DrugInteraction represents supplement-drug interactions
type DrugInteraction struct {
	DrugName    string `json:"drug_name"`
	Interaction string `json:"interaction"`
	Severity    string `json:"severity"`
	Management  string `json:"management"`
}

// SupplementRequest represents a request for supplement recommendations
type SupplementRequest struct {
	UserID           string   `json:"user_id" validate:"required"`
	WorkoutGoal      string   `json:"workout_goal" validate:"required"`
	WorkoutType      string   `json:"workout_type" validate:"oneof=gym home"`
	Duration         int      `json:"duration_minutes" validate:"min=15,max=180"`
	Intensity        string   `json:"intensity" validate:"oneof=low moderate high"`
	ExistingMeds     []string `json:"existing_medications"`
	Allergies        []string `json:"allergies"`
	HealthConditions []string `json:"health_conditions"`
}

// SupplementCategory represents supplement categories
type SupplementCategory struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Supplements []string `json:"supplements"`
}

// Available supplement categories
var SupplementCategories = []SupplementCategory{
	{
		Name:        "pre_workout",
		Description: "Supplements to enhance workout performance",
		Supplements: []string{"caffeine", "creatine", "beta_alanine", "citrulline"},
	},
	{
		Name:        "post_workout",
		Description: "Supplements for recovery and muscle building",
		Supplements: []string{"whey_protein", "casein", "glutamine", "bcaa"},
	},
	{
		Name:        "daily_health",
		Description: "Daily supplements for overall health",
		Supplements: []string{"multivitamin", "omega3", "vitamin_d", "magnesium"},
	},
}
