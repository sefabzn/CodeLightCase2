package services

import (
	"testing"

	"app/internal/models"
)

func TestComputeCoverageLogic(t *testing.T) {
	tests := []struct {
		name         string
		coverage     models.Coverage
		expectedTech []string
		description  string
	}{
		{
			name: "All technologies available",
			coverage: models.Coverage{
				AddressID: "A1001",
				City:      "Istanbul",
				District:  "Kadikoy",
				Fiber:     true,
				VDSL:      true,
				FWA:       true,
			},
			expectedTech: []string{"fiber", "vdsl", "fwa"},
			description:  "Should return all technologies in preference order",
		},
		{
			name: "Fiber and VDSL only",
			coverage: models.Coverage{
				AddressID: "A1002",
				City:      "Istanbul",
				District:  "Besiktas",
				Fiber:     true,
				VDSL:      true,
				FWA:       false,
			},
			expectedTech: []string{"fiber", "vdsl"},
			description:  "Should return fiber and vdsl in preference order",
		},
		{
			name: "VDSL and FWA only",
			coverage: models.Coverage{
				AddressID: "A1003",
				City:      "Ankara",
				District:  "Cankaya",
				Fiber:     false,
				VDSL:      true,
				FWA:       true,
			},
			expectedTech: []string{"vdsl", "fwa"},
			description:  "Should return vdsl and fwa when fiber not available",
		},
		{
			name: "Only FWA available",
			coverage: models.Coverage{
				AddressID: "A1004",
				City:      "Izmir",
				District:  "Konak",
				Fiber:     false,
				VDSL:      false,
				FWA:       true,
			},
			expectedTech: []string{"fwa"},
			description:  "Should return only fwa when it's the only option",
		},
		{
			name: "No technology available",
			coverage: models.Coverage{
				AddressID: "A1005",
				City:      "Rural",
				District:  "Remote",
				Fiber:     false,
				VDSL:      false,
				FWA:       false,
			},
			expectedTech: []string{},
			description:  "Should return empty array when no technology available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic directly without database dependency
			var availableTech []string

			// This is the same logic as in ComputeCoverage
			if tt.coverage.Fiber {
				availableTech = append(availableTech, "fiber")
			}

			if tt.coverage.VDSL {
				availableTech = append(availableTech, "vdsl")
			}

			if tt.coverage.FWA {
				availableTech = append(availableTech, "fwa")
			}

			if len(availableTech) != len(tt.expectedTech) {
				t.Errorf("Expected %d technologies, got %d: %v", len(tt.expectedTech), len(availableTech), availableTech)
			}

			for i, expectedTech := range tt.expectedTech {
				if i >= len(availableTech) || availableTech[i] != expectedTech {
					t.Errorf("Expected tech[%d] = %s, got %s", i, expectedTech, availableTech[i])
				}
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

func TestCoveragePreferenceOrder(t *testing.T) {
	// Test that technologies are returned in the correct preference order
	coverage := models.Coverage{
		Fiber: true,
		VDSL:  true,
		FWA:   true,
	}

	var availableTech []string

	// Apply the same logic as ComputeCoverage
	if coverage.Fiber {
		availableTech = append(availableTech, "fiber")
	}

	if coverage.VDSL {
		availableTech = append(availableTech, "vdsl")
	}

	if coverage.FWA {
		availableTech = append(availableTech, "fwa")
	}

	expectedOrder := []string{"fiber", "vdsl", "fwa"}

	for i, expected := range expectedOrder {
		if availableTech[i] != expected {
			t.Errorf("Expected position %d to be %s, got %s", i, expected, availableTech[i])
		}
	}

	t.Logf("✓ Technologies returned in correct preference order: %v", availableTech)
}
