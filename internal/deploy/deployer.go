package deploy

type Deployer interface {
	Deploy() error
}
