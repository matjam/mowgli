package mowgli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
)

// evalExpression evaluates an expression in the context of an object
// Supports expressions like:
//   - "fieldName == true"
//   - "fieldName != false"
//   - "fieldName > 0"
//   - "fieldName AND otherField"
//   - "field1 == value1 OR field2 == value2"
//   - "(field1 == value1) AND (field2 == value2 OR field3 == value3)"
//   - "fieldName != null"
func evalExpression(exprStr string, obj map[string]any) (bool, error) {
	exprStr = strings.TrimSpace(exprStr)
	if exprStr == "" {
		return false, fmt.Errorf("empty expression")
	}

	// Translate AND/OR to &&/|| for expr library compatibility
	// Use word boundaries to avoid replacing inside other words
	translatedExpr := translateExpression(exprStr)

	// Replace "null" and "nil" with nil for expr compatibility
	translatedExpr = strings.ReplaceAll(translatedExpr, " null ", " nil ")
	translatedExpr = strings.ReplaceAll(translatedExpr, " null", " nil")
	translatedExpr = strings.ReplaceAll(translatedExpr, "null ", "nil ")
	if translatedExpr == "null" {
		translatedExpr = "nil"
	}

	// Use Eval for dynamic map environments - this allows map keys to shadow built-in functions
	// when they exist in the provided object
	result, err := expr.Eval(translatedExpr, obj)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate expression '%s': %w", exprStr, err)
	}

	// Convert result to bool
	if boolResult, ok := result.(bool); ok {
		return boolResult, nil
	}

	return false, fmt.Errorf("expression '%s' did not evaluate to a boolean, got %T: %v", exprStr, result, result)
}

// translateExpression converts AND/OR to &&/|| while preserving word boundaries
func translateExpression(expr string) string {
	// Match AND/OR with word boundaries (whitespace, parentheses, operators, start/end of string)
	// Use negative lookbehind and lookahead to ensure we don't replace inside identifiers
	andPattern := regexp.MustCompile(`\bAND\b`)
	orPattern := regexp.MustCompile(`\bOR\b`)

	expr = andPattern.ReplaceAllString(expr, "&&")
	expr = orPattern.ReplaceAllString(expr, "||")

	return expr
}
