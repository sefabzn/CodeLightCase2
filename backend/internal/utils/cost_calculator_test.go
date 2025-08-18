package utils

import (
	"testing"

	"app/internal/models"
)

func TestCalcMobileOverage(t *testing.T) {
	// Test plan with specific overage rates
	testPlan := models.MobilePlan{
		PlanID:       1,
		PlanName:     "Test Plan",
		QuotaGB:      10.0,  // 10GB quota
		QuotaMin:     500.0, // 500 minutes quota
		MonthlyPrice: 100.0,
		OverageGB:    5.0, // 5 TL per GB over
		OverageMin:   0.5, // 0.5 TL per minute over
	}

	tests := []struct {
		name           string
		usage          LineUsage
		expectedResult OverageResult
		description    string
	}{
		{
			name: "No overage - under both quotas",
			usage: LineUsage{
				ExpectedGB:  8.0,   // Under 10GB quota
				ExpectedMin: 400.0, // Under 500min quota
			},
			expectedResult: OverageResult{
				OverageGB:   0.0,
				OverageMin:  0.0,
				OverageCost: 0.0,
			},
			description: "Usage under both quotas should result in no overage",
		},
		{
			name: "GB overage only",
			usage: LineUsage{
				ExpectedGB:  15.0,  // 5GB over quota
				ExpectedMin: 400.0, // Under minutes quota
			},
			expectedResult: OverageResult{
				OverageGB:   5.0, // 15 - 10 = 5GB over
				OverageMin:  0.0,
				OverageCost: 25.0, // 5GB * 5 TL/GB = 25 TL
			},
			description: "GB overage should be calculated correctly",
		},
		{
			name: "Minutes overage only",
			usage: LineUsage{
				ExpectedGB:  8.0,   // Under GB quota
				ExpectedMin: 650.0, // 150 minutes over quota
			},
			expectedResult: OverageResult{
				OverageGB:   0.0,
				OverageMin:  150.0, // 650 - 500 = 150min over
				OverageCost: 75.0,  // 150min * 0.5 TL/min = 75 TL
			},
			description: "Minutes overage should be calculated correctly",
		},
		{
			name: "Both GB and minutes overage",
			usage: LineUsage{
				ExpectedGB:  12.0,  // 2GB over quota
				ExpectedMin: 600.0, // 100 minutes over quota
			},
			expectedResult: OverageResult{
				OverageGB:   2.0,   // 12 - 10 = 2GB over
				OverageMin:  100.0, // 600 - 500 = 100min over
				OverageCost: 60.0,  // (2 * 5) + (100 * 0.5) = 10 + 50 = 60 TL
			},
			description: "Both overages should be calculated and summed correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcMobileOverage(tt.usage, testPlan)

			// Check GB overage
			if result.OverageGB != tt.expectedResult.OverageGB {
				t.Errorf("OverageGB = %v, want %v", result.OverageGB, tt.expectedResult.OverageGB)
			}

			// Check minutes overage
			if result.OverageMin != tt.expectedResult.OverageMin {
				t.Errorf("OverageMin = %v, want %v", result.OverageMin, tt.expectedResult.OverageMin)
			}

			// Check total overage cost
			if result.OverageCost != tt.expectedResult.OverageCost {
				t.Errorf("OverageCost = %v, want %v", result.OverageCost, tt.expectedResult.OverageCost)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

func TestCalcMobileLineCost(t *testing.T) {
	// Test plan
	testPlan := models.MobilePlan{
		PlanID:       1,
		PlanName:     "Test Plan",
		QuotaGB:      10.0,  // 10GB quota
		QuotaMin:     500.0, // 500 minutes quota
		MonthlyPrice: 100.0, // 100 TL base price
		OverageGB:    5.0,   // 5 TL per GB over
		OverageMin:   0.5,   // 0.5 TL per minute over
	}

	tests := []struct {
		name         string
		usage        LineUsage
		expectedCost float64
		description  string
	}{
		{
			name: "Exact quota usage",
			usage: LineUsage{
				ExpectedGB:  10.0,  // Exactly at quota
				ExpectedMin: 500.0, // Exactly at quota
			},
			expectedCost: 100.0, // Only monthly price, no overage
			description:  "Usage at quota should only charge monthly price",
		},
		{
			name: "Over GB only",
			usage: LineUsage{
				ExpectedGB:  15.0,  // 5GB over quota
				ExpectedMin: 400.0, // Under minutes quota
			},
			expectedCost: 125.0, // 100 (base) + 25 (5GB * 5 TL/GB)
			description:  "GB overage should be added to monthly price",
		},
		{
			name: "Over minutes only",
			usage: LineUsage{
				ExpectedGB:  8.0,   // Under GB quota
				ExpectedMin: 600.0, // 100 minutes over quota
			},
			expectedCost: 150.0, // 100 (base) + 50 (100min * 0.5 TL/min)
			description:  "Minutes overage should be added to monthly price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalcMobileLineCost(tt.usage, testPlan)

			if cost != tt.expectedCost {
				t.Errorf("CalcMobileLineCost() = %v, want %v", cost, tt.expectedCost)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

func TestCalcMobileTotal(t *testing.T) {
	// Test plans
	basicPlan := models.MobilePlan{
		PlanID:       1,
		PlanName:     "Basic Plan",
		QuotaGB:      5.0,
		QuotaMin:     300.0,
		MonthlyPrice: 99.90,
		OverageGB:    5.0,
		OverageMin:   0.5,
	}

	premiumPlan := models.MobilePlan{
		PlanID:       2,
		PlanName:     "Premium Plan",
		QuotaGB:      20.0,
		QuotaMin:     1000.0,
		MonthlyPrice: 199.90,
		OverageGB:    3.0,
		OverageMin:   0.3,
	}

	tests := []struct {
		name         string
		lines        []LineAssignment
		expectedCost float64
		description  string
	}{
		{
			name: "Single line no overage",
			lines: []LineAssignment{
				{
					LineID: "LINE001",
					Usage: LineUsage{
						ExpectedGB:  4.0,   // Under 5GB quota
						ExpectedMin: 250.0, // Under 300min quota
					},
					Plan: basicPlan,
				},
			},
			expectedCost: 99.90, // Only monthly price
			description:  "Single line with no overage should return plan monthly price",
		},
		{
			name: "Three lines with different plans and usage",
			lines: []LineAssignment{
				{
					LineID: "LINE001",
					Usage: LineUsage{
						ExpectedGB:  4.0,   // Under quota
						ExpectedMin: 250.0, // Under quota
					},
					Plan: basicPlan, // 99.90 TL
				},
				{
					LineID: "LINE002",
					Usage: LineUsage{
						ExpectedGB:  7.0,   // 2GB over quota (2 * 5 = 10 TL overage)
						ExpectedMin: 280.0, // Under quota
					},
					Plan: basicPlan, // 99.90 + 10 = 109.90 TL
				},
				{
					LineID: "LINE003",
					Usage: LineUsage{
						ExpectedGB:  18.0,  // Under 20GB quota
						ExpectedMin: 950.0, // Under 1000min quota
					},
					Plan: premiumPlan, // 199.90 TL
				},
			},
			expectedCost: 409.70, // 99.90 + 109.90 + 199.90 = 409.70 TL
			description:  "Multiple lines should sum all individual line costs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := CalcMobileTotal(tt.lines)

			// Use tolerance for floating point comparison
			tolerance := 0.01
			if diff := total - tt.expectedCost; diff < -tolerance || diff > tolerance {
				t.Errorf("CalcMobileTotal() = %v, want %v (diff: %v)", total, tt.expectedCost, diff)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
