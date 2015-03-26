package main

import "github.com/spf13/cobra"

var cmdHandin = &cobra.Command{
	Use:   "handin",
	Short: "hand in an assignment",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		common()
		// TODO(synful)
	}
	cmdHandin.Run = f
	addAllGlobalFlagsTo(cmdHandin.Flags())
	cmdMain.AddCommand(cmdHandin)
}
