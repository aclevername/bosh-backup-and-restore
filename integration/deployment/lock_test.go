package deployment

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/testcluster"
	"github.com/pivotal-cf-experimental/cf-webmock/mockbosh"
	"github.com/pivotal-cf-experimental/cf-webmock/mockhttp"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"time"

	"regexp"

	"strings"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unlock", func() {
	var director *mockhttp.Server
	var backupWorkspace string
	var session *gexec.Session
	var stdin io.WriteCloser
	var deploymentName string
	var downloadManifest bool
	var waitForBackupToFinish bool
	var verifyMocks bool
	var instance1 *testcluster.Instance
	manifest := `---
instance_groups:
- name: redis-dedicated-node
  instances: 1
  jobs:
  - name: redis
    release: redis
  - name: redis-writer
    release: redis
  - name: redis-broker
    release: redis
- name: redis-broker
  instances: 1
  jobs:
  - name: redis
    release: redis
  - name: redis-writer
    release: redis
  - name: redis-broker
    release: redis
`

	possibleBackupDirectories := func() []string {
		dirs, err := ioutil.ReadDir(backupWorkspace)
		Expect(err).NotTo(HaveOccurred())
		backupDirectoryPattern := regexp.MustCompile(`\b` + deploymentName + `_(\d){8}T(\d){6}Z\b`)

		matches := []string{}
		for _, dir := range dirs {
			dirName := dir.Name()
			if backupDirectoryPattern.MatchString(dirName) {
				matches = append(matches, dirName)
			}
		}
		return matches
	}

	backupDirectory := func() string {
		matches := possibleBackupDirectories()

		Expect(matches).To(HaveLen(1), "backup directory not found")
		return path.Join(backupWorkspace, matches[0])
	}

	metadataFile := func() string {
		return path.Join(backupDirectory(), "metadata")
	}

	artifactFile := func(name string) string {
		return path.Join(backupDirectory(), name)
	}

	BeforeEach(func() {
		deploymentName = "my-little-deployment"
		downloadManifest = false
		waitForBackupToFinish = true
		verifyMocks = true
		director = mockbosh.NewTLS()
		director.ExpectedBasicAuth("admin", "admin")
		var err error
		backupWorkspace, err = ioutil.TempDir(".", "backup-workspace-")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if verifyMocks {
			director.VerifyMocks()
		}
		director.Close()

		instance1.DieInBackground()
		Expect(os.RemoveAll(backupWorkspace)).To(Succeed())
	})

	JustBeforeEach(func() {
		env := []string{"BOSH_CLIENT_SECRET=admin"}

		params := []string{
			"deployment",
			"--ca-cert", sslCertPath,
			"--username", "admin",
			"--target", director.URL,
			"--deployment", deploymentName,
			"--debug",
			"unlock"}

		if waitForBackupToFinish {
			session = binary.Run(backupWorkspace, env, params...)
		} else {
			session, stdin = binary.Start(backupWorkspace, env, params...)
			Eventually(session).Should(gbytes.Say(".+"))
		}
	})

	Context("When there is a deployment which has one instance", func() {
		singleInstanceResponse := func(instanceGroupName string) []mockbosh.VMsOutput {
			return []mockbosh.VMsOutput{
				{
					IPs:     []string{"10.0.0.1"},
					JobName: instanceGroupName,
				},
			}
		}

		Context("and there is a plausible backup script", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				By("creating a dummy backup script")
				instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/backup", `#!/usr/bin/env sh

set -u
touch /tmp/backup-script-was-run
printf "backupcontent1" > $BBR_ARTIFACT_DIRECTORY/backupdump1
printf "backupcontent2" > $BBR_ARTIFACT_DIRECTORY/backupdump2
`)
			})

			Context("and we don't ask for the manifest to be downloaded", func() {
				BeforeEach(func() {
					MockDirectorWith(director,
						mockbosh.Info().WithAuthTypeBasic(),
						VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
						DownloadManifest(deploymentName, manifest),
						SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
						CleanupSSH(deploymentName, "redis-dedicated-node"))
				})

				It("successfully backs up the deployment", func() {
					By("not running non-existent pre-backup scripts")

					By("exiting zero", func() {
						Expect(session.ExitCode()).To(BeZero())
					})

					var redisNodeArchivePath string

					By("creating a backup directory which contains a backup artifact and a metadata file", func() {
						redisNodeArchivePath = artifactFile("redis-dedicated-node-0-redis.tar")
						Expect(backupDirectory()).To(BeADirectory())
						Expect(redisNodeArchivePath).To(BeARegularFile())
						Expect(metadataFile()).To(BeARegularFile())
					})

					By("having successfully run the backup script, using the $BBR_ARTIFACT_DIRECTORY variable", func() {
						archive := OpenTarArchive(redisNodeArchivePath)

						Expect(archive.Files()).To(ConsistOf("backupdump1", "backupdump2"))
						Expect(archive.FileContents("backupdump1")).To(Equal("backupcontent1"))
						Expect(archive.FileContents("backupdump2")).To(Equal("backupcontent2"))
					})

					By("correctly populating the metadata file", func() {
						metadataContents := ParseMetadata(metadataFile())

						currentTimezone, _ := time.Now().Zone()
						Expect(metadataContents.BackupActivityMetadata.StartTime).To(MatchRegexp(`^(\d{4})\/(\d{2})\/(\d{2}) (\d{2}):(\d{2}):(\d{2}) ` + currentTimezone + "$"))
						Expect(metadataContents.BackupActivityMetadata.FinishTime).To(MatchRegexp(`^(\d{4})\/(\d{2})\/(\d{2}) (\d{2}):(\d{2}):(\d{2}) ` + currentTimezone + "$"))

						Expect(metadataContents.InstancesMetadata).To(HaveLen(1))
						Expect(metadataContents.InstancesMetadata[0].InstanceName).To(Equal("redis-dedicated-node"))
						Expect(metadataContents.InstancesMetadata[0].InstanceIndex).To(Equal("0"))

						Expect(metadataContents.InstancesMetadata[0].Artifacts[0].Name).To(Equal("redis"))
						Expect(metadataContents.InstancesMetadata[0].Artifacts[0].Checksums).To(HaveLen(2))
						Expect(metadataContents.InstancesMetadata[0].Artifacts[0].Checksums["./backupdump1"]).To(Equal(ShaFor("backupcontent1")))
						Expect(metadataContents.InstancesMetadata[0].Artifacts[0].Checksums["./backupdump2"]).To(Equal(ShaFor("backupcontent2")))

						Expect(metadataContents.CustomArtifactsMetadata).To(BeEmpty())
					})

					By("printing the backup progress to the screen", func() {
						Expect(session.Out).To(gbytes.Say("INFO - Looking for scripts"))
						Expect(session.Out).To(gbytes.Say("INFO - redis-dedicated-node/fake-uuid/redis/backup"))
						Expect(session.Out).To(gbytes.Say(fmt.Sprintf("INFO - Running pre-checks for backup of %s...", deploymentName)))
						Expect(session.Out).To(gbytes.Say(fmt.Sprintf("INFO - Starting backup of %s...", deploymentName)))
						Expect(session.Out).To(gbytes.Say("INFO - Running pre-backup-lock scripts..."))
						Expect(session.Out).To(gbytes.Say("INFO - Finished running pre-backup-lock scripts."))
						Expect(session.Out).To(gbytes.Say("INFO - Running backup scripts..."))
						Expect(session.Out).To(gbytes.Say("INFO - Backing up redis on redis-dedicated-node/fake-uuid..."))
						Expect(session.Out).To(gbytes.Say("INFO - Finished running backup scripts."))
						Expect(session.Out).To(gbytes.Say("INFO - Running post-backup-unlock scripts..."))
						Expect(session.Out).To(gbytes.Say("INFO - Finished running post-backup-unlock scripts."))
						Expect(session.Out).To(gbytes.Say("INFO - Copying backup -- [^-]*-- from redis-dedicated-node/fake-uuid..."))
						Expect(session.Out).To(gbytes.Say("INFO - Finished copying backup -- from redis-dedicated-node/fake-uuid..."))
						Expect(session.Out).To(gbytes.Say("INFO - Starting validity checks -- from redis-dedicated-node/fake-uuid..."))
						Expect(session.Out).To(gbytes.Say(`DEBUG - Calculating shasum for local file ./backupdump[12]`))
						Expect(session.Out).To(gbytes.Say(`DEBUG - Calculating shasum for local file ./backupdump[12]`))
						Expect(session.Out).To(gbytes.Say("DEBUG - Calculating shasum for remote files"))
						Expect(session.Out).To(gbytes.Say("DEBUG - Comparing shasums"))
						Expect(session.Out).To(gbytes.Say("INFO - Finished validity checks -- from redis-dedicated-node/fake-uuid..."))
					})

					By("cleaning up backup artifacts from the remote", func() {
						Expect(instance1.FileExists("/var/vcap/store/bbr-backup")).To(BeFalse())
					})
				})

				Context("and the pre-backup-lock script is present", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/pre-backup-lock-script-was-run
`)
						instance1.CreateScript("/var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock", ``)
					})

					It("executes and logs the locks", func() {
						By("running the pre-backup-lock script", func() {
							Expect(instance1.FileExists("/tmp/pre-backup-lock-script-was-run")).To(BeTrue())
						})

						By("logging that it is locking the instance, and listing the scripts", func() {
							assertOutput(session, []string{
								`Locking redis on redis-dedicated-node/fake-uuid for backup`,
								"> /var/vcap/jobs/redis/bin/bbr/pre-backup-lock",
								"> /var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock",
							})
						})
					})

				})

				Context("when the pre-backup-lock script fails", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
echo 'ultra-bar'
(>&2 echo 'ultra-baz')
touch /tmp/pre-backup-lock-output
exit 1
`)
						instance1.CreateScript("/var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock", ``)
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/post-backup-unlock-output
`)
					})

					It("logs the failure, and unlocks the system", func() {
						By("running the pre-backup-lock scripts", func() {
							Expect(instance1.FileExists("/tmp/pre-backup-lock-output")).To(BeTrue())
						})

						By("not running the backup script", func() {
							Expect(instance1.FileExists("/tmp/backup-script-was-run")).NotTo(BeTrue())
						})

						By("exiting with the correct error code", func() {
							Expect(session.ExitCode()).To(Equal(4))
						})

						By("logging the error", func() {
							Expect(session.Err.Contents()).To(ContainSubstring(
								"Error attempting to run pre-backup-lock for job redis on redis-dedicated-node/fake-uuid"))
						})

						By("logging stderr", func() {
							Expect(session.Err.Contents()).To(ContainSubstring("ultra-baz"))
						})

						By("also running the post-backup-unlock scripts", func() {
							Expect(instance1.FileExists("/tmp/post-backup-unlock-output")).To(BeTrue())
						})

						By("not printing a recommendation to run bbr backup-cleanup", func() {
							Expect(string(session.Err.Contents())).NotTo(ContainSubstring(
								"It is recommended that you run `bbr backup-cleanup`"))
						})
					})
				})

				Context("when deployment has a post-backup-unlock script", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/post-backup-unlock-script-was-run
echo "Unlocking release"`)
					})

					It("prints unlock progress to the screen", func() {
						By("runs the pre-backup-lock scripts", func() {
							Expect(instance1.FileExists("/tmp/post-backup-unlock-script-was-run")).To(BeTrue())
						})

						By("logging the script action", func() {
							assertOutput(session, []string{
								"Unlocking redis on redis-dedicated-node/fake-uuid",
								"Finished unlocking redis on redis-dedicated-node/fake-uuid.",
							})
						})
					})

				})

				Context("when the post backup unlock script fails", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
echo 'ultra-bar'
(>&2 echo 'ultra-baz')
exit 1`)
					})

					It("exits and prints the error", func() {
						By("exits with the correct error code", func() {
							Expect(session).To(gexec.Exit(8))
						})

						By("prints stderr", func() {
							Expect(session.Err.Contents()).To(ContainSubstring("ultra-baz"))
						})

						By("prints an error", func() {
							Expect(session.Err.Contents()).To(ContainSubstring("Error attempting to run post-backup-unlock for job redis on redis-dedicated-node/fake-uuid"))
						})

						By("printing a recommendation to run bbr backup-cleanup", func() {
							Expect(string(session.Err.Contents())).To(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
						})
					})
				})

				Context("but /var/vcap/store is not world-accessible", func() {
					BeforeEach(func() {
						instance1.Run("sudo", "chmod", "700", "/var/vcap/store")
					})

					It("successfully backs up the deployment", func() {
						Expect(session.ExitCode()).To(BeZero())
					})
				})
			})
		})

		Context("when there are multiple plausible backup scripts", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				By("creating a dummy backup script")
				instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/backup", `#!/usr/bin/env sh

set -u

printf "backupcontent1" > $BBR_ARTIFACT_DIRECTORY/backupdump1
printf "backupcontent2" > $BBR_ARTIFACT_DIRECTORY/backupdump2
`)

				By("creating a dummy backup script")
				instance1.CreateScript("/var/vcap/jobs/redis-broker/bin/bbr/backup", `#!/usr/bin/env sh

set -u

printf "backupcontent1" > $BBR_ARTIFACT_DIRECTORY/backupdump1
printf "backupcontent2" > $BBR_ARTIFACT_DIRECTORY/backupdump2
`)

				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSH(deploymentName, "redis-dedicated-node"))
			})

			It("successfully backs up the deployment", func() {
				By("exiting zero", func() {
					Expect(session.ExitCode()).To(BeZero())
				})

				var redisNodeArchivePath, brokerArchivePath string
				By("creating a backup directory which contains the backup artifacts and a metadata file", func() {
					Expect(backupDirectory()).To(BeADirectory())
					redisNodeArchivePath = artifactFile("redis-dedicated-node-0-redis.tar")
					brokerArchivePath = artifactFile("redis-dedicated-node-0-redis-broker.tar")
					Expect(redisNodeArchivePath).To(BeARegularFile())
					Expect(brokerArchivePath).To(BeARegularFile())
					Expect(metadataFile()).To(BeARegularFile())
				})

				By("including the backup files from the instance", func() {
					redisNodeArchive := OpenTarArchive(redisNodeArchivePath)
					Expect(redisNodeArchive.Files()).To(ConsistOf("backupdump1", "backupdump2"))
					Expect(redisNodeArchive.FileContents("backupdump1")).To(Equal("backupcontent1"))
					Expect(redisNodeArchive.FileContents("backupdump2")).To(Equal("backupcontent2"))

					brokerArchive := OpenTarArchive(brokerArchivePath)
					Expect(brokerArchive.Files()).To(ConsistOf("backupdump1", "backupdump2"))
					Expect(brokerArchive.FileContents("backupdump1")).To(Equal("backupcontent1"))
					Expect(brokerArchive.FileContents("backupdump2")).To(Equal("backupcontent2"))
				})

				By("correctly populating the metadata file", func() {
					metadataContents := ParseMetadata(metadataFile())

					currentTimezone, _ := time.Now().Zone()
					Expect(metadataContents.BackupActivityMetadata.StartTime).To(MatchRegexp(`^(\d{4})\/(\d{2})\/(\d{2}) (\d{2}):(\d{2}):(\d{2}) ` + currentTimezone + "$"))
					Expect(metadataContents.BackupActivityMetadata.FinishTime).To(MatchRegexp(`^(\d{4})\/(\d{2})\/(\d{2}) (\d{2}):(\d{2}):(\d{2}) ` + currentTimezone + "$"))

					Expect(metadataContents.InstancesMetadata).To(HaveLen(1))
					Expect(metadataContents.InstancesMetadata[0].InstanceName).To(Equal("redis-dedicated-node"))
					Expect(metadataContents.InstancesMetadata[0].InstanceIndex).To(Equal("0"))

					redisArtifact := metadataContents.InstancesMetadata[0].FindArtifact("redis")
					Expect(redisArtifact.Name).To(Equal("redis"))
					Expect(redisArtifact.Checksums).To(HaveLen(2))
					Expect(redisArtifact.Checksums["./backupdump1"]).To(Equal(ShaFor("backupcontent1")))
					Expect(redisArtifact.Checksums["./backupdump2"]).To(Equal(ShaFor("backupcontent2")))

					brokerArtifact := metadataContents.InstancesMetadata[0].FindArtifact("redis-broker")
					Expect(brokerArtifact.Name).To(Equal("redis-broker"))
					Expect(brokerArtifact.Checksums).To(HaveLen(2))
					Expect(brokerArtifact.Checksums["./backupdump1"]).To(Equal(ShaFor("backupcontent1")))
					Expect(brokerArtifact.Checksums["./backupdump2"]).To(Equal(ShaFor("backupcontent2")))

					Expect(metadataContents.CustomArtifactsMetadata).To(BeEmpty())
				})

				By("cleaning up backup artifacts from the remote", func() {
					Expect(instance1.FileExists("/var/vcap/store/bbr-backup")).To(BeFalse())
				})
			})
		})

		Context("when a deployment can't be backed up", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSH(deploymentName, "redis-dedicated-node"),
				)

				instance1.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/ctl",
				)
			})
			It("exits and displays a message", func() {
				Expect(session.ExitCode()).NotTo(BeZero(), "returns a non-zero exit code")
				Expect(string(session.Err.Contents())).To(ContainSubstring("Deployment '"+deploymentName+"' has no backup scripts"),
					"prints an error")
				Expect(possibleBackupDirectories()).To(HaveLen(0), "does not create a backup on disk")

				By("not printing a recommendation to run bbr backup-cleanup", func() {
					Expect(string(session.Err.Contents())).NotTo(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
				})
			})
		})

		Context("when the instance backup script fails", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSH(deploymentName, "redis-dedicated-node"),
				)

				instance1.CreateScript(
					"/var/vcap/jobs/redis/bin/bbr/backup", "echo 'ultra-bar'; (>&2 echo 'ultra-baz'); exit 1",
				)

				instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/post-backup-unlock-script-was-run
echo "Unlocking release"`)
			})

			It("errors and exits gracefully", func() {
				By("returning exit code 1", func() {
					Expect(session.ExitCode()).To(Equal(1))
				})
				By("running the the post-backup-unlock scripts", func() {
					Expect(instance1.FileExists("/tmp/post-backup-unlock-script-was-run")).To(BeTrue())
				})

				By("not printing a recommendation to run bbr backup-cleanup", func() {
					Expect(string(session.Err.Contents())).NotTo(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
				})
			})

		})

		Context("when both the instance backup script and cleanup fail", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSHFails(deploymentName, "redis-dedicated-node", "ultra-foo"),
				)

				instance1.CreateScript(
					"/var/vcap/jobs/redis/bin/bbr/backup", "(>&2 echo 'ultra-baz'); exit 1",
				)
			})

			It("exits correctly and prints an error", func() {
				By("returning exit code 17 (16 + 1)", func() {
					Expect(session.ExitCode()).To(Equal(17))
				})

				By("printing an error", func() {
					assertErrorOutput(session, []string{
						"Error attempting to run backup for job redis on redis-dedicated-node/fake-uuid",
						"ultra-baz",
						"ultra-foo",
					})
				})

				By("printing a recommendation to run bbr backup-cleanup", func() {
					Expect(string(session.Err.Contents())).To(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
				})
			})

		})

		Context("when backup succeeds but cleanup fails", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSHFails(deploymentName, "redis-dedicated-node", "Can't do it mate"),
				)

				instance1.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)
			})

			It("exits correctly and prints the error", func() {
				By("returning the correct error code", func() {
					Expect(session.ExitCode()).To(Equal(16))
				})

				By("printing an error", func() {
					Expect(string(session.Err.Contents())).To(ContainSubstring("Deployment '" + deploymentName + "' failed while cleaning up with error: "))
				})

				By("including the failure message in error output", func() {
					Expect(string(session.Err.Contents())).To(ContainSubstring("Can't do it mate"))
				})

				By("creating a backup on disk", func() {
					Expect(backupDirectory()).To(BeADirectory())
				})

				By("printing a recommendation to run bbr backup-cleanup", func() {
					Expect(string(session.Err.Contents())).To(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
				})
			})

		})

		Context("when running the metadata script does not give valid yml", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata", `#!/usr/bin/env sh
touch /tmp/metadata-script-was-run-but-produces-invalid-yaml
echo "not valid yaml
"`)

				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
					DownloadManifest(deploymentName, manifest),
					SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1),
					CleanupSSH(deploymentName, "redis-dedicated-node"),
				)
			})

			It("attempts to use the metadata, and exits with an error", func() {
				By("running the metadata scripts", func() {
					Expect(instance1.FileExists("/tmp/metadata-script-was-run-but-produces-invalid-yaml")).To(BeTrue())
				})

				By("exiting with the correct error code", func() {
					Expect(session).To(gexec.Exit(1))
				})

				By("not printing a recommendation to run bbr backup-cleanup", func() {
					Expect(string(session.Err.Contents())).NotTo(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
				})
			})

		})
	})

	Context("When there is a deployment which has two instances", func() {
		twoInstancesResponse := func(firstInstanceGroupName, secondInstanceGroupName string) []mockbosh.VMsOutput {

			return []mockbosh.VMsOutput{
				{
					IPs:     []string{"10.0.0.1"},
					JobName: firstInstanceGroupName,
				},
				{
					IPs:     []string{"10.0.0.2"},
					JobName: secondInstanceGroupName,
				},
			}
		}

		Context("one backupable", func() {
			var firstReturnedInstance, secondReturnedInstance *testcluster.Instance

			BeforeEach(func() {
				deploymentName = "my-bigger-deployment"
				firstReturnedInstance = testcluster.NewInstance()
				secondReturnedInstance = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, twoInstancesResponse("redis-dedicated-node", "redis-broker")),
					DownloadManifest(deploymentName, manifest),
					append(SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, firstReturnedInstance),
						SetupSSH(deploymentName, "redis-broker", "fake-uuid-2", 0, secondReturnedInstance)...),
					append(CleanupSSH(deploymentName, "redis-dedicated-node"),
						CleanupSSH(deploymentName, "redis-broker")...),
				)
				firstReturnedInstance.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)
			})

			AfterEach(func() {
				firstReturnedInstance.DieInBackground()
				secondReturnedInstance.DieInBackground()
			})

			It("backs up deployment successfully", func() {
				Expect(session.ExitCode()).To(BeZero())
				Expect(backupDirectory()).To(BeADirectory())
				Expect(path.Join(backupDirectory(), "/redis-dedicated-node-0-redis.tar")).To(BeARegularFile())
				Expect(path.Join(backupDirectory(), "/redis-broker-0-redis.tar")).ToNot(BeAnExistingFile())
			})

			Context("with ordering on pre-backup-lock specified", func() {
				BeforeEach(func() {
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
					secondReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 0`)
					secondReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
				})

				It("locks in the specified order", func() {
					redisLockTime := firstReturnedInstance.GetCreatedTime("/tmp/redis-pre-backup-lock-called")
					redisWriterLockTime := secondReturnedInstance.GetCreatedTime("/tmp/redis-writer-pre-backup-lock-called")

					Expect(string(session.Out.Contents())).To(ContainSubstring("Detected order: redis-writer should be locked before redis/redis during backup"))

					Expect(redisWriterLockTime < redisLockTime).To(BeTrue(), fmt.Sprintf(
						"Writer locked at %s, which is after the server locked (%s)",
						strings.TrimSuffix(redisWriterLockTime, "\n"),
						strings.TrimSuffix(redisLockTime, "\n")))

				})
			})

			Context("with ordering on pre-backup-lock (where the default ordering would unlock in the wrong order)",
				func() {
					BeforeEach(func() {
						secondReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
						firstReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 0`)
						secondReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/redis-post-backup-unlock-called
exit 0`)
						firstReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis-writer/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/redis-writer-post-backup-unlock-called
exit 0`)
						firstReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
							`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
					})

					It("unlocks in the right order", func() {
						By("unlocking the redis job before unlocking the redis-writer job")
						redisUnlockTime := secondReturnedInstance.GetCreatedTime("/tmp/redis-post-backup-unlock-called")
						redisWriterUnlockTime := firstReturnedInstance.GetCreatedTime("/tmp/redis-writer-post-backup-unlock-called")

						Expect(redisUnlockTime < redisWriterUnlockTime).To(BeTrue(), fmt.Sprintf(
							"Writer unlocked at %s, which is before the server unlocked (%s)",
							strings.TrimSuffix(redisWriterUnlockTime, "\n"),
							strings.TrimSuffix(redisUnlockTime, "\n")))
					})
				})

			Context("but the pre-backup-lock ordering is cyclic", func() {
				BeforeEach(func() {
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 0`)
					firstReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
					firstReturnedInstance.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis-writer
  release: redis
"`)
				})

				It("Should fail", func() {
					By("exiting with an error", func() {
						Expect(session).To(gexec.Exit(1))
					})

					By("printing a helpful error message", func() {
						Expect(string(session.Err.Contents())).To(ContainSubstring("job locking dependency graph is cyclic"))
					})

					By("not creating a local backup artifact", func() {
						Expect(possibleBackupDirectories()).To(BeEmpty(),
							"Should quit before creating any local backup artifact.")
					})
				})
			})
		})

		Context("both backupable", func() {
			var backupableInstance1, backupableInstance2 *testcluster.Instance

			BeforeEach(func() {
				deploymentName = "my-two-instance-deployment"
				backupableInstance1 = testcluster.NewInstance()
				backupableInstance2 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, twoInstancesResponse("redis-dedicated-node", "redis-broker")),
					DownloadManifest(deploymentName, manifest),
					append(SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, backupableInstance1),
						SetupSSH(deploymentName, "redis-broker", "fake-uuid-2", 0, backupableInstance2)...),
					append(CleanupSSH(deploymentName, "redis-dedicated-node"),
						CleanupSSH(deploymentName, "redis-broker")...),
				)

				backupableInstance1.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)

				backupableInstance2.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)

			})

			AfterEach(func() {
				backupableInstance1.DieInBackground()
				backupableInstance2.DieInBackground()
			})

			It("backs up both instances and prints process to the screen", func() {
				By("backing up both instances successfully", func() {
					Expect(session.ExitCode()).To(BeZero())
					Expect(backupDirectory()).To(BeADirectory())
					Expect(path.Join(backupDirectory(), "/redis-dedicated-node-0-redis.tar")).To(BeARegularFile())
					Expect(path.Join(backupDirectory(), "/redis-broker-0-redis.tar")).To(BeARegularFile())
				})

				By("printing the backup progress to the screen", func() {
					assertOutput(session, []string{
						fmt.Sprintf("Starting backup of %s...", deploymentName),
						"Backing up redis on redis-dedicated-node/fake-uuid...",
						"Finished backing up redis on redis-dedicated-node/fake-uuid.",
						"Backing up redis on redis-broker/fake-uuid-2...",
						"Finished backing up redis on redis-broker/fake-uuid-2.",
						"Copying backup --",
						"from redis-dedicated-node/fake-uuid...",
						"from redis-broker/fake-uuid-2...",
						"Finished copying backup --",
						fmt.Sprintf("Backup created of %s on", deploymentName),
					})
				})
			})

			Context("and the backup artifact directory already exists on one of them", func() {
				BeforeEach(func() {
					backupableInstance2.CreateDir("/var/vcap/store/bbr-backup")
				})

				It("fails without destroying existing artifact", func() {
					By("failing", func() {
						Expect(session.ExitCode()).NotTo(BeZero())
					})

					By("not deleting the existing backup artifact directory", func() {
						Expect(backupableInstance2.FileExists("/var/vcap/store/bbr-backup")).To(BeTrue())
					})

					By("loging which instance has the extant artifact directory", func() {
						Expect(session.Err).To(gbytes.Say("Directory /var/vcap/store/bbr-backup already exists on instance redis-broker/fake-uuid-2"))
					})
				})
			})
		})

		Context("and both specify the same backup name in their metadata", func() {
			var backupableInstance1, backupableInstance2 *testcluster.Instance

			BeforeEach(func() {
				deploymentName = "my-two-instance-deployment"
				backupableInstance1 = testcluster.NewInstance()
				backupableInstance2 = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, twoInstancesResponse("redis-dedicated-node", "redis-broker")),
					DownloadManifest(deploymentName, manifest),
					append(SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, backupableInstance1),
						SetupSSH(deploymentName, "redis-broker", "fake-uuid-2", 0, backupableInstance2)...),
					append(CleanupSSH(deploymentName, "redis-dedicated-node"),
						CleanupSSH(deploymentName, "redis-broker")...),
				)

				backupableInstance1.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)

				backupableInstance2.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)

				backupableInstance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata", `#!/usr/bin/env sh
echo "---
backup_name: duplicate_name
"`)
				backupableInstance2.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata", `#!/usr/bin/env sh
echo "---
backup_name: duplicate_name
"`)
			})

			AfterEach(func() {
				backupableInstance1.DieInBackground()
				backupableInstance2.DieInBackground()
			})

			It("fails correctly, and doesn't create artifacts", func() {
				By("not creating a file with the duplicated backup name", func() {
					Expect(len(possibleBackupDirectories())).To(Equal(0))
				})

				By("refusing to perform backup", func() {
					Expect(session.Err.Contents()).To(ContainSubstring(
						"Multiple jobs in deployment 'my-two-instance-deployment' specified the same backup name",
					))
				})

				By("returning exit code 1", func() {
					Expect(session.ExitCode()).To(Equal(1))
				})
			})

		})

		Context("and one instance consumes restore custom name, which no instance provides", func() {
			var restoreInstance, backupableInstance *testcluster.Instance

			BeforeEach(func() {
				deploymentName = "my-two-instance-deployment"
				restoreInstance = testcluster.NewInstance()
				backupableInstance = testcluster.NewInstance()
				MockDirectorWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, twoInstancesResponse("redis-dedicated-node", "redis-broker")),
					DownloadManifest(deploymentName, manifest),
					append(SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, restoreInstance),
						SetupSSH(deploymentName, "redis-broker", "fake-uuid-2", 0, backupableInstance)...),
					append(CleanupSSH(deploymentName, "redis-dedicated-node"),
						CleanupSSH(deploymentName, "redis-broker")...),
				)

				restoreInstance.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/restore",
				)

				backupableInstance.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)

				restoreInstance.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata", `#!/usr/bin/env sh
echo "---
restore_name: name_1
"`)
				backupableInstance.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata", `#!/usr/bin/env sh
echo "---
backup_name: name_2
"`)
			})

			AfterEach(func() {
				restoreInstance.DieInBackground()
				backupableInstance.DieInBackground()
			})

			It("doesn't perform a backup", func() {
				By("refusing to perform backup", func() {
					Expect(string(session.Err.Contents())).To(ContainSubstring(
						"The redis-dedicated-node restore script expects a backup script which produces name_1 artifact which is not present in the deployment",
					))
				})
				By("returning exit code 1", func() {
					Expect(session.ExitCode()).To(Equal(1))
				})
			})

		})
	})

	Context("When deployment does not exist", func() {
		BeforeEach(func() {
			deploymentName = "my-non-existent-deployment"
			director.VerifyAndMock(
				mockbosh.Info().WithAuthTypeBasic(),
				mockbosh.VMsForDeployment(deploymentName).NotFound(),
			)
		})

		It("errors and exits", func() {
			By("returning exit code 1", func() {
				Expect(session.ExitCode()).To(Equal(1))
			})

			By("printing an error", func() {
				Expect(string(session.Err.Contents())).To(ContainSubstring("Director responded with non-successful status code"))
			})

			By("not printing a recommendation to run bbr backup-cleanup", func() {
				Expect(string(session.Err.Contents())).NotTo(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
			})
		})

	})
})
