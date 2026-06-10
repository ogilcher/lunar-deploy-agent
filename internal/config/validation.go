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
			// Native internal deployment steps.
			// No command validation required.
			case "git_pull":
			case "npm_install":
			case "npm_build":

			case "pm2_start_or_restart":
				if step.ProcessName == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].process_name is required",
						deploymentName,
						index,
					)
				}

				if step.StartCommand == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].start_command is required",
						deploymentName,
						index,
					)
				}

			case "pm2_restart":
				if step.ProcessName == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].process_name is required",
						deploymentName,
						index,
					)
				}

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
