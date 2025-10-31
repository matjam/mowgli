package examples

import (
	"testing"
)

func TestValidateUserPreferences(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "notifications enabled with valid email",
			data: map[string]any{
				"notificationsEnabled": true,
				"email":                "user@example.com",
				"age":                  25,
			},
			wantValid: true,
		},
		{
			name: "notifications disabled, email can be empty",
			data: map[string]any{
				"notificationsEnabled": false,
				"email":                "",
				"age":                  25,
			},
			wantValid: true,
		},
		{
			name: "notifications enabled but email missing",
			data: map[string]any{
				"notificationsEnabled": true,
				"email":                "",
				"age":                  25,
			},
			wantValid: false,
		},
		{
			name: "notifications enabled but invalid email",
			data: map[string]any{
				"notificationsEnabled": true,
				"email":                "not-an-email",
				"age":                  25,
			},
			wantValid: false,
		},
		{
			name: "minor with parental consent",
			data: map[string]any{
				"notificationsEnabled":   false,
				"age":                    16,
				"requireParentalConsent": true,
				"parentEmail":            "parent@example.com",
			},
			wantValid: true,
		},
		{
			name: "minor without parental consent",
			data: map[string]any{
				"notificationsEnabled":   false,
				"age":                    16,
				"requireParentalConsent": false,
			},
			wantValid: true, // Parental consent is optional, but parentEmail is required when consent is true
		},
		{
			name: "minor without parent email",
			data: map[string]any{
				"notificationsEnabled":   false,
				"age":                    16,
				"requireParentalConsent": true,
				"parentEmail":            "",
			},
			wantValid: false,
		},
		{
			name: "adult doesn't need parental consent",
			data: map[string]any{
				"notificationsEnabled": false,
				"age":                  25,
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateUserPreferences(tt.data)
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
