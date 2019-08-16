package list

import (
	"fmt"

	"github.com/redhat-developer/odo-fork/pkg/kdo/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended command name
const RecommendedCommandName = "list"

// NewCmdCatalogList implements the kdo catalog list command
func NewCmdCatalogList(name, fullName string) *cobra.Command {
	idp := NewCmdCatalogListIDPs(idpRecommendedCommandName, util.GetFullName(fullName, idpRecommendedCommandName))

	catalogListCmd := &cobra.Command{
		Use:     name,
		Short:   "List all available Iterative-Dev packs.",
		Long:    "List all available Iterative-Dev packs from GitHub",
		Example: fmt.Sprintf("%s\n\n", idp.Example),
	}

	catalogListCmd.AddCommand(
		idp,
	)

	return catalogListCmd
}
