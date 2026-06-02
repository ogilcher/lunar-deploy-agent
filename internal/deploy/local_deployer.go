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
		switch stepConfig.Type {
		case "git_pull":
			steps = append(steps, &GitPullStep{
				RepositoryPath: d.RepositoryPath,
			})

		case "shell":
			steps = append(steps, &ShellCommandStep{
				StepName:      stepConfig.Name,
				Command:       stepConfig.Command,
				Arguments:     stepConfig.Arguments,
				DirectoryPath: d.RepositoryPath,
			})
		}
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
