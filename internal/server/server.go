package server

import (
	"encoding/json"
	"net/http"
	"time"
	"os"

	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/config"
)

type HealthResponse struct {
	Status 		string 		`json:"status"`
	Service 	string 		`json:"service"`
	Timestamp 	time.Time 	`json:"timestamp"`
}

type DeploymentListItem struct {
	Name 			string	`json:"name"`
	RepositoryPath 	string 	`json:"repository_path"`
	Preset 			string	`json:"preset"`
	StepCount		int		`json:"step_count"`
}

type AgentStatusResponse struct {
	ConfigValid			bool	`json:"config_valid"`
	DeploymentCount 	int 	`json:"deployment_count"`
	HistoryExists		bool	`json:"history_exists"`
	HistoryCount		int 	`json:"history_count"`
	LastDeploymentSuccess bool 	`json:"last_deployment_success"`
}

func StartServer(address string, configPath string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/deployments", func(writer http.ResponseWriter, request *http.Request) {
		handleDeployments(writer, request, configPath)
	})
	mux.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		handleStatus(writer, request, configPath)
	})

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func handleHealth(
	writer http.ResponseWriter,
	request *http.Request,
) {
	response := HealthResponse{
		Status: 	"OK",
		Service: 	"lunar-deploy-agent",
		Timestamp: 	time.Now(),
	}

	writeJSON(writer, response)
}

func handleDeployments(
	writer http.ResponseWriter,
	request *http.Request,
	configPath string,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	deployments := []DeploymentListItem{}

	for name, deployment := range appConfig.Deployments {
		steps, err := config.ExpandedDeploymentSteps(deployment)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		deployments = append(deployments, DeploymentListItem{
			Name: 			name,
			RepositoryPath: deployment.RepositoryPath,
			Preset: 		deployment.Preset,
			StepCount: 		len(steps),
		})
	}

	writeJSON(writer, deployments)
}

func handleStatus(
	writer http.ResponseWriter,
	request *http.Request,
	configPath string,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := history.ReadDeploymentHistory(".lunar-deploy/history.jsonl")

	if err != nil {
		if os.IsNotExist(err) {
			writeJSON(writer, AgentStatusResponse {
				ConfigValid: 		true,
				DeploymentCount: 	len(appConfig.Deployments),
				HistoryExists: 		false,
				HistoryCount: 		0,
			})

			return
		}

		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	lastSuccess := false

	if len(results) > 0 {
		lastSuccess = results[len(results)-1].Success
	}

	writeJSON(writer, AgentStatusResponse{
		ConfigValid: 			true,
		DeploymentCount: 		len(appConfig.Deployments),
		HistoryExists: 			true,
		HistoryCount: 			len(results),
		LastDeploymentSuccess: 	lastSuccess,
	})
}

func writeJSON(
	writer http.ResponseWriter,
	value any,
) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(value)
}