package steps

import (
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// PM2RestartStep restarts a process managed by PM2.
type PM2RestartStep struct {
	ProcessName string
}

func (s *PM2RestartStep) Name() string {
	return "pm2_restart"
}

func (s *PM2RestartStep) Run(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Restarting PM2 process.",
		"deployment", context.DeploymentName,
		"process_name", s.ProcessName,
	)

	command := exec.Command(
		"pm2",
		"restart",
		s.ProcessName,
	)

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"PM2 restart output.",
		"process_name", s.ProcessName,
		"output", string(output),
	)

	return err
}
