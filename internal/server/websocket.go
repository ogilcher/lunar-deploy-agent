package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/ogilcher/lunar-deploy-agent/internal/events"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// handleEventsWebSocket streams deployment events to connected WebSocket clients.
func handleEventsWebSocket(
	writer http.ResponseWriter,
	request *http.Request,
) {
	connection, err := websocket.Accept(
		writer,
		request,
		nil,
	)

	if err != nil {
		logger.Log.Errorw("Failed to accept WebSocket connection.", "error", err)
		return
	}

	defer connection.Close(
		websocket.StatusNormalClosure,
		"connection closed",
	)

	eventChannel := events.GlobalEventBus.Subscribe()
	defer events.GlobalEventBus.Unsubscribe(eventChannel)

	for event := range eventChannel {
		eventJSON, err := json.Marshal(event)

		if err != nil {
			logger.Log.Errorw("Failed to serialize WebSocket event.", "error", err)
			continue
		}

		writeContext, cancel := context.WithTimeout(
			request.Context(),
			5*time.Second,
		)

		err = connection.Write(
			writeContext,
			websocket.MessageText,
			eventJSON,
		)

		cancel()

		if err != nil {
			logger.Log.Errorw("Failed to write WebSocket event.", "error", err)
			return
		}
	}
}
