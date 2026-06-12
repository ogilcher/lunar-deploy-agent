package health

type Registry struct {
	checks []Check
}

func NewRegistry() *Registry {
	return &Registry{
		checks: []Check{},
	}
}

func (r *Registry) Register(
	check Check,
) {
	r.checks = append(r.checks, check)
}

func (r *Registry) Run() Report {
	results := []CheckResult{}

	for _, check := range r.checks {
		results = append(
			results,
			check.Run(),
		)
	}

	return BuildReport(results)
}
