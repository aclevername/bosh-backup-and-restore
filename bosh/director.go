package bosh

import (
	"strings"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

func New(boshDirector director.Director,
	sshOptsGenerator SSHOptsGenerator,
	connectionFactory SSHConnectionFactory) backuper.BoshDirector {
	return client{
		Director:             boshDirector,
		SSHOptsGenerator:     sshOptsGenerator,
		SSHConnectionFactory: connectionFactory,
	}
}

//go:generate counterfeiter -o fakes/fake_opts_generator.go . SSHOptsGenerator
type SSHOptsGenerator func(uuidGen uuid.Generator) (director.SSHOpts, string, error)

//go:generate counterfeiter -o fakes/fake_ssh_connection_factory.go . SSHConnectionFactory
type SSHConnectionFactory func(host, user, privateKey string) (SSHConnection, error)

type client struct {
	director.Director
	SSHOptsGenerator
	SSHConnectionFactory
}

func (c client) FindInstances(deploymentName string) (backuper.Instances, error) {
	deployment, err := c.Director.FindDeployment(deploymentName)
	if err != nil {
		return nil, err
	}
	vms, err := deployment.VMInfos()
	if err != nil {
		return nil, err
	}
	sshOpts, privateKey, err := c.SSHOptsGenerator(uuid.NewGenerator())
	if err != nil {
		return nil, err
	}
	instances := backuper.Instances{}

	for _, jobName := range uniqueJobsFromVMs(vms) {
		allVmInstances, err := director.NewAllOrPoolOrInstanceSlugFromString(jobName)
		if err != nil {
			return nil, err
		}
		sshRes, err := deployment.SetUpSSH(allVmInstances, sshOpts)
		if err != nil {
			return nil, err
		}
		for _, host := range sshRes.Hosts {
			var sshConnection SSHConnection
			var err error
			sshConnection, err = c.SSHConnectionFactory(defaultToSSHPort(host.Host), host.Username, privateKey)

			if err != nil {
				return nil, err
			}
			instances = append(instances, NewBoshInstance(jobName, host.IndexOrID, sshConnection, deployment))
		}
	}

	return instances, nil
}

func defaultToSSHPort(host string) string {
	parts := strings.Split(host, ":")
	if len(parts) == 2 {
		return host
	} else {
		return host + ":22"
	}
}

func uniqueJobsFromVMs(vms []director.VMInfo) []string {
	var jobs []string
	for _, vm := range vms {
		if !contains(jobs, vm.JobName) {
			jobs = append(jobs, vm.JobName)
		}
	}
	return jobs
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
