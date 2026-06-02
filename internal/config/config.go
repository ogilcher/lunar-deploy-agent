package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentID       string `yaml:"agent_id"`
	Environment   string `yaml:"environment"`
	WorkspacePath string `yaml:"workspace_path"`
	LogLevel      string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var appConfig Config

	if err := yaml.Unmarshal(configFile, &appConfig); err != nil {
		return nil, err
	}

	return &appConfig, nil
}
