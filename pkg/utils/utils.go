package utils

import (
	"os"

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
