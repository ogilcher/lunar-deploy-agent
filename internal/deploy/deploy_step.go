package deploy

type DeployStep interface {
	Run() error
	Name() string
}
