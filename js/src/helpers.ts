import { validate, parseSpec, Spec, ValidationResult } from './index';

/**
 * Fetches a validation spec from an API endpoint
 * @param endpoint The API endpoint URL
 * @param options Optional configuration for the request
 * @returns Promise resolving to the parsed spec
 */
export async function fetchSpec(
  endpoint: string,
  options?: {
    headerName?: string;
    headerValue?: string;
    fetchOptions?: RequestInit;
  }
): Promise<Spec> {
  const {
    headerName = 'X-Mowgli',
    headerValue = 'Validation Request',
    fetchOptions = {},
  } = options || {};

  const headers: Record<string, string> = {
    ...(fetchOptions.headers as Record<string, string>),
    [headerName]: headerValue,
  };

  const response = await fetch(endpoint, {
    ...fetchOptions,
    headers,
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch spec: ${response.statusText}`);
  }

  const data = await response.json();
  return parseSpec(JSON.stringify(data));
}

/**
 * Validates a single field in a form
 * @param fieldName The name of the field to validate
 * @param fieldValue The current value of the field
 * @param allFormValues All current form values (for context in conditional validation)
 * @param spec The validation spec
 * @returns Object containing whether the field is valid and any error messages
 */
export function validateField(
  fieldName: string,
  fieldValue: any,
  allFormValues: Record<string, any>,
  spec: Spec,
  endpoint?: string
): {
  valid: boolean | Promise<boolean>;
  error: string | null | Promise<string | null>;
  errors: Array<{ path: string; message: string }> | Promise<Array<{ path: string; message: string }>>;
} {
  // Create a temporary object with the current field value in context
  const testData = { ...allFormValues, [fieldName]: fieldValue };

  const result = validate(testData, spec, endpoint ? { endpoint } : undefined);

  // Handle both sync and async results
  if (result instanceof Promise) {
    return {
      valid: result.then((r) => {
        const fieldErrors = r.errors.filter(
          (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
        );
        return fieldErrors.length === 0;
      }),
      error: result.then((r) => {
        const fieldErrors = r.errors.filter(
          (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
        );
        return fieldErrors[0]?.message || null;
      }),
      errors: result.then((r) => {
        const fieldErrors = r.errors.filter(
          (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
        );
        return fieldErrors.map((err) => ({
          path: err.path,
          message: err.message,
        }));
      }),
    };
  }

  // Client-side validation (sync)
  const fieldErrors = result.errors.filter(
    (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
  );

  return {
    valid: fieldErrors.length === 0,
    error: fieldErrors[0]?.message || null,
    errors: fieldErrors.map((err) => ({
      path: err.path,
      message: err.message,
    })),
  };
}

/**
 * Groups validation errors by field name
 * @param result The validation result
 * @returns Object mapping field names to their error messages
 */
export function groupErrorsByField(result: ValidationResult): Record<string, string> {
  const errorsByField: Record<string, string> = {};

  result.errors.forEach((err) => {
    const fieldName = err.path.split('.')[0]; // Get top-level field name
    if (!errorsByField[fieldName]) {
      errorsByField[fieldName] = err.message;
    }
  });

  return errorsByField;
}

/**
 * Groups validation errors by field name, preserving all errors per field
 * @param result The validation result
 * @returns Object mapping field names to arrays of error messages
 */
export function groupAllErrorsByField(
  result: ValidationResult
): Record<string, Array<{ path: string; message: string }>> {
  const errorsByField: Record<string, Array<{ path: string; message: string }>> = {};

  result.errors.forEach((err) => {
    const fieldName = err.path.split('.')[0];
    if (!errorsByField[fieldName]) {
      errorsByField[fieldName] = [];
    }
    errorsByField[fieldName].push({
      path: err.path,
      message: err.message,
    });
  });

  return errorsByField;
}

/**
 * Creates a form validator function that can be used with any framework
 * @param spec The validation spec
 * @param endpoint Optional endpoint URL for server-side validation when expressions are present
 * @returns An object with validation functions
 */
export function createFormValidator(
  spec: Spec,
  endpoint?: string
): {
  validateForm: (
    values: Record<string, any>
  ) => ValidationResult | Promise<ValidationResult>;
  validateField: (
    fieldName: string,
    fieldValue: any,
    allFormValues: Record<string, any>
  ) => ValidationResult | Promise<ValidationResult>;
  isValid: (
    values: Record<string, any>
  ) => boolean | Promise<boolean>;
} {
  return {
    /**
     * Validates all form values
     */
    validateForm: (values: Record<string, any>) => {
      return validate(values, spec, endpoint ? { endpoint } : undefined);
    },

    /**
     * Validates a single field
     */
    validateField: (
      fieldName: string,
      fieldValue: any,
      allFormValues: Record<string, any>
    ) => {
      // For field validation, we still need to validate the whole form context
      // but filter results for the specific field
      const testData = { ...allFormValues, [fieldName]: fieldValue };
      const result = validate(testData, spec, endpoint ? { endpoint } : undefined);
      
      // If result is a Promise (server-side validation), we need to handle it
      if (result instanceof Promise) {
        return result.then((r) => {
          const fieldErrors = r.errors.filter(
            (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
          );
          return {
            ...r,
            errors: fieldErrors,
            valid: fieldErrors.length === 0,
          } as ValidationResult;
        });
      }
      
      // Client-side validation
      const fieldErrors = result.errors.filter(
        (err) => err.path === fieldName || err.path.startsWith(`${fieldName}.`)
      );
      return {
        ...result,
        errors: fieldErrors,
        valid: fieldErrors.length === 0,
      } as ValidationResult;
    },

    /**
     * Checks if form is valid without returning detailed errors
     */
    isValid: (values: Record<string, any>) => {
      const result = validate(values, spec, endpoint ? { endpoint } : undefined);
      if (result instanceof Promise) {
        return result.then((r) => r.valid);
      }
      return result.valid;
    },
  };
}

