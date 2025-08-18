package utils

import (
	"app/internal/models"
	"fmt"
)

// LineUsage represents the expected usage for a single mobile line
type LineUsage struct {
	ExpectedGB  float64
	ExpectedMin float64
}

// OverageResult represents the overage calculation result
type OverageResult struct {
	OverageGB   float64 // GB over quota
	OverageMin  float64 // Minutes over quota
	OverageCost float64 // Total overage cost
}

// CalcMobileOverage calculates overage cost for a single mobile line
// Returns overage amounts and total overage cost
func CalcMobileOverage(usage LineUsage, plan models.MobilePlan) OverageResult {
	result := OverageResult{}

	// Calculate GB overage
	if usage.ExpectedGB > plan.QuotaGB {
		result.OverageGB = usage.ExpectedGB - plan.QuotaGB
		result.OverageCost += result.OverageGB * plan.OverageGB
	}

	// Calculate minute overage
	if usage.ExpectedMin > plan.QuotaMin {
		result.OverageMin = usage.ExpectedMin - plan.QuotaMin
		result.OverageCost += result.OverageMin * plan.OverageMin
	}

	return result
}

// CalcMobileLineCost calculates the total cost for a single mobile line
// Returns monthly_price + overage cost
func CalcMobileLineCost(usage LineUsage, plan models.MobilePlan) float64 {
	overage := CalcMobileOverage(usage, plan)
	return plan.MonthlyPrice + overage.OverageCost
}

// LineAssignment represents a mobile plan assignment to a specific line
type LineAssignment struct {
	LineID string
	Usage  LineUsage
	Plan   models.MobilePlan
}

// CalcMobileTotal calculates the total mobile cost for multiple lines
// Returns the sum of all line costs
func CalcMobileTotal(lines []LineAssignment) float64 {
	total := 0.0

	for _, line := range lines {
		lineCost := CalcMobileLineCost(line.Usage, line.Plan)
		total += lineCost
	}

	return total
}

// ApplyExtraLineDiscount applies discount for multiple mobile lines
// Rules: 2nd line -5%, 3+ lines -10% on mobile component only
// Returns the discounted mobile total and discount amount
func ApplyExtraLineDiscount(mobileTotal float64, lineCount int) (discountedTotal float64, discountAmount float64) {
	discountPercent := 0.0

	switch {
	case lineCount >= 3:
		discountPercent = 0.10 // 10% discount for 3+ lines
	case lineCount == 2:
		discountPercent = 0.05 // 5% discount for 2nd line
	default:
		discountPercent = 0.0 // No discount for single line
	}

	discountAmount = mobileTotal * discountPercent
	discountedTotal = mobileTotal - discountAmount

	return discountedTotal, discountAmount
}

// SelectHomePlan selects the best home plan based on available technology and needed speed
// Priority order: fiber > vdsl > fwa
// Returns the cheapest plan that meets speed requirements for the highest priority available tech
func SelectHomePlan(techAvailable []string, neededMbps int, homePlans []models.HomePlan) *models.HomePlan {
	// Define technology priority order
	techPriority := []string{"fiber", "vdsl", "fwa"}

	// Check each technology in priority order
	for _, preferredTech := range techPriority {
		// Skip if this tech is not available
		techFound := false
		for _, availableTech := range techAvailable {
			if availableTech == preferredTech {
				techFound = true
				break
			}
		}
		if !techFound {
			continue
		}

		// Find plans for this technology that meet speed requirements
		var candidatePlans []models.HomePlan
		for _, plan := range homePlans {
			if plan.Tech == preferredTech && plan.DownMbps >= neededMbps {
				candidatePlans = append(candidatePlans, plan)
			}
		}

		// If we found suitable plans, return the cheapest one
		if len(candidatePlans) > 0 {
			cheapest := candidatePlans[0]
			for _, plan := range candidatePlans[1:] {
				if plan.MonthlyPrice < cheapest.MonthlyPrice {
					cheapest = plan
				}
			}
			return &cheapest
		}
	}

	// No suitable plan found
	return nil
}

// SelectTvPlan selects the best TV plan based on required HD hours
// Returns the cheapest plan that covers the required HD hours
func SelectTvPlan(requiredHDHours float64, tvPlans []models.TVPlan) *models.TVPlan {
	if requiredHDHours <= 0 {
		return nil // No TV needed
	}

	var candidatePlans []models.TVPlan

	// Find all plans that cover the required HD hours
	for _, plan := range tvPlans {
		if plan.HDHoursIncluded >= requiredHDHours {
			candidatePlans = append(candidatePlans, plan)
		}
	}

	// If no plan covers the requirement, return nil
	if len(candidatePlans) == 0 {
		return nil
	}

	// Return the cheapest plan that meets requirements
	cheapest := candidatePlans[0]
	for _, plan := range candidatePlans[1:] {
		if plan.MonthlyPrice < cheapest.MonthlyPrice {
			cheapest = plan
		}
	}

	return &cheapest
}

// CalcBundleDiscount calculates bundle discount based on service combination
// Rules: mobile+home=10%; mobile+home+tv=15%; else 0%
// Returns discount percentage (as decimal, e.g., 0.10 for 10%)
func CalcBundleDiscount(hasMobile, hasHome, hasTV bool) float64 {
	if hasMobile && hasHome && hasTV {
		return 0.15 // 15% for triple bundle
	}

	if hasMobile && hasHome {
		return 0.10 // 10% for mobile + home bundle
	}

	return 0.0 // No bundle discount
}

// GrandTotalBreakdown represents the final cost calculation breakdown
type GrandTotalBreakdown struct {
	MobileTotal        float64 // Mobile cost after line discounts
	HomeTotal          float64 // Home plan cost
	TVTotal            float64 // TV plan cost
	SubTotal           float64 // Sum before bundle discount
	BundleDiscount     float64 // Bundle discount amount
	GrandTotal         float64 // Final total after all discounts
	BundleDiscountRate float64 // Bundle discount percentage applied
}

// CalcGrandTotal calculates the final total cost with all discounts applied
// Takes mobile total (after line discounts), home cost, TV cost, and bundle discount rate
func CalcGrandTotal(mobileAfterLineDisc, homeCost, tvCost, bundleDiscountRate float64) GrandTotalBreakdown {
	subTotal := mobileAfterLineDisc + homeCost + tvCost
	bundleDiscountAmount := subTotal * bundleDiscountRate
	grandTotal := subTotal - bundleDiscountAmount

	return GrandTotalBreakdown{
		MobileTotal:        mobileAfterLineDisc,
		HomeTotal:          homeCost,
		TVTotal:            tvCost,
		SubTotal:           subTotal,
		BundleDiscount:     bundleDiscountAmount,
		GrandTotal:         grandTotal,
		BundleDiscountRate: bundleDiscountRate,
	}
}

// RecommendationSummary represents all components used for generating reasoning
type RecommendationSummary struct {
	Lines              []LineAssignment
	HomePlan           *models.HomePlan
	TVPlan             *models.TVPlan
	LineDiscountAmount float64
	LineCount          int
	Breakdown          GrandTotalBreakdown
}

// BuildReasoning generates human-readable explanation of the recommendation
func BuildReasoning(summary RecommendationSummary) string {
	reasoning := "Selected plans: "

	// Mobile plans
	if len(summary.Lines) > 0 {
		reasoning += fmt.Sprintf("%d mobile line(s)", len(summary.Lines))

		// Add line discount info
		if summary.LineDiscountAmount > 0 {
			discountPercent := 5.0
			if summary.LineCount >= 3 {
				discountPercent = 10.0
			}
			reasoning += fmt.Sprintf(" (%.0f%% multi-line discount)", discountPercent)
		}
	}

	// Home plan
	if summary.HomePlan != nil {
		reasoning += fmt.Sprintf(", %s", summary.HomePlan.Name)
	}

	// TV plan
	if summary.TVPlan != nil {
		reasoning += fmt.Sprintf(", %s", summary.TVPlan.Name)
	}

	// Bundle discount
	if summary.Breakdown.BundleDiscountRate > 0 {
		bundlePercent := summary.Breakdown.BundleDiscountRate * 100
		reasoning += fmt.Sprintf(". Bundle discount: %.0f%% off total", bundlePercent)
	}

	reasoning += "."

	return reasoning
}
