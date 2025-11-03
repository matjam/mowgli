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

### Conditional Validation with Expressions

When your spec contains expressions in `conditions` (e.g., `"if": "enabled == true"`), the library automatically detects this and uses server-side validation. You need to provide an endpoint URL for validation:

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
    }
  ]
}`);

// With expressions, provide endpoint for server-side validation
const result = await validate(
  { enabled: true, value: "hello" },
  spec,
  { endpoint: "/api/users" }  // Server endpoint that handles X-Mowgli: validate header
);

// Note: validate() may return a Promise when expressions are present
if (result instanceof Promise) {
  const resolved = await result;
  console.log(resolved.valid);
} else {
  console.log(result.valid);
}
```

**Important:** When a spec contains expressions (conditions with `if` fields), you must provide an `endpoint` option. The library will automatically:
1. Detect that expressions are present
2. Make a POST request to your endpoint with `X-Mowgli: validate` header
3. Send the data as JSON in the request body
4. Parse validation errors from a 400 response, or treat 200 as valid

If no endpoint is provided, the library will perform basic constraint validation (required, min, max, etc.) but will skip expression evaluation.

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

- ✅ No external runtime dependencies (except for server-side expression evaluation)
- ✅ Full TypeScript support
- ✅ Automatic server-side validation for expressions (conditions)
- ✅ Client-side validation for basic constraints (required, min, max, minLength, etc.)
- ✅ All standard JSON types (string, number, integer, boolean, object, array, null)
- ✅ Comprehensive constraints (min, max, minLength, maxLength, pattern, enum, allowEmpty)
- ✅ Nested object and array validation
- ✅ Detailed error messages with field paths
- ✅ Transparent handling of sync/async validation (auto-detects when server-side is needed)

## License

MIT

