package mowgli

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
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
func evalExpression(expr string, obj map[string]any) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false, fmt.Errorf("empty expression")
	}

	// Tokenize the expression
	tokens, err := tokenizeExpression(expr)
	if err != nil {
		return false, err
	}

	// Parse and evaluate using recursive descent
	parser := &expressionParser{tokens: tokens, obj: obj, pos: 0}
	result, err := parser.parseExpression()
	if err != nil {
		return false, err
	}
	if parser.pos < len(tokens) {
		return false, fmt.Errorf("unexpected token at position %d: %v", parser.pos, tokens[parser.pos])
	}
	return result, nil
}

type tokenType int

const (
	tokenLiteral tokenType = iota
	tokenOperator
	tokenAnd
	tokenOr
	tokenLParen
	tokenRParen
)

type token struct {
	typ   tokenType
	value string
}

type expressionParser struct {
	tokens []token
	obj    map[string]any
	pos    int
}

func (p *expressionParser) peek() *token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.pos]
}

func (p *expressionParser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *expressionParser) parseExpression() (bool, error) {
	// Parse OR expression (lowest precedence)
	return p.parseOr()
}

func (p *expressionParser) parseOr() (bool, error) {
	left, err := p.parseAnd()
	if err != nil {
		return false, err
	}

	for tok := p.peek(); tok != nil && tok.typ == tokenOr; {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return false, err
		}
		left = left || right
		tok = p.peek()
	}

	return left, nil
}

func (p *expressionParser) parseAnd() (bool, error) {
	left, err := p.parseComparison()
	if err != nil {
		return false, err
	}

	for tok := p.peek(); tok != nil && tok.typ == tokenAnd; {
		p.advance()
		right, err := p.parseComparison()
		if err != nil {
			return false, err
		}
		left = left && right
		tok = p.peek()
	}

	return left, nil
}

func (p *expressionParser) parseComparison() (bool, error) {
	tok := p.peek()
	if tok == nil {
		return false, fmt.Errorf("unexpected end of expression")
	}

	// Handle parentheses
	if tok.typ == tokenLParen {
		p.advance()
		result, err := p.parseExpression()
		if err != nil {
			return false, err
		}
		next := p.peek()
		if next == nil || next.typ != tokenRParen {
			return false, fmt.Errorf("expected closing parenthesis")
		}
		p.advance()
		return result, nil
	}

	// Parse a literal or comparison
	if tok.typ != tokenLiteral {
		return false, fmt.Errorf("expected literal or comparison, got %v", tok)
	}

	// Try to parse as comparison first
	lit := tok.value
	p.advance()

	// Check if next token is an operator (comparison)
	next := p.peek()
	if next != nil && next.typ == tokenOperator {
		op := next.value
		p.advance()

		// Get the right operand
		rightTok := p.peek()
		if rightTok == nil || rightTok.typ != tokenLiteral {
			return false, fmt.Errorf("expected right operand for operator %s", op)
		}
		right := rightTok.value
		p.advance()

		return evaluateComparison(lit, op, right, p.obj)
	}

	// Not a comparison, treat as boolean literal or field reference
	return evaluateLiteral(lit, p.obj)
}

func tokenizeExpression(expr string) ([]token, error) {
	var tokens []token
	expr = strings.TrimSpace(expr)

	i := 0
	for i < len(expr) {
		// Skip whitespace
		for i < len(expr) && unicode.IsSpace(rune(expr[i])) {
			i++
		}
		if i >= len(expr) {
			break
		}

		// Check for parentheses
		if expr[i] == '(' {
			tokens = append(tokens, token{typ: tokenLParen, value: "("})
			i++
			continue
		}
		if expr[i] == ')' {
			tokens = append(tokens, token{typ: tokenRParen, value: ")"})
			i++
			continue
		}

		// Check for AND/OR (case insensitive, but require word boundaries)
		remaining := expr[i:]
		if i == 0 || unicode.IsSpace(rune(expr[i-1])) || expr[i-1] == '(' {
			if strings.HasPrefix(strings.ToUpper(remaining), "AND") &&
				(i+3 >= len(expr) || unicode.IsSpace(rune(expr[i+3])) || expr[i+3] == ')') {
				tokens = append(tokens, token{typ: tokenAnd, value: "AND"})
				i += 3
				continue
			}
			if strings.HasPrefix(strings.ToUpper(remaining), "OR") &&
				(i+2 >= len(expr) || unicode.IsSpace(rune(expr[i+2])) || expr[i+2] == ')') {
				tokens = append(tokens, token{typ: tokenOr, value: "OR"})
				i += 2
				continue
			}
		}

		// Check for comparison operators (longest first)
		operators := []string{"!=", "==", ">=", "<=", ">", "<"}
		foundOp := false
		for _, op := range operators {
			if strings.HasPrefix(remaining, op) {
				tokens = append(tokens, token{typ: tokenOperator, value: op})
				i += len(op)
				foundOp = true
				break
			}
		}
		if foundOp {
			continue
		}

		// Parse literal (field name, value, etc.)
		start := i
		inQuotes := false
		quoteChar := byte(0)

		for i < len(expr) {
			char := expr[i]

			// Handle quoted strings
			if (char == '"' || char == '\'') && !inQuotes {
				inQuotes = true
				quoteChar = char
				i++
				continue
			}

			if inQuotes {
				if char == quoteChar {
					// Check if it's escaped (only for same quote type)
					if i+1 < len(expr) && expr[i+1] == quoteChar {
						i += 2
						continue
					}
					// End of quoted string
					i++
					break
				}
				i++
				continue
			}

			// Stop at operators, parentheses, or whitespace
			if char == '(' || char == ')' || unicode.IsSpace(rune(char)) {
				break
			}

			// Check if we're hitting an operator
			remaining := expr[i:]
			hitOperator := false
			for _, op := range operators {
				if strings.HasPrefix(remaining, op) {
					hitOperator = true
					break
				}
			}
			if hitOperator {
				break
			}
			// Check for AND/OR (with word boundary check)
			if strings.HasPrefix(strings.ToUpper(remaining), "AND") &&
				(i == 0 || unicode.IsSpace(rune(expr[i-1])) || expr[i-1] == ')') &&
				(i+3 >= len(expr) || unicode.IsSpace(rune(expr[i+3])) || expr[i+3] == '(') {
				hitOperator = true
			}
			if strings.HasPrefix(strings.ToUpper(remaining), "OR") &&
				(i == 0 || unicode.IsSpace(rune(expr[i-1])) || expr[i-1] == ')') &&
				(i+2 >= len(expr) || unicode.IsSpace(rune(expr[i+2])) || expr[i+2] == '(') {
				hitOperator = true
			}
			if hitOperator {
				break
			}

			i++
		}

		literal := strings.TrimSpace(expr[start:i])
		if literal != "" {
			tokens = append(tokens, token{typ: tokenLiteral, value: literal})
		}
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens found in expression")
	}

	return tokens, nil
}

// evaluateLiteral evaluates a single literal (field reference or boolean value)
func evaluateLiteral(lit string, obj map[string]any) (bool, error) {
	lit = strings.TrimSpace(lit)

	// Check if it's a boolean literal
	if lit == "true" {
		return true, nil
	}
	if lit == "false" {
		return false, nil
	}

	// Check if it's a boolean field reference (shorthand for "fieldName == true")
	if val, exists := obj[lit]; exists {
		if b, ok := val.(bool); ok {
			return b, nil
		}
		// If it's not a boolean, treat as existence check
		return val != nil, nil
	}

	// Field doesn't exist
	return false, nil
}

// evaluateComparison evaluates a comparison expression
func evaluateComparison(left, op, right string, obj map[string]any) (bool, error) {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	leftVal, exists := obj[left]
	if !exists {
		// Field doesn't exist, evaluate based on operator and right side
		if right == "null" || right == "nil" {
			return op == "==", nil
		}
		return false, nil
	}

	rightVal, err := parseValue(right)
	if err != nil {
		return false, fmt.Errorf("invalid right operand in expression '%s %s %s': %w", left, op, right, err)
	}

	return compareValues(leftVal, rightVal, op)
}

func parseValue(s string) (any, error) {
	s = strings.TrimSpace(s)

	// null/nil
	if s == "null" || s == "nil" {
		return nil, nil
	}

	// boolean
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}

	// string (quoted)
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') && s[0] == s[len(s)-1] {
		// Remove quotes and unescape
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return s[1 : len(s)-1], nil // Fallback to simple unquote
		}
		return unquoted, nil
	}

	// number (int or float)
	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return f, nil
		}
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i, nil
	}

	return s, nil // Return as string if all else fails
}

func compareValues(left, right any, op string) (bool, error) {
	switch op {
	case "==":
		return deepEqual(left, right), nil
	case "!=":
		return !deepEqual(left, right), nil
	case ">", "<", ">=", "<=":
		return compareNumeric(left, right, op)
	default:
		return false, fmt.Errorf("unsupported operator: %s", op)
	}
}

func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Type assertion and comparison
	switch av := a.(type) {
	case bool:
		if bv, ok := b.(bool); ok {
			return av == bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av == bv
		}
	case float64:
		if bv, ok := b.(float64); ok {
			return av == bv
		}
		if bv, ok := b.(int); ok {
			return av == float64(bv)
		}
		if bv, ok := b.(int64); ok {
			return av == float64(bv)
		}
	case int:
		if bv, ok := b.(int); ok {
			return av == bv
		}
		if bv, ok := b.(int64); ok {
			return int64(av) == bv
		}
		if bv, ok := b.(float64); ok {
			return float64(av) == bv
		}
	case int64:
		if bv, ok := b.(int64); ok {
			return av == bv
		}
		if bv, ok := b.(int); ok {
			return av == int64(bv)
		}
		if bv, ok := b.(float64); ok {
			return float64(av) == bv
		}
	}

	return false
}

func compareNumeric(left, right any, op string) (bool, error) {
	leftNum := toFloat64(left)
	rightNum := toFloat64(right)

	if leftNum == nil || rightNum == nil {
		return false, fmt.Errorf("cannot compare non-numeric values with %s", op)
	}

	switch op {
	case ">":
		return *leftNum > *rightNum, nil
	case "<":
		return *leftNum < *rightNum, nil
	case ">=":
		return *leftNum >= *rightNum, nil
	case "<=":
		return *leftNum <= *rightNum, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", op)
	}
}

func toFloat64(v any) *float64 {
	switch val := v.(type) {
	case float64:
		return &val
	case int:
		f := float64(val)
		return &f
	case int64:
		f := float64(val)
		return &f
	case float32:
		f := float64(val)
		return &f
	default:
		return nil
	}
}
