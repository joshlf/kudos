package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/synful/kudos/lib/build"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kudos",
	Long:  "Print the version number of kudos",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		fmt.Println("kudos version", build.Version)
	}
	cmdVersion.Run = f
	cmdMain.AddCommand(cmdVersion)
}
