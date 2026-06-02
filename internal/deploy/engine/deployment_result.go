package engine

import "time"

// DeploymentStepResult represents the result of a single deployment step.
type DeploymentStepResult struct {
	StepName string
	Success  bool
	Started  time.Time
	Finished time.Time
	Error    string
}

// DeploymentResult represents the overall deployment execution result.
type DeploymentResult struct {
	DeploymentName string
	Success        bool
	Started        time.Time
	Finished       time.Time
	StepResults    []DeploymentStepResult
}
