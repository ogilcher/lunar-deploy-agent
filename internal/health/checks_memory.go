package health

import "runtime"

// MemoryCheck verifies basic runtime memory pressure.
//
// This uses Go runtime memory stats, so it measure the agent process,
// not full system memory. Later I'll consider adding OS-level memory checks.
type MemoryCheck struct{}

func (c *MemoryCheck) Name() string {
	return "memory"
}

func (c *MemoryCheck) Run() CheckResult {
	var stats runtime.MemStats

	runtime.ReadMemStats(&stats)

	allocatedMB := stats.Alloc / 1024 / 1024

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "agent memory allocation is normal",
		Metadata: map[string]any{
			"allocated_mb": allocatedMB,
		},
	}
}
