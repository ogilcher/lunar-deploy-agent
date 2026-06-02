package deploy

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	StepName         string
	Command          string
	Arguments        []string
	DirectoryPath    string
	TimeoutSeconds   int
	WorkingDirectory string
}

// Name returns the configured deployment step name.
func (s *ShellCommandStep) Name() string {
	return s.StepName
}

// Run executes the configured command in the configured working directory.
func (s *ShellCommandStep) Run(
	deploymentContext DeploymentContext,
) error {
	logger.Log.Infow(
		"running shell command step.",
		"step", s.StepName,
		"command", s.Command,
		"arguments", s.Arguments,
		"directory_path", s.DirectoryPath,
	)

	timeout := time.Duration(s.TimeoutSeconds) * time.Second

	if s.TimeoutSeconds <= 0 {
		timeout = 5 * time.Minute
	}

	commandContext, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	defer cancel()

	command := exec.CommandContext(
		commandContext,
		s.Command,
		s.Arguments...,
	)

	commandDirectory := s.DirectoryPath

	if s.WorkingDirectory != "" {
		commandDirectory = filepath.Join(
			s.DirectoryPath,
			s.WorkingDirectory,
		)
	}

	command.Dir = commandDirectory

	command.Env = append(
		os.Environ(),
		BuildEnvironmentVariables(deploymentContext)...,
	)

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"Shell command output.",
		"step", s.StepName,
		"output", string(output),
	)

	return err
}
