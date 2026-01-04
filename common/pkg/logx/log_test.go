// log_test.go
package logx

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(INFO, os.Stdout)
	Info("test message")
	if logger == nil {
		t.Error("NewLogger() should not return nil")
	}
	if logger.level != INFO {
		t.Errorf("Expected level to be INFO, got %v", logger.level)
	}
}

func TestSetLevel(t *testing.T) {
	oldLevel := defaultLogger.level
	defer SetLevel(oldLevel) // Restore original level after test

	SetLevel(DEBUG)
	if defaultLogger.level != DEBUG {
		t.Errorf("SetLevel() failed to set level. Expected %v, got %v", DEBUG, defaultLogger.level)
	}
}

func TestBasicLogFunctions(t *testing.T) {
	// Capture log output
	buf := &bytes.Buffer{}
	oldOutput := defaultLogger.logger
	oldLevel := defaultLogger.level
	defer func() {
		defaultLogger.logger = oldOutput
		SetLevel(oldLevel)
	}()

	// Create new logger with buffer output
	gkLogger := NewLogger(DEBUG, buf)
	defaultLogger = gkLogger

	Debug("debug message")
	Debugf("debug %s", "formatted")
	Info("info message")
	Infof("info %s", "formatted")
	Warn("warn message")
	Warnf("warn %s", "formatted")
	Error("error message")
	Errorf("error %s", "formatted")

	// Check that all messages are present in buffer
	output := buf.String()
	requiredMessages := []string{
		"debug message",
		"debug formatted",
		"info message",
		"info formatted",
		"warn message",
		"warn formatted",
		"error message",
		"error formatted",
	}

	for _, msg := range requiredMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected output to contain '%s', but it didn't. Output: %s", msg, output)
		}
	}
}

func TestLogLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	oldLevel := defaultLogger.level
	defer SetLevel(oldLevel)

	// Create new logger with buffer output
	gkLogger := NewLogger(WARN, buf)
	defaultLogger = gkLogger

	// These should NOT appear in output
	Debug("debug message")
	Info("info message")

	// These SHOULD appear in output
	Warn("warn message")
	Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") || strings.Contains(output, "info message") {
		t.Error("Lower level logs should have been filtered out")
	}
	if !strings.Contains(output, "warn message") || !strings.Contains(output, "error message") {
		t.Error("Higher or equal level logs should appear in output")
	}
}

func TestWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	oldLevel := defaultLogger.level
	defer SetLevel(oldLevel)

	// Create new logger with buffer output
	gkLogger := NewLogger(DEBUG, buf)
	defaultLogger = gkLogger

	fl := WithFields(map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	})

	fl.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected output to contain 'test message'")
	}
	// Note: go-kit/log formats fields differently than the previous implementation
	// Fields will be in key=value format but exact positioning may vary
	if !strings.Contains(output, "key1") || !strings.Contains(output, "value1") {
		t.Error("Expected output to contain field 'key1=value1'")
	}
	if !strings.Contains(output, "key2") || !strings.Contains(output, "123") {
		t.Error("Expected output to contain field 'key2=123'")
	}
}

func TestFatalFunctions(t *testing.T) {
	// Note: Testing fatal functions would normally exit the program
	// This test is just to ensure the functions exist and can be called
	// In a real scenario, you might use mock or test in a subprocess
	_ = Fatal
	_ = Fatalf
}
