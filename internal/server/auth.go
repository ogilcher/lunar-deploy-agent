package server

import "net/http"

// requireAPIToken protects API routes using a bearer token.
//
// Public endpoints like /health, /swagger, and /openapi.yaml can stay open.
// Operational endpoints like /deploy, /jobs, /history, and /queue should require auth.
func requireAPIToken(
	next http.HandlerFunc,
	apiToken string,
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if apiToken == "" {
			http.Error(writer, "api token is not configured", http.StatusUnauthorized)
			return
		}

		expectedHeader := "Bearer " + apiToken
		actualHeader := request.Header.Get("Authorization")

		if actualHeader != expectedHeader {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)
			return
		}

		next(writer, request)
	}
}
