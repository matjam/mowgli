import { Spec, ValidationResult, ValidationError } from './types';

export class ValidationErrorClass implements ValidationError {
  path: string;
  message: string;

  constructor(path: string, message: string) {
    this.path = path;
    this.message = message;
  }

  toString(): string {
    if (!this.path) {
      return this.message;
    }
    return `${this.path}: ${this.message}`;
  }
}

export class ValidationResultClass implements ValidationResult {
  valid: boolean;
  errors: ValidationError[];

  constructor() {
    this.valid = true;
    this.errors = [];
  }

  addError(path: string, message: string): void {
    this.valid = false;
    this.errors.push(new ValidationErrorClass(path, message));
  }

  validate(value: any, spec: Spec, path: string = ''): void {
    if (!spec) {
      this.addError(path, 'spec is nil');
      return;
    }

    // If no type specified, skip type validation (used in condition overrides)
    if (!spec.type) {
      return;
    }

    // Handle null values
    if (value === null || value === undefined) {
      if (spec.type !== 'null') {
        this.addError(path, `expected type ${spec.type}, got null`);
      }
      return;
    }

    switch (spec.type) {
      case 'string':
        this.validateString(path, value, spec);
        break;
      case 'number':
        this.validateNumber(path, value, spec);
        break;
      case 'integer':
        this.validateInteger(path, value, spec);
        break;
      case 'boolean':
        this.validateBoolean(path, value, spec);
        break;
      case 'object':
        this.validateObject(path, value, spec);
        break;
      case 'array':
        this.validateArray(path, value, spec);
        break;
      case 'null':
        // value is guaranteed to be non-null at this point (checked above)
        this.addError(path, 'expected null, got non-null value');
        break;
      default:
        this.addError(path, `unknown type: ${spec.type}`);
    }

    // Validate enum constraint if specified
    if (spec.enum && spec.enum.length > 0) {
      this.validateEnum(path, value, spec.enum);
    }
  }

  private validateString(path: string, value: any, spec: Spec): void {
    if (typeof value !== 'string') {
      this.addError(path, `expected string, got ${typeof value}`);
      return;
    }

    const str = value as string;

    // Handle allowEmpty - if true and string is empty, skip other validations
    if (spec.allowEmpty === true && str === '') {
      return;
    }

    // If allowEmpty is false or not set, empty strings must pass minLength check
    if (spec.minLength !== undefined && str.length < spec.minLength) {
      this.addError(
        path,
        `string length ${str.length} is less than minimum ${spec.minLength}`
      );
    }

    if (spec.maxLength !== undefined && str.length > spec.maxLength) {
      this.addError(
        path,
        `string length ${str.length} is greater than maximum ${spec.maxLength}`
      );
    }

    if (spec.pattern) {
      const regex = new RegExp(spec.pattern);
      if (!regex.test(str)) {
        this.addError(path, `string does not match pattern: ${spec.pattern}`);
      }
    }
  }

  private validateNumber(path: string, value: any, spec: Spec): void {
    const num = typeof value === 'number' ? value : parseFloat(value);
    if (isNaN(num)) {
      this.addError(path, `expected number, got ${typeof value}`);
      return;
    }

    if (spec.min !== undefined && num < spec.min) {
      this.addError(path, `number ${num} is less than minimum ${spec.min}`);
    }

    if (spec.max !== undefined && num > spec.max) {
      this.addError(path, `number ${num} is greater than maximum ${spec.max}`);
    }
  }

  private validateInteger(path: string, value: any, spec: Spec): void {
    const num = typeof value === 'number' ? value : parseFloat(value);
    if (isNaN(num)) {
      this.addError(path, `expected integer, got ${typeof value}`);
      return;
    }

    // Check if it's actually an integer (no fractional part)
    if (num !== Math.floor(num)) {
      this.addError(path, `expected integer, got float: ${num}`);
      return;
    }

    if (spec.min !== undefined && num < spec.min) {
      this.addError(path, `integer ${num} is less than minimum ${spec.min}`);
    }

    if (spec.max !== undefined && num > spec.max) {
      this.addError(path, `integer ${num} is greater than maximum ${spec.max}`);
    }
  }

  private validateBoolean(path: string, value: any, _spec: Spec): void {
    if (typeof value !== 'boolean') {
      this.addError(path, `expected boolean, got ${typeof value}`);
    }
  }

  private validateObject(path: string, value: any, spec: Spec): void {
    if (typeof value !== 'object' || Array.isArray(value) || value === null) {
      this.addError(path, `expected object, got ${typeof value}`);
      return;
    }

    const obj = value as Record<string, any>;

    // Check required fields
    if (spec.required) {
      for (const req of spec.required) {
        if (!(req in obj)) {
          this.addError(this.buildPath(path, req), 'required field is missing');
        }
      }
    }

    // Validate properties with conditional overrides
    if (spec.properties) {
      // Build effective specs for each property based on conditions
      const effectiveSpecs = this.buildEffectiveSpecs(obj, spec);

      for (const [key, propSpec] of Object.entries(spec.properties)) {
        const propValue = obj[key];
        const propExists = key in obj;

        // Get effective spec (with condition overrides applied)
        let effectiveSpec = propSpec;
        if (effectiveSpecs[key]) {
          effectiveSpec = effectiveSpecs[key];
        }

        if (propExists) {
          this.validate(propValue, effectiveSpec, this.buildPath(path, key));
        }
        // Field doesn't exist - check if it's required after conditions
        // This is handled by required check above, but conditions might add requirements
      }

      // Check for extra properties (warn but don't fail by default)
      for (const key of Object.keys(obj)) {
        if (spec.properties && !(key in spec.properties)) {
          // Extra property - silently allowed for now
        }
      }
    }
  }

  private validateArray(path: string, value: any, spec: Spec): void {
    if (!Array.isArray(value)) {
      this.addError(path, `expected array, got ${typeof value}`);
      return;
    }

    const arr = value as any[];

    if (spec.minLength !== undefined && arr.length < spec.minLength) {
      this.addError(
        path,
        `array length ${arr.length} is less than minimum ${spec.minLength}`
      );
    }

    if (spec.maxLength !== undefined && arr.length > spec.maxLength) {
      this.addError(
        path,
        `array length ${arr.length} is greater than maximum ${spec.maxLength}`
      );
    }

    if (spec.items) {
      for (let i = 0; i < arr.length; i++) {
        this.validate(arr[i], spec.items, this.buildArrayPath(path, i));
      }
    }
  }

  private validateEnum(path: string, value: any, enumValues: any[]): void {
    for (const allowed of enumValues) {
      if (deepEqual(value, allowed)) {
        return;
      }
    }

    // Format enum values for error message
    const enumStrs = enumValues.map((v) => String(v));
    this.addError(
      path,
      `value not in enum: ${value} (allowed: ${enumStrs.join(', ')})`
    );
  }

  private buildPath(base: string, field: string): string {
    if (!base) {
      return field;
    }
    if (!field) {
      return base;
    }
    return `${base}.${field}`;
  }

  private buildArrayPath(base: string, index: number): string {
    return `${base}[${index}]`;
  }

  private buildEffectiveSpecs(
    obj: Record<string, any>,
    spec: Spec
  ): Record<string, Spec> {
    const effectiveSpecs: Record<string, Spec> = {};

    if (!spec.conditions || spec.conditions.length === 0) {
      return effectiveSpecs;
    }

    // Client-side validation cannot evaluate expressions
    // This method will only be called for basic validation without expressions
    // When expressions are present, server-side validation should be used
    // We return empty effectiveSpecs to skip condition processing client-side
    return effectiveSpecs;
  }

  private mergeSpecs(base: Spec, override: Spec): Spec {
    if (!base) {
      return override;
    }
    if (!override) {
      return base;
    }

    const merged: Spec = {
      type: base.type,
      properties: base.properties,
      items: base.items,
      required: base.required,
      conditions: base.conditions,
      min: base.min,
      max: base.max,
      minLength: base.minLength,
      maxLength: base.maxLength,
      pattern: base.pattern,
      enum: base.enum,
      allowEmpty: base.allowEmpty,
    };

    // Apply overrides
    if (override.min !== undefined) {
      merged.min = override.min;
    }
    if (override.max !== undefined) {
      merged.max = override.max;
    }
    if (override.minLength !== undefined) {
      merged.minLength = override.minLength;
    }
    if (override.maxLength !== undefined) {
      merged.maxLength = override.maxLength;
    }
    if (override.pattern !== undefined) {
      merged.pattern = override.pattern;
    }
    if (override.enum !== undefined) {
      merged.enum = override.enum;
    }
    if (override.allowEmpty !== undefined) {
      merged.allowEmpty = override.allowEmpty;
    }
    if (override.type) {
      merged.type = override.type;
    }

    return merged;
  }
}

function deepEqual(a: any, b: any): boolean {
  if (a === null && b === null) {
    return true;
  }
  if (a === null || b === null) {
    return false;
  }
  if (a === b) {
    return true;
  }
  if (typeof a !== typeof b) {
    return false;
  }
  if (typeof a === 'object') {
    const keysA = Object.keys(a);
    const keysB = Object.keys(b);
    if (keysA.length !== keysB.length) {
      return false;
    }
    for (const key of keysA) {
      if (!keysB.includes(key) || !deepEqual(a[key], b[key])) {
        return false;
      }
    }
    return true;
  }
  return false;
}

/**
 * Checks if a spec contains expressions (conditions) that require server-side evaluation
 */
function hasExpressions(spec: Spec): boolean {
  if (!spec) {
    return false;
  }

  // Check if this spec has conditions
  if (spec.conditions && spec.conditions.length > 0) {
    return true;
  }

  // Recursively check nested specs
  if (spec.properties) {
    for (const propSpec of Object.values(spec.properties)) {
      if (hasExpressions(propSpec)) {
        return true;
      }
    }
  }

  if (spec.items && hasExpressions(spec.items)) {
    return true;
  }

  return false;
}

/**
 * Validates data using server-side validation endpoint
 * This is used when expressions are detected in the spec
 */
export async function validateWithServer(
  value: any,
  spec: Spec,
  endpoint: string,
  options?: {
    headerName?: string;
    headerValue?: string;
    fetchOptions?: RequestInit;
  }
): Promise<ValidationResult> {
  const {
    headerName = 'X-Mowgli',
    headerValue = 'validate',
    fetchOptions = {},
  } = options || {};

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(fetchOptions.headers as Record<string, string>),
    [headerName]: headerValue,
  };

  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      ...fetchOptions,
      headers,
      body: JSON.stringify(value),
    });

    if (!response.ok) {
      // If server returns validation errors, parse them
      if (response.status === 400) {
        try {
          const errors = await response.json();
          const result = new ValidationResultClass();
          if (Array.isArray(errors)) {
            for (const err of errors) {
              result.addError(err.path || '', err.message || String(err));
            }
          }
          return result;
        } catch {
          // Fall through to error handling
        }
      }

      throw new Error(`Server validation failed: ${response.statusText}`);
    }

    // If validation passes, server should return 200 with no body or a success response
    const result = new ValidationResultClass();
    return result;
  } catch (err) {
    const result = new ValidationResultClass();
    result.addError('', `Server validation error: ${err}`);
    return result;
  }
}

/**
 * Validates a value against a spec
 * If the spec contains expressions (conditions), this will attempt to use
 * server-side validation if an endpoint is provided via options.
 * Otherwise, it performs client-side validation (basic constraints only).
 */
export function validate(
  value: any,
  spec: Spec,
  options?: {
    /** Endpoint URL for server-side validation when expressions are present */
    endpoint?: string;
    /** Custom header name for validation request (default: 'X-Mowgli') */
    headerName?: string;
    /** Custom header value for validation request (default: 'validate') */
    headerValue?: string;
    /** Additional fetch options */
    fetchOptions?: RequestInit;
  }
): ValidationResult | Promise<ValidationResult> {
  // Check if spec has expressions that require server-side evaluation
  if (hasExpressions(spec)) {
    const endpoint = options?.endpoint;
    if (endpoint) {
      // Use server-side validation
      return validateWithServer(value, spec, endpoint, options);
    } else {
      // No endpoint provided, do basic validation only (expressions will be skipped)
      const result = new ValidationResultClass();
      result.validate(value, spec, '');
      return result;
    }
  }

  // No expressions, do full client-side validation
  const result = new ValidationResultClass();
  result.validate(value, spec, '');
  return result;
}

/**
 * Validates JSON data against a spec
 * See validate() for details on expression handling
 */
export function validateJSON(
  jsonData: string,
  spec: Spec,
  options?: {
    endpoint?: string;
    headerName?: string;
    headerValue?: string;
    fetchOptions?: RequestInit;
  }
): ValidationResult | Promise<ValidationResult> {
  try {
    const data = JSON.parse(jsonData);
    return validate(data, spec, options);
  } catch (err) {
    const result = new ValidationResultClass();
    result.addError('', `invalid JSON: ${err}`);
    return result;
  }
}

