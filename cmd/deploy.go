package cmd

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var repositoryPath string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run a deployment operation",
	Run: func(cmd *cobra.Command, args []string) {
		localDeployer := deploy.LocalDeployer{
			RepositoryPath: repositoryPath,
		}

		if err := localDeployer.Deploy(); err != nil {
			logger.Log.Errorw(
				"Deploy command failed.",
				"error", err,
			)

			return
		}

		logger.Log.Info("Deploy command completed.")
	},
}

func init() {
	deployCmd.Flags().StringVarP(
		&repositoryPath,
		"path",
		"p",
		".",
		"Path to the repository",
	)

	rootCmd.AddCommand(deployCmd)
}
