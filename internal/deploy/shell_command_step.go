package deploy

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

type ShellCommandStep struct {
	StepName      string
	Command       string
	Arguments     []string
	DirectoryPath string
}

func (s *ShellCommandStep) Name() string {
	return s.StepName
}

func (s *ShellCommandStep) Run() error {
	command := exec.Command(s.Command, s.Arguments...)
	command.Dir = s.DirectoryPath

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"Shell command output",
		"step", s.StepName,
		"command", s.Command,
		"arguments", s.Arguments,
		"directory_path", s.DirectoryPath,
		"output", string(output),
	)

	return err
}
