/**
 * Condition defines a conditional validation rule
 */
export interface Condition {
  /** Expression to evaluate, e.g., "enabled == true", "count > 0" */
  if: string;
  /** Spec overrides to apply when condition is true */
  then: Record<string, Spec>;
  /** Spec overrides to apply when condition is false (optional) */
  else?: Record<string, Spec>;
}

/**
 * Spec defines the validation specification structure
 */
export interface Spec {
  /** Type: string, number, integer, boolean, object, array, null */
  type?: string;
  /** For object type - map of property names to their specifications */
  properties?: Record<string, Spec>;
  /** For array type - specification for array items */
  items?: Spec;
  /** For object type - list of required property names */
  required?: string[];
  /** Conditional validation rules for object type */
  conditions?: Condition[];

  // Constraints
  /** For number/integer - minimum value */
  min?: number;
  /** For number/integer - maximum value */
  max?: number;
  /** For string/array - minimum length */
  minLength?: number;
  /** For string/array - maximum length */
  maxLength?: number;
  /** For string - regex pattern */
  pattern?: string;
  /** Array of allowed values */
  enum?: any[];
  /** For strings - allows empty string if true */
  allowEmpty?: boolean;
}

/**
 * ValidationError represents a validation error with a path to the field
 */
export interface ValidationError {
  /** JSON path to the field with the error */
  path: string;
  /** Error message describing the validation failure */
  message: string;
}

/**
 * ValidationResult contains the result of validation
 */
export interface ValidationResult {
  /** Whether validation passed */
  valid: boolean;
  /** List of validation errors */
  errors: ValidationError[];
}

