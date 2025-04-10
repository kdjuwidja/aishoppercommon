package os

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvInt(t *testing.T) {
	// Test with valid integer environment variable
	os.Setenv("TEST_INT", "42")
	assert.Equal(t, 42, GetEnvInt("TEST_INT", 0))

	// Test with invalid integer environment variable
	os.Setenv("TEST_INT", "not-a-number")
	assert.Equal(t, 0, GetEnvInt("TEST_INT", 0))

	// Test with non-positive integer environment variable
	os.Setenv("TEST_INT", "-1")
	assert.Equal(t, 0, GetEnvInt("TEST_INT", 0))

	// Test with missing environment variable
	os.Unsetenv("TEST_INT")
	assert.Equal(t, 0, GetEnvInt("TEST_INT", 0))

	// Test with custom default value
	assert.Equal(t, 100, GetEnvInt("NONEXISTENT_INT", 100))
}

func TestGetEnvString(t *testing.T) {
	// Test with valid string environment variable
	os.Setenv("TEST_STRING", "test-value")
	assert.Equal(t, "test-value", GetEnvString("TEST_STRING", "default"))

	// Test with empty string environment variable
	os.Setenv("TEST_STRING", "")
	assert.Equal(t, "default", GetEnvString("TEST_STRING", "default"))

	// Test with missing environment variable
	os.Unsetenv("TEST_STRING")
	assert.Equal(t, "default", GetEnvString("TEST_STRING", "default"))
}
