package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Helper function to capture log output
func captureLogOutput() *bytes.Buffer {
	var buf bytes.Buffer
	l.SetOutput(&buf)
	return &buf
}

// Helper function to parse log output
func parseLogOutput(buf *bytes.Buffer) (map[string]interface{}, error) {
	// Get the last line of output (in case there are multiple lines)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return nil, nil
	}
	lastLine := lines[len(lines)-1]

	var result map[string]interface{}
	err := json.Unmarshal([]byte(lastLine), &result)
	return result, err
}

func TestSetLevel(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()

	// Test valid level
	SetLevel("debug")
	assert.Equal(t, logrus.DebugLevel, l.GetLevel())

	// Test invalid level
	SetLevel("invalid")
	assert.Equal(t, logrus.InfoLevel, l.GetLevel())

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestGetLevel(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()

	// Test getting level
	l.SetLevel(logrus.DebugLevel)
	assert.Equal(t, logrus.DebugLevel, GetLevel())

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestSetServiceName(t *testing.T) {
	// Save original service name
	originalName := serviceName

	// Test setting service name
	SetServiceName("new-service")
	assert.Equal(t, "new-service", serviceName)

	// Restore original service name
	serviceName = originalName
}

func TestGetServiceName(t *testing.T) {
	// Save original service name
	originalName := serviceName

	// Test getting service name
	serviceName = "test-service"
	assert.Equal(t, "test-service", GetServiceName())

	// Restore original service name
	serviceName = originalName
}

func TestInfo(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()
	l.SetLevel(logrus.InfoLevel)

	// Setup
	buf := captureLogOutput()

	// Test
	Info("test message")

	// Verify
	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "test message", result["msg"])
	assert.Equal(t, "info", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Test Infof
	buf.Reset()
	Infof("test %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "test message", result["msg"])
	assert.Equal(t, "info", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestError(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()
	l.SetLevel(logrus.ErrorLevel)

	// Setup
	buf := captureLogOutput()

	// Test
	Error("error message")

	// Verify
	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "error message", result["msg"])
	assert.Equal(t, "error", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Test Errorf
	buf.Reset()
	Errorf("error %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "error message", result["msg"])
	assert.Equal(t, "error", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestDebug(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()
	l.SetLevel(logrus.DebugLevel)

	// Setup
	buf := captureLogOutput()

	// Test
	Debug("debug message")

	// Verify
	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "debug message", result["msg"])
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Test Debugf
	buf.Reset()
	Debugf("debug %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "debug message", result["msg"])
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestWarn(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()
	l.SetLevel(logrus.WarnLevel)

	// Setup
	buf := captureLogOutput()

	// Test
	Warn("warning message")

	// Verify
	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "warning message", result["msg"])
	assert.Equal(t, "warning", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Test Warnf
	buf.Reset()
	Warnf("warning %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "warning message", result["msg"])
	assert.Equal(t, "warning", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Restore original level
	l.SetLevel(originalLevel)
}

func TestTrace(t *testing.T) {
	// Save original level
	originalLevel := l.GetLevel()
	l.SetLevel(logrus.TraceLevel)

	// Setup
	buf := captureLogOutput()

	// Test
	Trace("trace message")

	// Verify
	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "trace message", result["msg"])
	assert.Equal(t, "trace", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Test Tracef
	buf.Reset()
	Tracef("trace %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "trace message", result["msg"])
	assert.Equal(t, "trace", result["level"])
	assert.Equal(t, serviceName, result["service"])

	// Restore original level
	l.SetLevel(originalLevel)
}

// Note: We can't fully test Fatal, Fatalf, Panic, and Panicf as they terminate the program
// But we can verify they exist and compile
func TestFatalExists(t *testing.T) {
	_ = Fatal
	_ = Fatalf
	_ = Panic
	_ = Panicf
}
