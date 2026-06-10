package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/cmd/util"
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy"
	"github.com/ogilcher/lunar-deploy-agent/internal/events"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
	"github.com/spf13/cobra"
)

var deployConfigPath string
var deploymentName string
var outputJSON bool
var printEvents bool
var printEventsJSON bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run a deployment operation",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := util.LoadAndValidateConfig(deployConfigPath)

		deploymentConfig, exists := appConfig.Deployments[deploymentName]

		if !exists {
			logger.Log.Errorw(
				"Deployment does not exist.",
				"deployment", deploymentName,
			)

			return
		}

		if printEvents {
			eventChannel := events.GlobalEventBus.Subscribe()

			go func() {
				for event := range eventChannel {
					if printEventsJSON {
						eventJSON, marshalErr := json.Marshal(event)

						if marshalErr != nil {
							logger.Log.Errorw(
								"Failed to serialize deployment event.",
								"error", marshalErr,
							)

							continue
						}

						fmt.Println(string(eventJSON))
						continue
					}

					fmt.Printf(
						"[EVENT] %s | deployment=%s | step=%s | message=%s\n",
						event.Type,
						event.Deployment,
						event.Step,
						event.Message,
					)
				}
			}()
		}

		expandedSteps, err := config.ExpandedDeploymentSteps(deploymentConfig)

		if err != nil {
			logger.Log.Errorw(
				"Failed to expand deployment steps.",
				"deployment", deploymentName,
				"error", err,
			)

			return
		}

		localDeployer := deploy.LocalDeployer{
			RepositoryPath: deploymentConfig.RepositoryPath,
			DeploymentName: deploymentName,
			Environment:    appConfig.Environment,
			StepConfigs:    expandedSteps,
		}

		result, err := localDeployer.Deploy()

		if saveErr := history.SaveDeploymentResult(
			".lunar-deploy",
			result,
		); saveErr != nil {
			logger.Log.Errorw(
				"Failed to save deployment history.",
				"error", saveErr,
			)
		}

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

	deployCmd.Flags().BoolVar(
		&printEvents,
		"events",
		false,
		"Print deployment events in real time",
	)

	deployCmd.Flags().BoolVar(
		&printEventsJSON,
		"event-json",
		false,
		"Print deployment events as JSON",
	)

	rootCmd.AddCommand(deployCmd)
}
