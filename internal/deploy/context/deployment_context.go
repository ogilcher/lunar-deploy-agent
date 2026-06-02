package context

// DeploymentContext contains runtime deployment information
// shared across deployment steps.
type DeploymentContext struct {
	DeploymentName string
	RepositoryPath string
	Environment    string
}
