package config

import "fmt"

func ValidateConfig(appConfig *Config) error {
	if appConfig.AgentID == "" {
		return fmt.Errorf("agend_id is required")
	}

	if appConfig.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if appConfig.Deployment.RepositoryPath == "" {
		return fmt.Errorf("deployment.repository_path is required")
	}

	for index, step := range appConfig.Deployment.Steps {
		if step.Name == "" {
			return fmt.Errorf("deployment.steps[%d].name is required", index)
		}

		if step.Command == "" {
			return fmt.Errorf("deployment.steps[%d].command is required", index)
		}
	}

	return nil
}
