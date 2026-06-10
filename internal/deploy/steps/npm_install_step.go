package steps

import (
	"os/exec"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// NPMInstallStep Installs Node.js dependencies using npm ci
type NPMInstallStep struct {
	DirectoryPath string
}

func (s *NPMInstallStep) Name() string {
	return "npm_install"
}

func (s *NPMInstallStep) Run(
	context deploycontext.DeploymentContext,
) error {
	logger.Log.Infow(
		"Installing npm dependencies...",
		"deployment", context.DeploymentName,
	)

	command := exec.Command(
		"npm",
		"ci",
	)

	command.Dir = s.DirectoryPath

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"NPM install output.",
		"output", string(output),
	)

	return err
}
