package mowgli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TestCase represents a single test case
type TestCase struct {
	Name          string         `json:"name"`
	Data          map[string]any `json:"data"`
	ExpectedValid bool           `json:"expectedValid"`
}

// TestCaseFile represents a file containing multiple test cases
type TestCaseFile struct {
	TestCases []TestCase `json:"testCases"`
}

// LoadSpec loads a spec JSON file from the testdata directory
func LoadSpec(filename string) (*Spec, error) {
	path := filepath.Join("testdata", "specs", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file %s: %w", path, err)
	}

	spec, err := ParseSpecString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec from %s: %w", path, err)
	}

	return spec, nil
}

// LoadTestCases loads test cases from a JSON file in the testdata directory
func LoadTestCases(filename string) ([]TestCase, error) {
	path := filepath.Join("testdata", "cases", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read test case file %s: %w", path, err)
	}

	var testFile TestCaseFile
	if err := json.Unmarshal(data, &testFile); err != nil {
		return nil, fmt.Errorf("failed to parse test case file %s: %w", path, err)
	}

	return testFile.TestCases, nil
}

