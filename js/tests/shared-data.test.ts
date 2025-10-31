import { loadSpec, loadTestCases, runTestCases } from './testdata-loader';

describe('Shared Test Data Validation', () => {
  describe('user_registration', () => {
    test('validates against shared spec and test cases', () => {
      const spec = loadSpec('user_registration.json');
      const testCases = loadTestCases('user_registration.json');
      const results = runTestCases(spec, testCases);

      expect(results.failed).toBe(0);
      expect(results.passed).toBe(testCases.length);
    });
  });

  describe('conditional_validation', () => {
    test('validates against shared spec and test cases', () => {
      const spec = loadSpec('conditional_validation.json');
      const testCases = loadTestCases('conditional_validation.json');
      const results = runTestCases(spec, testCases);

      expect(results.failed).toBe(0);
      expect(results.passed).toBe(testCases.length);
    });
  });

  describe('nested_objects', () => {
    test('validates against shared spec and test cases', () => {
      const spec = loadSpec('nested_objects.json');
      const testCases = loadTestCases('nested_objects.json');
      const results = runTestCases(spec, testCases);

      expect(results.failed).toBe(0);
      expect(results.passed).toBe(testCases.length);
    });
  });

  describe('advanced_conditional', () => {
    test('validates against shared spec and test cases', () => {
      const spec = loadSpec('advanced_conditional.json');
      const testCases = loadTestCases('advanced_conditional.json');
      const results = runTestCases(spec, testCases);

      expect(results.failed).toBe(0);
      expect(results.passed).toBe(testCases.length);
    });
  });
});

