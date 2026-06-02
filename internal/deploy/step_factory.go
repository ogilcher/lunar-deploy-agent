package deploy

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
			StepName:      stepConfig.Name,
			Command:       stepConfig.Command,
			Arguments:     stepConfig.Arguments,
			DirectoryPath: repositoryPath,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported deployment step type: %s", stepConfig.Type)

	}
}
