package server

import (
	"fmt"
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/config"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy"
	"github.com/ogilcher/lunar-deploy-agent/internal/events"
	"github.com/ogilcher/lunar-deploy-agent/internal/history"
	"github.com/ogilcher/lunar-deploy-agent/internal/jobs"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// StartDeploymentWorker starts a background worker that processes queued
// deployment jobs one at a time.
func StartDeploymentWorker(
	configPath string,
) {
	go func() {
		for job := range jobs.GlobalQueue.Jobs() {
			processDeploymentJob(configPath, job)
		}
	}()
}

func processDeploymentJob(
	configPath string,
	job *jobs.DeploymentJob,
) {
	if jobs.GlobalStore.IsCancelled(job.ID) {
		return
	}

	jobs.GlobalStore.MarkRunning(job.ID)

	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "job_running",
		JobID:      job.ID,
		Deployment: job.DeploymentName,
		Message:    "Deployment job started.",
		Timestamp:  time.Now(),
	})

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		jobs.GlobalStore.MarkFailed(job.ID, nil, err)
		return
	}

	if err := config.ValidateConfig(appConfig); err != nil {
		jobs.GlobalStore.MarkFailed(job.ID, nil, err)
		return
	}

	deploymentConfig, exists := appConfig.Deployments[job.DeploymentName]
	if !exists {
		jobs.GlobalStore.MarkFailed(
			job.ID,
			nil,
			fmt.Errorf("deployment %q not found", job.DeploymentName),
		)
		return
	}

	steps, err := config.ExpandedDeploymentSteps(deploymentConfig)
	if err != nil {
		jobs.GlobalStore.MarkFailed(job.ID, nil, err)
		return
	}

	localDeployer := deploy.LocalDeployer{
		RepositoryPath: deploymentConfig.RepositoryPath,
		DeploymentName: job.DeploymentName,
		Environment:    appConfig.Environment,
		StepConfigs:    steps,
	}

	result, err := localDeployer.Deploy()

	if saveErr := history.SaveDeploymentResult(
		".lunar-deploy",
		result,
	); saveErr != nil {
		logger.Log.Errorw(
			"Failed to save deployment history.",
			"error", saveErr,
		)
	}

	if err != nil {
		jobs.GlobalStore.MarkFailed(job.ID, result, err)

		events.GlobalEventBus.Publish(events.DeploymentEvent{
			Type:       "job_failed",
			JobID:      job.ID,
			Deployment: job.DeploymentName,
			Message:    err.Error(),
			Timestamp:  time.Now(),
		})

		return
	}

	jobs.GlobalStore.MarkSucceeded(job.ID, result)

	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "job_succeeded",
		JobID:      job.ID,
		Deployment: job.DeploymentName,
		Message:    "Deployment job succeeded.",
		Timestamp:  time.Now(),
	})
}
