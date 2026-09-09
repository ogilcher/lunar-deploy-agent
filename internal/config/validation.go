package config

import "fmt"

// ValidateConfig verifies that the loaded configuration contains all required values.
func ValidateConfig(appConfig *Config) error {
	if appConfig.AgentID == "" {
		return fmt.Errorf("agent_id is required")
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

		if deployment.Preset == "node_pm2" {
			if deployment.ProcessName == "" {
				return fmt.Errorf(
					"deployments.%s.process_name is required for node_pm2 preset",
					deploymentName,
				)
			}

			if len(deployment.Steps) > 0 {
				return fmt.Errorf(
					"deployments.%s cannot use both preset and manual steps",
					deploymentName,
				)
			}

			continue
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
			case "pm2_save":

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

			case "pm2_status":
				if step.ProcessName == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].process_name is required",
						deploymentName,
						index,
					)
				}

			case "http_health_check":
				if step.URL == "" {
					return fmt.Errorf(
						"deployments.%s.steps[%d].url is required",
						deploymentName,
						index,
					)
				}

				if step.ExpectedStatus == 0 {
					return fmt.Errorf(
						"deployments.%s.steps[%d].expected_status is required",
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

		if deployment.Preset != "" && deployment.Preset != "node_pm2" {
			return fmt.Errorf(
				"deployments.%s.preset %q is unsupported",
				deploymentName,
				deployment.Preset,
			)
		}
	}

	return nil
}
