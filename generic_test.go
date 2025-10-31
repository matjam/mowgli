package mowgli

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type User struct {
		Name string `json:"name" mowgli:"required,minLength=1,maxLength=100"`
		Age  int    `json:"age" mowgli:"required,min=0,max=150"`
	}

	tests := []struct {
		name      string
		data      any
		specJSON  string
		shouldErr bool
	}{
		{
			name: "valid data",
			data: map[string]any{
				"name": "John",
				"age":  30,
			},
			shouldErr: false,
		},
		{
			name: "missing required field",
			data: map[string]any{
				"name": "John",
			},
			shouldErr: true,
		},
		{
			name: "invalid age range",
			data: map[string]any{
				"name": "John",
				"age":  200,
			},
			shouldErr: true,
		},
		{
			name: "name too long",
			data: map[string]any{
				"name": "This is a very long name that definitely exceeds the maximum length of one hundred characters because it is way too long for the validation to pass",
				"age":  30,
			},
			shouldErr: true,
		},
		{
			name: "valid with JSON spec override",
			data: map[string]any{
				"name": "John",
				"age":  30,
			},
			specJSON: `{
				"properties": {
					"age": {
						"max": 200
					}
				}
			}`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user User
			result, typedResult, err := ValidateStructValue(user, tt.data, tt.specJSON)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid == tt.shouldErr {
				if tt.shouldErr {
					t.Errorf("expected validation to fail, but it passed")
					for _, err := range result.Errors {
						t.Logf("  error: %s", err)
					}
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
				}
			}

			if result.Valid {
				// Check that typed result is correct
				if typedResult.Name == "" && tt.data.(map[string]any)["name"] != nil {
					t.Errorf("typed result name is empty but should be set")
				}
				if typedResult.Age == 0 && tt.data.(map[string]any)["age"] != nil {
					t.Errorf("typed result age is 0 but should be set")
				}
			}
		})
	}
}

func TestValidateStructGeneric(t *testing.T) {
	type Product struct {
		Name    string  `json:"name" mowgli:"required,minLength=1"`
		Price   float64 `json:"price" mowgli:"required,min=0"`
		InStock bool    `json:"inStock"`
	}

	data := map[string]any{
		"name":    "Widget",
		"price":   29.99,
		"inStock": true,
	}

	result, product, err := ValidateStruct[Product](data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Valid {
		t.Fatalf("expected validation to pass, but it failed: %v", result.Errors)
	}

	if product.Name != "Widget" {
		t.Errorf("expected name Widget, got %s", product.Name)
	}

	if product.Price != 29.99 {
		t.Errorf("expected price 29.99, got %f", product.Price)
	}

	if !product.InStock {
		t.Errorf("expected inStock true, got %v", product.InStock)
	}
}

func TestValidateStructNested(t *testing.T) {
	type Address struct {
		Street string `json:"street" mowgli:"required"`
		City   string `json:"city" mowgli:"required"`
	}

	type User struct {
		Name    string  `json:"name" mowgli:"required"`
		Address Address `json:"address"`
	}

	data := map[string]any{
		"name": "John",
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "Springfield",
		},
	}

	result, user, err := ValidateStruct[User](data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Valid {
		t.Fatalf("expected validation to pass, but it failed: %v", result.Errors)
	}

	if user.Name != "John" {
		t.Errorf("expected name John, got %s", user.Name)
	}

	if user.Address.Street != "123 Main St" {
		t.Errorf("expected street 123 Main St, got %s", user.Address.Street)
	}

	if user.Address.City != "Springfield" {
		t.Errorf("expected city Springfield, got %s", user.Address.City)
	}
}

func TestValidateStructArray(t *testing.T) {
	type User struct {
		Name string   `json:"name" mowgli:"required"`
		Tags []string `json:"tags" mowgli:"minLength=0,maxLength=10"`
	}

	data := map[string]any{
		"name": "John",
		"tags": []any{"developer", "go"},
	}

	result, user, err := ValidateStruct[User](data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Valid {
		t.Fatalf("expected validation to pass, but it failed: %v", result.Errors)
	}

	if len(user.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(user.Tags))
	}

	if user.Tags[0] != "developer" {
		t.Errorf("expected first tag developer, got %s", user.Tags[0])
	}
}
