package orchestrator

//go:generate counterfeiter -o fakes/fake_deployment_manager.go . DeploymentManager
type DeploymentManager interface {
	Find(deploymentName string) (Deployment, error)
}

func NewBoshDeploymentManager(boshDirector BoshClient, logger Logger) DeploymentManager {
	return &BoshDeploymentManager{BoshClient: boshDirector, Logger: logger}
}

type BoshDeploymentManager struct {
	BoshClient
	Logger
}

func (b *BoshDeploymentManager) Find(deploymentName string) (Deployment, error) {
	instances, err := b.FindInstances(deploymentName)
	return NewBoshDeployment(b.Logger, instances), err
}
