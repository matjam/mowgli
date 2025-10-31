package examples

import (
	"fmt"

	"github.com/matjam/mowgli"
)

// Example: Validating user registration data
// This example shows how to validate a user registration form with
// required fields, email validation, and password requirements.

func ValidateUserRegistration(data map[string]any) (*mowgli.ValidationResult, error) {
	specJSON := `{
		"type": "object",
		"properties": {
			"username": {
				"type": "string",
				"minLength": 3,
				"maxLength": 30,
				"pattern": "^[a-zA-Z0-9_]+$"
			},
			"email": {
				"type": "string",
				"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
			},
			"password": {
				"type": "string",
				"minLength": 8,
				"maxLength": 128
			},
			"age": {
				"type": "integer",
				"min": 13,
				"max": 120
			}
		},
		"required": ["username", "email", "password"]
	}`

	spec, err := mowgli.ParseSpecString(specJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	result := mowgli.Validate(data, spec)
	return result, nil
}
