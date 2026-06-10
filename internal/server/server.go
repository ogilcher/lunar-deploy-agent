package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status 		string 		`json:"status"`
	Service 	string 		`json:"service"`
	Timestamp 	time.Time 	`json:"timestamp"`
}

func StartServer(address string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)

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
