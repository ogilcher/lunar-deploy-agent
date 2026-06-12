package server

import (
	"net/http"
	"runtime"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
)

type NodeInfoResponse struct {
	NodeID 			string 	 `json:"node_id"`
	NodeName 		string 	 `json:"node_name"`
	NodeRegion 		string 	 `json:"node_region"`
	Service			string 	 `json:"service"`
	Version			string 	 `json:"version"`
	OperatingSystem string 	 `json:"operating_system"`
	Architecture	string 	 `json:"architecture"`
	Capabilities 	[]string `json:"capabilities"`
}

func handleNodeInfo(
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

	writeJSON(writer, NodeInfoResponse{
		NodeID: 			appConfig.NodeID,
		NodeName: 			appConfig.NodeName,
		NodeRegion: 		appConfig.NodeRegion,
		Service: 			"lunar-deploy-agent",
		Version: 			"0.2.0",
		OperatingSystem: 	runtime.GOOS,
		Architecture: 		runtime.GOARCH,
		Capabilities: []string{
			"git_pull",
			"shell",
			"npm_install",
			"npm_build",
			"pm2_restart",
			"pm2_start_or_restart",
			"pm2_save",
			"pm2_status",
			"http_health_check",
			"websocket_events",
			"async_jobs",
			"swagger_ui",
		},
	})
}