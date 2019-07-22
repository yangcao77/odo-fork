package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// CmdPrintKdo External command
var CmdPrintKdo = &cobra.Command{
	Use:   "printkdo [string to print]",
	Short: "Print anything to the screen",
	Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Args: cobra.MinimumNArgs(1),
	Run:  kdoPrint,
}

var kdoPrint = func(cmd *cobra.Command, args []string) {
	fmt.Println("kdo: " + strings.Join(args, " "))
}
