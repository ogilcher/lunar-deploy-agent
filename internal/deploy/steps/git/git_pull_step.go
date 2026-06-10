package git

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// GitPullStep updates a local repository by running git pull.
type GitPullStep struct {
	RepositoryPath string
}

// Name returns the deployment step identifier used in logs.
func (s *GitPullStep) Name() string {
	return "git_pull"
}

// Run executes git pull inside the configured repository path.
func (s *GitPullStep) Run(
	context context.DeploymentContext,
) error {
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
