package utils

import (
	"testing"

	"app/internal/models"
)

func TestBuildReasoning(t *testing.T) {
	// Test plans
	fiberPlan := &models.HomePlan{
		HomeID: 1, Name: "Fiber 100Mbps", Tech: "fiber", DownMbps: 100, MonthlyPrice: 119.90,
	}

	tvPlan := &models.TVPlan{
		TVID: 2, Name: "Standard TV", HDHoursIncluded: 60.0, MonthlyPrice: 59.90,
	}

	tests := []struct {
		name              string
		summary           RecommendationSummary
		expectedReasoning string
		description       string
	}{
		{
			name: "Single mobile line only",
			summary: RecommendationSummary{
				Lines: []LineAssignment{
					{LineID: "LINE001", Usage: LineUsage{ExpectedGB: 5, ExpectedMin: 300}, Plan: models.MobilePlan{PlanName: "Basic Plan"}},
				},
				HomePlan:           nil,
				TVPlan:             nil,
				LineDiscountAmount: 0,
				LineCount:          1,
				Breakdown:          GrandTotalBreakdown{BundleDiscountRate: 0},
			},
			expectedReasoning: "Selected plans: 1 mobile line(s).",
			description:       "Should show single mobile line without discounts",
		},
		{
			name: "Triple bundle with all discounts",
			summary: RecommendationSummary{
				Lines: []LineAssignment{
					{LineID: "LINE001", Usage: LineUsage{ExpectedGB: 5, ExpectedMin: 300}, Plan: models.MobilePlan{PlanName: "Basic Plan"}},
					{LineID: "LINE002", Usage: LineUsage{ExpectedGB: 8, ExpectedMin: 400}, Plan: models.MobilePlan{PlanName: "Standard Plan"}},
					{LineID: "LINE003", Usage: LineUsage{ExpectedGB: 12, ExpectedMin: 600}, Plan: models.MobilePlan{PlanName: "Premium Plan"}},
				},
				HomePlan:           fiberPlan,
				TVPlan:             tvPlan,
				LineDiscountAmount: 30.0,
				LineCount:          3,
				Breakdown:          GrandTotalBreakdown{BundleDiscountRate: 0.15},
			},
			expectedReasoning: "Selected plans: 3 mobile line(s) (10% multi-line discount), Fiber 100Mbps, Standard TV. Bundle discount: 15% off total.",
			description:       "Should show all services with both line and bundle discounts",
		},
		{
			name: "Mobile + Home with bundle discount",
			summary: RecommendationSummary{
				Lines: []LineAssignment{
					{LineID: "LINE001", Usage: LineUsage{ExpectedGB: 5, ExpectedMin: 300}, Plan: models.MobilePlan{PlanName: "Basic Plan"}},
					{LineID: "LINE002", Usage: LineUsage{ExpectedGB: 8, ExpectedMin: 400}, Plan: models.MobilePlan{PlanName: "Standard Plan"}},
				},
				HomePlan:           fiberPlan,
				TVPlan:             nil,
				LineDiscountAmount: 10.0,
				LineCount:          2,
				Breakdown:          GrandTotalBreakdown{BundleDiscountRate: 0.10},
			},
			expectedReasoning: "Selected plans: 2 mobile line(s) (5% multi-line discount), Fiber 100Mbps. Bundle discount: 10% off total.",
			description:       "Should show mobile + home with appropriate discounts",
		},
		{
			name: "Home + TV only (no mobile)",
			summary: RecommendationSummary{
				Lines:              []LineAssignment{},
				HomePlan:           fiberPlan,
				TVPlan:             tvPlan,
				LineDiscountAmount: 0,
				LineCount:          0,
				Breakdown:          GrandTotalBreakdown{BundleDiscountRate: 0},
			},
			expectedReasoning: "Selected plans: , Fiber 100Mbps, Standard TV.",
			description:       "Should handle case with no mobile lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildReasoning(tt.summary)

			if result != tt.expectedReasoning {
				t.Errorf("BuildReasoning() = %q, want %q", result, tt.expectedReasoning)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
