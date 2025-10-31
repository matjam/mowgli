import { validate, validateJSON } from '../src/index';
import type { Spec } from '../src/types';

describe('validate', () => {
  describe('string validation', () => {
    test('valid string', () => {
      const spec: Spec = { type: 'string' };
      const result = validate('hello', spec);
      expect(result.valid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('string with minLength valid', () => {
      const spec: Spec = { type: 'string', minLength: 3 };
      const result = validate('hello', spec);
      expect(result.valid).toBe(true);
    });

    test('string with minLength invalid', () => {
      const spec: Spec = { type: 'string', minLength: 10 };
      const result = validate('hello', spec);
      expect(result.valid).toBe(false);
      expect(result.errors[0].message).toContain('less than minimum');
    });

    test('string with maxLength valid', () => {
      const spec: Spec = { type: 'string', maxLength: 10 };
      const result = validate('hello', spec);
      expect(result.valid).toBe(true);
    });

    test('string with maxLength invalid', () => {
      const spec: Spec = { type: 'string', maxLength: 3 };
      const result = validate('hello', spec);
      expect(result.valid).toBe(false);
    });

    test('string with pattern valid', () => {
      const spec: Spec = { type: 'string', pattern: '^[a-z]+$' };
      const result = validate('hello', spec);
      expect(result.valid).toBe(true);
    });

    test('string with pattern invalid', () => {
      const spec: Spec = { type: 'string', pattern: '^[a-z]+$' };
      const result = validate('Hello123', spec);
      expect(result.valid).toBe(false);
    });

    test('string with enum valid', () => {
      const spec: Spec = { type: 'string', enum: ['red', 'green', 'blue'] };
      const result = validate('red', spec);
      expect(result.valid).toBe(true);
    });

    test('string with enum invalid', () => {
      const spec: Spec = { type: 'string', enum: ['red', 'green', 'blue'] };
      const result = validate('yellow', spec);
      expect(result.valid).toBe(false);
    });

    test('wrong type', () => {
      const spec: Spec = { type: 'string' };
      const result = validate(123, spec);
      expect(result.valid).toBe(false);
    });
  });

  describe('number validation', () => {
    test('valid number', () => {
      const spec: Spec = { type: 'number' };
      const result = validate(3.14, spec);
      expect(result.valid).toBe(true);
    });

    test('number with min valid', () => {
      const spec: Spec = { type: 'number', min: 0 };
      const result = validate(5.5, spec);
      expect(result.valid).toBe(true);
    });

    test('number with min invalid', () => {
      const spec: Spec = { type: 'number', min: 10 };
      const result = validate(5.5, spec);
      expect(result.valid).toBe(false);
    });

    test('number with max valid', () => {
      const spec: Spec = { type: 'number', max: 100 };
      const result = validate(50.0, spec);
      expect(result.valid).toBe(true);
    });

    test('number with max invalid', () => {
      const spec: Spec = { type: 'number', max: 10 };
      const result = validate(50.0, spec);
      expect(result.valid).toBe(false);
    });
  });

  describe('integer validation', () => {
    test('valid integer', () => {
      const spec: Spec = { type: 'integer' };
      const result = validate(42, spec);
      expect(result.valid).toBe(true);
    });

    test('float not valid as integer', () => {
      const spec: Spec = { type: 'integer' };
      const result = validate(42.5, spec);
      expect(result.valid).toBe(false);
    });

    test('integer with min/max valid', () => {
      const spec: Spec = { type: 'integer', min: 0, max: 100 };
      const result = validate(50, spec);
      expect(result.valid).toBe(true);
    });

    test('integer with min/max invalid', () => {
      const spec: Spec = { type: 'integer', min: 0, max: 100 };
      const result = validate(150, spec);
      expect(result.valid).toBe(false);
    });
  });

  describe('boolean validation', () => {
    test('valid boolean true', () => {
      const spec: Spec = { type: 'boolean' };
      const result = validate(true, spec);
      expect(result.valid).toBe(true);
    });

    test('valid boolean false', () => {
      const spec: Spec = { type: 'boolean' };
      const result = validate(false, spec);
      expect(result.valid).toBe(true);
    });

    test('invalid boolean', () => {
      const spec: Spec = { type: 'boolean' };
      const result = validate('true', spec);
      expect(result.valid).toBe(false);
    });
  });

  describe('object validation', () => {
    test('simple object valid', () => {
      const spec: Spec = {
        type: 'object',
        properties: { name: { type: 'string' } },
      };
      const result = validate({ name: 'John' }, spec);
      expect(result.valid).toBe(true);
    });

    test('object with required field present', () => {
      const spec: Spec = {
        type: 'object',
        properties: { name: { type: 'string' } },
        required: ['name'],
      };
      const result = validate({ name: 'John' }, spec);
      expect(result.valid).toBe(true);
    });

    test('object with required field missing', () => {
      const spec: Spec = {
        type: 'object',
        properties: { name: { type: 'string' } },
        required: ['name'],
      };
      const result = validate({}, spec);
      expect(result.valid).toBe(false);
      expect(result.errors[0].message).toContain('required field is missing');
    });

    test('nested object valid', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          user: {
            type: 'object',
            properties: { name: { type: 'string' } },
          },
        },
      };
      const result = validate({ user: { name: 'John' } }, spec);
      expect(result.valid).toBe(true);
    });
  });

  describe('array validation', () => {
    test('simple array valid', () => {
      const spec: Spec = { type: 'array' };
      const result = validate([1, 2, 3], spec);
      expect(result.valid).toBe(true);
    });

    test('array with minLength valid', () => {
      const spec: Spec = { type: 'array', minLength: 2 };
      const result = validate([1, 2, 3], spec);
      expect(result.valid).toBe(true);
    });

    test('array with minLength invalid', () => {
      const spec: Spec = { type: 'array', minLength: 5 };
      const result = validate([1, 2, 3], spec);
      expect(result.valid).toBe(false);
    });

    test('array with maxLength valid', () => {
      const spec: Spec = { type: 'array', maxLength: 5 };
      const result = validate([1, 2, 3], spec);
      expect(result.valid).toBe(true);
    });

    test('array with maxLength invalid', () => {
      const spec: Spec = { type: 'array', maxLength: 2 };
      const result = validate([1, 2, 3], spec);
      expect(result.valid).toBe(false);
    });

    test('array with item spec valid', () => {
      const spec: Spec = {
        type: 'array',
        items: { type: 'string' },
      };
      const result = validate(['a', 'b', 'c'], spec);
      expect(result.valid).toBe(true);
    });

    test('array with item spec invalid', () => {
      const spec: Spec = {
        type: 'array',
        items: { type: 'string' },
      };
      const result = validate(['a', 123, 'c'], spec);
      expect(result.valid).toBe(false);
    });
  });

  describe('conditional validation', () => {
    test('conditional validation - enabled true, value present', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          enabled: { type: 'boolean' },
          value: { type: 'string' },
        },
        conditions: [
          {
            if: 'enabled == true',
            then: {
              value: { minLength: 1 },
            },
          },
          {
            if: 'enabled == false',
            then: {
              value: { allowEmpty: true },
            },
          },
        ],
      };
      const result = validate({ enabled: true, value: 'hello' }, spec);
      expect(result.valid).toBe(true);
    });

    test('conditional validation - enabled true, value missing', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          enabled: { type: 'boolean' },
          value: { type: 'string' },
        },
        conditions: [
          {
            if: 'enabled == true',
            then: {
              value: { minLength: 1 },
            },
          },
        ],
      };
      const result = validate({ enabled: true, value: '' }, spec);
      expect(result.valid).toBe(false);
    });

    test('conditional validation - enabled false, value empty allowed', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          enabled: { type: 'boolean' },
          value: { type: 'string' },
        },
        conditions: [
          {
            if: 'enabled == false',
            then: {
              value: { allowEmpty: true },
            },
          },
        ],
      };
      const result = validate({ enabled: false, value: '' }, spec);
      expect(result.valid).toBe(true);
    });

    test('conditional validation - numeric comparison', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          count: { type: 'integer' },
          message: { type: 'string' },
        },
        conditions: [
          {
            if: 'count > 0',
            then: {
              message: { minLength: 1 },
            },
          },
        ],
      };
      const result = validate({ count: 5, message: 'hello' }, spec);
      expect(result.valid).toBe(true);
    });

    test('conditional validation - string comparison', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          status: { type: 'string' },
          details: { type: 'string' },
        },
        conditions: [
          {
            if: 'status == "active"',
            then: {
              details: { minLength: 1 },
            },
          },
        ],
      };
      const result = validate({ status: 'active', details: 'some details' }, spec);
      expect(result.valid).toBe(true);
    });

    test('conditional validation - boolean field shorthand', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          enabled: { type: 'boolean' },
          value: { type: 'string' },
        },
        conditions: [
          {
            if: 'enabled',
            then: {
              value: { minLength: 1 },
            },
          },
        ],
      };
      const result = validate({ enabled: true, value: 'test' }, spec);
      expect(result.valid).toBe(true);
    });
  });

  describe('validateJSON', () => {
    test('valid JSON', () => {
      const spec: Spec = {
        type: 'object',
        properties: {
          name: { type: 'string', minLength: 1 },
          age: { type: 'integer', min: 0, max: 150 },
        },
        required: ['name', 'age'],
      };
      const result = validateJSON('{"name": "John", "age": 30}', spec);
      expect(result.valid).toBe(true);
    });

    test('invalid JSON', () => {
      const spec: Spec = { type: 'string' };
      const result = validateJSON('{invalid json}', spec);
      expect(result.valid).toBe(false);
      expect(result.errors[0].message).toContain('invalid JSON');
    });
  });
});

