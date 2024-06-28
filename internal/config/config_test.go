package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := getFileContent()
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
func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, but got nil")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Error(err)
		}
	}()

	_, err = tmpfile.Write([]byte("invalid: yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Load(tmpfile.Name())
	if err != nil {
		t.Error("Expected error for invalid YAML, but got nil : ", err)
	}
}

func TestLoadWithTilde(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	content := getFileContent()
	tmpfile, err := os.CreateTemp(homeDir, "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
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

	cfg, err := Load("~/" + filepath.Base(tmpfile.Name()))
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

func getFileContent() string {
	return `
sys:
  log_path:
  log_level: INFO
llm_engines:
  chatgpt:
    api_key: test_key
    model: gpt-3.5-turbo
    # base_url: https://api.openai.com/v1
    # organization_id:
  gemini:
    api_key: 
    model: gemini-1.5-flash
  ollama:
    api_key: 
    model: gemma2
    # base_url: http://127.0.0.1:11434
`
}
