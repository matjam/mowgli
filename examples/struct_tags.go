package examples

import (
	"github.com/matjam/mowgli"
)

// Example: Using struct tags for validation
// This example demonstrates how to use struct tags to define
// validation rules directly in your Go structs.

type Product struct {
	Name        string  `json:"name" mowgli:"required,minLength=1,maxLength=200"`
	Price       float64 `json:"price" mowgli:"required,min=0"`
	InStock     bool    `json:"inStock"`
	Quantity    int     `json:"quantity" mowgli:"min=0"`
	SKU         string  `json:"sku" mowgli:"required,pattern=^[A-Z0-9-]+$"`
	Category    string  `json:"category" mowgli:"enum=electronics,clothing,food,books"`
	Description string  `json:"description" mowgli:"maxLength=1000"`
}

type Order struct {
	OrderID    string          `json:"orderId" mowgli:"required"`
	CustomerID string          `json:"customerId" mowgli:"required"`
	Items      []OrderItem     `json:"items" mowgli:"minLength=1,maxLength=100"`
	Shipping   ShippingAddress `json:"shipping"`
}

type OrderItem struct {
	ProductID string  `json:"productId" mowgli:"required"`
	Quantity  int     `json:"quantity" mowgli:"required,min=1"`
	Price     float64 `json:"price" mowgli:"required,min=0"`
}

type ShippingAddress struct {
	Street  string `json:"street" mowgli:"required"`
	City    string `json:"city" mowgli:"required"`
	State   string `json:"state" mowgli:"required,minLength=2,maxLength=2"`
	ZipCode string `json:"zipCode" mowgli:"required,pattern=^[0-9]{5}(-[0-9]{4})?$"`
	Country string `json:"country" mowgli:"required,minLength=2,maxLength=2"`
}

func ValidateProduct(data map[string]any) (*mowgli.ValidationResult, Product, error) {
	result, typedProduct, err := mowgli.ValidateStruct[Product](data)
	return result, typedProduct, err
}

func ValidateOrder(data map[string]any) (*mowgli.ValidationResult, Order, error) {
	result, typedOrder, err := mowgli.ValidateStruct[Order](data)
	return result, typedOrder, err
}

func ValidateProductWithOverride(data map[string]any, overrideSpec string) (*mowgli.ValidationResult, Product, error) {
	result, typedProduct, err := mowgli.ValidateStruct[Product](data, overrideSpec)
	return result, typedProduct, err
}
