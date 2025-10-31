package examples

import (
	"fmt"

	"github.com/matjam/mowgli"
)

// Example: Validating API request payloads
// This example shows a more complex validation scenario with
// nested objects, arrays, and conditional validation based on
// request type.

func ValidateAPIRequest(data map[string]any) (*mowgli.ValidationResult, error) {
	specJSON := `{
		"type": "object",
		"properties": {
			"action": {
				"type": "string",
				"enum": ["create", "update", "delete"]
			},
			"resource": {
				"type": "string",
				"enum": ["user", "product", "order"]
			},
			"data": {
				"type": "object"
			},
			"metadata": {
				"type": "object",
				"properties": {
					"userId": {
						"type": "string",
						"minLength": 1
					},
					"timestamp": {
						"type": "integer",
						"min": 0
					},
					"tags": {
						"type": "array",
						"items": {
							"type": "string"
						},
						"maxLength": 20
					}
				},
				"required": ["userId"]
			},
			"validationLevel": {
				"type": "string",
				"enum": ["strict", "normal", "lenient"]
			}
		},
		"required": ["action", "resource"],
		"conditions": [
			{
				"if": "action == \"create\"",
				"then": {
					"data": {
						"type": "object",
						"properties": {}
					}
				}
			},
			{
				"if": "action == \"update\"",
				"then": {
					"data": {
						"type": "object",
						"properties": {}
					}
				}
			},
			{
				"if": "validationLevel == \"strict\"",
				"then": {
					"metadata": {
						"properties": {
							"timestamp": {
								"min": 1
							}
						},
						"required": ["userId", "timestamp"]
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
