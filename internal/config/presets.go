package config

import "fmt"

// ExpandedDeploymentSteps returns the final deployment step list.
//
// If a deployment uses a preset, the preset is expanded into concrete steps.
// If no preset is provided, the manually configured steps are returned.
func ExpandedDeploymentSteps(
	deployment DeploymentConfig,
) ([]DeployStepConfig, error) {
	switch deployment.Preset {
	case "":
		return deployment.Steps, nil

	case "node_pm2":
		return expandedNodePM2Preset(deployment)

	default:
		return nil, fmt.Errorf(
			"unsupported deployment preset: %s",
			deployment.Preset,
		)
	}
}

func expandedNodePM2Preset(
	deployment DeploymentConfig,
) ([]DeployStepConfig, error) {
	expectedStatus := deployment.HealthCheckExpectedStatus

	if expectedStatus == 0 {
		expectedStatus = 200
	}

	steps := []DeployStepConfig{
		{Type: "git_pull"},
		{Type: "npm_install"},
		{Type: "npm_build"},
		{
			Type:         "pm2_start_or_restart",
			ProcessName:  deployment.ProcessName,
			StartCommand: "npm",
			Arguments:    []string{"start"},
		},
		{Type: "pm2_save"},
		{
			Type:        "pm2_status",
			ProcessName: deployment.ProcessName,
		},
	}

	if deployment.HealthCheckURL != "" {
		steps = append(
			steps,
			DeployStepConfig{
				Type:           "http_health_check",
				URL:            deployment.HealthCheckURL,
				ExpectedStatus: expectedStatus,
				TimeoutSeconds: 10,
			},
		)
	}

	return steps, nil
}
