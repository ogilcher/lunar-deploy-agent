package cmd

import (
	"os"

	"github.com/ogilcher/lunar-deploy-agent/cmd/util"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var statusConfigPath string
var statusHistoryPath string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show local agent status",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := util.LoadAndValidateConfig(statusConfigPath)

		results, err := history.ReadDeploymentHistory(statusHistoryPath)

		if err != nil {
			if os.IsNotExist(err) {
				logger.Log.Infow(
					"Agent status.",
					"config_valid", true,
					"deployment_count", len(appConfig.Deployments),
					"history_exists", false,
					"history_count", 0,
				)

				return
			}

			logger.Log.Errorw(
				"Failed to read deployment history.",
				"error", err,
			)

			return
		}

		lastSuccess := false

		if len(results) > 0 {
			lastSuccess = results[len(results)-1].Success
		}

		logger.Log.Infow(
			"Agent status.",
			"config_valid", true,
			"deployment_count", len(appConfig.Deployments),
			"history_exists", true,
			"history_count", len(results),
			"last_deployment_success", lastSuccess,
		)
	},
}

func init() {
	statusCmd.Flags().StringVarP(
		&statusConfigPath,
		"config",
		"c",
		"config.example.yaml",
		"Path to the deployment config file",
	)

	statusCmd.Flags().StringVar(
		&statusHistoryPath,
		"history-file",
		".lunar-deploy/history.jsonl",
		"Path to deployment history file",
	)

	rootCmd.AddCommand(statusCmd)
}
