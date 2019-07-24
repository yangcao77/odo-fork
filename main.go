package main

import (
	"github.com/redhat-developer/odo-fork/common"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{Use: "kdo"}
	rootCmd.AddCommand(common.CmdPrintKdo)
	rootCmd.Execute()
}
