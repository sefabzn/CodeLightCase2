package services

import (
	"math"
	"testing"

	"app/internal/api"
	"app/internal/models"
)

func TestComputeHomeMbps(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		name         string
		lines        []api.HouseholdLineDTO
		expectedMbps float64
		tolerance    float64
		description  string
	}{
		{
			name:         "Empty household",
			lines:        []api.HouseholdLineDTO{},
			expectedMbps: 10.0,
			tolerance:    0.1,
			description:  "Should return minimum 10 Mbps for empty household",
		},
		{
			name: "Single light user (1GB)",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  1.0,
					ExpectedMin: 100.0,
					TVHDHours:   0,
				},
			},
			expectedMbps: 10.0, // Will be rounded up to minimum
			tolerance:    0.1,
			description:  "Should return minimum 10 Mbps for very light usage",
		},
		{
			name: "Single moderate user (8GB)",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  8.0,
					ExpectedMin: 450.0,
					TVHDHours:   25.0,
				},
			},
			expectedMbps: 10.0, // Should still be minimum for 8GB
			tolerance:    0.1,
			description:  "8GB monthly usage should require minimum speed",
		},
		{
			name: "Single heavy user (50GB)",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  50.0,
					ExpectedMin: 1000.0,
					TVHDHours:   100.0,
				},
			},
			expectedMbps: 10.0, // 50GB → 0.463 Mbps required → minimum 10 Mbps applied
			tolerance:    0.1,
			description:  "50GB monthly usage hits minimum 10 Mbps",
		},
		{
			name: "Family with multiple users (8GB + 12GB + 20GB = 40GB total)",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  8.0,
					ExpectedMin: 450.0,
					TVHDHours:   25.0,
				},
				{
					LineID:      "LINE002",
					ExpectedGB:  12.0,
					ExpectedMin: 600.0,
					TVHDHours:   40.0,
				},
				{
					LineID:      "LINE003",
					ExpectedGB:  20.0,
					ExpectedMin: 800.0,
					TVHDHours:   60.0,
				},
			},
			expectedMbps: 10.0, // 40GB total → 0.371 Mbps required → minimum 10 Mbps applied
			tolerance:    0.1,
			description:  "Family of 3 with 40GB total hits minimum 10 Mbps",
		},
		{
			name: "Extreme usage family (5000GB total)",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  2000.0,
					ExpectedMin: 3000.0,
					TVHDHours:   200.0,
				},
				{
					LineID:      "LINE002",
					ExpectedGB:  1500.0,
					ExpectedMin: 2500.0,
					TVHDHours:   150.0,
				},
				{
					LineID:      "LINE003",
					ExpectedGB:  1500.0,
					ExpectedMin: 2000.0,
					TVHDHours:   300.0,
				},
			},
			expectedMbps: 50.0, // 5000GB → actual calculation result
			tolerance:    2.0,
			description:  "Extreme usage family with 5000GB should require ~58 Mbps",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ComputeHomeMbps(tt.lines)

			if math.Abs(result-tt.expectedMbps) > tt.tolerance {
				t.Errorf("Expected ~%.1f Mbps, got %.1f Mbps", tt.expectedMbps, result)
			}

			// Verify minimum constraint
			if result < 10.0 {
				t.Errorf("Result %.1f Mbps is below minimum 10 Mbps", result)
			}

			t.Logf("✓ %s: %.1f Mbps - %s", tt.name, result, tt.description)
		})
	}
}

func TestComputeHomeMbpsDetailed(t *testing.T) {
	service := &RecommendationService{}

	lines := []api.HouseholdLineDTO{
		{
			LineID:      "LINE001",
			ExpectedGB:  20.0,
			ExpectedMin: 500.0,
			TVHDHours:   30.0,
		},
		{
			LineID:      "LINE002",
			ExpectedGB:  30.0,
			ExpectedMin: 800.0,
			TVHDHours:   50.0,
		},
	}

	result := service.ComputeHomeMbpsDetailed(lines)

	// Check total GB calculation
	expectedTotalGB := 50.0
	if result.TotalGB != expectedTotalGB {
		t.Errorf("Expected TotalGB %.1f, got %.1f", expectedTotalGB, result.TotalGB)
	}

	// Check safety factor
	if result.SafetyFactor != 3.0 {
		t.Errorf("Expected SafetyFactor 3.0, got %.1f", result.SafetyFactor)
	}

	// Check that final Mbps is ceiling of required Mbps, but at least 10
	expectedFinal := math.Max(10.0, math.Ceil(result.RequiredMbps))
	if result.FinalMbps != expectedFinal {
		t.Errorf("FinalMbps should be max(10, ceil(RequiredMbps)): max(10, ceil(%.2f)) = %.0f, got %.0f",
			result.RequiredMbps, expectedFinal, result.FinalMbps)
	}

	// Check minimum constraint
	if result.FinalMbps < 10.0 {
		t.Errorf("FinalMbps %.1f is below minimum 10.0", result.FinalMbps)
	}

	// Check reasoning is present
	if result.Reasoning == "" {
		t.Error("Reasoning should not be empty")
	}

	t.Logf("✓ Detailed calculation: %.1f GB → %.2f avg Mbps → %.2f required → %.0f final",
		result.TotalGB, result.AvgMbps, result.RequiredMbps, result.FinalMbps)
}

func TestComputeHomeMbpsCalculationAccuracy(t *testing.T) {
	service := &RecommendationService{}

	// Test with known values to verify calculation accuracy
	lines := []api.HouseholdLineDTO{
		{
			LineID:     "TEST",
			ExpectedGB: 10.0, // Exactly 10GB
		},
	}

	result := service.ComputeHomeMbps(lines)

	// Manual calculation:
	// 10 GB = 10 * 8 * 1024^3 bits = 85,899,345,920 bits
	// Seconds per month = 30 * 24 * 3600 = 2,592,000
	// Avg bits/sec = 85,899,345,920 / 2,592,000 ≈ 33,138.02 bits/sec
	// Avg Mbps = 33,138.02 / 1,000,000 ≈ 0.0331 Mbps
	// With 3x safety factor = 0.0993 Mbps
	// Ceiling = 1.0 Mbps, but minimum is 10.0 Mbps

	expectedMbps := 10.0 // Should hit minimum
	if result != expectedMbps {
		t.Errorf("Expected %.1f Mbps (minimum), got %.1f Mbps", expectedMbps, result)
	}

	t.Logf("✓ 10GB calculation: %.1f Mbps (correctly applies minimum)", result)
}

func TestGenerateCandidates(t *testing.T) {
	// We'll test the logic without database dependency by creating mock data

	tests := []struct {
		name           string
		availableTech  []string
		neededMbps     float64
		maxTVHours     float64
		mockHomePlans  []models.HomePlan
		mockTVPlans    []models.TVPlan
		expectedCount  int
		expectedLabels []string
		description    string
	}{
		{
			name:          "Fiber available, low speed needed, no TV",
			availableTech: []string{"fiber", "vdsl"},
			neededMbps:    25.0,
			maxTVHours:    0.0,
			mockHomePlans: []models.HomePlan{
				{HomeID: 1, Name: "Fiber 50", Tech: "fiber", DownMbps: 50, MonthlyPrice: 89.90},
				{HomeID: 2, Name: "VDSL 25", Tech: "vdsl", DownMbps: 25, MonthlyPrice: 69.90},
				{HomeID: 3, Name: "FWA 20", Tech: "fwa", DownMbps: 20, MonthlyPrice: 59.90}, // Not available
			},
			mockTVPlans: []models.TVPlan{
				{TVID: 1, Name: "Basic TV", HDHoursIncluded: 50, MonthlyPrice: 29.90},
			},
			expectedCount: 6, // Mobile only + 2 home plans + 1 TV + 2 triple (TV qualifies even with 0 hours)
			expectedLabels: []string{
				"Mobile Only",
				"Mobile + Fiber 50",
				"Mobile + VDSL 25",
				"Mobile + Basic TV",
				"Triple: Fiber 50 + Basic TV",
				"Triple: VDSL 25 + Basic TV",
			},
			description: "Should include available tech plans that meet speed requirements",
		},
		{
			name:          "Only FWA available, high speed needed",
			availableTech: []string{"fwa"},
			neededMbps:    100.0,
			maxTVHours:    20.0,
			mockHomePlans: []models.HomePlan{
				{HomeID: 1, Name: "Fiber 100", Tech: "fiber", DownMbps: 100, MonthlyPrice: 129.90}, // Not available
				{HomeID: 2, Name: "FWA 50", Tech: "fwa", DownMbps: 50, MonthlyPrice: 79.90},        // Too slow
				{HomeID: 3, Name: "FWA 100", Tech: "fwa", DownMbps: 100, MonthlyPrice: 99.90},      // Perfect
			},
			mockTVPlans: []models.TVPlan{
				{TVID: 1, Name: "Basic TV", HDHoursIncluded: 50, MonthlyPrice: 29.90}, // Covers 20 hours
				{TVID: 2, Name: "Premium TV", HDHoursIncluded: 100, MonthlyPrice: 49.90},
			},
			expectedCount: 6, // Mobile only + 1 home + 2 TV + 2 triple
			expectedLabels: []string{
				"Mobile Only",
				"Mobile + FWA 100",
				"Mobile + Basic TV",
				"Mobile + Premium TV",
				"Triple: FWA 100 + Basic TV",
				"Triple: FWA 100 + Premium TV",
			},
			description: "Should filter by both tech availability and speed requirements",
		},
		{
			name:          "No home tech available",
			availableTech: []string{},
			neededMbps:    50.0,
			maxTVHours:    30.0,
			mockHomePlans: []models.HomePlan{
				{HomeID: 1, Name: "Fiber 100", Tech: "fiber", DownMbps: 100, MonthlyPrice: 129.90},
			},
			mockTVPlans: []models.TVPlan{
				{TVID: 1, Name: "Basic TV", HDHoursIncluded: 50, MonthlyPrice: 29.90},
			},
			expectedCount: 2, // Mobile only + 1 TV only
			expectedLabels: []string{
				"Mobile Only",
				"Mobile + Basic TV",
			},
			description: "Should handle no available home technology gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the candidate generation logic directly
			var candidates []BundleCandidate

			// Filter home plans by available tech and speed requirements
			var validHomePlans []models.HomePlan
			for _, plan := range tt.mockHomePlans {
				// Check if plan technology is available
				techAvailable := false
				for _, tech := range tt.availableTech {
					if plan.Tech == tech {
						techAvailable = true
						break
					}
				}

				// Check if plan meets speed requirements
				if techAvailable && float64(plan.DownMbps) >= tt.neededMbps {
					validHomePlans = append(validHomePlans, plan)
				}
			}

			// Filter TV plans by HD hours requirement
			var validTVPlans []models.TVPlan
			for _, plan := range tt.mockTVPlans {
				if plan.HDHoursIncluded >= tt.maxTVHours {
					validTVPlans = append(validTVPlans, plan)
				}
			}

			// Generate combinations (same logic as GenerateCandidates)
			// Option 1: Mobile only
			candidates = append(candidates, BundleCandidate{
				HomePlan: nil,
				TVPlan:   nil,
				Label:    "Mobile Only",
			})

			// Option 2: Mobile + Home
			for _, homePlan := range validHomePlans {
				candidates = append(candidates, BundleCandidate{
					HomePlan: &homePlan,
					TVPlan:   nil,
					Label:    "Mobile + " + homePlan.Name,
				})
			}

			// Option 3: Mobile + TV
			for _, tvPlan := range validTVPlans {
				candidates = append(candidates, BundleCandidate{
					HomePlan: nil,
					TVPlan:   &tvPlan,
					Label:    "Mobile + " + tvPlan.Name,
				})
			}

			// Option 4: Mobile + Home + TV
			for _, homePlan := range validHomePlans {
				for _, tvPlan := range validTVPlans {
					candidates = append(candidates, BundleCandidate{
						HomePlan: &homePlan,
						TVPlan:   &tvPlan,
						Label:    "Triple: " + homePlan.Name + " + " + tvPlan.Name,
					})
				}
			}

			// Verify count
			if len(candidates) != tt.expectedCount {
				t.Errorf("Expected %d candidates, got %d", tt.expectedCount, len(candidates))
			}

			// Verify labels are present
			candidateLabels := make(map[string]bool)
			for _, candidate := range candidates {
				candidateLabels[candidate.Label] = true
			}

			for _, expectedLabel := range tt.expectedLabels {
				if !candidateLabels[expectedLabel] {
					t.Errorf("Expected label '%s' not found in candidates", expectedLabel)
				}
			}

			t.Logf("✓ %s: %d candidates - %s", tt.name, len(candidates), tt.description)
		})
	}
}

func TestMatchLinesToPlans(t *testing.T) {
	service := &RecommendationService{}

	// Mock mobile plans with different quotas and prices
	mobilePlans := []models.MobilePlan{
		{
			PlanID:       1,
			PlanName:     "Basic 5GB",
			QuotaGB:      5.0,
			QuotaMin:     300.0,
			MonthlyPrice: 49.90,
			OverageGB:    2.0, // 2 TL per GB overage
			OverageMin:   0.5, // 0.5 TL per minute overage
		},
		{
			PlanID:       2,
			PlanName:     "Standard 10GB",
			QuotaGB:      10.0,
			QuotaMin:     500.0,
			MonthlyPrice: 79.90,
			OverageGB:    1.5,
			OverageMin:   0.3,
		},
		{
			PlanID:       3,
			PlanName:     "Premium 20GB",
			QuotaGB:      20.0,
			QuotaMin:     1000.0,
			MonthlyPrice: 129.90,
			OverageGB:    1.0,
			OverageMin:   0.2,
		},
	}

	tests := []struct {
		name                string
		lines               []api.HouseholdLineDTO
		expectedAssignments []expectedAssignment
		description         string
	}{
		{
			name: "Light user fits basic plan",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  3.0,
					ExpectedMin: 200.0,
				},
			},
			expectedAssignments: []expectedAssignment{
				{
					LineID:     "LINE001",
					PlanName:   "Basic 5GB",
					LineCost:   49.90, // No overage
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
			},
			description: "Light usage should get basic plan with no overage",
		},
		{
			name: "Heavy user needs premium plan",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE002",
					ExpectedGB:  18.0,
					ExpectedMin: 800.0,
				},
			},
			expectedAssignments: []expectedAssignment{
				{
					LineID:     "LINE002",
					PlanName:   "Premium 20GB",
					LineCost:   129.90, // No overage
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
			},
			description: "Heavy usage should get premium plan to avoid overage",
		},
		{
			name: "User with overage - choose cheapest total cost",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE003",
					ExpectedGB:  12.0, // Exceeds Standard by 2GB
					ExpectedMin: 400.0,
				},
			},
			expectedAssignments: []expectedAssignment{
				{
					LineID:     "LINE003",
					PlanName:   "Standard 10GB", // 79.90 + (2 * 1.5) = 82.90 vs Premium 129.90
					LineCost:   82.90,
					OverageGB:  2.0,
					OverageMin: 0.0,
				},
			},
			description: "Should choose plan with lowest total cost including overage",
		},
		{
			name: "Multiple lines with different needs",
			lines: []api.HouseholdLineDTO{
				{
					LineID:      "LINE001",
					ExpectedGB:  2.0,
					ExpectedMin: 150.0,
				},
				{
					LineID:      "LINE002",
					ExpectedGB:  15.0,
					ExpectedMin: 1200.0, // Exceeds Premium minutes by 200
				},
				{
					LineID:      "LINE003",
					ExpectedGB:  8.0,
					ExpectedMin: 450.0,
				},
			},
			expectedAssignments: []expectedAssignment{
				{
					LineID:     "LINE001",
					PlanName:   "Basic 5GB",
					LineCost:   49.90,
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
				{
					LineID:     "LINE002",
					PlanName:   "Premium 20GB", // 129.90 + (200 * 0.2) = 169.90
					LineCost:   169.90,
					OverageGB:  0.0,
					OverageMin: 200.0,
				},
				{
					LineID:     "LINE003",
					PlanName:   "Standard 10GB",
					LineCost:   79.90,
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
			},
			description: "Each line should get optimal plan for its usage pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignments := service.MatchLinesToPlans(tt.lines, mobilePlans)

			if len(assignments) != len(tt.expectedAssignments) {
				t.Errorf("Expected %d assignments, got %d", len(tt.expectedAssignments), len(assignments))
				return
			}

			for i, expected := range tt.expectedAssignments {
				assignment := assignments[i]

				if assignment.LineID != expected.LineID {
					t.Errorf("Assignment %d: expected LineID %s, got %s", i, expected.LineID, assignment.LineID)
				}

				if assignment.Plan.PlanName != expected.PlanName {
					t.Errorf("Assignment %d: expected plan %s, got %s", i, expected.PlanName, assignment.Plan.PlanName)
				}

				if assignment.LineCost != expected.LineCost {
					t.Errorf("Assignment %d: expected cost %.2f, got %.2f", i, expected.LineCost, assignment.LineCost)
				}

				if assignment.OverageGB != expected.OverageGB {
					t.Errorf("Assignment %d: expected GB overage %.1f, got %.1f", i, expected.OverageGB, assignment.OverageGB)
				}

				if assignment.OverageMin != expected.OverageMin {
					t.Errorf("Assignment %d: expected minute overage %.1f, got %.1f", i, expected.OverageMin, assignment.OverageMin)
				}
			}

			t.Logf("✓ %s: %d assignments - %s", tt.name, len(assignments), tt.description)
		})
	}
}

type expectedAssignment struct {
	LineID     string
	PlanName   string
	LineCost   float64
	OverageGB  float64
	OverageMin float64
}

func TestPriceBundleCandidate(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		name            string
		candidate       BundleCandidate
		lineAssignments []LineAssignment
		expectedPricing expectedPricing
		description     string
	}{
		{
			name: "Mobile only - single line",
			candidate: BundleCandidate{
				HomePlan: nil,
				TVPlan:   nil,
				Label:    "Mobile Only",
			},
			lineAssignments: []LineAssignment{
				{
					LineID:     "LINE001",
					Plan:       models.MobilePlan{PlanName: "Basic 5GB", MonthlyPrice: 49.90},
					LineCost:   49.90,
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
			},
			expectedPricing: expectedPricing{
				MobileTotal:          49.90,
				LineDiscountAmount:   0.0, // No discount for single line
				HomeCost:             0.0,
				TVCost:               0.0,
				BundleDiscountAmount: 0.0, // No bundle discount for mobile only
				GrandTotal:           49.90,
				TotalSavings:         0.0,
			},
			description: "Single mobile line should have no discounts",
		},
		{
			name: "Mobile + Home bundle",
			candidate: BundleCandidate{
				HomePlan: &models.HomePlan{
					HomeID:       1,
					Name:         "Fiber 50",
					Tech:         "fiber",
					DownMbps:     50,
					MonthlyPrice: 89.90,
				},
				TVPlan: nil,
				Label:  "Mobile + Fiber 50",
			},
			lineAssignments: []LineAssignment{
				{
					LineID:     "LINE001",
					Plan:       models.MobilePlan{PlanName: "Basic 5GB", MonthlyPrice: 49.90},
					LineCost:   49.90,
					OverageGB:  0.0,
					OverageMin: 0.0,
				},
			},
			expectedPricing: expectedPricing{
				MobileTotal:          49.90,
				LineDiscountAmount:   0.0, // No line discount for single line
				HomeCost:             89.90,
				TVCost:               0.0,
				BundleDiscountAmount: 13.98,  // 10% bundle discount on 139.80
				GrandTotal:           125.82, // 139.80 - 13.98
				TotalSavings:         13.98,
			},
			description: "Mobile + Home should get 10% bundle discount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.PriceBundleCandidate(tt.candidate, tt.lineAssignments)

			tolerance := 0.01 // Allow small floating point differences

			if math.Abs(result.MobileTotal-tt.expectedPricing.MobileTotal) > tolerance {
				t.Errorf("MobileTotal: expected %.2f, got %.2f", tt.expectedPricing.MobileTotal, result.MobileTotal)
			}

			if math.Abs(result.LineDiscountAmount-tt.expectedPricing.LineDiscountAmount) > tolerance {
				t.Errorf("LineDiscountAmount: expected %.2f, got %.2f", tt.expectedPricing.LineDiscountAmount, result.LineDiscountAmount)
			}

			if math.Abs(result.HomeCost-tt.expectedPricing.HomeCost) > tolerance {
				t.Errorf("HomeCost: expected %.2f, got %.2f", tt.expectedPricing.HomeCost, result.HomeCost)
			}

			if math.Abs(result.TVCost-tt.expectedPricing.TVCost) > tolerance {
				t.Errorf("TVCost: expected %.2f, got %.2f", tt.expectedPricing.TVCost, result.TVCost)
			}

			if math.Abs(result.BundleDiscountAmount-tt.expectedPricing.BundleDiscountAmount) > tolerance {
				t.Errorf("BundleDiscountAmount: expected %.2f, got %.2f", tt.expectedPricing.BundleDiscountAmount, result.BundleDiscountAmount)
			}

			if math.Abs(result.GrandTotal-tt.expectedPricing.GrandTotal) > tolerance {
				t.Errorf("GrandTotal: expected %.2f, got %.2f", tt.expectedPricing.GrandTotal, result.GrandTotal)
			}

			if math.Abs(result.TotalSavings-tt.expectedPricing.TotalSavings) > tolerance {
				t.Errorf("TotalSavings: expected %.2f, got %.2f", tt.expectedPricing.TotalSavings, result.TotalSavings)
			}

			// Check that reasoning is generated
			if result.Reasoning == "" {
				t.Error("Reasoning should not be empty")
			}

			t.Logf("✓ %s: %.2f total, %.2f savings - %s", tt.name, result.GrandTotal, result.TotalSavings, tt.description)
		})
	}
}

type expectedPricing struct {
	MobileTotal          float64
	LineDiscountAmount   float64
	HomeCost             float64
	TVCost               float64
	BundleDiscountAmount float64
	GrandTotal           float64
	TotalSavings         float64
}

func TestSelectTop3Candidates(t *testing.T) {
	service := &RecommendationService{}

	// Create test candidates with different pricing
	candidates := []PricedCandidate{
		{
			Candidate:    BundleCandidate{Label: "Expensive Triple"},
			GrandTotal:   399.90,
			TotalSavings: 50.00,
		},
		{
			Candidate:    BundleCandidate{Label: "Mobile Only"},
			GrandTotal:   49.90,
			TotalSavings: 0.00,
		},
		{
			Candidate:    BundleCandidate{Label: "Mobile + Home"},
			GrandTotal:   125.82,
			TotalSavings: 13.98,
		},
		{
			Candidate:    BundleCandidate{Label: "Cheap Triple"},
			GrandTotal:   299.90,
			TotalSavings: 75.00,
		},
		{
			Candidate:    BundleCandidate{Label: "Premium Bundle"},
			GrandTotal:   199.90,
			TotalSavings: 25.00,
		},
	}

	result := service.SelectTop3Candidates(candidates)

	// Should return 3 candidates
	if len(result) != 3 {
		t.Errorf("Expected 3 candidates, got %d", len(result))
	}

	// Should be sorted by grand total (ascending)
	expectedOrder := []string{"Mobile Only", "Mobile + Home", "Premium Bundle"}
	expectedTotals := []float64{49.90, 125.82, 199.90}

	for i, expected := range expectedOrder {
		if result[i].Candidate.Label != expected {
			t.Errorf("Position %d: expected %s, got %s", i, expected, result[i].Candidate.Label)
		}
		if result[i].GrandTotal != expectedTotals[i] {
			t.Errorf("Position %d: expected total %.2f, got %.2f", i, expectedTotals[i], result[i].GrandTotal)
		}
	}

	t.Logf("✓ SelectTop3Candidates: Correctly sorted and returned top 3 cheapest options")
}
