package server

import (
	"encoding/json"
	"net/http"
	"time"

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

func StartServer(address string, configPath string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/deployments", func(writer http.ResponseWriter, request *http.Request) {
		handleDeployments(writer, request, configPath)
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

	writer.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(writer).Encode(response)
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

	writer.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(writer).Encode(deployments)
}
