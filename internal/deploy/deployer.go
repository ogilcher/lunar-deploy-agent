package deploy

// Deployer defines the behavior required for deployment implementations.
type Deployer interface {
	Deploy() error
}
