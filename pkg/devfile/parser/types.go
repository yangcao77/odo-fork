package parser

import (
	"github.com/openshift/odo/pkg/devfile/parser/data"
	devfileCtx "github.com/openshift/odo/pkg/devfile/parser/generic/context"
)

// DevfileObj is the runtime devfile object
type DevfileObj struct {
	// Ctx has devfile context info
	Ctx devfileCtx.DevfileCtx

	// Data has the devfile data
	Data data.DevfileData
}
