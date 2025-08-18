package utils

import "testing"

func TestCalcGrandTotal(t *testing.T) {
	tests := []struct {
		name                string
		mobileAfterLineDisc float64
		homeCost            float64
		tvCost              float64
		bundleDiscountRate  float64
		expectedBreakdown   GrandTotalBreakdown
		description         string
	}{
		{
			name:                "Mobile + Home bundle with 10% discount",
			mobileAfterLineDisc: 190.0, // Mobile after line discount
			homeCost:            89.90, // Home plan cost
			tvCost:              0.0,   // No TV
			bundleDiscountRate:  0.10,  // 10% bundle discount
			expectedBreakdown: GrandTotalBreakdown{
				MobileTotal:        190.0,
				HomeTotal:          89.90,
				TVTotal:            0.0,
				SubTotal:           279.90, // 190 + 89.90
				BundleDiscount:     27.99,  // 279.90 * 0.10
				GrandTotal:         251.91, // 279.90 - 27.99
				BundleDiscountRate: 0.10,
			},
			description: "Should apply 10% bundle discount on mobile + home combination",
		},
		{
			name:                "Triple bundle with 15% discount",
			mobileAfterLineDisc: 270.0,  // Mobile after line discount (3+ lines)
			homeCost:            119.90, // Home plan cost
			tvCost:              59.90,  // TV plan cost
			bundleDiscountRate:  0.15,   // 15% bundle discount
			expectedBreakdown: GrandTotalBreakdown{
				MobileTotal:        270.0,
				HomeTotal:          119.90,
				TVTotal:            59.90,
				SubTotal:           449.80, // 270 + 119.90 + 59.90
				BundleDiscount:     67.47,  // 449.80 * 0.15
				GrandTotal:         382.33, // 449.80 - 67.47
				BundleDiscountRate: 0.15,
			},
			description: "Should apply 15% bundle discount on triple bundle",
		},
		{
			name:                "Mobile only - no bundle discount",
			mobileAfterLineDisc: 99.90,
			homeCost:            0.0,
			tvCost:              0.0,
			bundleDiscountRate:  0.0, // No bundle discount
			expectedBreakdown: GrandTotalBreakdown{
				MobileTotal:        99.90,
				HomeTotal:          0.0,
				TVTotal:            0.0,
				SubTotal:           99.90,
				BundleDiscount:     0.0,
				GrandTotal:         99.90,
				BundleDiscountRate: 0.0,
			},
			description: "Should not apply bundle discount for mobile-only service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcGrandTotal(tt.mobileAfterLineDisc, tt.homeCost, tt.tvCost, tt.bundleDiscountRate)

			tolerance := 0.01

			// Check all fields with tolerance
			if diff := result.MobileTotal - tt.expectedBreakdown.MobileTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("MobileTotal = %v, want %v", result.MobileTotal, tt.expectedBreakdown.MobileTotal)
			}

			if diff := result.HomeTotal - tt.expectedBreakdown.HomeTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("HomeTotal = %v, want %v", result.HomeTotal, tt.expectedBreakdown.HomeTotal)
			}

			if diff := result.TVTotal - tt.expectedBreakdown.TVTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("TVTotal = %v, want %v", result.TVTotal, tt.expectedBreakdown.TVTotal)
			}

			if diff := result.SubTotal - tt.expectedBreakdown.SubTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("SubTotal = %v, want %v", result.SubTotal, tt.expectedBreakdown.SubTotal)
			}

			if diff := result.BundleDiscount - tt.expectedBreakdown.BundleDiscount; diff < -tolerance || diff > tolerance {
				t.Errorf("BundleDiscount = %v, want %v", result.BundleDiscount, tt.expectedBreakdown.BundleDiscount)
			}

			if diff := result.GrandTotal - tt.expectedBreakdown.GrandTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("GrandTotal = %v, want %v", result.GrandTotal, tt.expectedBreakdown.GrandTotal)
			}

			if result.BundleDiscountRate != tt.expectedBreakdown.BundleDiscountRate {
				t.Errorf("BundleDiscountRate = %v, want %v", result.BundleDiscountRate, tt.expectedBreakdown.BundleDiscountRate)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
