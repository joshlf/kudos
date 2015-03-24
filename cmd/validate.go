package main

import "github.com/spf13/cobra"

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "validate the configuration",
	Long: `Perform all of the steps necessary to determine kudos' configuration, 
reporting any errors encountered. This includes reading environment variables,
command line flags, and configuration files. Validate will check not only the 
syntax of these, but will also perform sanity checks to make sure that the 
configuration makes sense.`,
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		common()
		// TODO(synful)
	}
	cmdValidate.Run = f
	cmdMain.AddCommand(cmdValidate)
}
