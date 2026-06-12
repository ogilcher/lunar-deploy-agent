package health

import "github.com/ogilcher/lunar-deploy-agent/internal/jobs"

type QueueCheck struct{}

func (c *QueueCheck) Name() string {
	return "queue"
}

func (c *QueueCheck) Run() CheckResult {
	summary := jobs.GlobalStore.Summary()

	if summary.Queued > 100 {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: "queue backlog exceeds threshold",
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "queue operating normally",
	}
}
