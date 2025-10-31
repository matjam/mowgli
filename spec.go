package mowgli

import "encoding/json"

// Condition defines a conditional validation rule
type Condition struct {
	If   string           `json:"if"`             // Expression to evaluate, e.g., "enabled == true", "count > 0"
	Then map[string]*Spec `json:"then"`           // Spec overrides to apply when condition is true
	Else map[string]*Spec `json:"else,omitempty"` // Spec overrides to apply when condition is false
}

// Spec defines the validation specification structure
type Spec struct {
	Type       string           `json:"type"`                 // string, number, integer, boolean, object, array, null
	Properties map[string]*Spec `json:"properties,omitempty"` // For object type
	Items      *Spec            `json:"items,omitempty"`      // For array type
	Required   []string         `json:"required,omitempty"`   // For object type - list of required property names
	Conditions []Condition      `json:"conditions,omitempty"` // Conditional validation rules for object type

	// Constraints
	Min        *float64      `json:"min,omitempty"`        // For number/integer - minimum value
	Max        *float64      `json:"max,omitempty"`        // For number/integer - maximum value
	MinLength  *int          `json:"minLength,omitempty"`  // For string/array - minimum length
	MaxLength  *int          `json:"maxLength,omitempty"`  // For string/array - maximum length
	Pattern    *string       `json:"pattern,omitempty"`    // For string - regex pattern (future: could support regex validation)
	Enum       []any `json:"enum,omitempty"`       // Array of allowed values
	AllowEmpty *bool         `json:"allowEmpty,omitempty"` // For strings - allows empty string if true
}

// ParseSpec parses a JSON byte slice into a Spec
func ParseSpec(data []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

// ParseSpecString parses a JSON string into a Spec
func ParseSpecString(s string) (*Spec, error) {
	return ParseSpec([]byte(s))
}
