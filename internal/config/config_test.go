package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := `
llm_engines:
  chatgpt:
    api_key: test_key
    model: gpt-3.5-turbo
`
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// defer os.Remove(tmpfile.Name())
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Error(err)
		}
	}()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LLMEngines["chatgpt"].APIKey != "test_key" {
		t.Errorf("Expected API key 'test_key', got '%s'", cfg.LLMEngines["chatgpt"].APIKey)
	}

	if cfg.LLMEngines["chatgpt"].Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", cfg.LLMEngines["chatgpt"].Model)
	}
}
