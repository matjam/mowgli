package examples

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "valid minimal config",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 8080,
				},
			},
			wantValid: true,
		},
		{
			name: "valid full config",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "2.1.3",
					"debug":   true,
				},
				"server": map[string]any{
					"host":    "0.0.0.0",
					"port":    443,
					"timeout": 30,
				},
				"database": map[string]any{
					"type":           "postgres",
					"maxConnections": 20,
				},
				"features": []any{"auth", "logging", "metrics"},
			},
			wantValid: true,
		},
		{
			name: "missing required app name",
			data: map[string]any{
				"app": map[string]any{
					"version": "1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 8080,
				},
			},
			wantValid: false,
		},
		{
			name: "invalid version format",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "v1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 8080,
				},
			},
			wantValid: false,
		},
		{
			name: "port out of range",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 70000,
				},
			},
			wantValid: false,
		},
		{
			name: "invalid database type",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 8080,
				},
				"database": map[string]any{
					"type": "mongodb",
				},
			},
			wantValid: false,
		},
		{
			name: "too many features",
			data: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": "1.0.0",
				},
				"server": map[string]any{
					"host": "localhost",
					"port": 8080,
				},
				"features": make([]any, 51),
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateConfig(tt.data)
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
