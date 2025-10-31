package mowgli

import (
	"testing"
)

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
		check   func(*StructTagOptions) bool
	}{
		{
			name: "required",
			tag:  "required",
			check: func(opts *StructTagOptions) bool {
				return opts.Required == true
			},
		},
		{
			name: "min max",
			tag:  "min=0,max=100",
			check: func(opts *StructTagOptions) bool {
				return opts.Min != nil && *opts.Min == 0 && opts.Max != nil && *opts.Max == 100
			},
		},
		{
			name: "minLength maxLength",
			tag:  "minLength=1,maxLength=100",
			check: func(opts *StructTagOptions) bool {
				return opts.MinLength != nil && *opts.MinLength == 1 && opts.MaxLength != nil && *opts.MaxLength == 100
			},
		},
		{
			name: "pattern",
			tag:  "pattern=^[a-z]+$",
			check: func(opts *StructTagOptions) bool {
				return opts.Pattern != nil && *opts.Pattern == "^[a-z]+$"
			},
		},
		{
			name: "enum string",
			tag:  "enum=red,green,blue",
			check: func(opts *StructTagOptions) bool {
				return opts.Enum != nil && len(opts.Enum) == 3
			},
		},
		{
			name: "allowEmpty",
			tag:  "allowEmpty",
			check: func(opts *StructTagOptions) bool {
				return opts.AllowEmpty != nil && *opts.AllowEmpty == true
			},
		},
		{
			name:    "invalid format",
			tag:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseStructTag(tt.tag)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tt.check(opts) {
				t.Errorf("options check failed for tag: %s", tt.tag)
			}
		})
	}
}

func TestSpecFromStruct(t *testing.T) {
	type SimpleStruct struct {
		Name string `json:"name" mowgli:"required,minLength=1,maxLength=100"`
		Age  int    `json:"age" mowgli:"required,min=0,max=150"`
	}

	spec, err := SpecFromStruct(SimpleStruct{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if spec.Type != "object" {
		t.Errorf("expected type object, got %s", spec.Type)
	}

	if len(spec.Required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(spec.Required))
	}

	nameSpec, ok := spec.Properties["name"]
	if !ok {
		t.Fatal("name property not found")
	}
	if nameSpec.Type != "string" {
		t.Errorf("expected name type string, got %s", nameSpec.Type)
	}
	if nameSpec.MinLength == nil || *nameSpec.MinLength != 1 {
		t.Errorf("expected name minLength 1, got %v", nameSpec.MinLength)
	}

	ageSpec, ok := spec.Properties["age"]
	if !ok {
		t.Fatal("age property not found")
	}
	if ageSpec.Type != "integer" {
		t.Errorf("expected age type integer, got %s", ageSpec.Type)
	}
	if ageSpec.Min == nil || *ageSpec.Min != 0 {
		t.Errorf("expected age min 0, got %v", ageSpec.Min)
	}
}

func TestSpecFromStructNested(t *testing.T) {
	type Address struct {
		Street string `json:"street" mowgli:"required"`
		City   string `json:"city" mowgli:"required"`
	}

	type User struct {
		Name    string  `json:"name" mowgli:"required"`
		Address Address `json:"address"`
	}

	spec, err := SpecFromStruct(User{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	addressSpec, ok := spec.Properties["address"]
	if !ok {
		t.Fatal("address property not found")
	}

	if addressSpec.Type != "object" {
		t.Errorf("expected address type object, got %s", addressSpec.Type)
	}

	streetSpec, ok := addressSpec.Properties["street"]
	if !ok {
		t.Fatal("street property not found")
	}
	if streetSpec.Type != "string" {
		t.Errorf("expected street type string, got %s", streetSpec.Type)
	}
}

func TestSpecFromStructArray(t *testing.T) {
	type User struct {
		Tags []string `json:"tags" mowgli:"minLength=0,maxLength=10"`
	}

	spec, err := SpecFromStruct(User{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tagsSpec, ok := spec.Properties["tags"]
	if !ok {
		t.Fatal("tags property not found")
	}

	if tagsSpec.Type != "array" {
		t.Errorf("expected tags type array, got %s", tagsSpec.Type)
	}

	if tagsSpec.Items == nil {
		t.Fatal("tags items spec not found")
	}

	if tagsSpec.Items.Type != "string" {
		t.Errorf("expected tags items type string, got %s", tagsSpec.Items.Type)
	}
}

func TestMergeSpecs(t *testing.T) {
	base := &Spec{
		Type: "object",
		Properties: map[string]*Spec{
			"name": {
				Type:      "string",
				MinLength: intPtr(1),
			},
		},
		Required: []string{"name"},
	}

	override := &Spec{
		Properties: map[string]*Spec{
			"name": {
				MaxLength: intPtr(100),
			},
		},
	}

	merged := MergeSpecs(base, override)

	if merged.Type != "object" {
		t.Errorf("expected type object, got %s", merged.Type)
	}

	nameSpec := merged.Properties["name"]
	if nameSpec.MinLength == nil || *nameSpec.MinLength != 1 {
		t.Errorf("expected minLength to be preserved from base")
	}
	if nameSpec.MaxLength == nil || *nameSpec.MaxLength != 100 {
		t.Errorf("expected maxLength to be overridden")
	}

	if len(merged.Required) != 1 || merged.Required[0] != "name" {
		t.Errorf("expected required fields to be preserved")
	}
}

func intPtr(i int) *int {
	return &i
}
