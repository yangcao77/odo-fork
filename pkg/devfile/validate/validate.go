package validate

import (
	"fmt"

	"github.com/openshift/odo/pkg/devfile/parser/generic/data/common"
	"k8s.io/klog"

	v200 "github.com/openshift/odo/pkg/devfile/parser/generic/data/2.0.0"
)

// ValidateDevfileData validates whether sections of devfile are odo compatible
func ValidateDevfileData(data interface{}) error {
	var components []common.DevfileComponent
	var commandsMap map[string]common.DevfileCommand
	var events common.DevfileEvents

	switch d := data.(type) {
	case *v200.Devfile200:
		components = d.GetComponents()
		commandsMap = d.GetCommands()
		events = d.GetEvents()

		// Validate all the devfile components before validating commands
		if err := validateComponents(components); err != nil {
			return err
		}

		// Validate all the devfile commands before validating events
		if err := validateCommands(d.Commands, commandsMap, components); err != nil {
			return err
		}

		// Validate all the events
		if err := validateEvents(events, commandsMap, components); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown devfile type %T", d)
	}

	// Successful
	klog.V(2).Info("Successfully validated devfile sections")
	return nil

}
