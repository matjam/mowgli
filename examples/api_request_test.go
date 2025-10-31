package examples

import (
	"testing"
)

func TestValidateAPIRequest(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "valid create request",
			data: map[string]any{
				"action":   "create",
				"resource": "user",
				"data": map[string]any{
					"name": "John Doe",
				},
				"metadata": map[string]any{
					"userId": "user123",
				},
			},
			wantValid: true,
		},
		{
			name: "valid delete request without data",
			data: map[string]any{
				"action":   "delete",
				"resource": "product",
				"metadata": map[string]any{
					"userId": "user123",
				},
			},
			wantValid: true,
		},
		{
			name: "create with empty data object",
			data: map[string]any{
				"action":   "create",
				"resource": "user",
				"data":     map[string]any{},
				"metadata": map[string]any{
					"userId": "user123",
				},
			},
			wantValid: true, // Empty object is still valid - we can't enforce non-empty with current spec
		},
		{
			name: "invalid action",
			data: map[string]any{
				"action":   "invalid",
				"resource": "user",
			},
			wantValid: false,
		},
		{
			name: "strict validation with missing timestamp",
			data: map[string]any{
				"action":          "create",
				"resource":        "user",
				"validationLevel": "strict",
				"data": map[string]any{
					"name": "John Doe",
				},
				"metadata": map[string]any{
					"userId": "user123",
					// timestamp is required in strict mode but missing
				},
			},
			wantValid: false, // timestamp is required when validationLevel is strict
		},
		{
			name: "strict validation with timestamp",
			data: map[string]any{
				"action":          "create",
				"resource":        "user",
				"validationLevel": "strict",
				"data": map[string]any{
					"name": "John Doe",
				},
				"metadata": map[string]any{
					"userId":    "user123",
					"timestamp": 1234567890,
					"tags":      []any{"important", "urgent"},
				},
			},
			wantValid: true,
		},
		{
			name: "too many tags",
			data: map[string]any{
				"action":   "create",
				"resource": "user",
				"data": map[string]any{
					"name": "John Doe",
				},
				"metadata": map[string]any{
					"userId": "user123",
					"tags":   make([]any, 21),
				},
			},
			wantValid: false,
		},
		{
			name: "normal validation without timestamp",
			data: map[string]any{
				"action":          "create",
				"resource":        "user",
				"validationLevel": "normal",
				"data": map[string]any{
					"name": "John Doe",
				},
				"metadata": map[string]any{
					"userId": "user123",
				},
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAPIRequest(tt.data)
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
