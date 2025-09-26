package handlers

import (
	"github.com/labstack/echo/v4"
)

// SetupPDFRoutes configures all PDF-related routes
func SetupPDFRoutes(e *echo.Echo) {
	pdfHandler := NewPDFHandler()

	// PDF generation endpoint
	e.POST("/api/pdf/generate", pdfHandler.GeneratePDF)

	// Quota status endpoint
	e.GET("/api/pdf/quota", pdfHandler.GetQuotaStatus)

	// Alternative routes for frontend convenience
	e.POST("/pdf/generate", pdfHandler.GeneratePDF)
	e.GET("/pdf/quota", pdfHandler.GetQuotaStatus)
}
