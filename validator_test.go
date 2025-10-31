package mowgli

import (
	"encoding/json"
	"testing"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		value     any
		shouldErr bool
	}{
		{
			name:      "valid string",
			specJSON:  `{"type": "string"}`,
			value:     "hello",
			shouldErr: false,
		},
		{
			name:      "string with minLength valid",
			specJSON:  `{"type": "string", "minLength": 3}`,
			value:     "hello",
			shouldErr: false,
		},
		{
			name:      "string with minLength invalid",
			specJSON:  `{"type": "string", "minLength": 10}`,
			value:     "hello",
			shouldErr: true,
		},
		{
			name:      "string with maxLength valid",
			specJSON:  `{"type": "string", "maxLength": 10}`,
			value:     "hello",
			shouldErr: false,
		},
		{
			name:      "string with maxLength invalid",
			specJSON:  `{"type": "string", "maxLength": 3}`,
			value:     "hello",
			shouldErr: true,
		},
		{
			name:      "string with pattern valid",
			specJSON:  `{"type": "string", "pattern": "^[a-z]+$"}`,
			value:     "hello",
			shouldErr: false,
		},
		{
			name:      "string with pattern invalid",
			specJSON:  `{"type": "string", "pattern": "^[a-z]+$"}`,
			value:     "Hello123",
			shouldErr: true,
		},
		{
			name:      "wrong type",
			specJSON:  `{"type": "string"}`,
			value:     123,
			shouldErr: true,
		},
		{
			name:      "string with enum valid",
			specJSON:  `{"type": "string", "enum": ["red", "green", "blue"]}`,
			value:     "red",
			shouldErr: false,
		},
		{
			name:      "string with enum invalid",
			specJSON:  `{"type": "string", "enum": ["red", "green", "blue"]}`,
			value:     "yellow",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			result := Validate(tt.value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateNumber(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		value     any
		shouldErr bool
	}{
		{
			name:      "valid number",
			specJSON:  `{"type": "number"}`,
			value:     3.14,
			shouldErr: false,
		},
		{
			name:      "number with min valid",
			specJSON:  `{"type": "number", "min": 0}`,
			value:     5.5,
			shouldErr: false,
		},
		{
			name:      "number with min invalid",
			specJSON:  `{"type": "number", "min": 10}`,
			value:     5.5,
			shouldErr: true,
		},
		{
			name:      "number with max valid",
			specJSON:  `{"type": "number", "max": 100}`,
			value:     50.0,
			shouldErr: false,
		},
		{
			name:      "number with max invalid",
			specJSON:  `{"type": "number", "max": 10}`,
			value:     50.0,
			shouldErr: true,
		},
		{
			name:      "number with range valid",
			specJSON:  `{"type": "number", "min": 0, "max": 100}`,
			value:     50.0,
			shouldErr: false,
		},
		{
			name:      "number with range invalid (too low)",
			specJSON:  `{"type": "number", "min": 10, "max": 100}`,
			value:     5.0,
			shouldErr: true,
		},
		{
			name:      "number with range invalid (too high)",
			specJSON:  `{"type": "number", "min": 10, "max": 100}`,
			value:     150.0,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			result := Validate(tt.value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		value     any
		shouldErr bool
	}{
		{
			name:      "valid integer",
			specJSON:  `{"type": "integer"}`,
			value:     42,
			shouldErr: false,
		},
		{
			name:      "integer as float valid",
			specJSON:  `{"type": "integer"}`,
			value:     42.0,
			shouldErr: false,
		},
		{
			name:      "float not valid as integer",
			specJSON:  `{"type": "integer"}`,
			value:     42.5,
			shouldErr: true,
		},
		{
			name:      "integer with min/max valid",
			specJSON:  `{"type": "integer", "min": 0, "max": 100}`,
			value:     50,
			shouldErr: false,
		},
		{
			name:      "integer with min/max invalid",
			specJSON:  `{"type": "integer", "min": 0, "max": 100}`,
			value:     150,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			result := Validate(tt.value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		value     any
		shouldErr bool
	}{
		{
			name:      "valid boolean true",
			specJSON:  `{"type": "boolean"}`,
			value:     true,
			shouldErr: false,
		},
		{
			name:      "valid boolean false",
			specJSON:  `{"type": "boolean"}`,
			value:     false,
			shouldErr: false,
		},
		{
			name:      "invalid boolean",
			specJSON:  `{"type": "boolean"}`,
			value:     "true",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			result := Validate(tt.value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateObject(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		valueJSON string
		shouldErr bool
	}{
		{
			name:      "simple object valid",
			specJSON:  `{"type": "object", "properties": {"name": {"type": "string"}}}`,
			valueJSON: `{"name": "John"}`,
			shouldErr: false,
		},
		{
			name:      "object with required field present",
			specJSON:  `{"type": "object", "properties": {"name": {"type": "string"}}, "required": ["name"]}`,
			valueJSON: `{"name": "John"}`,
			shouldErr: false,
		},
		{
			name:      "object with required field missing",
			specJSON:  `{"type": "object", "properties": {"name": {"type": "string"}}, "required": ["name"]}`,
			valueJSON: `{}`,
			shouldErr: true,
		},
		{
			name:      "nested object valid",
			specJSON:  `{"type": "object", "properties": {"user": {"type": "object", "properties": {"name": {"type": "string"}}}}}`,
			valueJSON: `{"user": {"name": "John"}}`,
			shouldErr: false,
		},
		{
			name:      "object with multiple properties valid",
			specJSON:  `{"type": "object", "properties": {"name": {"type": "string"}, "age": {"type": "integer"}}, "required": ["name"]}`,
			valueJSON: `{"name": "John", "age": 30}`,
			shouldErr: false,
		},
		{
			name:      "object property validation fails",
			specJSON:  `{"type": "object", "properties": {"age": {"type": "integer", "min": 0, "max": 150}}}`,
			valueJSON: `{"age": 200}`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			var value any
			if err := json.Unmarshal([]byte(tt.valueJSON), &value); err != nil {
				t.Fatalf("Failed to parse value JSON: %v", err)
			}

			result := Validate(value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateArray(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		valueJSON string
		shouldErr bool
	}{
		{
			name:      "simple array valid",
			specJSON:  `{"type": "array"}`,
			valueJSON: `[1, 2, 3]`,
			shouldErr: false,
		},
		{
			name:      "array with minLength valid",
			specJSON:  `{"type": "array", "minLength": 2}`,
			valueJSON: `[1, 2, 3]`,
			shouldErr: false,
		},
		{
			name:      "array with minLength invalid",
			specJSON:  `{"type": "array", "minLength": 5}`,
			valueJSON: `[1, 2, 3]`,
			shouldErr: true,
		},
		{
			name:      "array with maxLength valid",
			specJSON:  `{"type": "array", "maxLength": 5}`,
			valueJSON: `[1, 2, 3]`,
			shouldErr: false,
		},
		{
			name:      "array with maxLength invalid",
			specJSON:  `{"type": "array", "maxLength": 2}`,
			valueJSON: `[1, 2, 3]`,
			shouldErr: true,
		},
		{
			name:      "array with item spec valid",
			specJSON:  `{"type": "array", "items": {"type": "string"}}`,
			valueJSON: `["a", "b", "c"]`,
			shouldErr: false,
		},
		{
			name:      "array with item spec invalid",
			specJSON:  `{"type": "array", "items": {"type": "string"}}`,
			valueJSON: `["a", 123, "c"]`,
			shouldErr: true,
		},
		{
			name:      "array with item constraints valid",
			specJSON:  `{"type": "array", "items": {"type": "integer", "min": 0, "max": 100}}`,
			valueJSON: `[10, 20, 30]`,
			shouldErr: false,
		},
		{
			name:      "array with item constraints invalid",
			specJSON:  `{"type": "array", "items": {"type": "integer", "min": 0, "max": 100}}`,
			valueJSON: `[10, 200, 30]`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			var value any
			if err := json.Unmarshal([]byte(tt.valueJSON), &value); err != nil {
				t.Fatalf("Failed to parse value JSON: %v", err)
			}

			result := Validate(value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateJSONString(t *testing.T) {
	specJSON := `{
		"type": "object",
		"properties": {
			"name": {
				"type": "string",
				"minLength": 1,
				"maxLength": 100
			},
			"age": {
				"type": "integer",
				"min": 0,
				"max": 150
			},
			"email": {
				"type": "string",
				"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
			}
		},
		"required": ["name", "age"]
	}`

	spec, err := ParseSpecString(specJSON)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	tests := []struct {
		name      string
		valueJSON string
		shouldErr bool
	}{
		{
			name:      "valid user object",
			valueJSON: `{"name": "John Doe", "age": 30, "email": "john@example.com"}`,
			shouldErr: false,
		},
		{
			name:      "missing required field",
			valueJSON: `{"name": "John Doe", "email": "john@example.com"}`,
			shouldErr: true,
		},
		{
			name:      "invalid age range",
			valueJSON: `{"name": "John Doe", "age": 200, "email": "john@example.com"}`,
			shouldErr: true,
		},
		{
			name:      "invalid email pattern",
			valueJSON: `{"name": "John Doe", "age": 30, "email": "not-an-email"}`,
			shouldErr: true,
		},
		{
			name:      "name too short",
			valueJSON: `{"name": "", "age": 30}`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateJSONString(tt.valueJSON, spec)
			if err != nil {
				t.Fatalf("Failed to validate JSON: %v", err)
			}

			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateComplexNested(t *testing.T) {
	specJSON := `{
		"type": "object",
		"properties": {
			"users": {
				"type": "array",
				"minLength": 1,
				"items": {
					"type": "object",
					"properties": {
						"name": {"type": "string", "minLength": 1},
						"age": {"type": "integer", "min": 0, "max": 150},
						"tags": {
							"type": "array",
							"items": {"type": "string"}
						}
					},
					"required": ["name"]
				}
			}
		},
		"required": ["users"]
	}`

	spec, err := ParseSpecString(specJSON)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	validJSON := `{
		"users": [
			{"name": "Alice", "age": 25, "tags": ["developer", "go"]},
			{"name": "Bob", "age": 30, "tags": ["designer"]}
		]
	}`

	result, err := ValidateJSONString(validJSON, spec)
	if err != nil {
		t.Fatalf("Failed to validate JSON: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
	}

	invalidJSON := `{
		"users": [
			{"name": "", "age": 25},
			{"age": 30}
		]
	}`

	result2, err := ValidateJSONString(invalidJSON, spec)
	if err != nil {
		t.Fatalf("Failed to validate JSON: %v", err)
	}

	if result2.Valid {
		t.Errorf("Expected validation to fail, but it passed")
	}
}

func TestValidateConditional(t *testing.T) {
	tests := []struct {
		name      string
		specJSON  string
		valueJSON string
		shouldErr bool
	}{
		{
			name: "conditional validation - enabled true, value present",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string"}
				},
				"conditions": [
					{
						"if": "enabled == true",
						"then": {
							"value": {"minLength": 1}
						}
					},
					{
						"if": "enabled == false",
						"then": {
							"value": {"allowEmpty": true}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": true, "value": "hello"}`,
			shouldErr: false,
		},
		{
			name: "conditional validation - enabled true, value missing",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string"}
				},
				"conditions": [
					{
						"if": "enabled == true",
						"then": {
							"value": {"minLength": 1}
						}
					},
					{
						"if": "enabled == false",
						"then": {
							"value": {"allowEmpty": true}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": true, "value": ""}`,
			shouldErr: true,
		},
		{
			name: "conditional validation - enabled false, value empty allowed",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string"}
				},
				"conditions": [
					{
						"if": "enabled == true",
						"then": {
							"value": {"minLength": 1}
						}
					},
					{
						"if": "enabled == false",
						"then": {
							"value": {"allowEmpty": true}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": false, "value": ""}`,
			shouldErr: false,
		},
		{
			name: "conditional validation - enabled false, value non-empty not allowed",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string", "maxLength": 0}
				},
				"conditions": [
					{
						"if": "enabled == false",
						"then": {
							"value": {"maxLength": 0}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": false, "value": "not allowed"}`,
			shouldErr: true,
		},
		{
			name: "conditional validation - numeric comparison",
			specJSON: `{
				"type": "object",
				"properties": {
					"count": {"type": "integer"},
					"message": {"type": "string"}
				},
				"conditions": [
					{
						"if": "count > 0",
						"then": {
							"message": {"minLength": 1}
						}
					},
					{
						"if": "count <= 0",
						"then": {
							"message": {"allowEmpty": true}
						}
					}
				]
			}`,
			valueJSON: `{"count": 5, "message": "hello"}`,
			shouldErr: false,
		},
		{
			name: "conditional validation - numeric comparison fails",
			specJSON: `{
				"type": "object",
				"properties": {
					"count": {"type": "integer"},
					"message": {"type": "string"}
				},
				"conditions": [
					{
						"if": "count > 0",
						"then": {
							"message": {"minLength": 1}
						}
					}
				]
			}`,
			valueJSON: `{"count": 5, "message": ""}`,
			shouldErr: true,
		},
		{
			name: "conditional validation - string comparison",
			specJSON: `{
				"type": "object",
				"properties": {
					"status": {"type": "string"},
					"details": {"type": "string"}
				},
				"conditions": [
					{
						"if": "status == \"active\"",
						"then": {
							"details": {"minLength": 1}
						}
					}
				]
			}`,
			valueJSON: `{"status": "active", "details": "some details"}`,
			shouldErr: false,
		},
		{
			name: "conditional validation - boolean field shorthand",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string"}
				},
				"conditions": [
					{
						"if": "enabled",
						"then": {
							"value": {"minLength": 1}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": true, "value": "test"}`,
			shouldErr: false,
		},
		{
			name: "conditional validation - boolean field shorthand false",
			specJSON: `{
				"type": "object",
				"properties": {
					"enabled": {"type": "boolean"},
					"value": {"type": "string"}
				},
				"conditions": [
					{
						"if": "enabled",
						"then": {
							"value": {"minLength": 1}
						}
					},
					{
						"if": "enabled == false",
						"then": {
							"value": {"allowEmpty": true}
						}
					}
				]
			}`,
			valueJSON: `{"enabled": false, "value": ""}`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := ParseSpecString(tt.specJSON)
			if err != nil {
				t.Fatalf("Failed to parse spec: %v", err)
			}

			var value any
			if err := json.Unmarshal([]byte(tt.valueJSON), &value); err != nil {
				t.Fatalf("Failed to parse value JSON: %v", err)
			}

			result := Validate(value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}

func TestValidateConditionalComplex(t *testing.T) {
	specJSON := `{
		"type": "object",
		"properties": {
			"hasEmail": {"type": "boolean"},
			"email": {"type": "string"},
			"hasPhone": {"type": "boolean"},
			"phone": {"type": "string"},
			"age": {"type": "integer"}
		},
		"conditions": [
			{
				"if": "hasEmail == true",
				"then": {
					"email": {"minLength": 1, "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"}
				}
			},
			{
				"if": "hasEmail == false",
				"then": {
					"email": {"allowEmpty": true}
				}
			},
			{
				"if": "hasPhone == true",
				"then": {
					"phone": {"minLength": 10}
				}
			},
			{
				"if": "age >= 18",
				"then": {
					"email": {"minLength": 1}
				}
			}
		]
	}`

	spec, err := ParseSpecString(specJSON)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	tests := []struct {
		name      string
		valueJSON string
		shouldErr bool
	}{
		{
			name:      "valid - hasEmail true with valid email",
			valueJSON: `{"hasEmail": true, "email": "test@example.com", "hasPhone": false, "age": 25}`,
			shouldErr: false,
		},
		{
			name:      "invalid - hasEmail true but email empty",
			valueJSON: `{"hasEmail": true, "email": "", "hasPhone": false, "age": 25}`,
			shouldErr: true,
		},
		{
			name:      "valid - hasEmail false with empty email",
			valueJSON: `{"hasEmail": false, "email": "", "hasPhone": false, "age": 25}`,
			shouldErr: false,
		},
		{
			name:      "invalid - hasEmail true but invalid email format",
			valueJSON: `{"hasEmail": true, "email": "not-an-email", "hasPhone": false, "age": 25}`,
			shouldErr: true,
		},
		{
			name:      "valid - hasPhone true with valid phone",
			valueJSON: `{"hasEmail": false, "hasPhone": true, "phone": "1234567890", "age": 25}`,
			shouldErr: false,
		},
		{
			name:      "invalid - hasPhone true but phone too short",
			valueJSON: `{"hasEmail": false, "hasPhone": true, "phone": "123", "age": 25}`,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value any
			if err := json.Unmarshal([]byte(tt.valueJSON), &value); err != nil {
				t.Fatalf("Failed to parse value JSON: %v", err)
			}

			result := Validate(value, spec)
			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("Expected validation to fail, but it passed")
				} else {
					t.Errorf("Expected validation to pass, but it failed: %v", result.Errors)
				}
			}
		})
	}
}
