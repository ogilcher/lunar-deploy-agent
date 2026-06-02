package deploy

// DeployStep represents one executable step in a deployment pipeline.
type DeployStep interface {
	Run() error
	Name() string
}
