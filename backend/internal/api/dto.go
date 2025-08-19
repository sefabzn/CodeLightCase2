package api

// RecommendationRequest represents the input for recommendation calculation
type RecommendationRequest struct {
	UserID     int                `json:"user_id" validate:"required"`
	AddressID  string             `json:"address_id" validate:"required"`
	Household  []HouseholdLineDTO `json:"household" validate:"required,min=1,dive"`
	PreferTech []string           `json:"prefer_tech,omitempty"`
}

// HouseholdLineDTO represents a single household line input
type HouseholdLineDTO struct {
	LineID      string  `json:"line_id" validate:"required"`
	ExpectedGB  float64 `json:"expected_gb" validate:"required,min=0"`
	ExpectedMin float64 `json:"expected_min" validate:"required,min=0"`
	TVHDHours   float64 `json:"tv_hd_hours" validate:"min=0"`
}

// RecommendationResponse represents the output of recommendation calculation
type RecommendationResponse struct {
	Top3 []RecommendationCandidateDTO `json:"top3"`
}

// RecommendationCandidateDTO represents a single bundle recommendation
type RecommendationCandidateDTO struct {
	ComboLabel   string                     `json:"combo_label"`
	Items        RecommendationItemsDTO     `json:"items"`
	MonthlyTotal float64                    `json:"monthly_total"`
	Savings      float64                    `json:"savings"`
	Reasoning    string                     `json:"reasoning"`
	Discounts    RecommendationDiscountsDTO `json:"discounts"`
}

// RecommendationItemsDTO represents the components of a recommendation
type RecommendationItemsDTO struct {
	Mobile []MobilePlanAssignmentDTO `json:"mobile"`
	Home   *HomePlanDTO              `json:"home,omitempty"`
	TV     *TVPlanDTO                `json:"tv,omitempty"`
}

// MobilePlanAssignmentDTO represents which mobile plan is assigned to which line
type MobilePlanAssignmentDTO struct {
	LineID     string        `json:"line_id"`
	Plan       MobilePlanDTO `json:"plan"`
	LineCost   float64       `json:"line_cost"` // including overage
	OverageGB  float64       `json:"overage_gb"`
	OverageMin float64       `json:"overage_min"`
}

// RecommendationDiscountsDTO represents applied discounts
type RecommendationDiscountsDTO struct {
	LineDiscount   float64 `json:"line_discount"`   // extra line discount amount
	BundleDiscount float64 `json:"bundle_discount"` // bundle discount amount
	TotalDiscount  float64 `json:"total_discount"`  // sum of all discounts
}

// Plan DTOs
type MobilePlanDTO struct {
	PlanID       int     `json:"plan_id"`
	PlanName     string  `json:"plan_name"`
	QuotaGB      float64 `json:"quota_gb"`
	QuotaMin     float64 `json:"quota_min"`
	MonthlyPrice float64 `json:"monthly_price"`
	OverageGB    float64 `json:"overage_gb"`
	OverageMin   float64 `json:"overage_min"`
}

type HomePlanDTO struct {
	HomeID       int     `json:"home_id"`
	Name         string  `json:"name"`
	Tech         string  `json:"tech"`
	DownMbps     int     `json:"down_mbps"`
	MonthlyPrice float64 `json:"monthly_price"`
	InstallFee   float64 `json:"install_fee"`
}

type TVPlanDTO struct {
	TVID            int     `json:"tv_id"`
	Name            string  `json:"name"`
	HDHoursIncluded float64 `json:"hd_hours_included"`
	MonthlyPrice    float64 `json:"monthly_price"`
}

// CheckoutRequest represents the checkout request
type CheckoutRequest struct {
	UserID        int                        `json:"user_id" validate:"required"`
	SelectedCombo RecommendationCandidateDTO `json:"selected_combo" validate:"required"`
	SlotID        string                     `json:"slot_id" validate:"required"`
	AddressID     string                     `json:"address_id" validate:"required"`
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}
