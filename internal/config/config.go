package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the root application configuration for the agent.
type Config struct {
	AgentID       string                      `yaml:"agent_id"`
	APIToken      string                      `yaml:"api_token"`
	NodeID        string                      `yaml:"node_id"`
	NodeName      string                      `yaml:"node_name"`
	NodeRegion    string                      `yaml:"node_region"`
	Environment   string                      `yaml:"environment"`
	WorkspacePath string                      `yaml:"workspace_path"`
	LogLevel      string                      `yaml:"log_level"`
	Deployments   map[string]DeploymentConfig `yaml:"deployments"`
}

// DeploymentConfig defines the repository and step configuration used during deployment.
type DeploymentConfig struct {
	RepositoryPath            string             `yaml:"repository_path"`
	Preset                    string             `yaml:"preset"`
	ProcessName               string             `yaml:"process_name"`
	HealthCheckURL            string             `yaml:"health_check_url"`
	HealthCheckExpectedStatus int                `yaml:"health_check_expected_status"`
	Steps                     []DeployStepConfig `yaml:"steps"`
}

// DeployStepConfig represents a single shell-based deployment step loaded from YAML.
type DeployStepConfig struct {
	Type              string   `yaml:"type"`
	Name              string   `yaml:"name"`
	Command           string   `yaml:"command"`
	Arguments         []string `yaml:"arguments"`
	TimeoutSeconds    int      `yaml:"timeout_seconds"`
	WorkingDirectory  string   `yaml:"working_directory"`
	Retries           int      `yaml:"retries"`
	RetryDelaySeconds int      `yaml:"retry_delay_seconds"`
	ProcessName       string   `yaml:"process_name"`
	StartCommand      string   `yaml:"start_command"`
	URL               string   `yaml:"url"`
	ExpectedStatus    int      `yaml:"expected_status"`
}

// LoadConfig reads and parses an agent configuration file from disk.
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
