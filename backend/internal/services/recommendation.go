package services

import (
	"context"
	"math"
	"sort"

	"app/internal/api"
	"app/internal/db"
	"app/internal/models"
	"app/internal/utils"
)

// RecommendationService handles recommendation calculations
type RecommendationService struct {
	db              db.DatabaseInterface
	coverageService *CoverageService
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(database db.DatabaseInterface, coverageService *CoverageService) *RecommendationService {
	return &RecommendationService{
		db:              database,
		coverageService: coverageService,
	}
}

// GetDB returns the database instance (for handler access)
func (s *RecommendationService) GetDB() db.DatabaseInterface {
	return s.db
}

// ComputeHomeMbps calculates the required home internet speed in Mbps
// based on total monthly data usage from all household lines
func (s *RecommendationService) ComputeHomeMbps(lines []api.HouseholdLineDTO) float64 {
	// Sum total GB usage across all lines
	totalGB := 0.0
	for _, line := range lines {
		totalGB += line.ExpectedGB
	}

	// If no data usage, return minimal speed
	if totalGB <= 0 {
		return 10.0 // 10 Mbps minimum
	}

	// Convert GB to bits: GB × 8 bits/byte × 1024³ bytes/GB
	totalBits := totalGB * 8 * 1024 * 1024 * 1024

	// Calculate average bits per second over a month
	// 30 days × 24 hours × 3600 seconds = 2,592,000 seconds per month
	secondsPerMonth := 30.0 * 24.0 * 3600.0
	avgBitsPerSecond := totalBits / secondsPerMonth

	// Convert to Mbps: bits/second ÷ 1,000,000
	avgMbps := avgBitsPerSecond / 1000000.0

	// Apply safety factor for peak usage (3x average)
	safetyFactor := 3.0
	requiredMbps := avgMbps * safetyFactor

	// Round up to next integer and ensure minimum
	result := math.Ceil(requiredMbps)
	if result < 10.0 {
		result = 10.0
	}

	return result
}

// ComputeHomeMbpsDetailed returns detailed calculation breakdown
func (s *RecommendationService) ComputeHomeMbpsDetailed(lines []api.HouseholdLineDTO) *MbpsCalculation {
	totalGB := 0.0
	for _, line := range lines {
		totalGB += line.ExpectedGB
	}

	if totalGB <= 0 {
		return &MbpsCalculation{
			TotalGB:      0,
			AvgMbps:      0,
			SafetyFactor: 3.0,
			RequiredMbps: 10.0,
			FinalMbps:    10.0,
			Reasoning:    "Minimum 10 Mbps for basic connectivity",
		}
	}

	// Same calculation as ComputeHomeMbps but with breakdown
	totalBits := totalGB * 8 * 1000 * 1000 * 1000
	secondsPerMonth := 30.0 * 24.0 * 3600.0
	avgBitsPerSecond := totalBits / secondsPerMonth
	avgMbps := avgBitsPerSecond / 1000000.0

	safetyFactor := 3.0
	requiredMbps := avgMbps * safetyFactor
	finalMbps := math.Ceil(requiredMbps)

	if finalMbps < 10.0 {
		finalMbps = 10.0
	}

	return &MbpsCalculation{
		TotalGB:      totalGB,
		AvgMbps:      avgMbps,
		SafetyFactor: safetyFactor,
		RequiredMbps: requiredMbps,
		FinalMbps:    finalMbps,
		Reasoning:    "Based on total monthly data usage with 3x safety factor",
	}
}

// MbpsCalculation represents detailed Mbps calculation breakdown
type MbpsCalculation struct {
	TotalGB      float64 `json:"total_gb"`
	AvgMbps      float64 `json:"avg_mbps"`
	SafetyFactor float64 `json:"safety_factor"`
	RequiredMbps float64 `json:"required_mbps"`
	FinalMbps    float64 `json:"final_mbps"`
	Reasoning    string  `json:"reasoning"`
}

// GenerateCandidates creates all valid combinations of home and TV plans
func (s *RecommendationService) GenerateCandidates(ctx context.Context, availableTech []string, neededMbps float64, maxTVHours float64) ([]BundleCandidate, error) {
	// Get all plans from database
	catalog, err := s.db.GetCatalog(ctx)
	if err != nil {
		return nil, err
	}

	var candidates []BundleCandidate

	// Filter home plans by available tech and speed requirements
	validHomePlans := []models.HomePlan{}
	for _, plan := range catalog.HomePlans {
		// Check if plan technology is available
		techAvailable := false
		for _, tech := range availableTech {
			if plan.Tech == tech {
				techAvailable = true
				break
			}
		}

		// Check if plan meets speed requirements
		if techAvailable && float64(plan.DownMbps) >= neededMbps {
			validHomePlans = append(validHomePlans, plan)
		}
	}

	// Filter TV plans by HD hours requirement
	validTVPlans := []models.TVPlan{}
	for _, plan := range catalog.TVPlans {
		if plan.HDHoursIncluded >= maxTVHours {
			validTVPlans = append(validTVPlans, plan)
		}
	}

	// Generate combinations
	// Option 1: Mobile only (no home, no TV)
	candidates = append(candidates, BundleCandidate{
		HomePlan: nil,
		TVPlan:   nil,
		Label:    "Mobile Only",
	})

	// Option 2: Mobile + Home (no TV)
	for _, homePlan := range validHomePlans {
		candidates = append(candidates, BundleCandidate{
			HomePlan: &homePlan,
			TVPlan:   nil,
			Label:    "Mobile + " + homePlan.Name,
		})
	}

	// Option 3: Mobile + TV (no home) - only if TV doesn't require home connection
	for _, tvPlan := range validTVPlans {
		candidates = append(candidates, BundleCandidate{
			HomePlan: nil,
			TVPlan:   &tvPlan,
			Label:    "Mobile + " + tvPlan.Name,
		})
	}

	// Option 4: Mobile + Home + TV (triple bundle)
	for _, homePlan := range validHomePlans {
		for _, tvPlan := range validTVPlans {
			candidates = append(candidates, BundleCandidate{
				HomePlan: &homePlan,
				TVPlan:   &tvPlan,
				Label:    "Triple: " + homePlan.Name + " + " + tvPlan.Name,
			})
		}
	}

	return candidates, nil
}

// BundleCandidate represents a potential bundle combination
type BundleCandidate struct {
	HomePlan *models.HomePlan `json:"home_plan,omitempty"`
	TVPlan   *models.TVPlan   `json:"tv_plan,omitempty"`
	Label    string           `json:"label"`
}

// MatchLinesToPlans assigns the optimal mobile plan to each household line
func (s *RecommendationService) MatchLinesToPlans(lines []api.HouseholdLineDTO, mobilePlans []models.MobilePlan) []LineAssignment {
	var assignments []LineAssignment

	for _, line := range lines {
		bestPlan := s.findBestMobilePlan(line, mobilePlans)

		// Calculate costs for this line assignment
		lineCost := s.calculateLineCost(line, bestPlan)
		overageGB, overageMin := s.calculateOverages(line, bestPlan)

		assignments = append(assignments, LineAssignment{
			LineID:     line.LineID,
			Plan:       bestPlan,
			LineCost:   lineCost,
			OverageGB:  overageGB,
			OverageMin: overageMin,
		})
	}

	return assignments
}

// findBestMobilePlan finds the cheapest plan that covers the line's usage
func (s *RecommendationService) findBestMobilePlan(line api.HouseholdLineDTO, plans []models.MobilePlan) models.MobilePlan {
	var bestPlan models.MobilePlan
	bestCost := float64(999999) // Very high initial cost

	for _, plan := range plans {
		// Calculate total cost for this plan including overages
		totalCost := s.calculateLineCost(line, plan)

		// Choose the plan with lowest total cost
		if totalCost < bestCost {
			bestCost = totalCost
			bestPlan = plan
		}
	}

	return bestPlan
}

// calculateLineCost calculates the total cost for a line with a given plan
func (s *RecommendationService) calculateLineCost(line api.HouseholdLineDTO, plan models.MobilePlan) float64 {
	// Start with monthly plan price
	totalCost := plan.MonthlyPrice

	// Add overage costs
	overageGB, overageMin := s.calculateOverages(line, plan)
	totalCost += overageGB * plan.OverageGB
	totalCost += overageMin * plan.OverageMin

	return totalCost
}

// calculateOverages calculates GB and minute overages for a line with a plan
func (s *RecommendationService) calculateOverages(line api.HouseholdLineDTO, plan models.MobilePlan) (float64, float64) {
	// Calculate GB overage
	overageGB := 0.0
	if line.ExpectedGB > plan.QuotaGB {
		overageGB = line.ExpectedGB - plan.QuotaGB
	}

	// Calculate minute overage
	overageMin := 0.0
	if line.ExpectedMin > plan.QuotaMin {
		overageMin = line.ExpectedMin - plan.QuotaMin
	}

	return overageGB, overageMin
}

// LineAssignment represents a mobile plan assigned to a specific line
type LineAssignment struct {
	LineID     string            `json:"line_id"`
	Plan       models.MobilePlan `json:"plan"`
	LineCost   float64           `json:"line_cost"`
	OverageGB  float64           `json:"overage_gb"`
	OverageMin float64           `json:"overage_min"`
}

// PriceBundleCandidate calculates the total price for a bundle candidate with all discounts
func (s *RecommendationService) PriceBundleCandidate(candidate BundleCandidate, lineAssignments []LineAssignment) PricedCandidate {
	// Convert LineAssignment to utils.LineAssignment for cost engine
	var utilsLines []utils.LineAssignment
	for _, assignment := range lineAssignments {
		utilsLines = append(utilsLines, utils.LineAssignment{
			LineID: assignment.LineID,
			Usage: utils.LineUsage{
				ExpectedGB:  assignment.OverageGB + assignment.Plan.QuotaGB, // Reconstruct total usage
				ExpectedMin: assignment.OverageMin + assignment.Plan.QuotaMin,
			},
			Plan: assignment.Plan,
		})
	}

	// Calculate mobile costs
	mobileTotal := utils.CalcMobileTotal(utilsLines)
	mobileAfterDiscount, lineDiscountAmount := utils.ApplyExtraLineDiscount(mobileTotal, len(lineAssignments))

	// Calculate home and TV costs
	homeCost := 0.0
	if candidate.HomePlan != nil {
		homeCost = candidate.HomePlan.MonthlyPrice
	}

	tvCost := 0.0
	if candidate.TVPlan != nil {
		tvCost = candidate.TVPlan.MonthlyPrice
	}

	// Calculate subtotal before bundle discount
	subtotal := mobileAfterDiscount + homeCost + tvCost

	// Calculate bundle discount
	bundleDiscountRate := utils.CalcBundleDiscount(true, candidate.HomePlan != nil, candidate.TVPlan != nil) // Always have mobile
	bundleDiscountAmount := subtotal * bundleDiscountRate

	// Calculate grand total using cost engine
	breakdown := utils.CalcGrandTotal(
		mobileAfterDiscount, // Mobile cost after line discount
		homeCost,
		tvCost,
		bundleDiscountRate, // Bundle discount rate
	)

	// Generate reasoning
	reasoning := utils.BuildReasoning(utils.RecommendationSummary{
		Lines:              utilsLines,
		HomePlan:           candidate.HomePlan,
		TVPlan:             candidate.TVPlan,
		LineDiscountAmount: lineDiscountAmount,
		LineCount:          len(lineAssignments),
		Breakdown:          breakdown,
	})

	return PricedCandidate{
		Candidate:            candidate,
		LineAssignments:      lineAssignments,
		MobileTotal:          mobileTotal,
		LineDiscountAmount:   lineDiscountAmount,
		HomeCost:             homeCost,
		TVCost:               tvCost,
		BundleDiscountAmount: bundleDiscountAmount,
		BundleDiscountRate:   bundleDiscountRate,
		GrandTotal:           breakdown.GrandTotal,
		TotalSavings:         lineDiscountAmount + bundleDiscountAmount,
		Reasoning:            reasoning,
		Breakdown:            breakdown,
	}
}

// PricedCandidate represents a bundle candidate with full pricing breakdown
type PricedCandidate struct {
	Candidate            BundleCandidate           `json:"candidate"`
	LineAssignments      []LineAssignment          `json:"line_assignments"`
	MobileTotal          float64                   `json:"mobile_total"`
	LineDiscountAmount   float64                   `json:"line_discount_amount"`
	HomeCost             float64                   `json:"home_cost"`
	TVCost               float64                   `json:"tv_cost"`
	BundleDiscountAmount float64                   `json:"bundle_discount_amount"`
	BundleDiscountRate   float64                   `json:"bundle_discount_rate"`
	GrandTotal           float64                   `json:"grand_total"`
	TotalSavings         float64                   `json:"total_savings"`
	Reasoning            string                    `json:"reasoning"`
	Breakdown            utils.GrandTotalBreakdown `json:"breakdown"`
}

// ProcessRecommendationRequest processes a full recommendation request
func (s *RecommendationService) ProcessRecommendationRequest(ctx context.Context, req *api.RecommendationRequest) (*api.RecommendationResponse, error) {
	// Step 1: Check coverage for the address
	availableTech, err := s.coverageService.ComputeCoverage(ctx, req.AddressID)
	if err != nil {
		return nil, err
	}

	// Step 2: Compute needed home Mbps
	neededMbps := s.ComputeHomeMbps(req.Household)

	// Step 3: Compute max TV hours needed
	maxTVHours := 0.0
	for _, line := range req.Household {
		if line.TVHDHours > maxTVHours {
			maxTVHours = line.TVHDHours
		}
	}

	// Step 4: Generate candidate combinations
	candidates, err := s.GenerateCandidates(ctx, availableTech, neededMbps, maxTVHours)
	if err != nil {
		return nil, err
	}

	// Get mobile plans catalog for line matching
	catalog, err := s.db.GetCatalog(ctx)
	if err != nil {
		return nil, err
	}

	// Step 5: Match lines to optimal mobile plans
	lineAssignments := s.MatchLinesToPlans(req.Household, catalog.MobilePlans)

	// Step 6: Price each candidate
	var pricedCandidates []PricedCandidate
	for _, candidate := range candidates {
		priced := s.PriceBundleCandidate(candidate, lineAssignments)
		pricedCandidates = append(pricedCandidates, priced)
	}

	// Step 7: Sort by best value (lowest grand total) and return top 3
	top3 := s.SelectTop3Candidates(pricedCandidates)

	// Convert to response DTOs
	response := s.ConvertToResponse(top3)

	return response, nil
}

// SelectTop3Candidates sorts candidates by grand total and returns the best 3
func (s *RecommendationService) SelectTop3Candidates(candidates []PricedCandidate) []PricedCandidate {
	// Sort by grand total (ascending - cheapest first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].GrandTotal < candidates[j].GrandTotal
	})

	// Return top 3 (or fewer if less than 3 candidates)
	if len(candidates) <= 3 {
		return candidates
	}

	return candidates[:3]
}

// ConvertToResponse converts PricedCandidates to API response format
func (s *RecommendationService) ConvertToResponse(pricedCandidates []PricedCandidate) *api.RecommendationResponse {
	var top3 []api.RecommendationCandidateDTO

	for _, candidate := range pricedCandidates {
		// Convert mobile plan assignments
		var mobileAssignments []api.MobilePlanAssignmentDTO
		for _, assignment := range candidate.LineAssignments {
			mobileAssignments = append(mobileAssignments, api.MobilePlanAssignmentDTO{
				LineID: assignment.LineID,
				Plan: api.MobilePlanDTO{
					PlanID:       assignment.Plan.PlanID,
					PlanName:     assignment.Plan.PlanName,
					QuotaGB:      assignment.Plan.QuotaGB,
					QuotaMin:     assignment.Plan.QuotaMin,
					MonthlyPrice: assignment.Plan.MonthlyPrice,
					OverageGB:    assignment.Plan.OverageGB,
					OverageMin:   assignment.Plan.OverageMin,
				},
				LineCost:   assignment.LineCost,
				OverageGB:  assignment.OverageGB,
				OverageMin: assignment.OverageMin,
			})
		}

		// Convert home plan (if present)
		var homePlan *api.HomePlanDTO
		if candidate.Candidate.HomePlan != nil {
			homePlan = &api.HomePlanDTO{
				HomeID:       candidate.Candidate.HomePlan.HomeID,
				Name:         candidate.Candidate.HomePlan.Name,
				Tech:         candidate.Candidate.HomePlan.Tech,
				DownMbps:     candidate.Candidate.HomePlan.DownMbps,
				MonthlyPrice: candidate.Candidate.HomePlan.MonthlyPrice,
				InstallFee:   candidate.Candidate.HomePlan.InstallFee,
			}
		}

		// Convert TV plan (if present)
		var tvPlan *api.TVPlanDTO
		if candidate.Candidate.TVPlan != nil {
			tvPlan = &api.TVPlanDTO{
				TVID:            candidate.Candidate.TVPlan.TVID,
				Name:            candidate.Candidate.TVPlan.Name,
				HDHoursIncluded: candidate.Candidate.TVPlan.HDHoursIncluded,
				MonthlyPrice:    candidate.Candidate.TVPlan.MonthlyPrice,
			}
		}

		// Create the recommendation candidate
		recommendationCandidate := api.RecommendationCandidateDTO{
			ComboLabel:   candidate.Candidate.Label,
			MonthlyTotal: candidate.GrandTotal,
			Savings:      candidate.TotalSavings,
			Reasoning:    candidate.Reasoning,
			Items: api.RecommendationItemsDTO{
				Mobile: mobileAssignments,
				Home:   homePlan,
				TV:     tvPlan,
			},
			Discounts: api.RecommendationDiscountsDTO{
				LineDiscount:   candidate.LineDiscountAmount,
				BundleDiscount: candidate.BundleDiscountAmount,
				TotalDiscount:  candidate.TotalSavings,
			},
		}

		top3 = append(top3, recommendationCandidate)
	}

	return &api.RecommendationResponse{
		Top3: top3,
	}
}
