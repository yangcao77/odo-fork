package catalog

import (
	"fmt"

	"github.com/redhat-developer/odo-fork/pkg/kdo/cli/catalog/list"
	"github.com/redhat-developer/odo-fork/pkg/kdo/cli/catalog/search"
	kdoutil "github.com/redhat-developer/odo-fork/pkg/kdo/util"

	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended catalog command name
const RecommendedCommandName = "catalog"

// NewCmdCatalog implements the odo catalog command
func NewCmdCatalog(name, fullName string) *cobra.Command {
	catalogListCmd := list.NewCmdCatalogList(list.RecommendedCommandName, kdoutil.GetFullName(fullName, list.RecommendedCommandName))
	catalogSearchCmd := search.NewCmdCatalogSearch(search.RecommendedCommandName, kdoutil.GetFullName(fullName, search.RecommendedCommandName))
	catalogCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [options]", name),
		Short: "Catalog related operations",
		Long:  "Catalog related operations",
		Example: fmt.Sprintf("%s\n%s\n",
			catalogListCmd.Example,
			catalogSearchCmd.Example),
	}

	catalogCmd.AddCommand(catalogListCmd, catalogSearchCmd)

	// Add a defined annotation in order to appear in the help menu
	catalogCmd.Annotations = map[string]string{"command": "main"}
	catalogCmd.SetUsageTemplate(kdoutil.CmdUsageTemplate)

	return catalogCmd
}
