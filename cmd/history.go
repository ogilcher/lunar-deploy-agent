package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var historyPath string
var historyOutputJSON bool
var historyLimit int
var historyDeploymentFilter string
var historySummary bool

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

		if historyDeploymentFilter != "" {
			filteredResults := []engine.DeploymentResult{}

			for _, result := range results {
				if result.DeploymentName == historyDeploymentFilter {
					filteredResults = append(
						filteredResults,
						result,
					)
				}
			}

			results = filteredResults
		}

		if historyLimit > 0 && historyLimit < len(results) {
			results = results[len(results)-historyLimit:]
		}

		if historySummary {
			total := len(results)
			successful := 0

			for _, result := range results {
				if result.Success {
					successful++
				}
			}

			failed := total - successful
			successRate := 0.0

			if total > 0 {
				successRate = float64(successful) / float64(total) * 100
			}

			logger.Log.Infow(
				"Deployment history summary.",
				"total", total,
				"successful", successful,
				"failed", failed,
				"success_rate", successRate,
			)

			return
		}

		if historyOutputJSON {
			resultsJSON, err := json.MarshalIndent(results, "", "  ")

			if err != nil {
				logger.Log.Errorw(
					"Failed to serialize deployment history.",
					"error", err,
				)

				return
			}

			fmt.Println(string(resultsJSON))
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

	historyCmd.Flags().BoolVar(
		&historyOutputJSON,
		"json",
		false,
		"Print deployment history as JSON",
	)

	historyCmd.Flags().IntVar(
		&historyLimit,
		"limit",
		0,
		"Maximum number of deployment history entries to show",
	)

	historyCmd.Flags().StringVar(
		&historyDeploymentFilter,
		"deployment",
		"",
		"Filter deployment history by deployment name",
	)

	historyCmd.Flags().BoolVar(
		&historySummary,
		"summary",
		false,
		"Show deployment history summary",
	)

	rootCmd.AddCommand(historyCmd)
}
