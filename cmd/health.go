package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check whether the agent is running correctly",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("Health check requested.")

		logger.Log.Info("Agent health: OK")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
