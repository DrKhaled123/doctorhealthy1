package handlers

import (
	"github.com/labstack/echo/v4"
)

// RegisterWarmupRoutes registers warmup and technique routes
func RegisterWarmupRoutes(e *echo.Echo, handler *WarmupHandler) {
	workoutGroup := e.Group("/api/v1/workouts")

	// Warmup and cooldown routes
	workoutGroup.POST("/warmup", handler.GenerateWarmup)
	workoutGroup.POST("/cooldown", handler.GenerateCooldown)

	// Exercise technique routes
	workoutGroup.GET("/exercises/:exerciseID/technique", handler.GetExerciseTechnique)
}
