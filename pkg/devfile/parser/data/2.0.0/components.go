package version200

import (
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
)

// GetMetadata returns the DevfileMetadata Object parsed from devfile
func (d *Devfile200) GetDevfileMetadata() common.DevfileMetadata {
	return d.DevfileMetadata
}
