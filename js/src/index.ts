export * from './types';
export * from './validator';
export * from './expression';
export { validateWithServer } from './validator';

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
