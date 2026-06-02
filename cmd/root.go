package cmd

import (
	"fmt"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lunar-client",
	Short: "Lunar Deploy Agent",
	Long:  "A lightweight deployment orchestration agent.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use --help to view available commands.")
	},
}

func Execute() {
	if err := logger.Initialize(); err != nil {
		fmt.Println("Failed to initialize logger:", err)

		os.Exit(1)
	}

	logger.Log.Info("Logger initialized.")

	if err := rootCmd.Execute(); err != nil {
		logger.Log.Error(err)

		os.Exit(1)
	}
}
