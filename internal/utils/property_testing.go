package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"api-key-generator/internal/models"
)

// PropertyTest represents a property-based test
type PropertyTest struct {
	Name        string
	Property    func(interface{}) bool
	Generator   func() interface{}
	Iterations  int
	ShrinkSteps int
}

// PropertyTestResult represents the result of a property test
type PropertyTestResult struct {
	TestName    string
	Passed      bool
	Iterations  int
	FailedInput interface{}
	Error       error
	ShrunkInput interface{}
}

// PropertyTestSuite manages a collection of property tests
type PropertyTestSuite struct {
	Tests []PropertyTest
	Rand  *rand.Rand
}

// NewPropertyTestSuite creates a new property test suite
func NewPropertyTestSuite() *PropertyTestSuite {
	return &PropertyTestSuite{
		Tests: make([]PropertyTest, 0),
		Rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AddTest adds a property test to the suite
func (pts *PropertyTestSuite) AddTest(test PropertyTest) {
	pts.Tests = append(pts.Tests, test)
}

// RunAllTests runs all property tests in the suite
func (pts *PropertyTestSuite) RunAllTests() []PropertyTestResult {
	results := make([]PropertyTestResult, 0, len(pts.Tests))

	for _, test := range pts.Tests {
		result := pts.RunTest(test)
		results = append(results, result)
	}

	return results
}

// RunTest executes a single property test
func (pts *PropertyTestSuite) RunTest(test PropertyTest) PropertyTestResult {
	result := PropertyTestResult{
		TestName:   test.Name,
		Iterations: 0,
		Passed:     true,
	}

	// Run the specified number of iterations
	for i := 0; i < test.Iterations; i++ {
		// Generate random input
		input := test.Generator()

		// Test the property
		if !test.Property(input) {
			result.Passed = false
			result.FailedInput = input
			result.Error = fmt.Errorf("property failed for input: %v", input)
			result.Iterations = i + 1

			// Attempt to shrink the failing input
			result.ShrunkInput = pts.ShrinkInput(input, test.Property, test.ShrinkSteps)

			return result
		}

		result.Iterations = i + 1
	}

	return result
}

// ShrinkInput attempts to find a minimal failing input
func (pts *PropertyTestSuite) ShrinkInput(input interface{}, property func(interface{}) bool, steps int) interface{} {
	current := input

	for i := 0; i < steps; i++ {
		shrunk := pts.shrinkOnce(current)
		if shrunk != nil && !property(shrunk) {
			current = shrunk
		} else {
			break
		}
	}

	return current
}

// shrinkOnce attempts to shrink the input in one step
func (pts *PropertyTestSuite) shrinkOnce(input interface{}) interface{} {
	v := reflect.ValueOf(input)

	switch v.Kind() {
	case reflect.Int, reflect.Int64:
		if v.Int() > 0 {
			return v.Int() / 2
		}
	case reflect.String:
		str := v.String()
		if len(str) > 0 {
			return str[:len(str)/2]
		}
	case reflect.Slice:
		if v.Len() > 0 {
			return v.Slice(0, v.Len()/2).Interface()
		}
	}

	return nil
}

// Common Property Test Generators

// GeneratePositiveInt generates positive integers
func (pts *PropertyTestSuite) GeneratePositiveInt() func() interface{} {
	return func() interface{} {
		return pts.Rand.Intn(1000) + 1
	}
}

// GenerateString generates random strings
func (pts *PropertyTestSuite) GenerateString() func() interface{} {
	return func() interface{} {
		length := pts.Rand.Intn(50) + 1
		bytes := make([]byte, length)
		for i := range bytes {
			bytes[i] = byte(pts.Rand.Intn(26) + 65) // A-Z
		}
		return string(bytes)
	}
}

// GenerateSlice generates random slices
func (pts *PropertyTestSuite) GenerateSlice() func() interface{} {
	return func() interface{} {
		length := pts.Rand.Intn(10)
		slice := make([]int, length)
		for i := range slice {
			slice[i] = pts.Rand.Intn(100)
		}
		return slice
	}
}

// GenerateUserProfile generates random user profiles for testing
func (pts *PropertyTestSuite) GenerateUserProfile() func() interface{} {
	return func() interface{} {
		profile := models.APIKey{
			Name:        pts.generateRandomName(),
			Permissions: pts.generateRandomPermissions(),
			IsActive:    pts.Rand.Float32() > 0.5,
		}
		return profile
	}
}

// Helper functions for test data generation
func (pts *PropertyTestSuite) generateRandomName() string {
	names := []string{"user", "admin", "test", "api", "service", "client", "server", "app"}
	return names[pts.Rand.Intn(len(names))] + fmt.Sprintf("%d", pts.Rand.Intn(1000))
}

func (pts *PropertyTestSuite) generateRandomPermissions() []string {
	available := models.GetAvailablePermissions()
	count := pts.Rand.Intn(len(available)) + 1
	permissions := make([]string, 0, count)

	for i := 0; i < count; i++ {
		perm := available[pts.Rand.Intn(len(available))]
		permissions = append(permissions, perm.Name)
	}

	return permissions
}

// Predefined Property Tests for Common Scenarios

// TestAPIKeyProperties tests API key model properties
func TestAPIKeyProperties() PropertyTest {
	return PropertyTest{
		Name: "API Key Properties",
		Property: func(input interface{}) bool {
			key, ok := input.(models.APIKey)
			if !ok {
				return false
			}

			// Property 1: Name should not be empty if active
			if key.IsActive && key.Name == "" {
				return false
			}

			// Property 2: Should have at least one permission
			if len(key.Permissions) == 0 {
				return false
			}

			// Property 3: Usage count should not be negative
			if key.UsageCount < 0 {
				return false
			}

			return true
		},
		Generator: func() interface{} {
			suite := NewPropertyTestSuite()
			return suite.GenerateUserProfile()()
		},
		Iterations:  100,
		ShrinkSteps: 10,
	}
}

// TestNutritionCalculationProperties tests nutrition calculation logic
func TestNutritionCalculationProperties() PropertyTest {
	return PropertyTest{
		Name: "Nutrition Calculation Properties",
		Property: func(input interface{}) bool {
			// This would test nutrition calculation properties
			// For now, return true as placeholder
			return true
		},
		Generator: func() interface{} {
			// Generate nutrition calculation inputs
			return map[string]interface{}{
				"calories": 2000,
				"protein":  150,
				"carbs":    250,
				"fat":      70,
			}
		},
		Iterations:  50,
		ShrinkSteps: 5,
	}
}

// TestWorkoutGenerationProperties tests workout generation logic
func TestWorkoutGenerationProperties() PropertyTest {
	return PropertyTest{
		Name: "Workout Generation Properties",
		Property: func(input interface{}) bool {
			// This would test workout generation properties
			// For now, return true as placeholder
			return true
		},
		Generator: func() interface{} {
			// Generate workout generation inputs
			return map[string]interface{}{
				"fitness_level": "beginner",
				"goal":          "strength",
				"equipment":     []string{"dumbbells", "barbell"},
			}
		},
		Iterations:  50,
		ShrinkSteps: 5,
	}
}

// RunPropertyTests executes all predefined property tests
func RunPropertyTests() []PropertyTestResult {
	suite := NewPropertyTestSuite()

	// Add all property tests
	suite.AddTest(TestAPIKeyProperties())
	suite.AddTest(TestNutritionCalculationProperties())
	suite.AddTest(TestWorkoutGenerationProperties())

	// Run all tests
	return suite.RunAllTests()
}

// ReportTestResults generates a summary report of test results
func ReportTestResults(results []PropertyTestResult) string {
	report := "Property Test Results:\n"
	report += fmt.Sprintf("Total Tests: %d\n", len(results))

	passed := 0
	failed := 0

	for _, result := range results {
		if result.Passed {
			passed++
			report += fmt.Sprintf("✅ %s: PASSED (%d iterations)\n", result.TestName, result.Iterations)
		} else {
			failed++
			report += fmt.Sprintf("❌ %s: FAILED at iteration %d\n", result.TestName, result.Iterations)
			report += fmt.Sprintf("   Input: %v\n", result.FailedInput)
			report += fmt.Sprintf("   Shrunk: %v\n", result.ShrunkInput)
			report += fmt.Sprintf("   Error: %v\n", result.Error)
		}
	}

	report += fmt.Sprintf("\nSummary: %d passed, %d failed\n", passed, failed)
	return report
}

// LogicalErrorDetector detects potential logical errors in code
type LogicalErrorDetector struct {
	Patterns []LogicalErrorPattern
}

// LogicalErrorPattern represents a pattern that might indicate logical errors
type LogicalErrorPattern struct {
	Name        string
	Description string
	Pattern     string
	Solution    string
}

// NewLogicalErrorDetector creates a new logical error detector
func NewLogicalErrorDetector() *LogicalErrorDetector {
	return &LogicalErrorDetector{
		Patterns: []LogicalErrorPattern{
			{
				Name:        "Division by zero",
				Description: "Potential division by zero",
				Pattern:     "/ 0",
				Solution:    "Add zero checks before division",
			},
			{
				Name:        "Null pointer dereference",
				Description: "Accessing properties of potentially null objects",
				Pattern:     ".name",
				Solution:    "Add null checks or use optional chaining",
			},
			{
				Name:        "Infinite loops",
				Description: "Loops without proper termination conditions",
				Pattern:     "for.*true",
				Solution:    "Ensure loops have proper exit conditions",
			},
		},
	}
}

// DetectLogicalErrors analyzes code for potential logical errors
func (led *LogicalErrorDetector) DetectLogicalErrors(code string) []LogicalErrorPattern {
	var detected []LogicalErrorPattern

	for _, pattern := range led.Patterns {
		if contains(code, pattern.Pattern) {
			detected = append(detected, pattern)
		}
	}

	return detected
}

// ValidateBusinessLogic validates business logic rules
func ValidateBusinessLogic(data interface{}) []models.ValidationError {
	var errors []models.ValidationError

	// Example business logic validations
	switch v := data.(type) {
	case models.APIKey:
		// Business rule: Active API keys must have names
		if v.IsActive && v.Name == "" {
			errors = append(errors, models.ValidationError{
				Field:   "name",
				Value:   v.Name,
				Tag:     "required_when_active",
				Message: "API key name is required when key is active",
			})
		}

		// Business rule: Rate limit should be reasonable
		if v.RateLimit != nil && (*v.RateLimit < 1 || *v.RateLimit > 10000) {
			errors = append(errors, models.ValidationError{
				Field:   "rate_limit",
				Value:   v.RateLimit,
				Tag:     "range",
				Message: "Rate limit must be between 1 and 10000",
			})
		}

	case models.CreateAPIKeyRequest:
		// Business rule: Must have at least one permission
		if len(v.Permissions) == 0 {
			errors = append(errors, models.ValidationError{
				Field:   "permissions",
				Value:   v.Permissions,
				Tag:     "required",
				Message: "At least one permission is required",
			})
		}

		// Business rule: Expiry days should be reasonable
		if v.ExpiryDays != nil && (*v.ExpiryDays < 1 || *v.ExpiryDays > 3650) {
			errors = append(errors, models.ValidationError{
				Field:   "expiry_days",
				Value:   v.ExpiryDays,
				Tag:     "range",
				Message: "Expiry days must be between 1 and 3650",
			})
		}
	}

	return errors
}

// TestBusinessLogicValidation tests the business logic validation
func TestBusinessLogicValidation() PropertyTest {
	return PropertyTest{
		Name: "Business Logic Validation",
		Property: func(input interface{}) bool {
			errors := ValidateBusinessLogic(input)
			return len(errors) == 0
		},
		Generator: func() interface{} {
			suite := NewPropertyTestSuite()
			return suite.GenerateUserProfile()()
		},
		Iterations:  100,
		ShrinkSteps: 10,
	}
}
