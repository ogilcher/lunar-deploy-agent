package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var lastHistoryPath string
var lastOutputJSON bool

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Show the most recent deployment result",
	Run: func(cmd *cobra.Command, args []string) {
		results, err := history.ReadDeploymentHistory(
			lastHistoryPath,
		)

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

		lastResult := results[len(results)-1]

		if lastOutputJSON {
			resultJSON, err := json.MarshalIndent(lastResult, "", "  ")

			if err != nil {
				logger.Log.Errorw(
					"Failed to serialize last deployment result.",
					"error", err,
				)

				return
			}

			fmt.Println(string(resultJSON))
			return
		}

		logger.Log.Infow(
			"Last deployment result.",
			"deployment", lastResult,
			"success", lastResult.Success,
			"started", lastResult.Started,
			"finished", lastResult.Finished,
			"step_count", len(lastResult.StepResults),
		)
	},
}

func init() {
	lastCmd.Flags().StringVar(
		&lastHistoryPath,
		"file",
		".lunar-deploy/history.jsonl",
		"Path to deployment history file",
	)

	lastCmd.Flags().BoolVar(
		&lastOutputJSON,
		"json",
		false,
		"Print last deployment result as JSON",
	)

	rootCmd.AddCommand(lastCmd)
}
