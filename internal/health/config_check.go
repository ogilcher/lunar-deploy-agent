package health

import "github.com/ogilcher/lunar-deploy-agent/internal/config"

type ConfigCheck struct {
	ConfigPath string
}

func (c *ConfigCheck) Name() string {
	return "config"
}

func (c *ConfigCheck) Run() CheckResult {
	appConfig, err := config.LoadConfig(
		c.ConfigPath,
	)

	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: err.Error(),
		}
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: err.Error(),
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "configuration is valid",
	}
}
