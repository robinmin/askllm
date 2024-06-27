package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LLMEngines map[string]LLMEngineConfig `yaml:"llm_engines"`
}

type LLMEngineConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
