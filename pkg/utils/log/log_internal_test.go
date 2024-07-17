package log

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"log/slog"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{
			name:     "Valid level",
			input:    "DEBUG",
			expected: slog.LevelDebug,
		},
		{
			name:     "Invalid level",
			input:    "INVALID",
			expected: slog.Level(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, result)
			}
			if err != nil && tt.expected != slog.Level(0) {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCloseLogger(t *testing.T) {
	t.Run("Close log file", func(t *testing.T) {
		fileName := "/tmp/test_f852b3ce_b97a_4e4d_aa8c_34a7bdc4381e.log"
		file, err := os.Create(fileName)
		assert.Nil(t, err)
		assert.NotNil(t, file)
		defer func() {
			if file != nil {
				_ = file.Close()
			}
			_ = os.Remove(fileName)
		}()

		CloseLogger(file)
		// assert.Nil(t, file)
	})
}
