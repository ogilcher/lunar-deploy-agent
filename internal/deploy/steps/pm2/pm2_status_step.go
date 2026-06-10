package pm2

import (
	"fmt"
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// PM2StatusStep checks whether a PM2 process exists and is visible to PM2.
type PM2StatusStep struct {
	ProcessName string
}

func (s *PM2StatusStep) Name() string {
	return "pm2_status"
}

func (s *PM2StatusStep) Run(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Checking PM2 process status...",
		"deployment", context.DeploymentName,
		"process_name", s.ProcessName,
	)

	command := exec.Command(
		"pm2",
		"describe",
		s.ProcessName,
	)

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"PM2 status output.",
		"process_name", s.ProcessName,
		"output", string(output),
	)

	if err != nil {
		return fmt.Errorf(
			"PM2 process %q was not found or is not healthy",
			s.ProcessName,
		)
	}

	return nil
}
