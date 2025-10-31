package examples

import (
	"fmt"

	"github.com/matjam/mowgli"
)

// Example: Advanced conditional validation with AND/OR operators and parentheses
// This example demonstrates complex conditional expressions using AND, OR,
// and parentheses to group expressions.

func ValidateAdvancedConditional(data map[string]any) (*mowgli.ValidationResult, error) {
	specJSON := `{
		"type": "object",
		"properties": {
			"hasEmail": {"type": "boolean"},
			"hasPhone": {"type": "boolean"},
			"email": {"type": "string"},
			"phone": {"type": "string"},
			"age": {"type": "integer"},
			"country": {"type": "string"},
			"requireVerification": {"type": "boolean"},
			"verificationCode": {"type": "string"}
		},
		"conditions": [
			{
				"if": "hasEmail == true OR hasPhone == true",
				"then": {
					"email": {
						"minLength": 1,
						"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
					},
					"phone": {
						"minLength": 10
					}
				}
			},
			{
				"if": "(age < 18 AND country == \"US\") OR requireVerification == true",
				"then": {
					"verificationCode": {
						"minLength": 1
					}
				}
			},
			{
				"if": "hasEmail == true AND hasPhone == true",
				"then": {
					"email": {
						"minLength": 1,
						"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
					},
					"phone": {
						"minLength": 10,
						"maxLength": 15
					}
				}
			}
		]
	}`

	spec, err := mowgli.ParseSpecString(specJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	result := mowgli.Validate(data, spec)
	return result, nil
}
