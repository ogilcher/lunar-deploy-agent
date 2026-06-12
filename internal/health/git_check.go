package health

import "os/exec"

// GitCheck verifies that the git command is installed and reachable.
type GitCheck struct{}

func (c *GitCheck) Name() string {
	return "git"
}

func (c *GitCheck) Run() CheckResult {
	command := exec.Command("git", "--version")

	output, err := command.CombinedOutput()

	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Healthy: false,
			Message: "git is not installed or not reachable",
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Healthy: true,
		Message: "git available: " + string(output),
	}
}
