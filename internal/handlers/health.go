package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func HealthCheckHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		checks := make(map[string]string)
		overallStatus := "healthy"

		// Check database connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		if err := db.PingContext(ctx); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			overallStatus = "unhealthy"
		} else {
			checks["database"] = "healthy"
		}
		cancel()

		// Check file system permissions
		testFile := "./data/health_check.tmp"
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			checks["filesystem"] = "unhealthy: " + err.Error()
			overallStatus = "unhealthy"
		} else {
			os.Remove(testFile)
			checks["filesystem"] = "healthy"
		}

		response := HealthResponse{
			Status:    overallStatus,
			Timestamp: time.Now(),
			Checks:    checks,
		}

		if overallStatus != "healthy" {
			return c.JSON(http.StatusServiceUnavailable, response)
		}
		return c.JSON(http.StatusOK, response)
	}
}
