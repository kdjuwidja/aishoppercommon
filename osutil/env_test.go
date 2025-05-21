package osutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvInt(t *testing.T) {
	// Save original environment variables
	originalEnv := make(map[string]string)
	for _, key := range []string{"TEST_INT_VALID", "TEST_INT_INVALID", "TEST_INT_ZERO", "TEST_INT_NEGATIVE"} {
		if val, exists := os.LookupEnv(key); exists {
			originalEnv[key] = val
		}
	}

	// Test cases
	testCases := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "Valid positive integer",
			envKey:       "TEST_INT_VALID",
			envValue:     "42",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Invalid integer (non-numeric)",
			envKey:       "TEST_INT_INVALID",
			envValue:     "not-a-number",
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "Zero value (should return default)",
			envKey:       "TEST_INT_ZERO",
			envValue:     "0",
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "Negative value (should return default)",
			envKey:       "TEST_INT_NEGATIVE",
			envValue:     "-5",
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "Empty environment variable",
			envKey:       "TEST_INT_EMPTY",
			envValue:     "",
			defaultValue: 100,
			expected:     100,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
			} else {
				os.Unsetenv(tc.envKey)
			}

			// Test
			result := GetEnvInt(tc.envKey, tc.defaultValue)
			assert.Equal(t, tc.expected, result)
		})
	}

	// Restore original environment variables
	for key, value := range originalEnv {
		os.Setenv(key, value)
	}
	for key := range originalEnv {
		if _, exists := os.LookupEnv(key); !exists {
			os.Unsetenv(key)
		}
	}
}

func TestGetEnvString(t *testing.T) {
	// Save original environment variables
	originalEnv := make(map[string]string)
	for _, key := range []string{"TEST_STR_VALID", "TEST_STR_EMPTY"} {
		if val, exists := os.LookupEnv(key); exists {
			originalEnv[key] = val
		}
	}

	// Test cases
	testCases := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "Valid string",
			envKey:       "TEST_STR_VALID",
			envValue:     "test-value",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "Empty string (should return default)",
			envKey:       "TEST_STR_EMPTY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Non-existent environment variable",
			envKey:       "TEST_STR_NONEXISTENT",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
			} else {
				os.Unsetenv(tc.envKey)
			}

			// Test
			result := GetEnvString(tc.envKey, tc.defaultValue)
			assert.Equal(t, tc.expected, result)
		})
	}

	// Restore original environment variables
	for key, value := range originalEnv {
		os.Setenv(key, value)
	}
	for key := range originalEnv {
		if _, exists := os.LookupEnv(key); !exists {
			os.Unsetenv(key)
		}
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "true lowercase",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "true uppercase",
			envValue:     "TRUE",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "true mixed case",
			envValue:     "TrUe",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "false value",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "empty value",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "invalid value",
			envValue:     "invalid",
			defaultValue: true,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("TEST_BOOL", tt.envValue)
				defer os.Unsetenv("TEST_BOOL")
			}

			// Test the function
			result := GetEnvBool("TEST_BOOL", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}
