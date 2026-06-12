package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/events"
	"github.com/ogilcher/lunar-deploy-agent/internal/health"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/jobs"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

type DeploymentListItem struct {
	Name           string `json:"name"`
	RepositoryPath string `json:"repository_path"`
	Preset         string `json:"preset"`
	StepCount      int    `json:"step_count"`
}

type AgentStatusResponse struct {
	ConfigValid           bool `json:"config_valid"`
	DeploymentCount       int  `json:"deployment_count"`
	HistoryExists         bool `json:"history_exists"`
	HistoryCount          int  `json:"history_count"`
	LastDeploymentSuccess bool `json:"last_deployment_success"`
}

type DeployRequest struct {
	Deployment string `json:"deployment"`
}

func StartServer(address string, configPath string) error {
	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/openapi.yaml", handleOpenAPI)
	mux.Handle("/swagger/", SwaggerHandler())

	protected := func(
		handler http.HandlerFunc,
	) http.HandlerFunc {
		return requireAPIToken(
			handler,
			appConfig.APIToken,
		)
	}

	mux.HandleFunc("/deployments", protected(func(writer http.ResponseWriter, request *http.Request) {
		handleDeployments(writer, request, configPath)
	}))
	mux.HandleFunc("/status", protected(func(writer http.ResponseWriter, request *http.Request) {
		handleStatus(writer, request, configPath)
	}))
	mux.HandleFunc("/deploy", protected(func(writer http.ResponseWriter, request *http.Request) {
		handleDeploy(writer, request, configPath)
	}))
	mux.HandleFunc("/node", protected(func(writer http.ResponseWriter, request *http.Request) {
		handleNodeInfo(writer, request, configPath)
	}))
	mux.HandleFunc("/heartbeat", protected(func(writer http.ResponseWriter, request *http.Request) {
		handleHeartbeat(writer, request, configPath)
	}))
	mux.HandleFunc("/history", protected(handleHistory))
	mux.HandleFunc("/jobs", protected(handleJobs))
	mux.HandleFunc("/jobs/", protected(handleJobByID))
	mux.HandleFunc("/queue", protected(handleQueue))
	mux.HandleFunc("/events", protected(handleEventsWebSocket))

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	StartDeploymentWorker(configPath)

	return server.ListenAndServe()
}

func handleHealth(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(writer, health.RunChecks())
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
			Name:           name,
			RepositoryPath: deployment.RepositoryPath,
			Preset:         deployment.Preset,
			StepCount:      len(steps),
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
			writeJSON(writer, AgentStatusResponse{
				ConfigValid:     true,
				DeploymentCount: len(appConfig.Deployments),
				HistoryExists:   false,
				HistoryCount:    0,
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
		ConfigValid:           true,
		DeploymentCount:       len(appConfig.Deployments),
		HistoryExists:         true,
		HistoryCount:          len(results),
		LastDeploymentSuccess: lastSuccess,
	})
}

func handleHistory(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results, err := history.ReadDeploymentHistory(".lunar-deploy/history.jsonl")

	if err != nil {
		if os.IsNotExist(err) {
			writeJSON(writer, []any{})
			return
		}

		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(writer, results)
}

func handleDeploy(
	writer http.ResponseWriter,
	request *http.Request,
	configPath string,
) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deployRequest DeployRequest

	if err := json.NewDecoder(request.Body).Decode(&deployRequest); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if deployRequest.Deployment == "" {
		http.Error(writer, "deployment is required", http.StatusBadRequest)
		return
	}

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, exists := appConfig.Deployments[deployRequest.Deployment]; !exists {
		http.Error(writer, "deployment not found", http.StatusNotFound)
		return
	}

	job := jobs.GlobalStore.CreateJob(deployRequest.Deployment)

	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "job_queued",
		JobID:      job.ID,
		Deployment: deployRequest.Deployment,
		Message:    "Deployment job queued.",
		Timestamp:  time.Now(),
	})

	jobs.GlobalQueue.Enqueue(job)

	writer.WriteHeader(http.StatusAccepted)
	writeJSON(writer, job)
}

func handleJobs(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(writer, jobs.GlobalStore.ListJobs())
}

func handleJobByID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id := request.URL.Path[len("/jobs/"):]

	if id == "" {
		http.Error(writer, "job id is required", http.StatusBadRequest)
		return
	}

	if request.Method == http.MethodDelete {
		cancelled := jobs.GlobalStore.MarkCancelled(id)

		if !cancelled {
			http.Error(writer, "job could not be cancelled", http.StatusBadRequest)
			return
		}

		events.GlobalEventBus.Publish(events.DeploymentEvent{
			Type:      "job_cancelled",
			JobID:     id,
			Message:   "Deployment job cancelled.",
			Timestamp: time.Now(),
		})
		writeJSON(writer, map[string]string{
			"status": "cancelled",
			"id":     id,
		})

		return
	}

	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	job, exists := jobs.GlobalStore.GetJob(id)

	if !exists {
		http.Error(writer, "job not found", http.StatusNotFound)
		return
	}

	writeJSON(writer, job)
}

func handleQueue(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(writer, jobs.GlobalStore.Summary())
}

func writeJSON(
	writer http.ResponseWriter,
	value any,
) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(value)
}
