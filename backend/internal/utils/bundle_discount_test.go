package utils

import "testing"

func TestCalcBundleDiscount(t *testing.T) {
	tests := []struct {
		name             string
		hasMobile        bool
		hasHome          bool
		hasTV            bool
		expectedDiscount float64
		description      string
	}{
		{
			name:             "Triple bundle - mobile + home + TV",
			hasMobile:        true,
			hasHome:          true,
			hasTV:            true,
			expectedDiscount: 0.15, // 15%
			description:      "Triple bundle should receive 15% discount",
		},
		{
			name:             "Pair bundle - mobile + home only",
			hasMobile:        true,
			hasHome:          true,
			hasTV:            false,
			expectedDiscount: 0.10, // 10%
			description:      "Mobile + home bundle should receive 10% discount",
		},
		{
			name:             "No bundle - mobile only",
			hasMobile:        true,
			hasHome:          false,
			hasTV:            false,
			expectedDiscount: 0.0, // 0%
			description:      "Mobile only should receive no bundle discount",
		},
		{
			name:             "No bundle - home only",
			hasMobile:        false,
			hasHome:          true,
			hasTV:            false,
			expectedDiscount: 0.0, // 0%
			description:      "Home only should receive no bundle discount",
		},
		{
			name:             "No bundle - TV only",
			hasMobile:        false,
			hasHome:          false,
			hasTV:            true,
			expectedDiscount: 0.0, // 0%
			description:      "TV only should receive no bundle discount",
		},
		{
			name:             "Mobile + TV (no home)",
			hasMobile:        true,
			hasHome:          false,
			hasTV:            true,
			expectedDiscount: 0.0, // 0%
			description:      "Mobile + TV without home should receive no bundle discount",
		},
		{
			name:             "Home + TV (no mobile)",
			hasMobile:        false,
			hasHome:          true,
			hasTV:            true,
			expectedDiscount: 0.0, // 0%
			description:      "Home + TV without mobile should receive no bundle discount",
		},
		{
			name:             "No services",
			hasMobile:        false,
			hasHome:          false,
			hasTV:            false,
			expectedDiscount: 0.0, // 0%
			description:      "No services should receive no bundle discount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discount := CalcBundleDiscount(tt.hasMobile, tt.hasHome, tt.hasTV)

			if discount != tt.expectedDiscount {
				t.Errorf("CalcBundleDiscount() = %v, want %v", discount, tt.expectedDiscount)
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
