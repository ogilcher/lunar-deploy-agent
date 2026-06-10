package deploy

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
	steps2 "github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// LocalDeployer runs deployment jobs against a locally accessible repository
type LocalDeployer struct {
	RepositoryPath string
	DeploymentName string
	Environment    string
	StepConfigs    []config.DeployStepConfig
}

// Deploy builds and executes the local deployment pipeline.
func (d *LocalDeployer) Deploy() (*engine.DeploymentResult, error) {
	logger.Log.Infow(
		"Starting local deployment.",
		"repository_path", d.RepositoryPath,
	)

	steps := []steps2.DeployStep{}

	for _, stepConfig := range d.StepConfigs {
		step, err := steps2.BuildStep(stepConfig, d.RepositoryPath)

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

	job := engine.DeployJob{
		Steps: steps,
	}

	context := deploycontext.DeploymentContext{
		DeploymentName: d.DeploymentName,
		RepositoryPath: d.RepositoryPath,
		Environment:    d.Environment,
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
