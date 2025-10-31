package examples

import (
	"fmt"

	"github.com/matjam/mowgli"
)

// Example: Validating application configuration
// This example demonstrates validating a configuration file with
// various types, nested objects, and optional fields.

func ValidateConfig(data map[string]any) (*mowgli.ValidationResult, error) {
	specJSON := `{
		"type": "object",
		"properties": {
			"app": {
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"minLength": 1
					},
					"version": {
						"type": "string",
						"pattern": "^\\d+\\.\\d+\\.\\d+$"
					},
					"debug": {
						"type": "boolean"
					}
				},
				"required": ["name", "version"]
			},
			"server": {
				"type": "object",
				"properties": {
					"host": {
						"type": "string",
						"minLength": 1
					},
					"port": {
						"type": "integer",
						"min": 1,
						"max": 65535
					},
					"timeout": {
						"type": "integer",
						"min": 1,
						"max": 3600
					}
				},
				"required": ["host", "port"]
			},
			"database": {
				"type": "object",
				"properties": {
					"type": {
						"type": "string",
						"enum": ["postgres", "mysql", "sqlite"]
					},
					"maxConnections": {
						"type": "integer",
						"min": 1,
						"max": 100
					}
				},
				"required": ["type"]
			},
			"features": {
				"type": "array",
				"items": {
					"type": "string"
				},
				"maxLength": 50
			}
		},
		"required": ["app", "server"]
	}`

	spec, err := mowgli.ParseSpecString(specJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	result := mowgli.Validate(data, spec)
	return result, nil
}
