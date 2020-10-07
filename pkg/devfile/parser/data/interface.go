package data

import (
	"github.com/openshift/odo/pkg/devfile/parser/generic/data"
)

// DevfileData is an interface that defines functions for Devfile data operations
type DevfileData interface {
	data.DevfileData
	GetDevfileMetadata()
}
