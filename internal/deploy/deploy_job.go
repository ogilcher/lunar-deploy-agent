package deploy

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

// DeployJob executes a sequence of deployment steps in order.
type DeployJob struct {
	Steps []DeployStep
}

// Run executes each deployment step and stops on the first failure.
func (j *DeployJob) Run(
	context DeploymentContext,
) error {
	for _, step := range j.Steps {
		logger.Log.Infow(
			"Running deployment step.",
			"step", step.Name(),
		)

		if err := step.Run(context); err != nil {
			logger.Log.Errorw(
				"Deployment step failed.",
				"step", step.Name(),
				"error", err,
			)

			return err
		}

		logger.Log.Infow(
			"Deployment step completed.",
			"step", step.Name(),
		)
	}

	return nil
}
