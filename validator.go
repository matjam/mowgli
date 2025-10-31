package mowgli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// ValidationError represents a validation error with a path to the field
type ValidationError struct {
	Path    string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationResult contains the result of validation
type ValidationResult struct {
	Valid  bool
	Errors []*ValidationError
}

// Validate validates a JSON value against a spec
func Validate(data any, spec *Spec) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []*ValidationError{},
	}

	if spec == nil {
		result.addError("", "spec is nil")
		return result
	}

	result.validate("", data, spec)
	return result
}

// ValidateJSON validates a JSON byte slice against a spec
func ValidateJSON(jsonData []byte, spec *Spec) (*ValidationResult, error) {
	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return Validate(data, spec), nil
}

// ValidateJSONString validates a JSON string against a spec
func ValidateJSONString(jsonStr string, spec *Spec) (*ValidationResult, error) {
	return ValidateJSON([]byte(jsonStr), spec)
}

func (r *ValidationResult) addError(path, message string) {
	r.Valid = false
	r.Errors = append(r.Errors, &ValidationError{
		Path:    path,
		Message: message,
	})
}

func buildPath(base, field string) string {
	if base == "" {
		return field
	}
	if field == "" {
		return base
	}
	return base + "." + field
}

func buildArrayPath(base string, index int) string {
	return base + "[" + strconv.Itoa(index) + "]"
}

func (r *ValidationResult) validate(path string, value any, spec *Spec) {
	if spec == nil {
		return
	}

	// Handle null values
	if value == nil {
		if spec.Type != "null" {
			r.addError(path, fmt.Sprintf("expected type %s, got null", spec.Type))
		}
		return
	}

	switch spec.Type {
	case "string":
		r.validateString(path, value, spec)
	case "number":
		r.validateNumber(path, value, spec)
	case "integer":
		r.validateInteger(path, value, spec)
	case "boolean":
		r.validateBoolean(path, value, spec)
	case "object":
		r.validateObject(path, value, spec)
	case "array":
		r.validateArray(path, value, spec)
	case "null":
		// value is guaranteed to be non-nil at this point (checked above)
		r.addError(path, "expected null, got non-null value")
	default:
		r.addError(path, fmt.Sprintf("unknown type: %s", spec.Type))
	}

	// Validate enum constraint if specified
	if len(spec.Enum) > 0 {
		r.validateEnum(path, value, spec.Enum)
	}
}

func (r *ValidationResult) validateString(path string, value any, spec *Spec) {
	str, ok := value.(string)
	if !ok {
		r.addError(path, fmt.Sprintf("expected string, got %T", value))
		return
	}

	// Handle allowEmpty - if true and string is empty, skip other validations
	if spec.AllowEmpty != nil && *spec.AllowEmpty && str == "" {
		return
	}

	// If allowEmpty is false or not set, empty strings must pass minLength check
	if spec.MinLength != nil && len(str) < *spec.MinLength {
		r.addError(path, fmt.Sprintf("string length %d is less than minimum %d", len(str), *spec.MinLength))
	}

	if spec.MaxLength != nil && len(str) > *spec.MaxLength {
		r.addError(path, fmt.Sprintf("string length %d is greater than maximum %d", len(str), *spec.MaxLength))
	}

	if spec.Pattern != nil {
		matched, err := regexp.MatchString(*spec.Pattern, str)
		if err != nil {
			r.addError(path, fmt.Sprintf("invalid pattern: %v", err))
		} else if !matched {
			r.addError(path, fmt.Sprintf("string does not match pattern: %s", *spec.Pattern))
		}
	}
}

func (r *ValidationResult) validateNumber(path string, value any, spec *Spec) {
	var num float64
	switch v := value.(type) {
	case float64:
		num = v
	case int:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	default:
		r.addError(path, fmt.Sprintf("expected number, got %T", value))
		return
	}

	if spec.Min != nil && num < *spec.Min {
		r.addError(path, fmt.Sprintf("number %g is less than minimum %g", num, *spec.Min))
	}

	if spec.Max != nil && num > *spec.Max {
		r.addError(path, fmt.Sprintf("number %g is greater than maximum %g", num, *spec.Max))
	}
}

func (r *ValidationResult) validateInteger(path string, value any, spec *Spec) {
	var num float64
	isInt := false

	switch v := value.(type) {
	case float64:
		num = v
		// Check if it's actually an integer (no fractional part)
		if num == float64(int64(num)) {
			isInt = true
		}
	case int:
		num = float64(v)
		isInt = true
	case int64:
		num = float64(v)
		isInt = true
	default:
		r.addError(path, fmt.Sprintf("expected integer, got %T", value))
		return
	}

	if !isInt {
		r.addError(path, fmt.Sprintf("expected integer, got float: %g", num))
		return
	}

	if spec.Min != nil && num < *spec.Min {
		r.addError(path, fmt.Sprintf("integer %g is less than minimum %g", num, *spec.Min))
	}

	if spec.Max != nil && num > *spec.Max {
		r.addError(path, fmt.Sprintf("integer %g is greater than maximum %g", num, *spec.Max))
	}
}

func (r *ValidationResult) validateBoolean(path string, value any, spec *Spec) {
	_, ok := value.(bool)
	if !ok {
		r.addError(path, fmt.Sprintf("expected boolean, got %T", value))
	}
}

func (r *ValidationResult) validateObject(path string, value any, spec *Spec) {
	obj, ok := value.(map[string]any)
	if !ok {
		r.addError(path, fmt.Sprintf("expected object, got %T", value))
		return
	}

	// Validate properties with conditional overrides
	// We need to do this first to get the effective specs for required field checking
	effectiveSpecs := make(map[string]*Spec)
	if spec.Properties != nil {
		effectiveSpecs = r.buildEffectiveSpecs(obj, spec)
	}

	// Check required fields (use base spec required fields)
	// Required fields from nested object overrides are handled when validating those nested objects
	if spec.Required != nil {
		for _, req := range spec.Required {
			if _, exists := obj[req]; !exists {
				r.addError(buildPath(path, req), "required field is missing")
			}
		}
	}

	// Validate properties with conditional overrides
	if spec.Properties != nil {
		for key, propSpec := range spec.Properties {
			propValue, exists := obj[key]

			// Get effective spec (with condition overrides applied)
			effectiveSpec := propSpec
			if override, ok := effectiveSpecs[key]; ok {
				effectiveSpec = override
			}

			if exists {
				r.validate(buildPath(path, key), propValue, effectiveSpec)
			}
		}

		// Check for extra properties (warn but don't fail by default)
		// This could be configurable in the future
		for key := range obj {
			if spec.Properties != nil {
				if _, exists := spec.Properties[key]; !exists {
					// Extra property - silently allowed for now
					_ = key // Placeholder for future validation
				}
			}
		}
	}
}

// buildEffectiveSpecs evaluates conditions and returns effective specs for each property
func (r *ValidationResult) buildEffectiveSpecs(obj map[string]any, spec *Spec) map[string]*Spec {
	effectiveSpecs := make(map[string]*Spec)

	if len(spec.Conditions) == 0 {
		return effectiveSpecs
	}

	// Collect all overrides first, then merge them all together
	for _, condition := range spec.Conditions {
		result, err := evalExpression(condition.If, obj)
		if err != nil {
			r.addError("", fmt.Sprintf("error evaluating condition '%s': %v", condition.If, err))
			continue
		}

		var overrides map[string]*Spec
		if result {
			overrides = condition.Then
		} else {
			overrides = condition.Else
		}

		if overrides != nil {
			for fieldName, overrideSpec := range overrides {
				// Get or create the effective spec for this field
				if existing, exists := effectiveSpecs[fieldName]; exists {
					// Merge with existing override
					effectiveSpecs[fieldName] = r.mergeSpecs(existing, overrideSpec)
				} else {
					// Start from base spec or create new
					if baseSpec, exists := spec.Properties[fieldName]; exists {
						merged := r.mergeSpecs(baseSpec, overrideSpec)
						effectiveSpecs[fieldName] = merged
					} else {
						// Condition defines a new validation for a field not in properties
						effectiveSpecs[fieldName] = overrideSpec
					}
				}
			}
		}
	}

	return effectiveSpecs
}

// mergeSpecs merges overrideSpec into baseSpec, with overrideSpec taking precedence
func (r *ValidationResult) mergeSpecs(base, override *Spec) *Spec {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := &Spec{
		Type:       base.Type,
		Properties: base.Properties,
		Items:      base.Items,
		Required:   base.Required,
		Conditions: base.Conditions,
		Min:        base.Min,
		Max:        base.Max,
		MinLength:  base.MinLength,
		MaxLength:  base.MaxLength,
		Pattern:    base.Pattern,
		Enum:       base.Enum,
		AllowEmpty: base.AllowEmpty,
	}

	// Apply overrides
	if override.Min != nil {
		merged.Min = override.Min
	}
	if override.Max != nil {
		merged.Max = override.Max
	}
	if override.MinLength != nil {
		merged.MinLength = override.MinLength
	}
	if override.MaxLength != nil {
		merged.MaxLength = override.MaxLength
	}
	if override.Pattern != nil {
		merged.Pattern = override.Pattern
	}
	if override.Enum != nil {
		merged.Enum = override.Enum
	}
	if override.AllowEmpty != nil {
		merged.AllowEmpty = override.AllowEmpty
	}
	if override.Type != "" {
		merged.Type = override.Type
	}
	// Merge properties
	if override.Properties != nil {
		if merged.Properties == nil {
			merged.Properties = make(map[string]*Spec)
		}
		for k, v := range base.Properties {
			merged.Properties[k] = v
		}
		for k, v := range override.Properties {
			if existing, exists := merged.Properties[k]; exists {
				merged.Properties[k] = r.mergeSpecs(existing, v)
			} else {
				merged.Properties[k] = v
			}
		}
	}
	// Merge required fields (override replaces, as that's more intuitive for conditions)
	if override.Required != nil {
		merged.Required = override.Required
	}

	return merged
}

func (r *ValidationResult) validateArray(path string, value any, spec *Spec) {
	arr, ok := value.([]any)
	if !ok {
		// Try to handle arrays of other types
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			r.addError(path, fmt.Sprintf("expected array, got %T", value))
			return
		}

		// Convert to []any
		arr = make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			arr[i] = val.Index(i).Interface()
		}
	}

	if spec.MinLength != nil && len(arr) < *spec.MinLength {
		r.addError(path, fmt.Sprintf("array length %d is less than minimum %d", len(arr), *spec.MinLength))
	}

	if spec.MaxLength != nil && len(arr) > *spec.MaxLength {
		r.addError(path, fmt.Sprintf("array length %d is greater than maximum %d", len(arr), *spec.MaxLength))
	}

	if spec.Items != nil {
		for i, item := range arr {
			r.validate(buildArrayPath(path, i), item, spec.Items)
		}
	}
}

func (r *ValidationResult) validateEnum(path string, value any, enum []any) {
	for _, allowed := range enum {
		if reflect.DeepEqual(value, allowed) {
			return
		}
	}

	// Format enum values for error message
	enumStrs := make([]string, len(enum))
	for i, v := range enum {
		enumStrs[i] = fmt.Sprintf("%v", v)
	}
	r.addError(path, fmt.Sprintf("value not in enum: %v (allowed: %v)", value, enumStrs))
}
