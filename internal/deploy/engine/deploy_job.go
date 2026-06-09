package engine

import (
	"time"

	deploycontext "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"
	"github.com/ogilcher/lunar-deploy-agent/internal/deploy/steps"
	"github.com/ogilcher/lunar-deploy-agent/internal/events"
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// DeployJob executes a sequence of deployment steps in order.
//
// A job is made up of ordered deployment steps. Each step runs one at a time,
// and the job stops immediately if any step fails.
type DeployJob struct {
	Steps []steps.DeployStep
}

// Run executes each deployment step and stops on the first failure.
func (j *DeployJob) Run(
	context deploycontext.DeploymentContext,
) (*DeploymentResult, error) {
	result := newDeploymentResult(context)

	publishDeploymentStarted(context, result.Started)

	for _, step := range j.Steps {
		stepResult, err := j.runStep(context, step)

		result.StepResults = append(result.StepResults, stepResult)

		if err != nil {
			result.Success = false
			result.Finished = time.Now()

			publishDeploymentFailed(context, result, err)

			return result, err
		}
	}

	result.Success = true
	result.Finished = time.Now()

	publishDeploymentCompleted(context, result)

	logger.Log.Infow(
		"Deployment job completed successfully",
		"deployment", context.DeploymentName,
		"step_count", len(result.StepResults),
	)

	return result, nil
}

func (j *DeployJob) runStep(
	context deploycontext.DeploymentContext,
	step steps.DeployStep,
) (DeploymentStepResult, error) {
	stepResult := DeploymentStepResult {
		StepName: step.Name(),
		Success: false,
		Started: time.Now(),
	}

	publishStepStarted(context, step, stepResult.Started)

	logger.Log.Infow(
		"Running deployment step.",
		"deployment", context.DeploymentName,
		"step", step.Name(),
	)

	err := step.Run(context)

	stepResult.Finished = time.Now()

	if err != nil {
		stepResult.Error = err.Error()

		logger.Log.Errorw(
			"Deployment step failed.",
			"deployment", context.DeploymentName,
			"step", step.Name(),
			"error", err,
		)

		publishStepFailed(context, step, stepResult, err)

		return stepResult, err
	}

	stepResult.Success = true

	logger.Log.Infow(
		"Deployment step completed successfully.",
		"deployment", context.DeploymentName,
		"step", step.Name(),
	)

	publishStepCompleted(context, step, stepResult)

	return stepResult, nil
}

func newDeploymentResult(
	context deploycontext.DeploymentContext,
) *DeploymentResult {
	return &DeploymentResult{
		DeploymentName: context.DeploymentName,
		Success:        false,
		Started:        time.Now(),
		StepResults:    []DeploymentStepResult{},
	}
}

func publishDeploymentStarted(
	context deploycontext.DeploymentContext,
	started time.Time,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "deployment_started",
		Deployment: context.DeploymentName,
		Message:    "Deployment started.",
		Timestamp:  started,
	})
}

func publishDeploymentCompleted(
	context deploycontext.DeploymentContext,
	result *DeploymentResult,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "deployment_completed",
		Deployment: context.DeploymentName,
		Message:    "Deployment completed successfully.",
		DurationMilliseconds: result.Finished.Sub(
			result.Started,
		).Milliseconds(),
	})
}

func publishDeploymentFailed(
	context deploycontext.DeploymentContext,
	result *DeploymentResult,
	err error,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "deployment_failed",
		Deployment: context.DeploymentName,
		Message:    err.Error(),
		Timestamp:  result.Finished,
		DurationMilliseconds: result.Finished.Sub(
			result.Started,
		).Milliseconds(),
	})
}

func publishStepStarted(
	context deploycontext.DeploymentContext,
	step steps.DeployStep,
	started time.Time,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "step_started",
		Deployment: context.DeploymentName,
		Step:       step.Name(),
		Message:    "Deployment step started.",
		Timestamp:  started,
	})
}

func publishStepCompleted(
	context deploycontext.DeploymentContext,
	step steps.DeployStep,
	result DeploymentStepResult,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "step_completed",
		Deployment: context.DeploymentName,
		Step:       step.Name(),
		Message:    "Deployment step completed successfully.",
		Timestamp:  result.Finished,
		DurationMilliseconds: result.Finished.Sub(
			result.Started,
		).Milliseconds(),
	})
}

func publishStepFailed(
	context deploycontext.DeploymentContext,
	step steps.DeployStep,
	result DeploymentStepResult,
	err error,
) {
	events.GlobalEventBus.Publish(events.DeploymentEvent{
		Type:       "step_failed",
		Deployment: context.DeploymentName,
		Step:       step.Name(),
		Message:    err.Error(),
		Timestamp:  result.Finished,
		DurationMilliseconds: result.Finished.Sub(
			result.Started,
		).Milliseconds(),
	})
}
