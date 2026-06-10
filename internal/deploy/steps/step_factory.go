package steps

import (
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
)

func BuildStep(
	stepConfig config.DeployStepConfig,
	repositoryPath string,
) (DeployStep, error) {
	switch stepConfig.Type {
	case "git_pull":
		return &GitPullStep{
			RepositoryPath: repositoryPath,
		}, nil

	case "shell":
		return &ShellCommandStep{
			StepName:          stepConfig.Name,
			Command:           stepConfig.Command,
			Arguments:         stepConfig.Arguments,
			DirectoryPath:     repositoryPath,
			TimeoutSeconds:    stepConfig.TimeoutSeconds,
			WorkingDirectory:  stepConfig.WorkingDirectory,
			Retries:           stepConfig.Retries,
			RetryDelaySeconds: stepConfig.RetryDelaySeconds,
		}, nil

	case "pm2_restart":
		return &PM2RestartStep{
			ProcessName: stepConfig.ProcessName,
		}, nil

	case "npm_install":
		return &NPMInstallStep{
			DirectoryPath: repositoryPath,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported deployment step type: %s", stepConfig.Type)
	}
}
