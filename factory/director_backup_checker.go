package factory

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orderer"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/standalone"
)

func BuildDirectorBackupChecker(host, username, privateKeyPath string, hasDebug bool) *orchestrator.BackupChecker {
	logger := BuildLogger(hasDebug)
	deploymentManager := standalone.NewDeploymentManager(logger,
		host,
		username,
		privateKeyPath,
		instance.NewJobFinder(logger),
		ssh.NewSshRemoteRunner,
	)

	return orchestrator.NewBackupChecker(logger, deploymentManager, orderer.NewDirectorLockOrderer())
}
