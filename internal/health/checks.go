package health

func RunChecks() Report {
	registry := NewRegistry()

	registry.Register(
		&QueueCheck{},
	)

	registry.Register(
		&PM2Check{},
	)

	registry.Register(
		&GitCheck{},
	)

	return registry.Run()
}
