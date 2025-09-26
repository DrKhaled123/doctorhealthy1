package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"api-key-generator/internal/models"

	"github.com/google/uuid"
)

// UserService handles user operations
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// CreateUser creates a new user profile
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Input validation
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create user model
	user := &models.User{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Email:            req.Email,
		Age:              req.Age,
		Weight:           req.Weight,
		Height:           req.Height,
		Gender:           req.Gender,
		ActivityLevel:    req.ActivityLevel,
		MetabolicRate:    req.MetabolicRate,
		Goal:             req.Goal,
		FoodDislikes:     req.FoodDislikes,
		Allergies:        req.Allergies,
		Diseases:         req.Diseases,
		Medications:      req.Medications,
		PreferredCuisine: req.PreferredCuisine,
		Language:         req.Language,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Set default language if not provided
	if user.Language == "" {
		user.Language = "en"
	}

	// Serialize arrays to JSON
	foodDislikesJSON, err := json.Marshal(user.FoodDislikes)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize food dislikes: %w", err)
	}

	allergiesJSON, err := json.Marshal(user.Allergies)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize allergies: %w", err)
	}

	diseasesJSON, err := json.Marshal(user.Diseases)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize diseases: %w", err)
	}

	medicationsJSON, err := json.Marshal(user.Medications)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize medications: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO users (
			id, name, email, age, weight, height, gender, activity_level, 
			metabolic_rate, goal, food_dislikes, allergies, diseases, 
			medications, preferred_cuisine, language, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Age,
		user.Weight,
		user.Height,
		user.Gender,
		user.ActivityLevel,
		user.MetabolicRate,
		user.Goal,
		string(foodDislikesJSON),
		string(allergiesJSON),
		string(diseasesJSON),
		string(medicationsJSON),
		user.PreferredCuisine,
		user.Language,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	query := `
		SELECT id, name, email, age, weight, height, gender, activity_level,
			   metabolic_rate, goal, food_dislikes, allergies, diseases,
			   medications, preferred_cuisine, language, created_at, updated_at
		FROM users 
		WHERE id = ?
	`

	var user models.User
	var foodDislikesJSON, allergiesJSON, diseasesJSON, medicationsJSON string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Weight,
		&user.Height,
		&user.Gender,
		&user.ActivityLevel,
		&user.MetabolicRate,
		&user.Goal,
		&foodDislikesJSON,
		&allergiesJSON,
		&diseasesJSON,
		&medicationsJSON,
		&user.PreferredCuisine,
		&user.Language,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Deserialize JSON arrays
	if err := json.Unmarshal([]byte(foodDislikesJSON), &user.FoodDislikes); err != nil {
		return nil, fmt.Errorf("failed to deserialize food dislikes: %w", err)
	}
	if err := json.Unmarshal([]byte(allergiesJSON), &user.Allergies); err != nil {
		return nil, fmt.Errorf("failed to deserialize allergies: %w", err)
	}
	if err := json.Unmarshal([]byte(diseasesJSON), &user.Diseases); err != nil {
		return nil, fmt.Errorf("failed to deserialize diseases: %w", err)
	}
	if err := json.Unmarshal([]byte(medicationsJSON), &user.Medications); err != nil {
		return nil, fmt.Errorf("failed to deserialize medications: %w", err)
	}

	return &user, nil
}

// CalculateCalories calculates daily calorie needs based on user data
func (s *UserService) CalculateCalories(ctx context.Context, user *models.User) (*models.CalorieCalculation, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}

	// Calculate BMI
	heightInMeters := user.Height / 100
	bmi := user.Weight / (heightInMeters * heightInMeters)

	// Calculate BMR using Mifflin-St Jeor Equation
	var bmr float64
	if user.Gender == "male" {
		bmr = (10 * user.Weight) + (6.25 * user.Height) - (5 * float64(user.Age)) + 5
	} else {
		bmr = (10 * user.Weight) + (6.25 * user.Height) - (5 * float64(user.Age)) - 161
	}

	// Activity multipliers
	activityMultipliers := map[string]float64{
		"sedentary":   1.2,
		"light":       1.375,
		"moderate":    1.55,
		"active":      1.725,
		"very_active": 1.9,
	}

	multiplier, exists := activityMultipliers[user.ActivityLevel]
	if !exists {
		multiplier = 1.55 // default to moderate
	}

	// Calculate TDEE
	tdee := bmr * multiplier

	// Determine calories per kg based on BMI and goals
	var caloriesPerKg int
	var method string
	var explanation string

	switch {
	case bmi >= 18 && bmi <= 30:
		caloriesPerKg = 20
		method = "20 calories per kg"
		explanation = "Suitable for weight loss or maintenance (BMI 18-30)"
	case bmi >= 15 && bmi <= 17:
		caloriesPerKg = 25
		method = "25 calories per kg"
		explanation = "For thin individuals wanting to maintain or gain weight (BMI 15-17)"
	case user.MetabolicRate == "high" || user.Goal == "build_muscle" || user.Goal == "improve_strength":
		caloriesPerKg = 30
		method = "30 calories per kg"
		explanation = "For high metabolic rate or muscle building goals"
	default:
		caloriesPerKg = 25
		method = "25 calories per kg"
		explanation = "Default calculation for your profile"
	}

	// Adjust based on goal
	recommendedCalories := int(user.Weight * float64(caloriesPerKg))

	switch user.Goal {
	case "lose_weight":
		recommendedCalories = int(float64(recommendedCalories) * 0.85) // 15% deficit
	case "gain_weight", "build_muscle":
		recommendedCalories = int(float64(recommendedCalories) * 1.15) // 15% surplus
	}

	// Calculate macros
	var proteinGrams float64
	if user.ActivityLevel == "active" || user.ActivityLevel == "very_active" {
		proteinGrams = user.Weight * 1.7 // 1.5-1.7g per kg for active individuals
	} else {
		proteinGrams = user.Weight * 1.2 // 1-1.5g per kg for sedentary
	}

	// Standard macro distribution (can be adjusted based on plan type)
	proteinCalories := proteinGrams * 4
	fatCalories := float64(recommendedCalories) * 0.25 // 25% from fats
	carbCalories := float64(recommendedCalories) - proteinCalories - fatCalories

	fatsGrams := fatCalories / 9
	carbsGrams := carbCalories / 4

	return &models.CalorieCalculation{
		BMI:                 bmi,
		BMR:                 bmr,
		TDEE:                tdee,
		CaloriesPerKg:       caloriesPerKg,
		RecommendedCalories: recommendedCalories,
		ProteinGrams:        proteinGrams,
		CarbsGrams:          carbsGrams,
		FatsGrams:           fatsGrams,
		Method:              method,
		Explanation:         explanation,
	}, nil
}

// Private helper methods

func (s *UserService) validateCreateUserRequest(req *models.CreateUserRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Age < 13 || req.Age > 120 {
		return fmt.Errorf("age must be between 13 and 120")
	}
	if req.Weight < 30 || req.Weight > 300 {
		return fmt.Errorf("weight must be between 30 and 300 kg")
	}
	if req.Height < 100 || req.Height > 250 {
		return fmt.Errorf("height must be between 100 and 250 cm")
	}
	if req.Gender != "male" && req.Gender != "female" {
		return fmt.Errorf("gender must be 'male' or 'female'")
	}

	validActivityLevels := map[string]bool{
		"sedentary": true, "light": true, "moderate": true, "active": true, "very_active": true,
	}
	if !validActivityLevels[req.ActivityLevel] {
		return fmt.Errorf("invalid activity level")
	}

	validMetabolicRates := map[string]bool{
		"low": true, "medium": true, "high": true,
	}
	if !validMetabolicRates[req.MetabolicRate] {
		return fmt.Errorf("invalid metabolic rate")
	}

	if req.Goal == "" {
		return fmt.Errorf("goal is required")
	}

	return nil
}
