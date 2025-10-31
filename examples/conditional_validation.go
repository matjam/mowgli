package examples

import (
	"fmt"

	"github.com/matjam/mowgli"
)

// Example: Conditional validation for optional fields
// This example shows how to use conditional validation where
// certain fields become required or have different constraints
// based on other field values.

func ValidateUserPreferences(data map[string]any) (*mowgli.ValidationResult, error) {
	specJSON := `{
		"type": "object",
		"properties": {
			"notificationsEnabled": {
				"type": "boolean"
			},
			"email": {
				"type": "string"
			},
			"phone": {
				"type": "string"
			},
			"age": {
				"type": "integer",
				"min": 0,
				"max": 150
			},
			"requireParentalConsent": {
				"type": "boolean"
			},
			"parentEmail": {
				"type": "string"
			}
		},
		"conditions": [
			{
				"if": "notificationsEnabled == true",
				"then": {
					"email": {
						"minLength": 1,
						"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
					}
				},
				"else": {
					"email": {
						"allowEmpty": true
					}
				}
			},
			{
				"if": "age < 18",
				"then": {
					"requireParentalConsent": {
						"type": "boolean"
					},
					"parentEmail": {
						"minLength": 1,
						"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
					}
				}
			},
			{
				"if": "requireParentalConsent == true",
				"then": {
					"parentEmail": {
						"minLength": 1,
						"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
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
