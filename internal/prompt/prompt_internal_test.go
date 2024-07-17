package prompt

import (
	// "os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/robinmin/askllm/pkg/utils"
)

func TestIsValidFilePath(t *testing.T) {
	t.Run("ValidFilePath", func(t *testing.T) {
		// Create a temporary file
		data := []byte("")
		tmpFile, err := utils.WriteTempFile("testfile", "", data)
		assert.NoError(t, err)

		// Defer cleanup to ensure it happens even if the function returns early
		defer func() {
			err := utils.CleanupTempFile(tmpFile)
			assert.NoError(t, err)
		}()

		if !isValidFilePath(tmpFile) {
			t.Errorf("Expected valid file path, got invalid")
		}
	})

	t.Run("InvalidFilePath", func(t *testing.T) {
		if isValidFilePath("/invalid/path/to/file") {
			t.Errorf("Expected invalid file path, got valid")
		}
	})
}

func TestIsQueryString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"ValidQueryString", "key1=value1&key2=value2", true},
		{"InvalidQueryString", "key1value1&key2value2", false},
		{"EmptyString", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := isQueryString(tt.input); result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseQueryString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
		hasError bool
	}{
		{
			"ValidQueryString",
			"key1=value1&key2=value2",
			map[string]any{"key1": "value1", "key2": "value2"},
			false,
		},
		{
			"InvalidQueryString",
			"key1value1&key2value2",
			nil,
			true,
		},
		{
			"EmptyString",
			"",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseQueryString(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("Expected error: %v, got: %v", tt.hasError, err)
			}
			if !tt.hasError && !equalMaps(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func equalMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
