package mowgli

import (
	"testing"
)

func TestLoadSpec(t *testing.T) {
	spec, err := LoadSpec("user_registration.json")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	if spec.Type != "object" {
		t.Errorf("expected type object, got %s", spec.Type)
	}

	if len(spec.Required) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(spec.Required))
	}
}

func TestLoadTestCases(t *testing.T) {
	testCases, err := LoadTestCases("user_registration.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	if len(testCases) == 0 {
		t.Error("expected at least one test case, got none")
	}

	// Verify first test case structure
	if testCases[0].Name == "" {
		t.Error("test case name should not be empty")
	}
}

func TestSharedDataValidation(t *testing.T) {
	// Load spec and test cases
	spec, err := LoadSpec("user_registration.json")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	testCases, err := LoadTestCases("user_registration.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := Validate(tc.Data, spec)
			if result.Valid != tc.ExpectedValid {
				if !tc.ExpectedValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestSharedConditionalValidation(t *testing.T) {
	spec, err := LoadSpec("conditional_validation.json")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	testCases, err := LoadTestCases("conditional_validation.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := Validate(tc.Data, spec)
			if result.Valid != tc.ExpectedValid {
				if !tc.ExpectedValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestSharedNestedObjects(t *testing.T) {
	spec, err := LoadSpec("nested_objects.json")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	testCases, err := LoadTestCases("nested_objects.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := Validate(tc.Data, spec)
			if result.Valid != tc.ExpectedValid {
				if !tc.ExpectedValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestSharedAdvancedConditional(t *testing.T) {
	spec, err := LoadSpec("advanced_conditional.json")
	if err != nil {
		t.Fatalf("failed to load spec: %v", err)
	}

	testCases, err := LoadTestCases("advanced_conditional.json")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := Validate(tc.Data, spec)
			if result.Valid != tc.ExpectedValid {
				if !tc.ExpectedValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

