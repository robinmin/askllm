package log_test

import (
	"testing"

	// "github.com/dusted-go/logging/prettylog"
	testee "github.com/robinmin/askllm/pkg/utils/log"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name       string
		logpath    string
		appID      string
		level      string
		verbose    bool
		expect_nil bool
	}{
		{
			name:       "Create log file",
			logpath:    "logs",
			appID:      "test",
			level:      "DEBUG",
			verbose:    true,
			expect_nil: false,
		},
		{
			name:       "No log file created",
			logpath:    "",
			appID:      "test",
			level:      "INFO",
			verbose:    false,
			expect_nil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testee.InitLogger(tt.logpath, tt.appID, tt.level, tt.verbose)
			if (result == nil && !tt.expect_nil) || (result != nil && tt.expect_nil) {
				t.Errorf("Expected return nil: %v, got %v", tt.expect_nil, result)
			}
		})
	}
}

func TestLoggingFunctions(t *testing.T) {
	testee.InitLogger("logs", "test", "DEBUG", true)

	tests := []struct {
		name     string
		function func(string, ...interface{})
	}{
		{
			name:     "Debug",
			function: testee.Debug,
		},
		{
			name:     "Info",
			function: testee.Info,
		},
		{
			name:     "Warn",
			function: testee.Warn,
		},
		{
			name:     "Error",
			function: testee.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function("Test message")
			// Add assertions based on the logging output
		})
	}
}
