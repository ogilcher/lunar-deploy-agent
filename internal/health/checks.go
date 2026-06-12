package health

func RunChecks(
	configPath string,
) Report {
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

	registry.Register(
		&ConfigCheck{
			ConfigPath: configPath,
		},
	)

	registry.Register(
		&DiskCheck{
			Path: ".",
		},
	)

	registry.Register(
		&MemoryCheck{},
	)

	registry.Register(
		&CapabilitiesCheck{
			ConfigPath: configPath,
		},
	)

	registry.Register(
		&UptimeCheck{},
	)

	return registry.Run()
}
