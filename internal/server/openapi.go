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
  /health:
	get:
	  summary: Check agent health
	  responses:
		"200":
		  description: Agent is healthy

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
        "500":
          description: Deployment failed
`
