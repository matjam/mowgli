/**
 * Evaluates an expression in the context of an object
 * Supports expressions like:
 *   - "fieldName == true"
 *   - "fieldName != false"
 *   - "fieldName > 0"
 *   - "fieldName AND otherField"
 *   - "field1 == value1 OR field2 == value2"
 *   - "(field1 == value1) AND (field2 == value2 OR field3 == value3)"
 *   - "fieldName != null"
 */
export function evalExpression(expr: string, obj: Record<string, any>): boolean {
  expr = expr.trim();
  if (expr === '') {
    throw new Error('empty expression');
  }

  // Tokenize the expression
  const tokens = tokenizeExpression(expr);

  // Parse and evaluate using recursive descent
  const parser = new ExpressionParser(tokens, obj);
  const result = parser.parseExpression();
  if (parser.pos < tokens.length) {
    throw new Error(`unexpected token at position ${parser.pos}: ${JSON.stringify(tokens[parser.pos])}`);
  }
  return result;
}

enum TokenType {
  Literal,
  Operator,
  And,
  Or,
  LParen,
  RParen,
}

interface Token {
  type: TokenType;
  value: string;
}

class ExpressionParser {
  tokens: Token[];
  obj: Record<string, any>;
  pos: number;

  constructor(tokens: Token[], obj: Record<string, any>) {
    this.tokens = tokens;
    this.obj = obj;
    this.pos = 0;
  }

  peek(): Token | null {
    if (this.pos >= this.tokens.length) {
      return null;
    }
    return this.tokens[this.pos];
  }

  advance(): void {
    if (this.pos < this.tokens.length) {
      this.pos++;
    }
  }

  parseExpression(): boolean {
    // Parse OR expression (lowest precedence)
    return this.parseOr();
  }

  parseOr(): boolean {
    let left = this.parseAnd();

    let tok = this.peek();
    while (tok !== null && tok.type === TokenType.Or) {
      this.advance();
      const right = this.parseAnd();
      left = left || right;
      tok = this.peek();
    }

    return left;
  }

  parseAnd(): boolean {
    let left = this.parseComparison();

    let tok = this.peek();
    while (tok !== null && tok.type === TokenType.And) {
      this.advance();
      const right = this.parseComparison();
      left = left && right;
      tok = this.peek();
    }

    return left;
  }

  parseComparison(): boolean {
    const tok = this.peek();
    if (tok === null) {
      throw new Error('unexpected end of expression');
    }

    // Handle parentheses
    if (tok.type === TokenType.LParen) {
      this.advance();
      const result = this.parseExpression();
      const next = this.peek();
      if (next === null || next.type !== TokenType.RParen) {
        throw new Error('expected closing parenthesis');
      }
      this.advance();
      return result;
    }

    // Parse a literal or comparison
    if (tok.type !== TokenType.Literal) {
      throw new Error(`expected literal or comparison, got ${JSON.stringify(tok)}`);
    }

    // Try to parse as comparison first
    const lit = tok.value;
    this.advance();

    // Check if next token is an operator (comparison)
    const next = this.peek();
    if (next !== null && next.type === TokenType.Operator) {
      const op = next.value;
      this.advance();

      // Get the right operand
      const rightTok = this.peek();
      if (rightTok === null || rightTok.type !== TokenType.Literal) {
        throw new Error(`expected right operand for operator ${op}`);
      }
      const right = rightTok.value;
      this.advance();

      return evaluateComparison(lit, op, right, this.obj);
    }

    // Not a comparison, treat as boolean literal or field reference
    return evaluateLiteral(lit, this.obj);
  }
}

function tokenizeExpression(expr: string): Token[] {
  const tokens: Token[] = [];
  expr = expr.trim();

  let i = 0;
  while (i < expr.length) {
    // Skip whitespace
    while (i < expr.length && /\s/.test(expr[i])) {
      i++;
    }
    if (i >= expr.length) {
      break;
    }

    // Check for parentheses
    if (expr[i] === '(') {
      tokens.push({ type: TokenType.LParen, value: '(' });
      i++;
      continue;
    }
    if (expr[i] === ')') {
      tokens.push({ type: TokenType.RParen, value: ')' });
      i++;
      continue;
    }

    // Check for AND/OR (case insensitive, but require word boundaries)
    const remaining = expr.substring(i);
    if (
      (i === 0 || /\s|\)/.test(expr[i - 1])) &&
      remaining.toUpperCase().startsWith('AND') &&
      (i + 3 >= expr.length || /\s|\(/.test(expr[i + 3]))
    ) {
      tokens.push({ type: TokenType.And, value: 'AND' });
      i += 3;
      continue;
    }
    if (
      (i === 0 || /\s|\)/.test(expr[i - 1])) &&
      remaining.toUpperCase().startsWith('OR') &&
      (i + 2 >= expr.length || /\s|\(/.test(expr[i + 2]))
    ) {
      tokens.push({ type: TokenType.Or, value: 'OR' });
      i += 2;
      continue;
    }

    // Check for comparison operators (longest first)
    const operators = ['!=', '==', '>=', '<=', '>', '<'];
    let foundOp = false;
    for (const op of operators) {
      if (remaining.startsWith(op)) {
        tokens.push({ type: TokenType.Operator, value: op });
        i += op.length;
        foundOp = true;
        break;
      }
    }
    if (foundOp) {
      continue;
    }

    // Parse literal (field name, value, etc.)
    const start = i;
    let inQuotes = false;
    let quoteChar: string | null = null;

    while (i < expr.length) {
      const char = expr[i];

      // Handle quoted strings
      if ((char === '"' || char === "'") && !inQuotes) {
        inQuotes = true;
        quoteChar = char;
        i++;
        continue;
      }

      if (inQuotes) {
        if (char === quoteChar) {
          // Check if it's escaped
          if (i + 1 < expr.length && expr[i + 1] === quoteChar) {
            i += 2;
            continue;
          }
          // End of quoted string
          i++;
          break;
        }
        i++;
        continue;
      }

      // Stop at operators, parentheses, or whitespace
      if (char === '(' || char === ')' || /\s/.test(char)) {
        break;
      }

      // Check if we're hitting an operator
      const remaining2 = expr.substring(i);
      let hitOperator = false;
      for (const op of operators) {
        if (remaining2.startsWith(op)) {
          hitOperator = true;
          break;
        }
      }
      if (hitOperator) {
        break;
      }
      // Check for AND/OR (with word boundary check)
      if (
        remaining2.toUpperCase().startsWith('AND') &&
        (i === 0 || /\s|\)/.test(expr[i - 1])) &&
        (i + 3 >= expr.length || /\s|\(/.test(expr[i + 3]))
      ) {
        hitOperator = true;
      }
      if (
        remaining2.toUpperCase().startsWith('OR') &&
        (i === 0 || /\s|\)/.test(expr[i - 1])) &&
        (i + 2 >= expr.length || /\s|\(/.test(expr[i + 2]))
      ) {
        hitOperator = true;
      }
      if (hitOperator) {
        break;
      }

      i++;
    }

    const literal = expr.substring(start, i).trim();
    if (literal !== '') {
      tokens.push({ type: TokenType.Literal, value: literal });
    }
  }

  if (tokens.length === 0) {
    throw new Error('no tokens found in expression');
  }

  return tokens;
}

// evaluateLiteral evaluates a single literal (field reference or boolean value)
function evaluateLiteral(lit: string, obj: Record<string, any>): boolean {
  lit = lit.trim();

  // Check if it's a boolean literal
  if (lit === 'true') {
    return true;
  }
  if (lit === 'false') {
    return false;
  }

  // Check if it's a boolean field reference (shorthand for "fieldName == true")
  if (obj.hasOwnProperty(lit)) {
    const val = obj[lit];
    if (typeof val === 'boolean') {
      return val;
    }
    // If it's not a boolean, treat as existence check
    return val != null;
  }

  // Field doesn't exist
  return false;
}

// evaluateComparison evaluates a comparison expression
function evaluateComparison(
  left: string,
  op: string,
  right: string,
  obj: Record<string, any>
): boolean {
  left = left.trim();
  right = right.trim();

  const leftVal = obj[left];
  if (!obj.hasOwnProperty(left)) {
    // Field doesn't exist, evaluate based on operator and right side
    if (right === 'null' || right === 'nil') {
      return op === '==';
    }
    return false;
  }

  const rightVal = parseValue(right);
  return compareValues(leftVal, rightVal, op);
}

function parseValue(s: string): any {
  s = s.trim();

  // null/nil
  if (s === 'null' || s === 'nil') {
    return null;
  }

  // boolean
  if (s === 'true') {
    return true;
  }
  if (s === 'false') {
    return false;
  }

  // string (quoted)
  if (s.length >= 2 && (s[0] === '"' || s[0] === "'") && s[0] === s[s.length - 1]) {
    // Remove quotes and unescape
    try {
      return JSON.parse(s);
    } catch {
      return s.slice(1, -1); // Fallback to simple unquote
    }
  }

  // number (int or float)
  if (s.includes('.')) {
    const f = parseFloat(s);
    if (!isNaN(f)) {
      return f;
    }
  }

  const i = parseInt(s, 10);
  if (!isNaN(i)) {
    return i;
  }

  return s; // Return as string if all else fails
}

function compareValues(left: any, right: any, op: string): boolean {
  switch (op) {
    case '==':
      return deepEqual(left, right);
    case '!=':
      return !deepEqual(left, right);
    case '>':
    case '<':
    case '>=':
    case '<=':
      return compareNumeric(left, right, op);
    default:
      throw new Error(`Unsupported operator: ${op}`);
  }
}

function deepEqual(a: any, b: any): boolean {
  if (a === null && b === null) {
    return true;
  }
  if (a === null || b === null) {
    return false;
  }

  // Handle type coercion for numbers
  if (typeof a === 'number' && typeof b === 'number') {
    return a === b;
  }
  if (typeof a === 'number' && typeof b === 'string') {
    const bNum = parseFloat(b);
    return !isNaN(bNum) && a === bNum;
  }
  if (typeof a === 'string' && typeof b === 'number') {
    const aNum = parseFloat(a);
    return !isNaN(aNum) && aNum === b;
  }

  return a === b;
}

function compareNumeric(left: any, right: any, op: string): boolean {
  const leftNum = toNumber(left);
  const rightNum = toNumber(right);

  if (leftNum === null || rightNum === null) {
    throw new Error(`Cannot compare non-numeric values with ${op}`);
  }

  switch (op) {
    case '>':
      return leftNum > rightNum;
    case '<':
      return leftNum < rightNum;
    case '>=':
      return leftNum >= rightNum;
    case '<=':
      return leftNum <= rightNum;
    default:
      throw new Error(`Unsupported operator: ${op}`);
  }
}

function toNumber(v: any): number | null {
  if (typeof v === 'number') {
    return v;
  }
  if (typeof v === 'string') {
    const num = parseFloat(v);
    return isNaN(num) ? null : num;
  }
  return null;
}
