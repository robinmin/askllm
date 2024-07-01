package config

import (
	"os"
	"path/filepath"

	"github.com/robinmin/askllm/pkg/utils"
)

const (
	VERSION = "0.1.1"
)

type Config struct {
	Sys struct {
		LogPath  string `yaml:"log_path,omitempty"`
		LogLevel string `yaml:"log_level,omitempty"`
	} `yaml:"sys"`
	LLMEngines map[string]LLMEngineConfig `yaml:"llm_engines"`
}

type LLMEngineConfig struct {
	APIKey        string `yaml:"api_key"`
	Model         string `yaml:"model"`
	BaseURL       string `yaml:"base_url,omitempty"`
	OrgnizationId string `yaml:"organization_id,omitempty"` // So far, only avaliable for chatgpt
}

func Load(filename string) (*Config, error) {
	// Expand the tilde to the user's home directory
	absolutePath, err := expandTilde(filename)
	if err != nil {
		return nil, err
	}

	return utils.LoadConfig[Config](absolutePath)
}

func expandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil // Path doesn't start with '~', return as is
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err // Error getting home directory
	}

	return filepath.Join(homeDir, path[1:]), nil
}
