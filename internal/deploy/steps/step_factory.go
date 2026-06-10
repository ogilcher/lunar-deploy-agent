package steps

import (
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps/git"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps/http"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps/npm"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps/pm2"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps/shell"
)

func BuildStep(
	stepConfig config.DeployStepConfig,
	repositoryPath string,
) (DeployStep, error) {
	switch stepConfig.Type {
	case "git_pull":
		return &git.GitPullStep{
			RepositoryPath: repositoryPath,
		}, nil

	case "shell":
		return &shell.ShellCommandStep{
			StepName:          stepConfig.Name,
			Command:           stepConfig.Command,
			Arguments:         stepConfig.Arguments,
			DirectoryPath:     repositoryPath,
			TimeoutSeconds:    stepConfig.TimeoutSeconds,
			WorkingDirectory:  stepConfig.WorkingDirectory,
			Retries:           stepConfig.Retries,
			RetryDelaySeconds: stepConfig.RetryDelaySeconds,
		}, nil

	case "http_health_check":
		return &http.HTTPHealthCheckStep{
			URL:            stepConfig.URL,
			ExpectedStatus: stepConfig.ExpectedStatus,
			TimeoutSeconds: stepConfig.TimeoutSeconds,
		}, nil

	case "pm2_restart":
		return &pm2.PM2RestartStep{
			ProcessName: stepConfig.ProcessName,
		}, nil

	case "pm2_start_or_restart":
		return &pm2.PM2StartOrRestartStep{
			ProcessName:   stepConfig.ProcessName,
			StartCommand:  stepConfig.StartCommand,
			Arguments:     stepConfig.Arguments,
			DirectoryPath: repositoryPath,
		}, nil

	case "pm2_save":
		return &pm2.PM2SaveStep{}, nil

	case "npm_install":
		return &npm.NPMInstallStep{
			DirectoryPath: repositoryPath,
		}, nil

	case "npm_build":
		return &npm.NPMBuildStep{
			DirectoryPath: repositoryPath,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported deployment step type: %s", stepConfig.Type)
	}
}
