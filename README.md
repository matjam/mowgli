# Mowgli - JSON Validation Library

Mowgli is a validation library for Go and JavaScript that validates JSON data using a JSON-based specification format. The spec format is designed to work in both languages, making it useful for projects that need consistent validation across backend and frontend.

## Features

Mowgli supports standard JSON types (string, number, integer, boolean, object, array, null) with constraints like min/max values, string lengths, regex patterns, and enums. It also includes conditional validation where field requirements can depend on other field values.

The library is still being developed and missing some features, but provides a solid foundation for JSON validation needs. Future versions will expand on these capabilities.

## Quick Start

### Go - Using JSON Specs

```go
import "github.com/matjam/mowgli"

specJSON := `{
  "type": "object",
  "properties": {
    "name": {"type": "string", "minLength": 1, "maxLength": 100},
    "age": {"type": "integer", "min": 0, "max": 150}
  },
  "required": ["name"]
}`

spec, _ := mowgli.ParseSpecString(specJSON)
data := map[string]any{"name": "John", "age": 30}
result := mowgli.Validate(data, spec)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("%s: %s\n", err.Path, err.Message)
    }
}
```

### Go - Using Struct Tags

```go
type User struct {
    Name string `json:"name" mowgli:"required,minLength=1,maxLength=100"`
    Age  int    `json:"age" mowgli:"min=0,max=150"`
}

data := map[string]any{"name": "John", "age": 30}
result, user, err := mowgli.ValidateStruct[User](data)
// user is of type User, fully typed and validated
```

### JavaScript/TypeScript

```typescript
import { validate, parseSpec } from 'mowgli';

const spec = parseSpec(`{
  "type": "object",
  "properties": {
    "name": {"type": "string", "minLength": 1}
  },
  "required": ["name"]
}`);

const result = validate({name: "John"}, spec);
if (!result.valid) {
    result.errors.forEach(err => console.error(`${err.path}: ${err.message}`));
}
```

## Examples

The `examples/` directory contains working examples demonstrating various use cases:

- **[User Registration](examples/user_registration.go)** - Validating registration forms with email, password, and age requirements
- **[Configuration Files](examples/config_validation.go)** - Validating application config with nested objects
- **[Conditional Validation](examples/conditional_validation.go)** - Fields that become required based on other field values
- **[Struct Tags](examples/struct_tags.go)** - Using struct tags for type-safe validation
- **[API Requests](examples/api_request.go)** - Complex API payload validation with conditional rules

Each example includes comprehensive tests showing both valid and invalid inputs.

## Specification Format

Specs are JSON objects that define type and constraints. See [example_spec.json](example_spec.json) for a complete example.

**Supported constraints:**
- Strings: `minLength`, `maxLength`, `pattern`, `enum`, `allowEmpty`
- Numbers/Integers: `min`, `max`, `enum`
- Arrays: `minLength`, `maxLength`, `items` (for item validation)
- Objects: `properties`, `required`, `conditions` (for conditional validation)

**Conditional validation** uses simple expressions in the `if` field (e.g., `"enabled == true"`, `"count > 0"`) and applies spec overrides in `then`/`else` blocks.

## Installation

**Go:**
```bash
go get github.com/matjam/mowgli
```

**JavaScript/TypeScript:**
```bash
yarn add mowgli
# or
npm install mowgli
```

## License

MIT

## Acknowledgments

An LLM was used extensively to help with the implementation of this library. The code was carefully inspected, tested, and adjusted by humans to ensure a clean and accurate implementation.
