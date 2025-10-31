import { evalExpression } from '../src/expression';

describe('evalExpression with AND/OR and parentheses', () => {
  describe('AND operator', () => {
    test('both true', () => {
      const obj = { enabled: true, active: true };
      expect(evalExpression('enabled AND active', obj)).toBe(true);
    });

    test('first false', () => {
      const obj = { enabled: false, active: true };
      expect(evalExpression('enabled AND active', obj)).toBe(false);
    });

    test('comparison AND comparison', () => {
      const obj = { count: 5, status: 'active' };
      expect(evalExpression('count > 0 AND status == "active"', obj)).toBe(true);
    });
  });

  describe('OR operator', () => {
    test('both true', () => {
      const obj = { enabled: true, active: true };
      expect(evalExpression('enabled OR active', obj)).toBe(true);
    });

    test('first true', () => {
      const obj = { enabled: true, active: false };
      expect(evalExpression('enabled OR active', obj)).toBe(true);
    });

    test('comparison OR comparison', () => {
      const obj = { count: 5, status: 'inactive' };
      expect(evalExpression('count > 10 OR status == "active"', obj)).toBe(false);
    });
  });

  describe('AND and OR together', () => {
    test('AND before OR', () => {
      const obj = { a: true, b: false, c: true };
      expect(evalExpression('a AND b OR c', obj)).toBe(true);
    });

    test('complex expression', () => {
      const obj = { count: 0, status: 'inactive', enabled: true };
      expect(evalExpression('count > 0 AND status == "active" OR enabled', obj)).toBe(true);
    });
  });

  describe('parentheses', () => {
    test('simple parentheses', () => {
      const obj = { enabled: true };
      expect(evalExpression('(enabled)', obj)).toBe(true);
    });

    test('parentheses change precedence', () => {
      const obj = { a: true, b: false, c: true };
      expect(evalExpression('(a OR b) AND c', obj)).toBe(true);
    });

    test('nested parentheses', () => {
      const obj = { a: false, b: true, c: true };
      expect(evalExpression('((a AND b) OR c)', obj)).toBe(true);
    });

    test('complex with parentheses', () => {
      const obj = { count: 0, status: 'inactive', enabled: true, verified: true };
      expect(
        evalExpression(
          '(count > 0 AND status == "active") OR (enabled AND verified)',
          obj
        )
      ).toBe(true);
    });

    test('multiple parentheses groups', () => {
      const obj = { a: 1, b: 3, c: 3 };
      expect(evalExpression('(a == 1) AND (b == 2 OR c == 3)', obj)).toBe(true);
    });
  });

  describe('complex expressions', () => {
    test('complex nested expression', () => {
      const obj = { a: 5, b: 15, c: 'inactive', enabled: true };
      expect(
        evalExpression('((a > 0 AND b < 10) OR (c == "active")) AND enabled', obj)
      ).toBe(false);
    });

    test('complex nested expression true', () => {
      const obj = { a: 5, b: 5, c: 'inactive', enabled: true };
      expect(
        evalExpression('((a > 0 AND b < 10) OR (c == "active")) AND enabled', obj)
      ).toBe(true);
    });
  });
});

