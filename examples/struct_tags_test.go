package examples

import (
	"testing"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "valid product",
			data: map[string]any{
				"name":        "Laptop",
				"price":       1299.99,
				"inStock":     true,
				"quantity":    50,
				"sku":         "LAPTOP-001",
				"category":    "electronics",
				"description": "A high-performance laptop",
			},
			wantValid: true,
		},
		{
			name: "missing required name",
			data: map[string]any{
				"price": 1299.99,
			},
			wantValid: false,
		},
		{
			name: "negative price",
			data: map[string]any{
				"name":  "Laptop",
				"price": -10,
			},
			wantValid: false,
		},
		{
			name: "invalid SKU format",
			data: map[string]any{
				"name":  "Laptop",
				"price": 1299.99,
				"sku":   "invalid_sku",
			},
			wantValid: false,
		},
		{
			name: "invalid category",
			data: map[string]any{
				"name":     "Laptop",
				"price":    1299.99,
				"sku":      "LAPTOP-001",
				"category": "invalid",
			},
			wantValid: false,
		},
		{
			name: "description too long",
			data: map[string]any{
				"name":        "Laptop",
				"price":       1299.99,
				"sku":         "LAPTOP-001",
				"category":    "electronics",
				"description": string(make([]byte, 1001)),
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, product, err := ValidateProduct(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.wantValid {
				if !tt.wantValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
					for _, err := range result.Errors {
						t.Logf("  error: %s", err)
					}
				}
			}

			// If valid, verify typed result
			if result.Valid {
				if product.Name == "" && tt.data["name"] != nil {
					t.Errorf("typed result name is empty but should be set")
				}
			}
		})
	}
}

func TestValidateOrder(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]any
		wantValid bool
	}{
		{
			name: "valid order",
			data: map[string]any{
				"orderId":    "ORD-123",
				"customerId": "CUST-456",
				"items": []any{
					map[string]any{
						"productId": "PROD-001",
						"quantity":  2,
						"price":     29.99,
					},
				},
				"shipping": map[string]any{
					"street":  "123 Main St",
					"city":    "Springfield",
					"state":   "IL",
					"zipCode": "62701",
					"country": "US",
				},
			},
			wantValid: true,
		},
		{
			name: "order without items",
			data: map[string]any{
				"orderId":    "ORD-123",
				"customerId": "CUST-456",
				"items":      []any{},
			},
			wantValid: false,
		},
		{
			name: "invalid zip code",
			data: map[string]any{
				"orderId":    "ORD-123",
				"customerId": "CUST-456",
				"items": []any{
					map[string]any{
						"productId": "PROD-001",
						"quantity":  2,
						"price":     29.99,
					},
				},
				"shipping": map[string]any{
					"street":  "123 Main St",
					"city":    "Springfield",
					"state":   "IL",
					"zipCode": "invalid",
					"country": "US",
				},
			},
			wantValid: false,
		},
		{
			name: "item with zero quantity",
			data: map[string]any{
				"orderId":    "ORD-123",
				"customerId": "CUST-456",
				"items": []any{
					map[string]any{
						"productId": "PROD-001",
						"quantity":  0,
						"price":     29.99,
					},
				},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, order, err := ValidateOrder(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.wantValid {
				if !tt.wantValid {
					t.Errorf("expected validation to fail, but it passed")
				} else {
					t.Errorf("expected validation to pass, but it failed: %v", result.Errors)
					for _, err := range result.Errors {
						t.Logf("  error: %s", err)
					}
				}
			}

			// If valid, verify typed result
			if result.Valid {
				if order.OrderID == "" {
					t.Errorf("typed result OrderID is empty but should be set")
				}
			}
		})
	}
}

func TestValidateProductWithOverride(t *testing.T) {
	data := map[string]any{
		"name":     "Laptop",
		"price":    1299.99,
		"sku":      "LAPTOP-001",
		"category": "electronics",
	}

	// Override to allow higher price
	overrideSpec := `{
		"properties": {
			"price": {
				"max": 10000
			}
		}
	}`

	result, product, err := ValidateProductWithOverride(data, overrideSpec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Valid {
		t.Fatalf("expected validation to pass, but it failed: %v", result.Errors)
	}

	if product.Name != "Laptop" {
		t.Errorf("expected name Laptop, got %s", product.Name)
	}
}
