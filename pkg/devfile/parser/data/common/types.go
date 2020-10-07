package common

import "github.com/openshift/odo/pkg/devfile/parser/generic/data/common"

// DevfileMetadata contains odo specific metadata for devfile
type DevfileMetadata struct {
	*common.Metadata

	AlphaBuildDockerfile string `json:"alpha.build-dockerfile,omitempty" yaml:"alpha.build-dockerfile,omitempty"`

	AlphaDeploymentManifest string `json:"alpha.deployment-manifest,omitempty" yaml:"alpha.deployment-manifest,omitempty"`
}
