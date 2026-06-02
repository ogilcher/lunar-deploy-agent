package deploy

import (
	"time"

	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// DeployJob executes a sequence of deployment steps in order.
//
// A job is made up of ordered deployment steps. Each step runs one at a time,
// and the job stops immediately if any step fails.
type DeployJob struct {
	Steps []DeployStep
}

// Run executes each deployment step and stops on the first failure.
func (j *DeployJob) Run(
	context DeploymentContext,
) (*DeploymentResult, error) {
	result := &DeploymentResult{
		DeploymentName: context.DeploymentName,
		Success:        false,
		Started:        time.Now(),
		StepResults:    []DeploymentStepResult{},
	}

	for _, step := range j.Steps {
		// Step starting
		stepResult := DeploymentStepResult{
			StepName: step.Name(),
			Success:  false,
			Started:  time.Now(),
		}

		logger.Log.Infow(
			"Running deployment step.",
			"deployment", context.DeploymentName,
			"step", step.Name(),
		)

		err := step.Run(context)

		stepResult.Finished = time.Now()

		if err != nil {
			stepResult.Error = err.Error()

			result.StepResults = append(
				result.StepResults,
				stepResult,
			)

			result.Success = false
			result.Finished = time.Now()

			logger.Log.Errorw(
				"Deployment step failed.",
				"deployment", context.DeploymentName,
				"step", step.Name(),
				"error", err,
			)

			return result, err
		}

		stepResult.Success = true

		result.StepResults = append(
			result.StepResults,
			stepResult,
		)

		logger.Log.Infow(
			"Deployment step completed successfully.",
			"deployment", context.DeploymentName,
			"step", step.Name(),
		)

	}

	result.Success = true
	result.Finished = time.Now()

	logger.Log.Infow(
		"Deployment job completed successfully.",
		"deployment", context.DeploymentName,
		"success", result.Success,
		"step_count", len(result.StepResults),
	)

	return result, nil
}
