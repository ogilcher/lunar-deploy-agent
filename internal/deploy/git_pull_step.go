package deploy

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

type GitPullStep struct {
	RepositoryPath string
}

func (s *GitPullStep) Name() string {
	return "git_pull"
}

func (s *GitPullStep) Run() error {
	command := exec.Command(
		"git",
		"-C",
		s.RepositoryPath,
		"pull",
	)

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"Git pull output.",
		"output", string(output),
	)

	return err
}
