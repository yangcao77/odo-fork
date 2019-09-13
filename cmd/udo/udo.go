package main

import (
	"flag"
	"os"

	"github.com/redhat-developer/odo-fork/pkg/kdo/cli"
)

func main() {
	var root = cli.NewCmdUdo(cli.UdoRecommendedName, cli.UdoRecommendedName)

	// override usage so that flag.Parse uses root command's usage instead of default one when invoked with -h
	root.Flags().AddGoFlagSet(flag.CommandLine)
	flag.Usage = func() {
		_ = root.Help()
	}

	// parse the flags but hack around to avoid exiting with error code 2 on help
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	args := os.Args[1:]
	if err := flag.CommandLine.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
	}

	root.Execute()
}
