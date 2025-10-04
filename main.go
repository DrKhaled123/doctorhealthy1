package main

import (
	"context"
	"embed"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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

// Embed frontend at build time for resilient serving
//
//go:embed frontend/*
var embeddedFrontend embed.FS

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Debug: Log configuration values (only in debug/development mode)
	if cfg.Logging.EnableDebug {
		log.Printf("[DEBUG] Server Port: %s", cfg.Server.Port)
		log.Printf("[DEBUG] Database Path: %s", cfg.Database.Path)
		log.Printf("[DEBUG] JWT Secret present: %t", os.Getenv("JWT_SECRET") != "")
		log.Printf("[DEBUG] CORS Origins: %s", os.Getenv("CORS_ORIGINS"))
		log.Printf("[DEBUG] Log Level: %s", cfg.Logging.Level)
	}

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
	apiKeyService, err := services.NewAPIKeyService(db, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize API key service: %v", err)
	}
	defer func() {
		if err := apiKeyService.Close(); err != nil {
			log.Printf("Error closing API key service: %v", err)
		}
	}()

	userService := services.NewUserService(db)
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

	// Production-ready CORS configuration
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{
			"https://my.doctorhealthy1.com",
			"https://doctorhealthy1.com",
			"https://www.doctorhealthy1.com",
			"http://localhost:3000", // Local development
			"http://localhost:8080", // Local testing
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-API-Key",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			echo.HeaderContentLength,
			echo.HeaderContentType,
			"X-Request-ID",
		},
		MaxAge: 300, // 5 minutes
	}))

	e.Use(middleware.Security())
	e.Use(middleware.RateLimit(cfg))
	// Optional JWT (sets user_id when Authorization: Bearer is present)
	e.Use(middleware.OptionalJWT())
	// Enforce per-user monthly quotas on generate endpoints (free/pro/lifetime)
	e.Use(middleware.UserQuotaMiddleware())

	// Centralized error handler
	middleware.SetupErrorHandler(e)

	// Health check
	e.GET("/health", handlers.HealthCheckHandler(db))

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

	// Health management routes (handlers already registered in individual route files)

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

	// PDF generation routes with quota management
	handlers.SetupPDFRoutes(e)

	// URL slug generation routes
	handlers.RegisterURLSlugRoutes(api, apiKeyService)

	// Static frontend serving with dual-source fallback (disk first, then embedded)
	// Serve root index
	e.GET("/", func(c echo.Context) error {
		return serveStaticFile(c, "index.html", embeddedFrontend)
	})

	// Serve frontend pages
	e.GET("/diet.html", func(c echo.Context) error {
		return serveStaticFile(c, "diet.html", embeddedFrontend)
	})

	e.GET("/workouts.html", func(c echo.Context) error {
		return serveStaticFile(c, "workouts.html", embeddedFrontend)
	})

	e.GET("/recipes.html", func(c echo.Context) error {
		return serveStaticFile(c, "recipes.html", embeddedFrontend)
	})

	e.GET("/lifestyle.html", func(c echo.Context) error {
		return serveStaticFile(c, "lifestyle.html", embeddedFrontend)
	})

	e.GET("/pricing.html", func(c echo.Context) error {
		return serveStaticFile(c, "pricing.html", embeddedFrontend)
	})

	e.GET("/home.html", func(c echo.Context) error {
		return serveStaticFile(c, "home.html", embeddedFrontend)
	})

	// Serve URL slug generator page
	e.GET("/url-slug", func(c echo.Context) error {
		return serveStaticFile(c, "url-slug.html", embeddedFrontend)
	})

	// Serve CSS and JS assets
	e.GET("/frontend/css/*", func(c echo.Context) error {
		path := strings.TrimPrefix(c.Request().URL.Path, "/frontend/css/")
		return serveStaticFile(c, "css/"+path, embeddedFrontend)
	})

	e.GET("/frontend/js/*", func(c echo.Context) error {
		path := strings.TrimPrefix(c.Request().URL.Path, "/frontend/js/")
		return serveStaticFile(c, "js/"+path, embeddedFrontend)
	})

	// Catch-all fallback to index for SPA-like navigation (only for non-API routes)
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path
		// Don't serve index.html for API routes
		if strings.HasPrefix(path, "/api/") {
			return echo.NewHTTPError(http.StatusNotFound, "endpoint not found")
		}
		return serveStaticFile(c, "index.html", embeddedFrontend)
	})

	// Check if port is available
	if cfg.Logging.EnableDebug {
		log.Printf("[DEBUG] Attempting to start server on port %s", cfg.Server.Port)
	}

	// Start server
	errorChan := make(chan error, 1)
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			if cfg.Logging.EnableDebug {
				log.Printf("[DEBUG] Server start error: %v", err)
			}
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

// serveStaticFile serves a file from disk if present under ./frontend, otherwise from the embedded FS.
func serveStaticFile(c echo.Context, relPath string, embedded embed.FS) error {
	// Prefer on-disk to allow hot updates in deployments mounting /frontend
	diskPath := filepath.Join("frontend", relPath)
	if f, err := os.Open(diskPath); err == nil { // nosec G304 - controlled path within frontend directory
		defer func() { _ = f.Close() }()
		// Determine content type
		ext := filepath.Ext(relPath)
		ct := mime.TypeByExtension(ext)
		if ct == "" {
			ct = "text/plain; charset=utf-8"
		}
		c.Response().Header().Set(echo.HeaderContentType, ct)
		_, _ = io.Copy(c.Response(), f)
		return nil
	}

	// Fallback to embedded
	data, err := embedded.ReadFile(filepath.ToSlash(filepath.Join("frontend", relPath))) // nosec G304 - controlled path within embedded assets
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "asset not found")
	}
	ext := filepath.Ext(relPath)
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		ct = "text/plain; charset=utf-8"
	}
	return c.Blob(http.StatusOK, ct, data)
}
