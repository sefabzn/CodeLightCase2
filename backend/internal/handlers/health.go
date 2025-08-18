package handlers

import (
	"net/http"

	"app/internal/db"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db db.DatabaseInterface
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(database db.DatabaseInterface) *HealthHandler {
	return &HealthHandler{
		db: database,
	}
}

// GetHealth handles GET /health
func (h *HealthHandler) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()

	// Check database connectivity
	if err := h.db.Health(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status":   "error",
			"database": "unhealthy",
			"error":    "database connection failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "ok",
		"database": "connected",
		"service":  "recommendation-api",
		"version":  "1.0.0",
	})
}
