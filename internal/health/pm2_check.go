package health

import "os/exec"

// PM2Check verifies that the pm2 command is installed and reachable.
type PM2Check struct{}

func (c *PM2Check) Name() string {
	return "pm2"
}

func (c *PM2Check) Run() CheckResult {
	command := exec.Command("pm2", "--version")

	output, err := command.CombinedOutput()

	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: "pm2 is not installed or not reachable",
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "pm2 available: " + string(output),
	}
}
