package bosh

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cppforlife/go-patch/patch"
	"gopkg.in/yaml.v1"
	"github.com/pkg/errors"
)

type ManifestReleaseMapping struct {
	manifest   interface{}
	v2Manifest bool
}

func (rm ManifestReleaseMapping) FindReleaseName(instanceGroupName, jobName string) (string, error) {
	var releasePath string
	if rm.v2Manifest {
		releasePath = fmt.Sprintf("/instance_groups/name=%s/jobs/name=%s/release", instanceGroupName, jobName)
	} else {
		releasePath = fmt.Sprintf("/jobs/name=%s/templates/name=%s/release", instanceGroupName, jobName)
	}

	releasePointer, err := patch.NewPointerFromString(releasePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error finding release name for job %s in instance group %s", jobName, instanceGroupName))
	}

	release, err := patch.FindOp{Path: releasePointer}.Apply(rm.manifest)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error finding release name for job %s in instance group %s", jobName, instanceGroupName))
	}

	return release.(string), nil
}

func NewBoshManifestReleaseMapping(manifest string) (instance.ReleaseMapping, error) {
	var parsedManifest interface{}

	err := yaml.Unmarshal([]byte(manifest), &parsedManifest)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling manifest yaml")
	}

	v2Manifest := isV2Manifest(parsedManifest)

	return ManifestReleaseMapping{manifest: parsedManifest, v2Manifest: v2Manifest}, nil
}

func isV2Manifest(manifest interface{}) bool {
	instanceGroupPath := patch.MustNewPointerFromString(fmt.Sprintf("/instance_groups"))
	_, err := patch.FindOp{Path: instanceGroupPath}.Apply(manifest)
	if err != nil {
		return false
	}
	return true
}
