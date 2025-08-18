package models

import "time"

// User represents a user in the system
type User struct {
	UserID             int       `json:"user_id" db:"user_id"`
	Name               string    `json:"name" db:"name"`
	AddressID          string    `json:"address_id" db:"address_id"`
	CurrentBundleLabel *string   `json:"current_bundle_label" db:"current_bundle_label"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// Coverage represents technology availability by address
type Coverage struct {
	AddressID string `json:"address_id" db:"address_id"`
	City      string `json:"city" db:"city"`
	District  string `json:"district" db:"district"`
	Fiber     bool   `json:"fiber" db:"fiber"`
	VDSL      bool   `json:"vdsl" db:"vdsl"`
	FWA       bool   `json:"fwa" db:"fwa"`
}

// Household represents a household member and their usage patterns
type Household struct {
	ID          int     `json:"id" db:"id"`
	UserID      int     `json:"user_id" db:"user_id"`
	LineID      string  `json:"line_id" db:"line_id"`
	ExpectedGB  float64 `json:"expected_gb" db:"expected_gb"`
	ExpectedMin float64 `json:"expected_min" db:"expected_min"`
	TVHDHours   float64 `json:"tv_hd_hours" db:"tv_hd_hours"`
}

// CurrentServices represents current services the user has
type CurrentServices struct {
	ID            int     `json:"id" db:"id"`
	UserID        int     `json:"user_id" db:"user_id"`
	HasHome       bool    `json:"has_home" db:"has_home"`
	HomeTech      *string `json:"home_tech" db:"home_tech"`
	HomeSpeed     *int    `json:"home_speed" db:"home_speed"`
	HasTV         bool    `json:"has_tv" db:"has_tv"`
	MobilePlanIDs string  `json:"mobile_plan_ids" db:"mobile_plan_ids"` // JSON array
}
