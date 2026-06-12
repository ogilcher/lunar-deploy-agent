package health

func RunChecks() Report {
	registry := NewRegistry()

	registry.Register(
		&QueueCheck{},
	)

	return registry.Run()
}
