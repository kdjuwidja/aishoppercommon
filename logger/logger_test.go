package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a logger with a buffer output
func setupTestLogger(serviceName, level string) (*Logger, *bytes.Buffer) {
	logger := &Logger{}
	logger.Initialize(serviceName, level)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	return logger, &buf
}

// Helper function to parse log output
func parseLogOutput(buf *bytes.Buffer) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &result)
	return result, err
}

func TestInitialize(t *testing.T) {
	logger := &Logger{}

	// Test with valid parameters
	logger.Initialize("test-service", "debug")
	assert.Equal(t, "test-service", logger.GetServiceName())
	assert.Equal(t, logrus.DebugLevel, logger.GetLevel())

	// Test with empty level (should default to info)
	logger = &Logger{}
	logger.Initialize("test-service", "")
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())

	// Test with invalid level (should default to info)
	logger = &Logger{}
	logger.Initialize("test-service", "invalid-level")
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
}

func TestSetLevel(t *testing.T) {
	logger := &Logger{}
	logger.Initialize("test-service", "info")

	// Test setting level
	logger.SetLevel(logrus.DebugLevel)
	assert.Equal(t, logrus.DebugLevel, logger.GetLevel())
}

func TestInfo(t *testing.T) {
	logger, buf := setupTestLogger("test-service", "info")

	// Test Info
	logger.Info("test message")

	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "test message", result["msg"])
	assert.Equal(t, "info", result["level"])
	assert.Equal(t, "test-service", result["service"])

	// Test Infof
	buf.Reset()
	logger.Infof("test %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "test message", result["msg"])
	assert.Equal(t, "info", result["level"])
	assert.Equal(t, "test-service", result["service"])
}

func TestError(t *testing.T) {
	logger, buf := setupTestLogger("test-service", "info")

	// Test Error
	logger.Error("error message")

	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "error message", result["msg"])
	assert.Equal(t, "error", result["level"])
	assert.Equal(t, "test-service", result["service"])

	// Test Errorf
	buf.Reset()
	logger.Errorf("error %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "error message", result["msg"])
	assert.Equal(t, "error", result["level"])
	assert.Equal(t, "test-service", result["service"])
}

func TestDebug(t *testing.T) {
	logger, buf := setupTestLogger("test-service", "debug")

	// Test Debug
	logger.Debug("debug message")

	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "debug message", result["msg"])
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, "test-service", result["service"])

	// Test Debugf
	buf.Reset()
	logger.Debugf("debug %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "debug message", result["msg"])
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, "test-service", result["service"])
}

func TestWarn(t *testing.T) {
	logger, buf := setupTestLogger("test-service", "info")

	// Test Warn
	logger.Warn("warning message")

	result, err := parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "warning message", result["msg"])
	assert.Equal(t, "warning", result["level"])
	assert.Equal(t, "test-service", result["service"])

	// Test Warnf
	buf.Reset()
	logger.Warnf("warning %s", "message")

	result, err = parseLogOutput(buf)
	assert.NoError(t, err)
	assert.Equal(t, "warning message", result["msg"])
	assert.Equal(t, "warning", result["level"])
	assert.Equal(t, "test-service", result["service"])
}

func TestClose(t *testing.T) {
	logger := &Logger{}
	logger.Initialize("test-service", "info")

	// Test Close
	logger.Close()
	assert.Nil(t, logger.logger)
}
