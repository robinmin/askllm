package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

// create instance and load default values which defined in the struct definition
func NewInstance[T any]() *T {
	obj := new(T)
	if err := defaults.Set(obj); err != nil {
		return nil
	}
	return obj
}

// LoadConfig 从指定的YAML文件中加载配置信息
func LoadConfig[T any](yamlFile string) (*T, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 将配置信息保存到指定的YAML文件中
func SaveConfig[T any](cfg *T, yamlFile string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFile, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// WriteTempFile creates a temporary file with the given name prefix,
// writes the provided data to it, and returns the full file path.
// The caller is responsible for calling CleanupTempFile when done.
func WriteTempFile(fileNamePrefix string, data []byte) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", fileNamePrefix)
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %w", err)
	}
	// defer tempFile.Close()
	defer func() {
		_ = tempFile.Close()
	}()

	// Write data to the file
	if _, err := tempFile.Write(data); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %w", err)
	}

	// Get the full file path
	fullPath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error getting absolute file path: %w", err)
	}

	return fullPath, nil
}

// CleanupTempFile removes the specified temporary file.
// It's safe to call this function multiple times on the same file.
func CleanupTempFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing temporary file: %w", err)
	}
	return nil
}
