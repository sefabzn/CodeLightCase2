package models

// RecommendationRequest represents the input for recommendation calculation
type RecommendationRequest struct {
	UserID     int         `json:"user_id"`
	AddressID  string      `json:"address_id"`
	Household  []Household `json:"household"`
	PreferTech []string    `json:"prefer_tech,omitempty"` // optional tech preference order
}

// RecommendationResponse represents the output of recommendation calculation
type RecommendationResponse struct {
	Top3 []RecommendationCandidate `json:"top3"`
}

// RecommendationCandidate represents a single bundle recommendation
type RecommendationCandidate struct {
	ComboLabel   string                  `json:"combo_label"`
	Items        RecommendationItems     `json:"items"`
	MonthlyTotal float64                 `json:"monthly_total"`
	Savings      float64                 `json:"savings"`
	Reasoning    string                  `json:"reasoning"`
	Discounts    RecommendationDiscounts `json:"discounts"`
}

// RecommendationItems represents the components of a recommendation
type RecommendationItems struct {
	Mobile []MobilePlanAssignment `json:"mobile"`
	Home   *HomePlan              `json:"home,omitempty"`
	TV     *TVPlan                `json:"tv,omitempty"`
}

// MobilePlanAssignment represents which mobile plan is assigned to which line
type MobilePlanAssignment struct {
	LineID     string     `json:"line_id"`
	Plan       MobilePlan `json:"plan"`
	LineCost   float64    `json:"line_cost"` // including overage
	OverageGB  float64    `json:"overage_gb"`
	OverageMin float64    `json:"overage_min"`
}

// RecommendationDiscounts represents applied discounts
type RecommendationDiscounts struct {
	LineDiscount   float64 `json:"line_discount"`   // extra line discount
	BundleDiscount float64 `json:"bundle_discount"` // bundle discount
	TotalDiscount  float64 `json:"total_discount"`  // sum of all discounts
}

// CheckoutRequest represents the checkout request
type CheckoutRequest struct {
	UserID        int                     `json:"user_id"`
	SelectedCombo RecommendationCandidate `json:"selected_combo"`
	SlotID        int                     `json:"slot_id"`
	AddressID     string                  `json:"address_id"`
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}
