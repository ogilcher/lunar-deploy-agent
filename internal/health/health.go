package health

type CheckResult struct {
	Name     string         `json:"name"`
	Healthy  bool           `json:"healthy"`
	Message  string         `json:"message"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Report struct {
	Status string        `json:"status"`
	Checks []CheckResult `json:"checks"`
}

func BuildReport(
	checks []CheckResult,
) Report {
	status := "healthy"

	for _, check := range checks {
		if !check.Healthy {
			status = "unhealthy"
			break
		}
	}

	return Report{
		Status: status,
		Checks: checks,
	}
}
