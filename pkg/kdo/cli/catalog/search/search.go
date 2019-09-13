package search

import (
	"fmt"

	"github.com/redhat-developer/odo-fork/pkg/kdo/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended command name
const RecommendedCommandName = "search"

// NewCmdCatalogSearch implements the udo catalog search command
func NewCmdCatalogSearch(name, fullName string) *cobra.Command {
	component := NewCmdCatalogSearchIDP(idpRecommendedCommandName, util.GetFullName(fullName, idpRecommendedCommandName))
	catalogSearchCmd := &cobra.Command{
		Use:   name,
		Short: "Search available Iterative-Dev Packs",
		Long: `Search available  Iterative-Dev Packs..

This searches for a partial match for the given search term in all the available
Iterative-Dev Packs.
`,
		Example: fmt.Sprintf("%s\n", component.Example),
	}
	catalogSearchCmd.AddCommand(component)

	return catalogSearchCmd
}
