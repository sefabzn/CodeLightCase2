package handlers

import (
	"net/http"

	"app/internal/api"
	"app/internal/services"
	"app/internal/utils"

	"github.com/labstack/echo/v4"
)

// RecommendationHandler handles recommendation-related HTTP requests
type RecommendationHandler struct {
	recommendationService *services.RecommendationService
	validator             *utils.Validator
}

// NewRecommendationHandler creates a new recommendation handler
func NewRecommendationHandler(recommendationService *services.RecommendationService, validator *utils.Validator) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
		validator:             validator,
	}
}

// GetRecommendations handles POST /api/recommendation
func (h *RecommendationHandler) GetRecommendations(c echo.Context) error {
	// Parse request body
	var req api.RecommendationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "INVALID_REQUEST_BODY",
				Message: "Failed to parse request body",
				Details: []string{err.Error()},
			},
		})
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Request validation failed",
				Details: validationErrors,
			},
		})
	}

	// Process recommendation request
	response, err := h.recommendationService.ProcessRecommendationRequest(c.Request().Context(), &req)
	if err != nil {
		// Log error for debugging
		c.Logger().Errorf("Recommendation processing failed: %v", err)

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "RECOMMENDATION_FAILED",
				Message: "Failed to generate recommendations",
				Details: []string{"An internal error occurred while processing your request"},
			},
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetCoverage handles GET /api/coverage/:address_id
func (h *RecommendationHandler) GetCoverage(c echo.Context) error {
	addressID := c.Param("address_id")
	if addressID == "" {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "MISSING_ADDRESS_ID",
				Message: "Address ID is required",
			},
		})
	}

	// Get coverage info
	coverageService := services.NewCoverageService(h.recommendationService.GetDB())
	coverageInfo, err := coverageService.GetCoverageInfo(c.Request().Context(), addressID)
	if err != nil {
		c.Logger().Errorf("Coverage lookup failed for address %s: %v", addressID, err)

		return c.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "COVERAGE_NOT_FOUND",
				Message: "Coverage information not found for the specified address",
				Details: []string{addressID},
			},
		})
	}

	return c.JSON(http.StatusOK, coverageInfo)
}

// GetInstallSlots handles GET /api/install-slots/:address_id?tech=fiber
func (h *RecommendationHandler) GetInstallSlots(c echo.Context) error {
	addressID := c.Param("address_id")
	if addressID == "" {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "MISSING_ADDRESS_ID",
				Message: "Address ID is required",
			},
		})
	}

	tech := c.QueryParam("tech")
	if tech == "" {
		tech = "fiber" // Default to fiber
	}

	// Get install slots
	slots, err := h.recommendationService.GetDB().GetInstallSlots(c.Request().Context(), addressID, tech)
	if err != nil {
		c.Logger().Errorf("Install slots lookup failed for address %s, tech %s: %v", addressID, tech, err)

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "SLOTS_LOOKUP_FAILED",
				Message: "Failed to retrieve install slots",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"address_id": addressID,
		"tech":       tech,
		"slots":      slots,
	})
}

// PostCheckout handles POST /api/checkout
func (h *RecommendationHandler) PostCheckout(c echo.Context) error {
	// Parse request body
	var req api.CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "INVALID_REQUEST_BODY",
				Message: "Failed to parse checkout request",
				Details: []string{err.Error()},
			},
		})
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Checkout validation failed",
				Details: validationErrors,
			},
		})
	}

	// For now, return a mock successful checkout
	// In a real implementation, this would process payment, create orders, etc.
	response := api.CheckoutResponse{
		Status:  "success",
		OrderID: "ORD-" + generateOrderID(),
	}

	c.Logger().Infof("Checkout completed for user %d: %s", req.UserID, response.OrderID)

	return c.JSON(http.StatusOK, response)
}

// generateOrderID generates a simple order ID (in production, use a proper ID generator)
func generateOrderID() string {
	// Simple timestamp-based ID for demo purposes
	return "12345678"
}
