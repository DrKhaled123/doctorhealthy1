package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-key-generator/internal/config"
	"api-key-generator/internal/database"
	"api-key-generator/internal/handlers"
	"api-key-generator/internal/middleware"
	"api-key-generator/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Debug: Log configuration values
	log.Printf("DEBUG: Server Port: %s", cfg.Server.Port)
	log.Printf("DEBUG: Database Path: %s", cfg.Database.Path)
	log.Printf("DEBUG: JWT Secret present: %t", os.Getenv("JWT_SECRET") != "")
	log.Printf("DEBUG: CORS Origins: %s", os.Getenv("CORS_ORIGINS"))

	// Initialize database
	db, err := database.Initialize(cfg.Database.Path)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize services
	apiKeyService := services.NewAPIKeyService(db, cfg)
	userService := services.NewUserService(db)
	dataLoader := services.NewDataLoader(".")
	nutritionService := services.NewNutritionService(db, userService)
	workoutService := services.NewWorkoutService(db, userService)
	healthService := services.NewHealthService(db, userService, dataLoader)
	recipeService := services.NewRecipeService(db)
	vipService := services.NewVIPIntegrationService(".")
	ultimateService := services.NewUltimateDataService()

	// Initialize Enhanced Health Service
	enhancedService := services.NewEnhancedHealthService(".")

	// Initialize Echo
	e := echo.New()

	// Set validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.Security())
	e.Use(middleware.RateLimit(cfg))
	// Optional JWT (sets user_id when Authorization: Bearer is present)
	e.Use(middleware.OptionalJWT())

	// Centralized error handler
	middleware.SetupErrorHandler(e)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   getVersion(),
			"features":  []string{"api_keys", "nutrition", "workouts", "health", "recipes", "vip_data_integration", "ultimate_comprehensive_database", "sets_reps_management", "warmup_technique_system", "comprehensive_workout_system", "enhanced_health_system"},
		})
	})

	// Readiness check
	e.GET("/ready", func(c echo.Context) error {
		// DB ping
		if err := db.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"ready": false, "reason": "db_unreachable"})
		}
		// Optionally check core data loaded later
		return c.JSON(http.StatusOK, map[string]interface{}{"ready": true})
	})

	// API routes
	api := e.Group("/api/v1")

	// API Key management routes
	handlers.RegisterAPIKeyRoutes(api, apiKeyService)

	// Health management routes
	handlers.RegisterHealthRoutes(api, userService, nutritionService, workoutService, healthService, recipeService)

	// Recipe API routes with security integration
	handlers.SetupRecipeRoutes(e, recipeService, apiKeyService)

	// Supplement service and routes
	supplementService := services.NewSupplementService(db, userService)
	handlers.SetupSupplementRoutes(e, supplementService, apiKeyService)

	// VIP Data integration routes
	handlers.RegisterVIPDataRoutes(api, vipService)

	// Ultimate comprehensive database routes
	handlers.RegisterUltimateRoutes(api, ultimateService, apiKeyService)

	// Sets, reps, and rest period management routes
	setsRepsHandler := handlers.NewSetsRepsHandler()
	handlers.RegisterSetsRepsRoutes(e, setsRepsHandler)

	// Warmup and technique routes
	warmupHandler := handlers.NewWarmupHandler()
	handlers.RegisterWarmupRoutes(e, warmupHandler)

	// Comprehensive workout system routes
	comprehensiveWorkoutHandler := handlers.NewComprehensiveWorkoutHandler()
	handlers.RegisterComprehensiveWorkoutRoutes(e, comprehensiveWorkoutHandler, apiKeyService)

	// Enhanced Health System routes
	handlers.RegisterEnhancedHealthRoutes(e, enhancedService, apiKeyService)

	// Check if port is available
	log.Printf("DEBUG: Attempting to start server on port %s", cfg.Server.Port)

	// Start server
	errorChan := make(chan error, 1)
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			log.Printf("DEBUG: Server start error: %v", err)
			errorChan <- err
		}
	}()

	log.Printf("🚀 Health Management System server started on port %s", cfg.Server.Port)
	log.Printf("📊 Available features: API Keys, Nutrition Plans, Workout Plans, Health Management, Recipes, Ultimate Comprehensive Database, Sets/Reps Management, Warmup/Technique System, Comprehensive Workout System, Enhanced Health System")
	log.Printf("🏆 Ultimate Database: 118+ items across 10 categories with full authentication")
	log.Printf("🌟 Enhanced Health System: Diet Plans, Workout Plans, Lifestyle Management, Recipes, Injury Management, Supplement Guidance with PDF generation")

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errorChan:
		log.Printf("Failed to start server: %v", err)
		// Attempt graceful shutdown if possible
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if shutdownErr := e.Shutdown(ctx); shutdownErr != nil {
			log.Printf("Server forced to shutdown after start error: %v", shutdownErr)
		}
		os.Exit(1)
	case <-quit:
		log.Println("Received shutdown signal")
	}

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Server exited")
}

// Version can be set at build time with -ldflags "-X main.version=x.x.x"
var version = "1.0.0"

// getVersion returns the application version
func getVersion() string {
	if v := os.Getenv("APP_VERSION"); v != "" {
		return v
	}
	return version
}

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
