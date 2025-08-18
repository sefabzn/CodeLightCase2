package handlers

import (
	"app/internal/db"
	"app/internal/services"
	"app/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRoutes configures all HTTP routes and middleware
func SetupRoutes(e *echo.Echo, database db.DatabaseInterface) {
	// Create services
	coverageService := services.NewCoverageService(database)
	recommendationService := services.NewRecommendationService(database, coverageService)
	validator := utils.NewValidator()

	// Create handlers
	healthHandler := NewHealthHandler(database)
	recommendationHandler := NewRecommendationHandler(recommendationService, validator)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Health check
	e.GET("/health", healthHandler.GetHealth)

	// API routes
	api := e.Group("/api")
	{
		// Recommendation endpoints
		api.POST("/recommendation", recommendationHandler.GetRecommendations)
		api.POST("/checkout", recommendationHandler.PostCheckout)

		// Utility endpoints
		api.GET("/coverage/:address_id", recommendationHandler.GetCoverage)
		api.GET("/install-slots/:address_id", recommendationHandler.GetInstallSlots)
	}
}
