package cmd

import (
	"bufio"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var historyPath string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show saved deployment history",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(historyPath)
		if err != nil {
			logger.Log.Errorw("Failed to open deployment history.", "error", err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		
		for scanner.Scan() {
			logger.Log.Info(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			logger.Log.Errorw("Failed to read deployment history.", "error", err)
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
