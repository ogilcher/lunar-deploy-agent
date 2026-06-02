package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var configPath string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Load and validate the agent config",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig, err := config.LoadConfig(configPath)

		if err != nil {
			logger.Log.Errorw("Failed to load config.", "error", err)
			return
		}

		if err := config.ValidateConfig(appConfig); err != nil {
			logger.Log.Errorw("Invalid config.", "error", err)
			return
		}

		logger.Log.Infow(
			"Config loaded successfully",
			"agent_id", appConfig.AgentID,
			"environment", appConfig.Environment,
			"workspace_path", appConfig.WorkspacePath,
			"log_level", appConfig.LogLevel,
		)
	},
}

func init() {
	configCmd.Flags().StringVarP(&configPath, "file", "f", "config.example.yaml", "Path to the config file")
	rootCmd.AddCommand(configCmd)
}
