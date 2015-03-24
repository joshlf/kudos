package main

import "github.com/spf13/cobra"

var cmdMain = &cobra.Command{
	Use:   "kudos",
	Short: "kudos is a simple grading system",
	Long: `kudos is a simple grading system made out of love and frustration by m, ezr, 
and jliebowf`,
}

var verboseFlag bool
var quietFlag bool
var configFlag string

func main() {
	cmdMain.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "be more verbose than normal")
	cmdMain.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "be more quiet than normal; overrides --verbose")
	cmdMain.PersistentFlags().StringVarP(&configFlag, "config", "c", DefaultConfigFile, "location of the global config file")
	cmdMain.Execute()
}

func common() {
	if quietFlag {
		SetLoggingLevel(Warn)
	} else if verboseFlag {
		SetLoggingLevel(Verbose)
	}

	InitConfig(configFlag)
}
