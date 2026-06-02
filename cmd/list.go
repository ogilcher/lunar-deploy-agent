package cmd

import (
	"sort"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var listConfigPath string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured deployments",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig, err := config.LoadConfig(listConfigPath)
		if err != nil {
			logger.Log.Error("Failed to load config.", "error", err)
			return
		}

		if err := config.ValidateConfig(appConfig); err != nil {
			logger.Log.Error("Invalid deployment config.", "error", err)
			return
		}

		deploymentNames := make([]string, 0, len(appConfig.Deployments))

		for deploymentName := range appConfig.Deployments {
			deploymentNames = append(deploymentNames, deploymentName)
		}

		sort.Strings(deploymentNames)

		for _, deploymentName := range deploymentNames {
			deploymentConfig := appConfig.Deployments[deploymentName]

			logger.Log.Infow(
				"Configured deployment.",
				"name", deploymentName,
				"repository_path", deploymentConfig.RepositoryPath,
				"step_count", len(deploymentConfig.Steps),
			)
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(
		&listConfigPath,
		"config",
		"c",
		"config.example.yaml",
		"Path to the deployment config file",
	)

	rootCmd.AddCommand(listCmd)
}