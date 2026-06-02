package deploy

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

type LocalDeployer struct {
	RepositoryPath string
	StepConfigs    []config.DeployStepConfig
}

func (d *LocalDeployer) Deploy() error {
	logger.Log.Infow(
		"Starting local deployment.",
		"repository_path", d.RepositoryPath,
	)

	steps := []DeployStep{
		&GitPullStep{
			RepositoryPath: d.RepositoryPath,
		},
	}

	for _, stepConfig := range d.StepConfigs {
		steps = append(steps, &ShellCommandStep{
			StepName:      stepConfig.Name,
			Command:       stepConfig.Command,
			Arguments:     stepConfig.Arguments,
			DirectoryPath: d.RepositoryPath,
		})
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
