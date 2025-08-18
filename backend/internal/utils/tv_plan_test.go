package utils

import (
	"testing"

	"app/internal/models"
)

func TestSelectTvPlan(t *testing.T) {
	// Test TV plans
	tvPlans := []models.TVPlan{
		{TVID: 1, Name: "Basic TV", HDHoursIncluded: 30.0, MonthlyPrice: 39.90},
		{TVID: 2, Name: "Standard TV", HDHoursIncluded: 60.0, MonthlyPrice: 59.90},
		{TVID: 3, Name: "Premium TV", HDHoursIncluded: 120.0, MonthlyPrice: 89.90},
		{TVID: 4, Name: "Sports Package", HDHoursIncluded: 150.0, MonthlyPrice: 119.90},
	}

	tests := []struct {
		name           string
		requiredHours  float64
		expectedPlanID *int // nil if no plan expected
		description    string
	}{
		{
			name:           "Low usage - 20 hours",
			requiredHours:  20.0,
			expectedPlanID: intPtr(1), // Basic TV (30 hours, cheapest that covers 20)
			description:    "Should select Basic TV for low HD hours requirement",
		},
		{
			name:           "Medium usage - 60 hours",
			requiredHours:  60.0,
			expectedPlanID: intPtr(2), // Standard TV (exactly 60 hours)
			description:    "Should select Standard TV for medium HD hours requirement",
		},
		{
			name:           "High usage - 120 hours",
			requiredHours:  120.0,
			expectedPlanID: intPtr(3), // Premium TV (exactly 120 hours, cheaper than Sports)
			description:    "Should select Premium TV for high HD hours requirement",
		},
		{
			name:           "Very high usage - 140 hours",
			requiredHours:  140.0,
			expectedPlanID: intPtr(4), // Sports Package (150 hours, only one that covers 140)
			description:    "Should select Sports Package for very high HD hours requirement",
		},
		{
			name:           "No TV needed - 0 hours",
			requiredHours:  0.0,
			expectedPlanID: nil, // No TV plan needed
			description:    "Should return nil when no TV hours are required",
		},
		{
			name:           "Requirement exceeds all plans",
			requiredHours:  200.0,
			expectedPlanID: nil, // No plan can cover 200 hours
			description:    "Should return nil when requirement exceeds all available plans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectTvPlan(tt.requiredHours, tvPlans)
			
			if tt.expectedPlanID == nil {
				if result != nil {
					t.Errorf("Expected no plan, got plan ID %d", result.TVID)
				}
			} else {
				if result == nil {
					t.Errorf("Expected plan ID %d, got nil", *tt.expectedPlanID)
				} else if result.TVID != *tt.expectedPlanID {
					t.Errorf("Expected plan ID %d, got %d", *tt.expectedPlanID, result.TVID)
				}
			}
			
			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
