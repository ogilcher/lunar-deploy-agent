package pm2

import (
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// PM2SaveStep saves the current PM2 process list for reboot persistence.
type PM2SaveStep struct{}

func (s *PM2SaveStep) Name() string {
	return "pm2_save"
}

func (s *PM2SaveStep) Run(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Saving PM2 process list...",
		"deployment", context.DeploymentName,
	)

	command := exec.Command("pm2", "save")

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"PM2 save output.",
		"output", string(output),
	)

	return err
}
