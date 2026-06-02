package context

// BuildEnvironmentVariables converts deployed context
// into shell environment variables for deployment steps.
func BuildEnvironmentVariables(
	context DeploymentContext,
) []string {
	return []string{
		"LUNAR_DEPLOYMENT=" + context.DeploymentName,
		"LUNAR_ENVIRONMENT=" + context.Environment,
		"LUNAR_REPOSITORY_PATH=" + context.RepositoryPath,
	}
}
