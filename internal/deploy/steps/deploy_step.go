package steps

import "github.com/ogilcher/lunar-deploy-agent/internal/deploy/context"

// DeployStep represents one executable step in a deployment pipeline.
type DeployStep interface {
	Run(context context.DeploymentContext) error
	Name() string
}
