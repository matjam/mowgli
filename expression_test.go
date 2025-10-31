package mowgli

import (
	"testing"
)

func TestEvalExpressionAND(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		obj      map[string]any
		expected bool
		wantErr  bool
	}{
		{
			name:     "both true",
			expr:     "enabled AND active",
			obj:      map[string]any{"enabled": true, "active": true},
			expected: true,
		},
		{
			name:     "first false",
			expr:     "enabled AND active",
			obj:      map[string]any{"enabled": false, "active": true},
			expected: false,
		},
		{
			name:     "second false",
			expr:     "enabled AND active",
			obj:      map[string]any{"enabled": true, "active": false},
			expected: false,
		},
		{
			name:     "both false",
			expr:     "enabled AND active",
			obj:      map[string]any{"enabled": false, "active": false},
			expected: false,
		},
		{
			name:     "comparison AND comparison",
			expr:     "count > 0 AND status == \"active\"",
			obj:      map[string]any{"count": 5, "status": "active"},
			expected: true,
		},
		{
			name:     "comparison AND comparison false",
			expr:     "count > 0 AND status == \"active\"",
			obj:      map[string]any{"count": 0, "status": "active"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalExpression(tt.expr, tt.obj)
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
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalExpressionOR(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		obj      map[string]any
		expected bool
		wantErr  bool
	}{
		{
			name:     "both true",
			expr:     "enabled OR active",
			obj:      map[string]any{"enabled": true, "active": true},
			expected: true,
		},
		{
			name:     "first true",
			expr:     "enabled OR active",
			obj:      map[string]any{"enabled": true, "active": false},
			expected: true,
		},
		{
			name:     "second true",
			expr:     "enabled OR active",
			obj:      map[string]any{"enabled": false, "active": true},
			expected: true,
		},
		{
			name:     "both false",
			expr:     "enabled OR active",
			obj:      map[string]any{"enabled": false, "active": false},
			expected: false,
		},
		{
			name:     "comparison OR comparison",
			expr:     "count > 10 OR status == \"active\"",
			obj:      map[string]any{"count": 5, "status": "active"},
			expected: true,
		},
		{
			name:     "comparison OR comparison both false",
			expr:     "count > 10 OR status == \"inactive\"",
			obj:      map[string]any{"count": 5, "status": "active"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalExpression(tt.expr, tt.obj)
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
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalExpressionANDOR(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		obj      map[string]any
		expected bool
		wantErr  bool
	}{
		{
			name:     "AND before OR",
			expr:     "a AND b OR c",
			obj:      map[string]any{"a": true, "b": false, "c": true},
			expected: true, // (a AND b) OR c = false OR true = true
		},
		{
			name:     "OR before AND (need parentheses)",
			expr:     "a OR b AND c",
			obj:      map[string]any{"a": false, "b": true, "c": true},
			expected: true, // a OR (b AND c) = false OR true = true (AND has higher precedence)
		},
		{
			name:     "multiple AND",
			expr:     "a AND b AND c",
			obj:      map[string]any{"a": true, "b": true, "c": true},
			expected: true,
		},
		{
			name:     "multiple AND one false",
			expr:     "a AND b AND c",
			obj:      map[string]any{"a": true, "b": false, "c": true},
			expected: false,
		},
		{
			name:     "multiple OR",
			expr:     "a OR b OR c",
			obj:      map[string]any{"a": false, "b": false, "c": true},
			expected: true,
		},
		{
			name:     "complex expression",
			expr:     "count > 0 AND status == \"active\" OR enabled",
			obj:      map[string]any{"count": 0, "status": "inactive", "enabled": true},
			expected: true, // (count > 0 AND status == "active") OR enabled = false OR true = true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalExpression(tt.expr, tt.obj)
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
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalExpressionParentheses(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		obj      map[string]any
		expected bool
		wantErr  bool
	}{
		{
			name:     "simple parentheses",
			expr:     "(enabled)",
			obj:      map[string]any{"enabled": true},
			expected: true,
		},
		{
			name:     "parentheses with AND",
			expr:     "(a AND b)",
			obj:      map[string]any{"a": true, "b": true},
			expected: true,
		},
		{
			name:     "parentheses change precedence",
			expr:     "(a OR b) AND c",
			obj:      map[string]any{"a": true, "b": false, "c": true},
			expected: true, // (true OR false) AND true = true AND true = true
		},
		{
			name:     "parentheses change precedence 2",
			expr:     "a OR (b AND c)",
			obj:      map[string]any{"a": false, "b": true, "c": true},
			expected: true, // false OR (true AND true) = false OR true = true
		},
		{
			name:     "nested parentheses",
			expr:     "((a AND b) OR c)",
			obj:      map[string]any{"a": false, "b": true, "c": true},
			expected: true, // ((false AND true) OR true) = (false OR true) = true
		},
		{
			name:     "complex with parentheses",
			expr:     "(count > 0 AND status == \"active\") OR (enabled AND verified)",
			obj:      map[string]any{"count": 0, "status": "inactive", "enabled": true, "verified": true},
			expected: true, // (false AND false) OR (true AND true) = false OR true = true
		},
		{
			name:     "comparison in parentheses",
			expr:     "(count > 5)",
			obj:      map[string]any{"count": 10},
			expected: true,
		},
		{
			name:     "multiple parentheses groups",
			expr:     "(a == 1) AND (b == 2 OR c == 3)",
			obj:      map[string]any{"a": 1, "b": 3, "c": 3},
			expected: true, // (1 == 1) AND (3 == 2 OR 3 == 3) = true AND (false OR true) = true AND true = true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalExpression(tt.expr, tt.obj)
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
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalExpressionComplex(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		obj      map[string]any
		expected bool
		wantErr  bool
	}{
		{
			name:     "complex nested expression",
			expr:     "((a > 0 AND b < 10) OR (c == \"active\")) AND enabled",
			obj:      map[string]any{"a": 5, "b": 15, "c": "inactive", "enabled": true},
			expected: false, // ((5 > 0 AND 15 < 10) OR ("inactive" == "active")) AND true = (false OR false) AND true = false
		},
		{
			name:     "complex nested expression true",
			expr:     "((a > 0 AND b < 10) OR (c == \"active\")) AND enabled",
			obj:      map[string]any{"a": 5, "b": 5, "c": "inactive", "enabled": true},
			expected: true, // ((5 > 0 AND 5 < 10) OR false) AND true = (true OR false) AND true = true
		},
		{
			name:     "chained ORs and ANDs",
			expr:     "a AND b AND c OR d",
			obj:      map[string]any{"a": false, "b": false, "c": false, "d": true},
			expected: true, // ((false AND false) AND false) OR true = false OR true = true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalExpression(tt.expr, tt.obj)
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
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalExpressionErrors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		obj     map[string]any
		wantErr bool
	}{
		{
			name:    "mismatched parentheses - too many closing",
			expr:    "(a AND b))",
			obj:     map[string]any{"a": true, "b": true},
			wantErr: true,
		},
		{
			name:    "mismatched parentheses - too many opening",
			expr:    "((a AND b)",
			obj:     map[string]any{"a": true, "b": true},
			wantErr: true,
		},
		{
			name:    "AND without operand",
			expr:    "a AND",
			obj:     map[string]any{"a": true},
			wantErr: true,
		},
		{
			name:    "OR without operand",
			expr:    "OR a",
			obj:     map[string]any{"a": true},
			wantErr: true,
		},
		{
			name:    "empty expression",
			expr:    "",
			obj:     map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evalExpression(tt.expr, tt.obj)
			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
