package deploy

import (
	"github.com/ogilcher/lunar-deploy-agent/internal/logger"
)

type DeployJob struct {
	Steps []DeployStep
}

func (j *DeployJob) Run() error {
	for _, step := range j.Steps {
		logger.Log.Infow(
			"Running deployment step.",
			"step", step.Name(),
		)

		if err := step.Run(); err != nil {
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
