package deployment

import (
	"time"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Suite")
}

var (
	commandPath    string
	boshCaCertPath string
	err            error
)

var fixturesPath = "../../fixtures/redis-backup/"

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)

	//var wg sync.WaitGroup

	//wg.Add(3)

	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//	By("deploying the Redis test release")
	//	RedisDeployment.Deploy()
	//
	//	By("deploying the Redis with metadata")
	//	RedisWithMetadataDeployment.Deploy()
	//}()
	//
	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//	By("deploying the Redis with missing backup script")
	//	RedisWithMissingScriptDeployment.Deploy()
	//
	//	By("deploying the slow backup Redis test release")
	//	RedisSlowBackupDeployment.Deploy()
	//}()

	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//By("deploying the Redis with locking order release")
	//RedisWithLockingOrderDeployment.Deploy()
	//
	//By("deploying the jump box")
	//JumpboxDeployment.Deploy()

	By("deploying the many-bbr-jobs deployment")
	ManyBbrJobsDeployment.Deploy()
	//}()

	//wg.Wait()

	By("building bbr")
	commandPath, err = gexec.Build("github.com/cloudfoundry-incubator/bosh-backup-and-restore/cmd/bbr")
	Expect(err).NotTo(HaveOccurred())
	//
	//By("setting up the jump box")
	//Eventually(JumpboxInstance.RunCommand(
	//	fmt.Sprintf("sudo mkdir %s && sudo chown vcap:vcap %s && sudo chmod 0777 %s", workspaceDir, workspaceDir, workspaceDir))).Should(gexec.Exit(0))
	//
	//By("writing $BOSH_CA_CERT to a temp file")
	//boshCaCertPath, err = WriteEnvVarToTempFile("BOSH_CA_CERT")
	//Expect(err).NotTo(HaveOccurred())
	//
	//By("copying bbr and bosh.crt to the jumpbox")
	//JumpboxInstance.Copy(commandPath, workspaceDir)
	//JumpboxInstance.Copy(boshCaCertPath, workspaceDir+"/bosh.crt")
})

var _ = AfterSuite(func() {
	//var wg sync.WaitGroup
	//
	//wg.Add(3)
	//
	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//	By("tearing down the redis release")
	//	RedisDeployment.Delete()
	//
	//	By("tearing down the other redis release")
	//	RedisWithMetadataDeployment.Delete()
	//}()
	//
	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//	By("tearing down the other redis release")
	//	RedisWithMissingScriptDeployment.Delete()
	//
	//	By("tearing down the slow backup Redis test release")
	//	RedisSlowBackupDeployment.Delete()
	//}()
	//
	//go func() {
	//	defer GinkgoRecover()
	//	defer wg.Done()
	//
	//	By("tearing down the Redis with locking order release")
	//	RedisWithLockingOrderDeployment.Delete()
	//
	//	By("tearing down the jump box")
	//	JumpboxDeployment.Delete()

	By("tearing down the many-bbr-jobs deployment")
	ManyBbrJobsDeployment.Delete()
	//}()

	//wg.Wait()
})

func runOnInstances(instanceCollection map[string][]string, f func(string, string)) {
	for instanceGroup, instances := range instanceCollection {
		for _, instanceIndex := range instances {
			f(instanceGroup, instanceIndex)
		}
	}
}
