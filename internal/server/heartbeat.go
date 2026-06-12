package server

import (
	"net/http"
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/jobs"
)

var serverStartedAt = time.Now()

type HeartbeatResponse struct {
	NodeID        string       `json:"node_id"`
	NodeName      string       `json:"node_name"`
	NodeRegion    string       `json:"node_region"`
	Service       string       `json:"service"`
	Version       string       `json:"version"`
	Status        string       `json:"status"`
	StartedAt     time.Time    `json:"started_at"`
	Timestamp     time.Time    `json:"timestamp"`
	UptimeSeconds int64        `json:"uptime_seconds"`
	Queue         jobs.Summary `json:"queue"`
}

func handleHeartbeat(
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

	now := time.Now()

	writeJSON(writer, HeartbeatResponse{
		NodeID:        appConfig.NodeID,
		NodeName:      appConfig.NodeName,
		NodeRegion:    appConfig.NodeRegion,
		Service:       "lunar-deploy-agent",
		Version:       "0.2.0",
		Status:        "online",
		StartedAt:     serverStartedAt,
		Timestamp:     now,
		UptimeSeconds: int64(now.Sub(serverStartedAt).Seconds()),
		Queue:         jobs.GlobalStore.Summary(),
	})
}
