package main

import (
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "initialize a course's kudos installation",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		common()
		// TODO(synful)
	}
	cmdInit.Run = f
	addAllGlobalFlagsTo(cmdInit.Flags())
	cmdMain.AddCommand(cmdInit)
}
