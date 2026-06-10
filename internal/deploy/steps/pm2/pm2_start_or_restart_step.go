package pm2

import (
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// PM2StartOrRestartStep restarts an existing PM2 process or starts it
// if it does not exist yet.
type PM2StartOrRestartStep struct {
	ProcessName   string
	StartCommand  string
	Arguments     []string
	DirectoryPath string
}

func (s *PM2StartOrRestartStep) Name() string {
	return "pm2_start_or_restart"
}

func (s *PM2StartOrRestartStep) Run(
	context deploycontext.DeploymentContext,
) error {
	if s.processExists() {
		return s.restartProcess(context)
	}

	return s.startProcess(context)
}

func (s *PM2StartOrRestartStep) processExists() bool {
	command := exec.Command(
		"pm2",
		"describe",
		s.ProcessName,
	)

	return command.Run() == nil
}

func (s *PM2StartOrRestartStep) restartProcess(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Restarting existing PM2 process...",
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

func (s *PM2StartOrRestartStep) startProcess(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Starting new Pm2 process...",
		"deployment", context.DeploymentName,
		"process_name", s.ProcessName,
		"command", s.StartCommand,
		"arguments", s.Arguments,
	)

	args := append(
		[]string{
			"start",
			s.StartCommand,
			"--name",
			s.ProcessName,
			"--",
		},
		s.Arguments...,
	)

	command := exec.Command(
		"pm2",
		args...,
	)

	command.Dir = s.DirectoryPath

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"PM2 start output.",
		"process_name", s.ProcessName,
		"output", string(output),
	)

	return err
}
