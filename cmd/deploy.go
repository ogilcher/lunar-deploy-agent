package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var deployConfigPath string
var deploymentName string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run a deployment operation",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig, err := config.LoadConfig(deployConfigPath)
		if err != nil {
			logger.Log.Errorw("Failed to load config.", "error", err)
			return
		}

		if err := config.ValidateConfig(appConfig); err != nil {
			logger.Log.Errorw("Invalid deploy config.", "error", err)
			return
		}

		deploymentConfig, exists := appConfig.Deployments[deploymentName]

		if !exists {
			logger.Log.Errorw(
				"Deployment does not exist.",
				"deployment", deploymentName,
			)

			return
		}

		localDeployer := deploy.LocalDeployer{
			RepositoryPath: deploymentConfig.RepositoryPath,
			StepConfigs:    deploymentConfig.Steps,
		}

		if err := localDeployer.Deploy(); err != nil {
			logger.Log.Errorw("Deploy command failed.", "error", err)
			return
		}

		logger.Log.Info("Deploy command completed.")
	},
}

func init() {
	deployCmd.Flags().StringVarP(
		&deployConfigPath,
		"config",
		"c",
		"config.example.yaml",
		"Path to the deployment config file",
	)

	deployCmd.Flags().StringVarP(
		&deploymentName,
		"deployment",
		"d",
		"",
		"Deployment name to execute",
	)

	rootCmd.AddCommand(deployCmd)
}
