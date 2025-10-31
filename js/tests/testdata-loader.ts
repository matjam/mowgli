import { readFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { parseSpec, validate, Spec } from '../src/index';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

export interface TestCase {
  name: string;
  data: Record<string, any>;
  expectedValid: boolean;
}

export interface TestCaseFile {
  testCases: TestCase[];
}

/**
 * Loads a spec JSON file from the testdata directory
 */
export function loadSpec(filename: string): Spec {
  // Try relative to js directory first, then parent directory
  let path = join(__dirname, '..', 'testdata', 'specs', filename);
  try {
    readFileSync(path, 'utf-8');
  } catch {
    path = join(__dirname, '..', '..', 'testdata', 'specs', filename);
  }
  const data = readFileSync(path, 'utf-8');
  return parseSpec(data);
}

/**
 * Loads test cases from a JSON file in the testdata directory
 */
export function loadTestCases(filename: string): TestCase[] {
  // Try relative to js directory first, then parent directory
  let path = join(__dirname, '..', 'testdata', 'cases', filename);
  try {
    readFileSync(path, 'utf-8');
  } catch {
    path = join(__dirname, '..', '..', 'testdata', 'cases', filename);
  }
  const data = readFileSync(path, 'utf-8');
  const testFile: TestCaseFile = JSON.parse(data);
  return testFile.testCases;
}

/**
 * Runs all test cases against a spec and returns results
 */
export function runTestCases(spec: Spec, testCases: TestCase[]): {
  passed: number;
  failed: number;
  failures: Array<{ name: string; expected: boolean; actual: boolean; errors?: any[] }>;
} {
  let passed = 0;
  let failed = 0;
  const failures: Array<{ name: string; expected: boolean; actual: boolean; errors?: any[] }> = [];

  for (const testCase of testCases) {
    const result = validate(testCase.data, spec);
    if (result.valid === testCase.expectedValid) {
      passed++;
    } else {
      failed++;
      failures.push({
        name: testCase.name,
        expected: testCase.expectedValid,
        actual: result.valid,
        errors: result.errors,
      });
    }
  }

  return { passed, failed, failures };
}

