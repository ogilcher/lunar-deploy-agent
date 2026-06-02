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
func (d *LocalDeployer) Deploy() (*DeploymentResult, error) {
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

			return nil, err
		}

		steps = append(steps, step)
	}

	job := DeployJob{
		Steps: steps,
	}

	context := DeploymentContext{
		DeploymentName: "local",
		RepositoryPath: d.RepositoryPath,
		Environment:    "development",
	}

	result, err := job.Run(context)

	if err != nil {
		return result, err
	}

	logger.Log.Info(
		"Deployment completed successfully.",
		"step_count", len(result.StepResults),
	)

	return result, nil
}
