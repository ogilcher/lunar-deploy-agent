package health

import "time"

var agentStartedAt = time.Now()

// UptimeCheck reports how long the agent process has been running.
type UptimeCheck struct {}

func (c *UptimeCheck) Name() string {
	return "uptime"
}

func (c *UptimeCheck) Run() CheckResult {
	uptimeSeconds := int64(time.Since(agentStartedAt).Seconds())

	return CheckResult{
		Name: 		c.Name(),
		Healthy: 	true,
		Message: 	"agent process is running",
		Metadata: map[string]any{
			"started_at": agentStartedAt,
			"uptime_seconds": uptimeSeconds,
		},
	}
}
