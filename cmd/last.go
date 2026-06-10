package cmd

import (
	"bufio"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var lastHistoryPath string

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Show the most recent deployment result",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(lastHistoryPath)
		if err != nil {
			logger.Log.Errorw("Failed to open deployment history.", "error", err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		lastLine := ""

		for scanner.Scan() {
			lastLine = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			logger.Log.Errorw("Failed to read deployment history.", "error", err)
			return
		}

		if lastLine == "" {
			logger.Log.Info("No deployment history found.")
			return
		}

		logger.Log.Info(lastLine)
	},
}

func init() {
	lastCmd.Flags().StringVar(
		&lastHistoryPath,
		"file",
		".lunar-deploy/history.jsonl",
		"Path to deployment history file",
	)

	rootCmd.AddCommand(lastCmd)
}