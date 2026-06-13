package health

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
)

type CapabilitiesCheck struct {
	ConfigPath string
}

func (c *CapabilitiesCheck) Name() string {
	return "capabilities"
}

func (c *CapabilitiesCheck) Run() CheckResult {
	appConfig, err := config.LoadConfig(c.ConfigPath)
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: err.Error(),
		}
	}

	requiredCommands := map[string]bool{
		"git": false,
		"npm": false,
		"pm2": false,
	}

	for _, deployment := range appConfig.Deployments {
		steps, err := config.ExpandedDeploymentSteps(deployment)
		if err != nil {
			return CheckResult{
				Name:    c.Name(),
				Healthy: false,
				Message: err.Error(),
			}
		}

		for _, step := range steps {
			switch step.Type {
			case "git_pull":
				requiredCommands["git"] = true
			case "npm_install":
				requiredCommands["npm"] = true
			case "pm2_restart", "pm2_start_or_restart", "pm2_save", "pm2_staus":
				requiredCommands["pm2"] = true
			}
		}
	}

	missing := []string{}

	for command, required := range requiredCommands {
		if required {
			if _, err := exec.LookPath(command); err != nil {
				missing = append(missing, command)
			}
		}
	}

	if len(missing) > 0 {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: "missing required commands",
			Metadata: map[string]any{
				"missing": missing,
			},
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "required deployment capabilities are available",
	}
}
