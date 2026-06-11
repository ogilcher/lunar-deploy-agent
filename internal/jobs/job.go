package jobs

import (
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/engine"
)

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobSucceeded JobStatus = "succeeded"
	JobFailed    JobStatus = "failed"
)

type DeploymentJob struct {
	ID             string                   `json:"id"`
	DeploymentName string                   `json:"deployment_name"`
	Status         JobStatus                `json:"status"`
	CreatedAt      time.Time                `json:"created_at"`
	StartedAt      *time.Time               `json:"started_at,omitempty"`
	FinishedAt     *time.Time               `json:"finished_at,omitempty"`
	Result         *engine.DeploymentResult `json:"result,omitempty"`
	Error          string                   `json:"error,omitempty"`
}
