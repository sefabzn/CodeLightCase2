package utils

import (
	"testing"

	"app/internal/api"
)

func TestValidator(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name           string
		input          interface{}
		expectedErrors []string
		description    string
	}{
		{
			name: "Valid recommendation request",
			input: api.RecommendationRequest{
				UserID:    1,
				AddressID: "A1001",
				Household: []api.HouseholdLineDTO{
					{
						LineID:      "LINE001",
						ExpectedGB:  8.0,
						ExpectedMin: 450.0,
						TVHDHours:   25.0,
					},
				},
				PreferTech: []string{"fiber", "vdsl"},
			},
			expectedErrors: nil,
			description:    "Valid request should pass validation",
		},
		{
			name: "Missing required fields",
			input: api.RecommendationRequest{
				// Missing UserID and AddressID
				Household: []api.HouseholdLineDTO{},
			},
			expectedErrors: []string{
				"user_id is required",
				"address_id is required",
				"household must have at least 1 item(s)",
			},
			description: "Should catch all missing required fields",
		},
		{
			name: "Invalid household line",
			input: api.RecommendationRequest{
				UserID:    1,
				AddressID: "A1001",
				Household: []api.HouseholdLineDTO{
					{
						// Missing LineID, invalid negative values
						ExpectedGB:  -5.0,
						ExpectedMin: -100.0,
						TVHDHours:   -10.0,
					},
				},
			},
			expectedErrors: []string{
				"line_id is required",
				"expected_gb must be at least 0",
				"expected_min must be at least 0",
				"tv_hd_hours must be at least 0",
			},
			description: "Should validate household line fields",
		},
		{
			name: "Valid checkout request",
			input: api.CheckoutRequest{
				UserID:    1,
				SlotID:    123,
				AddressID: "A1001",
				SelectedCombo: api.RecommendationCandidateDTO{
					ComboLabel:   "Triple Bundle",
					MonthlyTotal: 299.90,
					Savings:      50.0,
					Reasoning:    "Best value triple bundle",
				},
			},
			expectedErrors: nil,
			description:    "Valid checkout request should pass validation",
		},
		{
			name:  "Invalid checkout request",
			input: api.CheckoutRequest{
				// Missing all required fields
			},
			expectedErrors: []string{
				"user_id is required",
				"slot_id is required",
				"address_id is required",
			},
			description: "Should catch missing checkout fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateStruct(tt.input)

			if tt.expectedErrors == nil {
				if len(errors) != 0 {
					t.Errorf("Expected no errors, got: %v", errors)
				}
			} else {
				if len(errors) != len(tt.expectedErrors) {
					t.Errorf("Expected %d errors, got %d: %v", len(tt.expectedErrors), len(errors), errors)
				}

				// Check if all expected errors are present
				errorMap := make(map[string]bool)
				for _, err := range errors {
					errorMap[err] = true
				}

				for _, expectedErr := range tt.expectedErrors {
					if !errorMap[expectedErr] {
						t.Errorf("Expected error '%s' not found in: %v", expectedErr, errors)
					}
				}
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}
