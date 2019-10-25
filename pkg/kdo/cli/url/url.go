package url

import (
	"fmt"

	projectCmd "github.com/redhat-developer/odo-fork/pkg/kdo/cli/project"
	ktemplates "k8s.io/kubectl/pkg/util/templates"

	kdoutil "github.com/redhat-developer/odo-fork/pkg/kdo/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended url command name
const RecommendedCommandName = "url"

var (
	urlShortDesc = `Expose component to the outside world`
	urlLongDesc  = ktemplates.LongDesc(`Expose component to the outside world.
		
		The URLs that are generated using this command, can be used to access the deployed components from outside the cluster.`)
)

// NewCmdURL returns the top-level url command
func NewCmdURL(name, fullName string) *cobra.Command {
	urlCreateCmd := NewCmdURLCreate(createRecommendedCommandName, kdoutil.GetFullName(fullName, createRecommendedCommandName))
	urlDeleteCmd := NewCmdURLDelete(deleteRecommendedCommandName, kdoutil.GetFullName(fullName, deleteRecommendedCommandName))
	urlListCmd := NewCmdURLList(listRecommendedCommandName, kdoutil.GetFullName(fullName, listRecommendedCommandName))
	urlCmd := &cobra.Command{
		Use:   name,
		Short: urlShortDesc,
		Long:  urlLongDesc,
		Example: fmt.Sprintf("%s\n%s\n%s",
			urlCreateCmd.Example,
			urlDeleteCmd.Example,
			urlListCmd.Example),
	}

	// Add a defined annotation in order to appear in the help menu
	urlCmd.Annotations = map[string]string{"command": "main"}
	urlCmd.SetUsageTemplate(kdoutil.CmdUsageTemplate)
	urlCmd.AddCommand(urlCreateCmd, urlDeleteCmd, urlListCmd)

	//Adding `--project` flag
	projectCmd.AddProjectFlag(urlListCmd)
	projectCmd.AddProjectFlag(urlCreateCmd)
	projectCmd.AddProjectFlag(urlDeleteCmd)

	//Adding `--application` flag
	// appCmd.AddApplicationFlag(urlListCmd)
	// appCmd.AddApplicationFlag(urlDeleteCmd)
	// appCmd.AddApplicationFlag(urlCreateCmd)

	//Adding `--component` flag
	// componentCmd.AddComponentFlag(urlDeleteCmd)
	// componentCmd.AddComponentFlag(urlListCmd)
	// componentCmd.AddComponentFlag(urlCreateCmd)

	return urlCmd
}
