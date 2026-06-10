package npm

import (
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// NPMBuildStep builds a Node.js project using npm run build.
type NPMBuildStep struct {
	DirectoryPath string
}

func (s *NPMBuildStep) Name() string {
	return "npm_build"
}

func (s *NPMBuildStep) Run(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Building npm project...",
		"deployment", context.DeploymentName,
	)

	command := exec.Command(
		"npm",
		"run",
		"build",
	)

	command.Dir = s.DirectoryPath

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"NPM build output.",
		"output", string(output),
	)

	return err
}
