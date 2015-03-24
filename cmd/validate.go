package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "validate the configuration",
	Long:  ``,
}

func init() {
	var verbose bool
	cmdValidate.Flags().BoolVarP(&verbose, "verbose", "v", false, "display extra information about the configuration")

	f := func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Verbose!")
		} else {
			fmt.Println("Not verbose!")
		}
	}
	cmdValidate.Run = f
	cmdMain.AddCommand(cmdValidate)
}
