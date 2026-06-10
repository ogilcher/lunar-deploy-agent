package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	core "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
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
	StepName          string
	Command           string
	Arguments         []string
	DirectoryPath     string
	TimeoutSeconds    int
	WorkingDirectory  string
	Retries           int
	RetryDelaySeconds int
}

// Name returns the configured deployment step name.
func (s *ShellCommandStep) Name() string {
	return s.StepName
}

// Run executes the configured command in the configured working directory.
func (s *ShellCommandStep) Run(
	deploymentContext core.DeploymentContext,
) error {
	attempts := s.Retries + 1

	for attempt := 1; attempt <= attempts; attempt++ {
		logger.Log.Infow(
			"Running shell command step.",
			"step", s.StepName,
			"attempt", attempt,
			"max_attempts", attempts,
			"command", s.Command,
		)

		timeout := time.Duration(s.TimeoutSeconds) * time.Second

		if s.TimeoutSeconds <= 0 {
			timeout = 5 * time.Minute
		}

		commandContext, cancel := context.WithTimeout(
			context.Background(),
			timeout,
		)

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
			core.BuildEnvironmentVariables(
				deploymentContext,
			)...,
		)

		output, err := command.CombinedOutput()

		cancel()

		logger.Log.Infow(
			"Shell command output.",
			"step", s.StepName,
			"attempt", attempt,
			"output", string(output),
		)

		if err == nil {
			return nil
		}

		logger.Log.Errorw(
			"Shell command attempt failed.",
			"step", s.StepName,
			"attempt", attempt,
			"error", err,
		)

		if attempt < attempts {
			delay := time.Duration(
				s.RetryDelaySeconds,
			) * time.Second

			if delay <= 0 {
				delay = 1 * time.Second
			}

			time.Sleep(delay)
		}
	}

	return fmt.Errorf(
		"shell command step failed after %d attempts",
		attempts,
	)
}
