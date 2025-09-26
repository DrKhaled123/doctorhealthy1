package handlers

import (
	"github.com/labstack/echo/v4"
)

// RegisterSetsRepsRoutes registers all sets, reps, and rest period management routes
func RegisterSetsRepsRoutes(e *echo.Echo, handler *SetsRepsHandler) {
	// Sets, reps, and rest period management routes
	workoutGroup := e.Group("/api/v1/workouts")

	// Exercise recommendations
	workoutGroup.POST("/recommendations", handler.GetExerciseRecommendations)
	workoutGroup.GET("/sets-reps", handler.GetSetsRepsRecommendation)
	workoutGroup.GET("/rest-period", handler.GetRestPeriodRecommendation)

	// Workout session tracking
	workoutGroup.POST("/sessions", handler.TrackWorkoutSession)
	workoutGroup.GET("/progress/:userID/:programID", handler.GetProgressReport)

	// Progression adjustments
	workoutGroup.POST("/progression-adjustment", handler.GetProgressionAdjustment)
}
