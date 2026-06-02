package deploy

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// ShellCommandStep represents a configurable deployment step that runs
// a shell command inside a target project directory.
//
// This is used for commands like:
// - npm ci
// - npm run build
// - pm2 restart app
type ShellCommandStep struct {
	StepName      string
	Command       string
	Arguments     []string
	DirectoryPath string
}

// Name returns the configured deployment step name.
func (s *ShellCommandStep) Name() string {
	return s.StepName
}

// Run executes the configured command in the configured working directory.
func (s *ShellCommandStep) Run(
	context DeploymentContext,
) error {
	logger.Log.Infow(
		"running shell command step.",
		"step", s.StepName,
		"command", s.Command,
		"arguments", s.Arguments,
		"directory_path", s.DirectoryPath,
	)

	command := exec.Command(s.Command, s.Arguments...)
	command.Dir = s.DirectoryPath

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"Shell command output.",
		"step", s.StepName,
		"output", string(output),
	)

	return err
}
