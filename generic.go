package mowgli

import (
	"encoding/json"
)

// ValidateStruct validates data against a struct type and returns a typed result
// T is the struct type to validate against
// data is the raw data (map[string]any or compatible)
// specJSON is an optional JSON spec that will be merged on top of the struct-generated spec
func ValidateStruct[T any](data any, specJSON ...string) (*ValidationResult, T, error) {
	var zero T

	// Generate spec from struct type
	structSpec, err := SpecFromStruct(zero)
	if err != nil {
		return nil, zero, err
	}

	// Merge with JSON spec if provided
	if len(specJSON) > 0 && specJSON[0] != "" {
		jsonSpec, err := ParseSpecString(specJSON[0])
		if err != nil {
			return nil, zero, err
		}
		structSpec = MergeSpecs(structSpec, jsonSpec)
	}

	// Validate the data
	result := Validate(data, structSpec)
	if !result.Valid {
		return result, zero, nil
	}

	// Convert data to the struct type
	typedResult, err := convertToType(data, zero)
	if err != nil {
		return nil, zero, err
	}

	return result, typedResult, nil
}

// ValidateStructValue validates data against a struct value (empty instance) and returns a typed result
// structValue should be an empty struct instance like MyStruct{}
// data is the raw data (map[string]any or compatible)
// specJSON is an optional JSON spec that will be merged on top of the struct-generated spec
func ValidateStructValue[T any](structValue T, data any, specJSON ...string) (*ValidationResult, T, error) {
	return ValidateStruct[T](data, specJSON...)
}

// convertToType converts the data to the target type T
func convertToType[T any](data any, target T) (T, error) {
	var result T

	// If data is already of the correct type, return it directly
	if typedData, ok := data.(T); ok {
		return typedData, nil
	}

	// Use JSON marshaling/unmarshaling for conversion
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}

// ValidateAndConvert validates data and converts it to the target type
// This is a convenience function that combines validation and conversion
func ValidateAndConvert[T any](data any, spec *Spec) (*ValidationResult, T, error) {
	var zero T

	// Validate the data
	result := Validate(data, spec)
	if !result.Valid {
		return result, zero, nil
	}

	// Convert data to the struct type
	typedResult, err := convertToType(data, zero)
	if err != nil {
		return nil, zero, err
	}

	return result, typedResult, nil
}

// GetSpecFromStruct generates a Spec from a struct type
// This is a convenience function that can be used to get the spec for a struct type
func GetSpecFromStruct[T any](structValue T) (*Spec, error) {
	return SpecFromStruct(structValue)
}
