package main

import (
	"github.com/openshift/odo/common"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{Use: "kdo"}
	rootCmd.AddCommand(common.CmdPrintKdo)
	rootCmd.Execute()
}
