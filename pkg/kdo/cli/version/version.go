package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/redhat-developer/odo-fork/pkg/kdo/genericclioptions"
	"github.com/redhat-developer/odo-fork/pkg/kdo/util"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

var (
	// VERSION  is version number that will be displayed when running ./odo version
	VERSION = "v0.0.1"

	// GITCOMMIT is hash of the commit that will be displayed when running ./odo version
	// this will be overwritten when running  build like this: go build -ldflags="-X github.com/openshift/odo/cmd.GITCOMMIT=$(GITCOMMIT)"
	// HEAD is default indicating that this was not set during build
	GITCOMMIT = "HEAD"
)

// RecommendedCommandName is the recommended version command name
const RecommendedCommandName = "version"

var versionLongDesc = ktemplates.LongDesc("Print the client version information")

var versionExample = ktemplates.Examples(`
# Print the client version of odo
%[1]s`,
)

// VersionOptions encapsulates all options for odo version command
type VersionOptions struct {
	// clientFlag indicates if the user only wants client information
	clientFlag bool
}

// NewVersionOptions creates a new VersionOptions instance
func NewVersionOptions() *VersionOptions {
	return &VersionOptions{}
}

// Complete completes VersionOptions after they have been created
func (o *VersionOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	return
}

// Validate validates the VersionOptions based on completed values
func (o *VersionOptions) Validate() (err error) {
	return
}

// Run contains the logic for the odo version command
func (o *VersionOptions) Run() (err error) {
	// If verbose mode is enabled, dump all KUBECLT_* env variables
	// this is usefull for debuging oc plugin integration
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "KUBECTL_") {
			glog.V(4).Info(v)
		}
	}

	fmt.Println("udo " + VERSION + " (" + GITCOMMIT + ")")

	return
}

// NewCmdVersion implements the version odo command
func NewCmdVersion(name, fullName string) *cobra.Command {
	o := NewVersionOptions()
	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:     name,
		Short:   versionLongDesc,
		Long:    versionLongDesc,
		Example: fmt.Sprintf(versionExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	// Add a defined annotation in order to appear in the help menu
	versionCmd.Annotations = map[string]string{"command": "utility"}
	versionCmd.SetUsageTemplate(util.CmdUsageTemplate)
	versionCmd.Flags().BoolVar(&o.clientFlag, "client", false, "Client version only (no server required).")

	return versionCmd
}
