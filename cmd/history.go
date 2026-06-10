package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var historyPath string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show saved deployment history",
	Run: func(cmd *cobra.Command, args []string) {
		results, err := history.ReadDeploymentHistory(historyPath)

		if err != nil {
			logger.Log.Errorw(
				"Failed to read deployment history.",
				"error", err,
			)

			return
		}

		if len(results) == 0 {
			logger.Log.Info("No deployment history found.")
			return
		}

		for _, result := range results {
			logger.Log.Infow(
				"Deployment result.",
				"deployment", result.DeploymentName,
				"success", result.Success,
				"started", result.Started,
				"finished", result.Finished,
				"step_count", len(result.StepResults),
			)
		}
	},
}

func init() {
	historyCmd.Flags().StringVar(
		&historyPath,
		"file",
		".lunar-deploy/history.jsonl",
		"Path to deployment history file",
	)

	rootCmd.AddCommand(historyCmd)
}
