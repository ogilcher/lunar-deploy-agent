package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/cmd/util"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var statusConfigPath string
var statusHistoryPath string
var statusOutputJSON bool

type AgentStatus struct {
	ConfigValid           bool `json:"config_valid"`
	DeploymentCount       int  `json:"deployment_count"`
	HistoryExists         bool `json:"history_exists"`
	HistoryCount          int  `json:"history_count"`
	LastDeploymentSuccess bool `json:"last_deployment_success"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show local agent status",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := util.LoadAndValidateConfig(statusConfigPath)

		results, err := history.ReadDeploymentHistory(statusHistoryPath)

		if err != nil {
			if os.IsNotExist(err) {
				printAgentStatus(
					AgentStatus{
						ConfigValid:     true,
						DeploymentCount: len(appConfig.Deployments),
						HistoryExists:   false,
						HistoryCount:    0,
					},
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

		printAgentStatus(
			AgentStatus{
				ConfigValid:           true,
				DeploymentCount:       len(appConfig.Deployments),
				HistoryExists:         true,
				HistoryCount:          len(results),
				LastDeploymentSuccess: lastSuccess,
			})
	},
}

func printAgentStatus(
	status AgentStatus,
) {
	if statusOutputJSON {
		statusJSON, err := json.MarshalIndent(status, "", "	")

		if err != nil {
			logger.Log.Errorw(
				"Failed to serialize agent status.",
				"error", err,
			)

			return
		}

		fmt.Println(string(statusJSON))
		return
	}

	logger.Log.Infow(
		"Agent status.",
		"config_valid", status.ConfigValid,
		"deployment_count", status.DeploymentCount,
		"history_exists", status.HistoryExists,
		"history_count", status.HistoryCount,
		"last_deployment_success", status.LastDeploymentSuccess,
	)
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

	statusCmd.Flags().BoolVar(
		&statusOutputJSON,
		"json",
		false,
		"Print local agent status as JSON",
	)

	rootCmd.AddCommand(statusCmd)
}
