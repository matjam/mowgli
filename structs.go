package mowgli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// StructTagOptions holds parsed validation options from struct tags
type StructTagOptions struct {
	Required   bool
	Min        *float64
	Max        *float64
	MinLength  *int
	MaxLength  *int
	Pattern    *string
	Enum       []any
	AllowEmpty *bool
}

// ParseStructTag parses a mowgli struct tag and returns validation options
// Example: `mowgli:"required,min=0,max=100,minLength=1"`
func ParseStructTag(tag string) (*StructTagOptions, error) {
	options := &StructTagOptions{}

	if tag == "" {
		return options, nil
	}

	// Handle enum specially since values can contain commas
	// First check if there's an enum tag and extract it
	if enumIdx := strings.Index(tag, "enum="); enumIdx != -1 {
		// Extract everything before enum=
		before := strings.TrimSpace(tag[:enumIdx])
		// Extract enum value (everything after enum=)
		after := tag[enumIdx+5:] // 5 = len("enum=")

		// Parse everything before enum
		if before != "" {
			// Remove trailing comma if present
			before = strings.TrimSuffix(before, ",")
			if before != "" {
				beforeOpts, err := ParseStructTag(before)
				if err != nil {
					return nil, err
				}
				// Merge options
				options.Required = beforeOpts.Required
				options.Min = beforeOpts.Min
				options.Max = beforeOpts.Max
				options.MinLength = beforeOpts.MinLength
				options.MaxLength = beforeOpts.MaxLength
				options.Pattern = beforeOpts.Pattern
				options.AllowEmpty = beforeOpts.AllowEmpty
			}
		}

		// Parse enum value (may contain commas, so don't split)
		enumValues, err := parseEnumValues(after)
		if err != nil {
			return nil, fmt.Errorf("invalid enum value: %s: %w", after, err)
		}
		options.Enum = enumValues
		return options, nil
	}

	// No enum, proceed with normal comma-separated parsing
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if part == "required" {
			options.Required = true
			continue
		}

		if part == "allowEmpty" {
			trueVal := true
			options.AllowEmpty = &trueVal
			continue
		}

		// Parse key=value pairs
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag format: %s", part)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "min":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid min value: %s", value)
			}
			options.Min = &val
		case "max":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid max value: %s", value)
			}
			options.Max = &val
		case "minLength":
			val, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid minLength value: %s", value)
			}
			options.MinLength = &val
		case "maxLength":
			val, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid maxLength value: %s", value)
			}
			options.MaxLength = &val
		case "pattern":
			options.Pattern = &value
		default:
			return nil, fmt.Errorf("unknown tag option: %s", key)
		}
	}

	return options, nil
}

func parseEnumValues(value string) ([]any, error) {
	// Try parsing as JSON array first
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var result []any
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result, nil
		}
	}

	// Otherwise, parse as comma-separated values
	parts := strings.Split(value, ",")
	result := make([]any, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		// Try to parse as number, otherwise keep as string
		if num, err := strconv.ParseFloat(part, 64); err == nil {
			result[i] = num
		} else if part == "true" {
			result[i] = true
		} else if part == "false" {
			result[i] = false
		} else {
			// Remove quotes if present
			if len(part) >= 2 && ((part[0] == '"' && part[len(part)-1] == '"') || (part[0] == '\'' && part[len(part)-1] == '\'')) {
				part = part[1 : len(part)-1]
			}
			result[i] = part
		}
	}
	return result, nil
}

// SpecFromStruct generates a Spec from a struct type using reflection and struct tags
func SpecFromStruct(v any) (*Spec, error) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("SpecFromStruct expects a struct type, got %s", rt.Kind())
	}

	spec := &Spec{
		Type:       "object",
		Properties: make(map[string]*Spec),
		Required:   []string{},
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag for field name
		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		if jsonTag != "" && jsonTag != "-" {
			jsonParts := strings.Split(jsonTag, ",")
			if jsonParts[0] != "" {
				fieldName = jsonParts[0]
			}
		}

		// Skip if JSON tag is "-"
		if jsonTag == "-" {
			continue
		}

		// Get mowgli tag
		mowgliTag := field.Tag.Get("mowgli")
		options, err := ParseStructTag(mowgliTag)
		if err != nil {
			return nil, fmt.Errorf("error parsing struct tag for field %s: %w", fieldName, err)
		}

		// Determine field type
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		fieldSpec := &Spec{}

		// Map Go types to JSON types
		switch fieldType.Kind() {
		case reflect.String:
			fieldSpec.Type = "string"
			fieldSpec.MinLength = options.MinLength
			fieldSpec.MaxLength = options.MaxLength
			fieldSpec.Pattern = options.Pattern
			fieldSpec.AllowEmpty = options.AllowEmpty
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldSpec.Type = "integer"
			fieldSpec.Min = options.Min
			fieldSpec.Max = options.Max
		case reflect.Float32, reflect.Float64:
			fieldSpec.Type = "number"
			fieldSpec.Min = options.Min
			fieldSpec.Max = options.Max
		case reflect.Bool:
			fieldSpec.Type = "boolean"
		case reflect.Slice, reflect.Array:
			fieldSpec.Type = "array"
			fieldSpec.MinLength = options.MinLength
			fieldSpec.MaxLength = options.MaxLength
			// Handle array item types
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			switch elemType.Kind() {
			case reflect.String:
				fieldSpec.Items = &Spec{Type: "string"}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldSpec.Items = &Spec{Type: "integer"}
			case reflect.Float32, reflect.Float64:
				fieldSpec.Items = &Spec{Type: "number"}
			case reflect.Bool:
				fieldSpec.Items = &Spec{Type: "boolean"}
			case reflect.Struct:
				// Recursively generate spec for array item type
				nestedSpec, err := SpecFromStruct(reflect.New(elemType).Interface())
				if err != nil {
					return nil, fmt.Errorf("error generating spec for array item type %s: %w", fieldName, err)
				}
				fieldSpec.Items = nestedSpec
			}
		case reflect.Struct:
			// Recursively generate spec for nested struct
			nestedSpec, err := SpecFromStruct(reflect.New(fieldType).Interface())
			if err != nil {
				return nil, fmt.Errorf("error generating spec for nested struct %s: %w", fieldName, err)
			}
			fieldSpec = nestedSpec
		case reflect.Map:
			fieldSpec.Type = "object"
			// For maps, we treat them as generic objects
		default:
			return nil, fmt.Errorf("unsupported field type for %s: %s", fieldName, fieldType.Kind())
		}

		// Apply common options
		if options.Min != nil && fieldSpec.Min == nil {
			fieldSpec.Min = options.Min
		}
		if options.Max != nil && fieldSpec.Max == nil {
			fieldSpec.Max = options.Max
		}
		if options.Enum != nil {
			fieldSpec.Enum = options.Enum
		}

		spec.Properties[fieldName] = fieldSpec

		if options.Required {
			spec.Required = append(spec.Required, fieldName)
		}
	}

	return spec, nil
}

// MergeSpecs merges two specs, with override taking precedence
func MergeSpecs(base, override *Spec) *Spec {
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

	// Merge properties
	if override.Properties != nil {
		if merged.Properties == nil {
			merged.Properties = make(map[string]*Spec)
		}
		for k, v := range base.Properties {
			merged.Properties[k] = v
		}
		for k, v := range override.Properties {
			merged.Properties[k] = MergeSpecs(merged.Properties[k], v)
		}
	}

	// Apply overrides
	if override.Type != "" {
		merged.Type = override.Type
	}
	if override.Items != nil {
		merged.Items = override.Items
	}
	if override.Required != nil {
		merged.Required = override.Required
	}
	if override.Conditions != nil {
		merged.Conditions = override.Conditions
	}
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

	return merged
}
