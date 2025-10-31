export * from './types';
export * from './validator';
export * from './expression';

/**
 * Parses a JSON spec string into a Spec object
 */
export function parseSpec(specJSON: string): import('./types').Spec {
  return JSON.parse(specJSON);
}

// Framework-agnostic helpers
export {
  fetchSpec,
  validateField,
  groupErrorsByField,
  groupAllErrorsByField,
  createFormValidator,
} from './helpers';
