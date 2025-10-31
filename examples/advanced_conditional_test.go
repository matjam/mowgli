package examples

import (
	"testing"
)

func TestValidateAdvancedConditional(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "OR condition - hasEmail true, email required",
			data: map[string]any{
				"hasEmail": true,
				"hasPhone": false,
				"email":    "user@example.com",
				"age":      25,
				"country":  "US",
			},
			wantValid: true,
		},
		{
			name: "OR condition - hasPhone true, phone required",
			data: map[string]any{
				"hasEmail": false,
				"hasPhone": true,
				"phone":    "1234567890",
				"age":      25,
				"country":  "US",
			},
			wantValid: true,
		},
		{
			name: "OR condition - both true, both validated",
			data: map[string]any{
				"hasEmail": true,
				"hasPhone": true,
				"email":    "user@example.com",
				"phone":    "1234567890",
				"age":      25,
				"country":  "US",
			},
			wantValid: true,
		},
		{
			name: "OR condition - email missing",
			data: map[string]any{
				"hasEmail": true,
				"hasPhone": false,
				"email":    "",
				"age":      25,
				"country":  "US",
			},
			wantValid: false,
		},
		{
			name: "Parentheses - minor in US needs verification",
			data: map[string]any{
				"hasEmail":            false,
				"hasPhone":            false,
				"age":                 16,
				"country":             "US",
				"requireVerification": false,
				"verificationCode":    "ABC123",
			},
			wantValid: true,
		},
		{
			name: "Parentheses - minor in US without verification code",
			data: map[string]any{
				"hasEmail":            false,
				"hasPhone":            false,
				"age":                 16,
				"country":             "US",
				"requireVerification": false,
				"verificationCode":    "",
			},
			wantValid: false,
		},
		{
			name: "Parentheses - adult in US doesn't need verification",
			data: map[string]any{
				"hasEmail":            false,
				"hasPhone":            false,
				"age":                 25,
				"country":             "US",
				"requireVerification": false,
			},
			wantValid: true,
		},
		{
			name: "Parentheses - requireVerification true needs code",
			data: map[string]any{
				"hasEmail":            false,
				"hasPhone":            false,
				"age":                 25,
				"country":             "CA",
				"requireVerification": true,
				"verificationCode":    "XYZ789",
			},
			wantValid: true,
		},
		{
			name: "Complex - AND condition both email and phone validated",
			data: map[string]any{
				"hasEmail": true,
				"hasPhone": true,
				"email":    "user@example.com",
				"phone":    "1234567890123456", // Too long when both are required
				"age":      25,
				"country":  "US",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAdvancedConditional(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.wantValid {
				if !tt.wantValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
					for _, err := range result.Errors {
						t.Logf("  error: %s", err)
					}
				}
			}
		})
	}
}
