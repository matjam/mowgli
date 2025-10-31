# Mowgli - JavaScript/TypeScript Validation Library

A self-contained JavaScript/TypeScript validation library that validates JSON data against a JSON-based validation specification. Compatible with the Go implementation.

## Installation

```bash
yarn add mowgli
```

Or using npm:

```bash
npm install mowgli
```

## Usage

### Basic Example

```typescript
import { validate, validateJSON, parseSpec } from 'mowgli';

// Parse a spec
const specJSON = `{
  "type": "object",
  "properties": {
    "name": {"type": "string", "minLength": 1},
    "age": {"type": "integer", "min": 0, "max": 150}
  },
  "required": ["name"]
}`;

const spec = parseSpec(specJSON);

// Validate a value
const data = { name: "John", age: 30 };
const result = validate(data, spec);

if (!result.valid) {
  result.errors.forEach(err => {
    console.error(`${err.path}: ${err.message}`);
  });
}

// Validate JSON string directly
const jsonResult = validateJSON('{"name": "John", "age": 30}', spec);
```

### Conditional Validation

```typescript
const spec = parseSpec(`{
  "type": "object",
  "properties": {
    "enabled": {"type": "boolean"},
    "value": {"type": "string"}
  },
  "conditions": [
    {
      "if": "enabled == true",
      "then": {
        "value": {"minLength": 1}
      }
    },
    {
      "if": "enabled == false",
      "then": {
        "value": {"allowEmpty": true}
      }
    }
  ]
}`);

// When enabled is true, value must have minLength 1
const result1 = validate({ enabled: true, value: "hello" }, spec);
console.log(result1.valid); // true

// When enabled is false, value can be empty
const result2 = validate({ enabled: false, value: "" }, spec);
console.log(result2.valid); // true
```

### TypeScript Support

Full TypeScript types are included:

```typescript
import type { Spec, ValidationResult, ValidationError } from 'mowgli';

const spec: Spec = {
  type: 'object',
  properties: {
    name: { type: 'string' }
  }
};

const result: ValidationResult = validate({ name: 'John' }, spec);
```

## API

### `validate(value: any, spec: Spec): ValidationResult`

Validates a value against a spec.

### `validateJSON(jsonData: string, spec: Spec): ValidationResult`

Validates a JSON string against a spec.

### `parseSpec(specJSON: string): Spec`

Parses a JSON spec string into a Spec object.

### Types

- `Spec` - Validation specification
- `ValidationResult` - Result of validation with `valid` boolean and `errors` array
- `ValidationError` - Individual error with `path` and `message`
- `Condition` - Conditional validation rule

## Validation Spec Format

See the main [README.md](../README.md) for the complete validation spec format documentation.

## Building

```bash
yarn build
```

This will:
1. Compile TypeScript to JavaScript
2. Generate TypeScript declaration files (.d.ts)
3. Create a minified bundle (dist/index.min.js)

## Development

```bash
# Install dependencies
yarn install

# Run tests
yarn test

# Build
yarn build

# Clean build artifacts
yarn clean
```

## Features

- ✅ No external runtime dependencies
- ✅ Full TypeScript support
- ✅ Conditional validation with expression evaluation
- ✅ All standard JSON types (string, number, integer, boolean, object, array, null)
- ✅ Comprehensive constraints (min, max, minLength, maxLength, pattern, enum, allowEmpty)
- ✅ Nested object and array validation
- ✅ Detailed error messages with field paths

## License

MIT

