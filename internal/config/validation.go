package config

import "fmt"

// ValidateConfig verifies that the loaded configuration contains all required values.
func ValidateConfig(appConfig *Config) error {
	if appConfig.AgentID == "" {
		return fmt.Errorf("agend_id is required")
	}

	if appConfig.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	if len(appConfig.Deployments) == 0 {
		return fmt.Errorf("at least one deployment is required")
	}

	for deploymentName, deployment := range appConfig.Deployments {
		if deployment.RepositoryPath == "" {
			return fmt.Errorf(
				"deployments.%s.repository_path is required",
				deploymentName,
			)
		}

		for index, step := range deployment.Steps {
			if step.Type == "" {
				return fmt.Errorf(
					"deployment.%s.steps[%d].type is required",
					deploymentName,
					index,
				)
			}

			switch step.Type {
			case "git_pull":
				// Native internal deployment step.
				// No command validation required.

			case "shell":
				if step.Name == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].name is required",
						deploymentName,
						index,
					)
				}

				if step.Command == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].command is required",
						deploymentName,
						index,
					)
				}

			default:
				return fmt.Errorf(
					"deployments.%s.steps[%d].type %q is unsupported",
					deploymentName,
					index,
					step.Type,
				)
			}

		}
	}

	return nil
}
