package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/ogilcher/lunar-deploy-agent/internal/server"
	"github.com/spf13/cobra"
)

var serveAddress string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Lunar Deploy Agent HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Infow(
			"Starting Lunar Deploy Agent server.",
			"address", serveAddress,
		)

		if err := server.StartServer(serveAddress); err != nil {
			logger.Log.Errorw(
				"Server stopped with error.",
				"error", err,
			)
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(
		&serveAddress,
		"address",
		":8080",
		"HTTP server listen address",
	)

	rootCmd.AddCommand(serveCmd)
}
