package version

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	type tc struct {
		name        string
		input       string
		expected    Version
		expectError bool
	}

	tests := map[string]tc{
		"valid version": {
			input:    "1.2.3",
			expected: Version{Major: 1, Minor: 2, Patch: 3},
		},
		"version with v prefix": {
			input:    "v2.3.4",
			expected: Version{Major: 2, Minor: 3, Patch: 4},
		},
		"invalid format": {
			input:       "1.2",
			expectError: true,
		},
		"invalid major": {
			input:       "a.2.3",
			expectError: true,
		},
		"invalid minor": {
			input:       "1.b.3",
			expectError: true,
		},
		"invalid patch": {
			input:       "1.2.c",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := parseVersion(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := map[string]struct {
		version  Version
		expected string
	}{
		"standard version": {
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		"zero version": {
			version:  Version{Major: 0, Minor: 0, Patch: 0},
			expected: "0.0.0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.version.String())
		})
	}
}

func TestVersionIncrement(t *testing.T) {
	type tc struct {
		name        string
		initial     Version
		level       string
		expected    Version
		expectError bool
	}

	tests := map[string]tc{
		"increment major": {
			initial:  Version{Major: 1, Minor: 2, Patch: 3},
			level:    "major",
			expected: Version{Major: 2, Minor: 0, Patch: 0},
		},
		"increment minor": {
			initial:  Version{Major: 1, Minor: 2, Patch: 3},
			level:    "minor",
			expected: Version{Major: 1, Minor: 3, Patch: 0},
		},
		"increment patch": {
			initial:  Version{Major: 1, Minor: 2, Patch: 3},
			level:    "patch",
			expected: Version{Major: 1, Minor: 2, Patch: 4},
		},
		"invalid level": {
			initial:     Version{Major: 1, Minor: 2, Patch: 3},
			level:       "invalid",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := tc.initial
			err := v.Increment(tc.level)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, v)
			}
		})
	}
}

func TestReadWrite(t *testing.T) {
	// Create temp file path
	tmpFile, err := os.CreateTemp("", "version_test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Clean up after test
	defer os.Remove(tmpPath)

	tests := map[string]struct {
		initialVersion Version
		expectedError  bool
	}{
		"write and read version": {
			initialVersion: Version{Major: 1, Minor: 2, Patch: 3},
		},
		"read non-existent file": {
			expectedError: false, // Should return 0.0.1
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Remove any existing file
			os.Remove(tmpPath)

			if !tc.expectedError {
				if err := WriteToFile(tc.initialVersion, tmpPath); err != nil {
					t.Fatal(err)
				}

				readVersion, err := ReadFromFile(tmpPath)
				assert.NoError(t, err)
				if tc.initialVersion.Major == 0 && tc.initialVersion.Minor == 0 && tc.initialVersion.Patch == 0 {
					// Expect default version for non-existent file
					assert.Equal(t, Version{Major: 0, Minor: 0, Patch: 1}, readVersion)
				} else {
					assert.Equal(t, tc.initialVersion, readVersion)
				}
			}
		})
	}
}
