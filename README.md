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

### Client-Side Form Validation

Mowgli provides framework-agnostic helpers that work with any JavaScript framework (React, Vue, Svelte, Angular, vanilla JS, etc.):

```typescript
import { fetchSpec, createFormValidator, groupErrorsByField } from 'mowgli';

// Fetch validation spec from API endpoint
const spec = await fetchSpec('/api/users', {
  headerName: 'X-Mowgli',
  headerValue: 'Validation Request'
});

// Create a reusable validator
const validator = createFormValidator(spec);

// Validate entire form on submit
const formValues = { username: 'john', email: 'john@example.com', password: 'pass123' };
const result = validator.validateForm(formValues);
if (!result.valid) {
  const errorsByField = groupErrorsByField(result);
  // Display errors to user
  // { username: "string length 4 is less than minimum 5", ... }
}

// Validate single field on change (for real-time feedback)
const fieldResult = validator.validateField('username', 'ab', formValues);
if (!fieldResult.valid) {
  // Show error for this field: fieldResult.error
  // "string length 2 is less than minimum 3"
}
```

**Available helpers:**
- **`fetchSpec`** - Fetches and parses a validation spec from an API endpoint
- **`validateField`** - Validates a single field in isolation (with form context for conditional validation)
- **`groupErrorsByField`** - Groups validation errors by field name (first error per field)
- **`groupAllErrorsByField`** - Groups all validation errors by field name (all errors per field)
- **`createFormValidator`** - Creates a reusable validator with helper methods for a spec

**Example: Using with a framework**

Here's how you might use these helpers in a React component (or similar pattern for Vue/Svelte/etc.):

```typescript
import { useState, useEffect } from 'react';
import { fetchSpec, createFormValidator, groupErrorsByField } from 'mowgli';

function UserRegistrationForm() {
  const [values, setValues] = useState({ username: '', email: '', password: '' });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [validator, setValidator] = useState<any>(null);

  useEffect(() => {
    fetchSpec('/api/users').then(spec => {
      setValidator(createFormValidator(spec));
    });
  }, []);

  const handleChange = (field: string, value: string) => {
    const newValues = { ...values, [field]: value };
    setValues(newValues);
    
    // Real-time field validation
    if (validator) {
      const result = validator.validateField(field, value, newValues);
      setErrors(prev => ({
        ...prev,
        [field]: result.error || ''
      }));
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validator) return;
    
    const result = validator.validateForm(values);
    if (!result.valid) {
      setErrors(groupErrorsByField(result));
      return;
    }
    
    // Submit form...
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={values.username}
        onChange={(e) => handleChange('username', e.target.value)}
      />
      {errors.username && <span>{errors.username}</span>}
      {/* ... more fields ... */}
    </form>
  );
}
```

These helpers work with any framework—you just integrate them with your framework's state management and event handling patterns.

## Examples

The `examples/` directory contains working examples demonstrating various use cases:

- **[User Registration](examples/user_registration.go)** - Validating registration forms with email, password, and age requirements
- **[Configuration Files](examples/config_validation.go)** - Validating application config with nested objects
- **[Conditional Validation](examples/conditional_validation.go)** - Fields that become required based on other field values
- **[Struct Tags](examples/struct_tags.go)** - Using struct tags for type-safe validation
- **[API Requests](examples/api_request.go)** - Complex API payload validation with conditional rules

Each example includes comprehensive tests showing both valid and invalid inputs.

## API Service Usage

Mowgli is designed to work seamlessly in API services where you want to validate request payloads and provide schema discovery. The recommended pattern is:

1. **POST/PUT requests** - Validate the request body against your spec
2. **GET requests with schema discovery** - Return the validation spec when requested with a special header or query parameter

### Example Pattern

```go
// Example endpoint handler
func handleUserEndpoint(w http.ResponseWriter, r *http.Request) {
    // Schema discovery: GET with X-Mowgli header returns the spec
    if r.Method == "GET" && r.Header.Get("X-Mowgli") == "Validation Request" {
        spec, _ := LoadSpec("user_registration.json")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(spec)
        return
    }
    
    // Validation: POST/PUT validates the request body
    if r.Method == "POST" || r.Method == "PUT" {
        var data map[string]any
        json.NewDecoder(r.Body).Decode(&data)
        
        spec, _ := LoadSpec("user_registration.json")
        result := mowgli.Validate(data, spec)
        
        if !result.Valid {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(result.Errors)
            return
        }
        
        // Process valid data...
    }
}
```

This pattern allows API clients to:
- Discover validation requirements by making a GET request with the `X-Mowgli: Validation Request` header (or your preferred header/query parameter)
- Validate payloads on the server side before processing
- Use the same spec format for both server-side validation and client-side validation (with the JavaScript library)

The header name and format are flexible—you can use any header name or query parameter that suits your API design.

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
