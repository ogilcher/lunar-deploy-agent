package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var deployConfigPath string
var deploymentName string
var outputJSON bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run a deployment operation",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig, err := config.LoadConfig(deployConfigPath)
		if err != nil {
			logger.Log.Errorw("Failed to load config.", "error", err)
			return
		}

		if err := config.ValidateConfig(appConfig); err != nil {
			logger.Log.Errorw("Invalid deploy config.", "error", err)
			return
		}

		deploymentConfig, exists := appConfig.Deployments[deploymentName]

		if !exists {
			logger.Log.Errorw(
				"Deployment does not exist.",
				"deployment", deploymentName,
			)

			return
		}

		localDeployer := deploy.LocalDeployer{
			RepositoryPath: deploymentConfig.RepositoryPath,
			StepConfigs:    deploymentConfig.Steps,
		}

		result, err := localDeployer.Deploy()

		if outputJSON && result != nil {
			resultJSON, marshalErr := json.MarshalIndent(
				result,
				"",
				"	",
			)

			if marshalErr != nil {
				logger.Log.Errorw(
					"Failed to serialize deployment result.",
					"error", marshalErr,
				)

				return
			}

			fmt.Println(string(resultJSON))
		}

		if err != nil {
			logger.Log.Errorw(
				"Deploy command failed.",
				"error", err,
			)

			os.Exit(1)
		}

		logger.Log.Info("Deploy command completed successfully.")
	},
}

func init() {
	deployCmd.Flags().StringVarP(
		&deployConfigPath,
		"config",
		"c",
		"config.example.yaml",
		"Path to the deployment config file",
	)

	deployCmd.Flags().StringVarP(
		&deploymentName,
		"deployment",
		"d",
		"",
		"Deployment name to execute",
	)

	deployCmd.Flags().BoolVar(
		&outputJSON,
		"json",
		false,
		"Print deployment result as JSON",
	)

	rootCmd.AddCommand(deployCmd)
}
