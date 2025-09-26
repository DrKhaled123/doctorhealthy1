package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

// SetsRepsHandler handles sets, reps, and rest period management endpoints
type SetsRepsHandler struct {
	setsRepsManager *services.SetsRepsManager
}

// NewSetsRepsHandler creates a new sets/reps handler
func NewSetsRepsHandler() *SetsRepsHandler {
	return &SetsRepsHandler{
		setsRepsManager: services.NewSetsRepsManager(),
	}
}

// GetExerciseRecommendations generates dynamic sets/reps recommendations
// POST /api/v1/workouts/recommendations
func (h *SetsRepsHandler) GetExerciseRecommendations(c echo.Context) error {
	var req services.ExerciseRecommendationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	recommendations, err := h.setsRepsManager.GenerateExerciseRecommendations(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate recommendations")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    recommendations,
	})
}

// TrackWorkoutSession records workout performance and progress
// POST /api/v1/workouts/sessions
func (h *SetsRepsHandler) TrackWorkoutSession(c echo.Context) error {
	var req services.WorkoutSessionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := h.setsRepsManager.TrackWorkoutSession(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to track workout session")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "Workout session tracked successfully",
	})
}

// GetProgressReport generates comprehensive progress analytics
// GET /api/v1/workouts/progress/{userID}/{programID}
func (h *SetsRepsHandler) GetProgressReport(c echo.Context) error {
	userID := c.Param("userID")
	programID := c.Param("programID")

	if userID == "" || programID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID and Program ID are required")
	}

	// Parse optional time range parameters
	timeRange := services.TimeRange{
		StartDate: time.Now().AddDate(0, -1, 0), // Default: last month
		EndDate:   time.Now(),
	}

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			timeRange.StartDate = startDate
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			timeRange.EndDate = endDate
		}
	}

	report, err := h.setsRepsManager.GenerateProgressReport(c.Request().Context(), userID, programID, timeRange)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate progress report")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

// GetRestPeriodRecommendation calculates optimal rest period for an exercise
// GET /api/v1/workouts/rest-period
func (h *SetsRepsHandler) GetRestPeriodRecommendation(c echo.Context) error {
	exerciseType := c.QueryParam("exercise_type")
	intensityStr := c.QueryParam("intensity")
	goals := c.QueryParams()["goals"]

	if exerciseType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Exercise type is required")
	}

	intensity := 70.0 // default
	if intensityStr != "" {
		if parsed, err := strconv.ParseFloat(intensityStr, 64); err == nil {
			intensity = parsed
		}
	}

	// Create a temporary progression engine to calculate rest period
	progressionEngine := services.NewProgressionEngine()
	restSeconds := progressionEngine.CalculateRestPeriod(exerciseType, intensity, goals)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"exercise_type": exerciseType,
			"intensity":     intensity,
			"goals":         goals,
			"rest_seconds":  restSeconds,
			"rest_display":  formatRestDisplay(restSeconds),
		},
	})
}

// GetSetsRepsRecommendation generates sets/reps recommendation for specific parameters
// GET /api/v1/workouts/sets-reps
func (h *SetsRepsHandler) GetSetsRepsRecommendation(c echo.Context) error {
	goals := c.QueryParams()["goals"]
	exerciseType := c.QueryParam("exercise_type")
	fitnessLevel := c.QueryParam("fitness_level")
	weekNumberStr := c.QueryParam("week_number")

	if len(goals) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one goal is required")
	}

	if exerciseType == "" {
		exerciseType = "isolation" // default
	}

	if fitnessLevel == "" {
		fitnessLevel = "intermediate" // default
	}

	weekNumber := 1 // default
	if weekNumberStr != "" {
		if parsed, err := strconv.Atoi(weekNumberStr); err == nil && parsed > 0 {
			weekNumber = parsed
		}
	}

	// Create a temporary progression engine to generate recommendation
	progressionEngine := services.NewProgressionEngine()
	recommendation, err := progressionEngine.GenerateSetsRepsRecommendation(
		c.Request().Context(),
		goals,
		exerciseType,
		fitnessLevel,
		weekNumber,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate recommendation")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"goals":         goals,
			"exercise_type": exerciseType,
			"fitness_level": fitnessLevel,
			"week_number":   weekNumber,
			"sets":          recommendation.Sets,
			"reps_min":      recommendation.RepsMin,
			"reps_max":      recommendation.RepsMax,
			"reps_display":  formatRepsDisplay(recommendation.RepsMin, recommendation.RepsMax),
			"rest_seconds":  recommendation.RestSeconds,
			"rest_display":  formatRestDisplay(recommendation.RestSeconds),
			"intensity":     recommendation.Intensity,
			"notes":         recommendation.Notes,
		},
	})
}

// GetProgressionAdjustment suggests progression adjustments based on performance
// POST /api/v1/workouts/progression-adjustment
func (h *SetsRepsHandler) GetProgressionAdjustment(c echo.Context) error {
	var req struct {
		ExerciseProgress      *ExerciseProgressData       `json:"exercise_progress"`
		CurrentRecommendation *SetsRepsRecommendationData `json:"current_recommendation"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if req.ExerciseProgress == nil || req.CurrentRecommendation == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Exercise progress and current recommendation are required")
	}

	// This is a simplified example - in reality you'd convert the request data properly
	adjustment := map[string]interface{}{
		"suggested_intensity": req.CurrentRecommendation.Intensity * 1.05, // 5% increase example
		"notes":               "Increase intensity based on performance",
		"next_week_sets":      req.CurrentRecommendation.Sets,
		"next_week_reps":      fmt.Sprintf("%d-%d", req.CurrentRecommendation.RepsMin, req.CurrentRecommendation.RepsMax),
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    adjustment,
	})
}

// Helper functions

func formatRestDisplay(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60

	if remainingSeconds == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

func formatRepsDisplay(min, max int) string {
	if min == max {
		return strconv.Itoa(min)
	}
	return fmt.Sprintf("%d-%d", min, max)
}

// Request/Response helper types for the API

type ExerciseProgressData struct {
	ExerciseID    string    `json:"exercise_id"`
	Weight        float64   `json:"weight"`
	Reps          int       `json:"reps"`
	Sets          int       `json:"sets"`
	RPE           int       `json:"rpe"`
	LastPerformed time.Time `json:"last_performed"`
}

type SetsRepsRecommendationData struct {
	Sets        int     `json:"sets"`
	RepsMin     int     `json:"reps_min"`
	RepsMax     int     `json:"reps_max"`
	RestSeconds int     `json:"rest_seconds"`
	Intensity   float64 `json:"intensity"`
	Notes       string  `json:"notes"`
}
