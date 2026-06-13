package health

import "syscall"

// DiskCheck verifies that the current filesystem has enough free space.
type DiskCheck struct {
	Path string
}

func (c *DiskCheck) Name() string {
	return "disk"
}

func (c *DiskCheck) Run() CheckResult {
	path := c.Path

	if path == "" {
		path = "."
	}

	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: err.Error(),
		}
	}

	availableBytes := stat.Bavail * uint64(stat.Bsize)
	availableGB := availableBytes / 1024 / 1024 / 1024

	if availableGB < 1 {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: "available disk space is below 1GB",
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "available disk space is sufficient",
	}
}
