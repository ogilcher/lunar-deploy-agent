package http

import (
	"fmt"
	"net/http"
	"time"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// HTTPHealthCheckStep verifies that an HTTP endpoint responds with the
// expected status code after deployment
type HTTPHealthCheckStep struct {
	URL            string
	ExpectedStatus int
	TimeoutSeconds int
}

func (s *HTTPHealthCheckStep) Name() string {
	return "http_health_check"
}

func (s *HTTPHealthCheckStep) Run(
	context deploycontext.DeploymentContext,
) error {
	timeout := time.Duration(s.TimeoutSeconds) * time.Second

	if s.TimeoutSeconds <= 0 {
		timeout = 10 * time.Second
	}

	client := http.Client{
		Timeout: timeout,
	}

	logger.Log.Infow(
		"Running HTTP health check...",
		"deployment", context.DeploymentName,
		"url", s.URL,
		"expected_status", s.ExpectedStatus,
	)

	response, err := client.Get(s.URL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != s.ExpectedStatus {
		return fmt.Errorf(
			"health check failed: expected status %d but got %d",
			s.ExpectedStatus,
			response.StatusCode,
		)
	}

	logger.Log.Infow(
		"HTTP health check passed successfully.",
		"deployment", context.DeploymentName,
		"url", s.URL,
		"status", response.StatusCode,
	)

	return nil
}
