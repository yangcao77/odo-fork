package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// CmdPrintUdo External command
var CmdPrintUdo = &cobra.Command{
	Use:   "printudo [string to print]",
	Short: "Print anything to the screen",
	Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Args: cobra.MinimumNArgs(1),
	Run:  udoPrint,
}

var udoPrint = func(cmd *cobra.Command, args []string) {
	fmt.Println("udo: " + strings.Join(args, " "))
}
