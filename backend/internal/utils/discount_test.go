package utils

import "testing"

func TestApplyExtraLineDiscount(t *testing.T) {
	tests := []struct {
		name                    string
		mobileTotal             float64
		lineCount               int
		expectedDiscountedTotal float64
		expectedDiscountAmount  float64
		description             string
	}{
		{
			name:                    "Single line - no discount",
			mobileTotal:             100.0,
			lineCount:               1,
			expectedDiscountedTotal: 100.0, // No discount
			expectedDiscountAmount:  0.0,
			description:             "Single line should not receive any discount",
		},
		{
			name:                    "Two lines - 5% discount",
			mobileTotal:             200.0,
			lineCount:               2,
			expectedDiscountedTotal: 190.0, // 200 - (200 * 0.05) = 190
			expectedDiscountAmount:  10.0,  // 200 * 0.05 = 10
			description:             "Second line should receive 5% discount",
		},
		{
			name:                    "Three lines - 10% discount",
			mobileTotal:             300.0,
			lineCount:               3,
			expectedDiscountedTotal: 270.0, // 300 - (300 * 0.10) = 270
			expectedDiscountAmount:  30.0,  // 300 * 0.10 = 30
			description:             "Three lines should receive 10% discount",
		},
		{
			name:                    "Five lines - 10% discount",
			mobileTotal:             500.0,
			lineCount:               5,
			expectedDiscountedTotal: 450.0, // 500 - (500 * 0.10) = 450
			expectedDiscountAmount:  50.0,  // 500 * 0.10 = 50
			description:             "More than three lines should still receive 10% discount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discountedTotal, discountAmount := ApplyExtraLineDiscount(tt.mobileTotal, tt.lineCount)

			// Use tolerance for floating point comparison
			tolerance := 0.01

			if diff := discountedTotal - tt.expectedDiscountedTotal; diff < -tolerance || diff > tolerance {
				t.Errorf("Discounted total = %v, want %v (diff: %v)", discountedTotal, tt.expectedDiscountedTotal, diff)
			}

			if diff := discountAmount - tt.expectedDiscountAmount; diff < -tolerance || diff > tolerance {
				t.Errorf("Discount amount = %v, want %v (diff: %v)", discountAmount, tt.expectedDiscountAmount, diff)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
