package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentID       string           `yaml:"agent_id"`
	Environment   string           `yaml:"environment"`
	WorkspacePath string           `yaml:"workspace_path"`
	LogLevel      string           `yaml:"log_level"`
	Deployment    DeploymentConfig `yaml:"deployment"`
}

type DeploymentConfig struct {
	RepositoryPath string             `yaml:"repository_path"`
	Steps          []DeployStepConfig `yaml:"steps"`
}

type DeployStepConfig struct {
	Name      string   `yaml:"name"`
	Command   string   `yaml:"command"`
	Arguments []string `yaml:"arguments"`
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
