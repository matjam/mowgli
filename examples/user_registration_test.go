package examples

import (
	"testing"
)

func TestValidateUserRegistration(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "valid registration",
			data: map[string]any{
				"username": "johndoe",
				"email":    "john@example.com",
				"password": "SecurePass123",
				"age":      25,
			},
			wantValid: true,
		},
		{
			name: "missing required field",
			data: map[string]any{
				"username": "johndoe",
				"email":    "john@example.com",
				// password missing
			},
			wantValid: false,
		},
		{
			name: "invalid email format",
			data: map[string]any{
				"username": "johndoe",
				"email":    "not-an-email",
				"password": "SecurePass123",
			},
			wantValid: false,
		},
		{
			name: "password too short",
			data: map[string]any{
				"username": "johndoe",
				"email":    "john@example.com",
				"password": "short",
			},
			wantValid: false,
		},
		// Note: Password complexity validation (uppercase, lowercase, numbers)
		// would require lookahead regex which Go doesn't support.
		// This is handled via length requirements only in this example.
		{
			name: "username with invalid characters",
			data: map[string]any{
				"username": "john-doe",
				"email":    "john@example.com",
				"password": "SecurePass123",
			},
			wantValid: false,
		},
		{
			name: "age too young",
			data: map[string]any{
				"username": "johndoe",
				"email":    "john@example.com",
				"password": "SecurePass123",
				"age":      10,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateUserRegistration(tt.data)
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
