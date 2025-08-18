package utils

import (
	"testing"

	"app/internal/models"
)

func TestSelectHomePlan(t *testing.T) {
	// Test home plans
	homePlans := []models.HomePlan{
		{HomeID: 1, Name: "Fiber 50Mbps", Tech: "fiber", DownMbps: 50, MonthlyPrice: 89.90},
		{HomeID: 2, Name: "Fiber 100Mbps", Tech: "fiber", DownMbps: 100, MonthlyPrice: 119.90},
		{HomeID: 3, Name: "Fiber 200Mbps", Tech: "fiber", DownMbps: 200, MonthlyPrice: 159.90},
		{HomeID: 4, Name: "VDSL 25Mbps", Tech: "vdsl", DownMbps: 25, MonthlyPrice: 69.90},
		{HomeID: 5, Name: "VDSL 50Mbps", Tech: "vdsl", DownMbps: 50, MonthlyPrice: 89.90},
		{HomeID: 6, Name: "FWA 30Mbps", Tech: "fwa", DownMbps: 30, MonthlyPrice: 79.90},
		{HomeID: 7, Name: "FWA 50Mbps", Tech: "fwa", DownMbps: 50, MonthlyPrice: 99.90},
	}

	tests := []struct {
		name           string
		techAvailable  []string
		neededMbps     int
		expectedPlanID *int // nil if no plan expected
		description    string
	}{
		{
			name:           "Fiber available - choose cheapest fiber plan meeting requirements",
			techAvailable:  []string{"fiber", "vdsl", "fwa"},
			neededMbps:     75,        // Needs more than 50, but less than 100
			expectedPlanID: intPtr(2), // Fiber 100Mbps (cheapest that meets 75+ Mbps)
			description:    "Should select cheapest fiber plan that meets speed requirement",
		},
		{
			name:           "Fiber missing - fallback to VDSL",
			techAvailable:  []string{"vdsl", "fwa"},
			neededMbps:     40,        // Needs more than 25, but less than 50
			expectedPlanID: intPtr(5), // VDSL 50Mbps (meets requirement)
			description:    "Should fallback to VDSL when fiber not available",
		},
		{
			name:           "Only FWA available - choose FWA",
			techAvailable:  []string{"fwa"},
			neededMbps:     35,        // Needs more than 30, but less than 50
			expectedPlanID: intPtr(7), // FWA 50Mbps (meets requirement)
			description:    "Should select FWA when only FWA is available",
		},
		{
			name:           "No technology can meet requirement",
			techAvailable:  []string{"vdsl"},
			neededMbps:     100, // Needs 100+ Mbps, but VDSL max is 50
			expectedPlanID: nil, // No plan can meet requirement
			description:    "Should return nil when no plan can meet speed requirement",
		},
		{
			name:           "Prefer fiber over higher speed VDSL",
			techAvailable:  []string{"fiber", "vdsl"},
			neededMbps:     25,        // Both fiber 50 and VDSL 50 can meet this
			expectedPlanID: intPtr(1), // Fiber 50Mbps (fiber preferred over VDSL)
			description:    "Should prefer fiber technology over VDSL even if VDSL is available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectHomePlan(tt.techAvailable, tt.neededMbps, homePlans)

			if tt.expectedPlanID == nil {
				if result != nil {
					t.Errorf("Expected no plan, got plan ID %d", result.HomeID)
				}
			} else {
				if result == nil {
					t.Errorf("Expected plan ID %d, got nil", *tt.expectedPlanID)
				} else if result.HomeID != *tt.expectedPlanID {
					t.Errorf("Expected plan ID %d, got %d", *tt.expectedPlanID, result.HomeID)
				}
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
