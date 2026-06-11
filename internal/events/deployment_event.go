package events

import "time"

// DeploymentEvent represents a real-time deployment lifecycle event.
type DeploymentEvent struct {
	Type                 string    `json:"type"`
	Deployment           string    `json:"deployment"`
	Step                 string    `json:"step"`
	Message              string    `json:"message"`
	Timestamp            time.Time `json:"timestamp"`
	DurationMilliseconds int64     `json:"duration_milliseconds,omitempty"`
	JobID                string    `json:"job_id,omitempty"`
}
