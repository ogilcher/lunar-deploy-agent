package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var validateConfigPath string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the deployment config file",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig, err := config.LoadConfig(validateConfigPath)
		if err != nil {
			logger.Log.Errorw("Failed to load config.", "error", err)
			return
		}

		if err := config.ValidateConfig(appConfig); err != nil {
			logger.Log.Errorw("Invalid deployment config.", "error", err)
			return
		}

		logger.Log.Info("Deployment config is valid.")
	},
}

func init() {
	validateCmd.Flags().StringVarP(
		&validateConfigPath,
		"config",
		"c",
		"config.example.yaml",
		"Path to the deployment config file",
	)

	rootCmd.AddCommand(validateCmd)
}
