package deploy

import (
	"os/exec"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

type LocalDeployer struct {
	RepositoryPath string
}

func (d *LocalDeployer) Deploy() error {
	logger.Log.Infow(
		"Starting local deployment.",
		"repository_path", d.RepositoryPath,
	)

	command := exec.Command("git", "-C", d.RepositoryPath, "pull")

	output, err := command.CombinedOutput()

	logger.Log.Infow(
		"Git pull completed.",
		"output", string(output),
	)

	if err != nil {
		logger.Log.Errorw(
			"Deployment failed.",
			"error", err,
		)

		return err
	}

	logger.Log.Info("Deployment completed successfully.")

	return nil
}
