package handlers

import (
	"net/http"

	"api-key-generator/internal/services"

	"github.com/labstack/echo/v4"
)

// WarmupHandler handles warmup and technique endpoints
type WarmupHandler struct {
	warmupService *services.WarmupService
}

// NewWarmupHandler creates a new warmup handler
func NewWarmupHandler() *WarmupHandler {
	return &WarmupHandler{
		warmupService: services.NewWarmupService(),
	}
}

// GenerateWarmup creates a customized warmup routine
// POST /api/v1/workouts/warmup
func (h *WarmupHandler) GenerateWarmup(c echo.Context) error {
	var req services.WarmupRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Set defaults
	if req.Language == "" {
		req.Language = "en"
	}
	if req.Intensity == "" {
		req.Intensity = "moderate"
	}

	routine, err := h.warmupService.GenerateWarmupRoutine(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate warmup routine")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    routine,
	})
}

// GenerateCooldown creates a post-workout cooldown routine
// POST /api/v1/workouts/cooldown
func (h *WarmupHandler) GenerateCooldown(c echo.Context) error {
	var req services.CooldownRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if req.Language == "" {
		req.Language = "en"
	}

	routine, err := h.warmupService.GenerateCooldownRoutine(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate cooldown routine")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    routine,
	})
}

// GetExerciseTechnique provides comprehensive technique guidance
// GET /api/v1/workouts/exercises/{exerciseID}/technique
func (h *WarmupHandler) GetExerciseTechnique(c echo.Context) error {
	exerciseID := c.Param("exerciseID")
	language := c.QueryParam("language")

	if exerciseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Exercise ID is required")
	}

	if language == "" {
		language = "en"
	}

	technique, err := h.warmupService.GetExerciseTechnique(c.Request().Context(), exerciseID, language)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Exercise technique not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    technique,
	})
}
