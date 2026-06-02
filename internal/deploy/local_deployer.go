package deploy

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// LocalDeployer runs deployment jobs against a locally accessible repository
type LocalDeployer struct {
	RepositoryPath string
	StepConfigs    []config.DeployStepConfig
}

// Deploy builds and executes the local deployment pipeline.
func (d *LocalDeployer) Deploy() error {
	logger.Log.Infow(
		"Starting local deployment.",
		"repository_path", d.RepositoryPath,
	)

	steps := []DeployStep{}

	for _, stepConfig := range d.StepConfigs {
		step, err := BuildStep(stepConfig, d.RepositoryPath)

		if err != nil {
			logger.Log.Errorw(
				"Failed to build deployment step.",
				"type", stepConfig.Type,
				"error", err,
			)

			return err
		}

		steps = append(steps, step)
	}

	job := DeployJob{
		Steps: steps,
	}

	if err := job.Run(); err != nil {
		return err
	}

	logger.Log.Info("Deployment completed successfully.")

	return nil
}
