package models

import "time"

// MobilePlan represents a mobile plan in the catalog
type MobilePlan struct {
	PlanID       int     `json:"plan_id" db:"plan_id"`
	PlanName     string  `json:"plan_name" db:"plan_name"`
	QuotaGB      float64 `json:"quota_gb" db:"quota_gb"`
	QuotaMin     float64 `json:"quota_min" db:"quota_min"`
	MonthlyPrice float64 `json:"monthly_price" db:"monthly_price"`
	OverageGB    float64 `json:"overage_gb" db:"overage_gb"`
	OverageMin   float64 `json:"overage_min" db:"overage_min"`
}

// HomePlan represents a home internet plan in the catalog
type HomePlan struct {
	HomeID       int     `json:"home_id" db:"home_id"`
	Name         string  `json:"name" db:"name"`
	Tech         string  `json:"tech" db:"tech"` // fiber, vdsl, fwa
	DownMbps     int     `json:"down_mbps" db:"down_mbps"`
	MonthlyPrice float64 `json:"monthly_price" db:"monthly_price"`
	InstallFee   float64 `json:"install_fee" db:"install_fee"`
}

// TVPlan represents a TV plan in the catalog
type TVPlan struct {
	TVID            int     `json:"tv_id" db:"tv_id"`
	Name            string  `json:"name" db:"name"`
	HDHoursIncluded float64 `json:"hd_hours_included" db:"hd_hours_included"`
	MonthlyPrice    float64 `json:"monthly_price" db:"monthly_price"`
}

// BundlingRule represents a bundling rule for discounts
type BundlingRule struct {
	RuleID          int     `json:"rule_id" db:"rule_id"`
	RuleType        string  `json:"rule_type" db:"rule_type"` // line_discount, bundle_discount
	Description     string  `json:"description" db:"description"`
	DiscountPercent float64 `json:"discount_percent" db:"discount_percent"`
	AppliesTo       string  `json:"applies_to" db:"applies_to"` // mobile, home, tv, total
}

// InstallSlot represents an installation time slot
type InstallSlot struct {
	SlotID    int       `json:"slot_id" db:"slot_id"`
	AddressID string    `json:"address_id" db:"address_id"`
	SlotStart time.Time `json:"slot_start" db:"slot_start"`
	SlotEnd   time.Time `json:"slot_end" db:"slot_end"`
	Tech      string    `json:"tech" db:"tech"` // fiber, vdsl, fwa
	Available bool      `json:"available" db:"available"`
}
