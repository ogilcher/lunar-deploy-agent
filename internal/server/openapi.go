package server

import "net/http"

func handleOpenAPI(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/yaml")
	_, _ = writer.Write([]byte(openAPISpec))
}

const openAPISpec = `
openapi: 3.0.3
info:
  title: Lunar Deploy Agent API
  version: 0.1.0
  description: HTTP API for the Lunar Deploy Agent node service.

paths:
  /status:
    get:
      summary: Get local agent status
      responses:
        "200":
          description: Agent status response

  /deployments:
    get:
      summary: List configured deployments
      responses:
        "200":
          description: Configured deployments

  /history:
    get:
      summary: Get deployment history
      responses:
        "200":
          description: Deployment history

  /jobs:
    get:
      summary: List deployment jobs
      responses:
        "200":
          description: Deployment jobs

  /jobs/{id}:
    get:
      summary: Get a deployment job by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          description: Deployment job
        "404":
          description: Job not found

    delete:
      summary: Cancel a deployment job
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          description: Job cancelled
        "400":
          description: Job could not be cancelled
        "404":
          description: Job not found

  /events:
    get:
      summary: Stream deployment events over WebSocket
      responses:
        "101":
          description: WebSocket connection established

  /health:
    get:
      summary: Get agent health
      responses:
        "405":
          description: Method not allowed
        "500":
          description: Internal server error

  /deploy:
    post:
      summary: Trigger a deployment
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - deployment
              properties:
                deployment:
                  type: string
      responses:
        "200":
          description: Deployment completed
        "202":
          description: Deployment job queued
        "500":
          description: Deployment failed
`
