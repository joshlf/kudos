package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdMain = &cobra.Command{
	Use:   "kudos",
	Short: "kudos is a simple grading system",
	Long: `kudos is a simple grading system made out of love and frustration by m, ezr,
and jliebowf`,
}

var verboseFlag bool
var quietFlag bool
var configFlag string
var courseFlag string

func main() {
	cmdMain.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "be more verbose than normal")
	cmdMain.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "be more quiet than normal; overrides --verbose")
	cmdMain.PersistentFlags().StringVarP(&configFlag, "config", "", DefaultConfigFile, "location of the global config file")
	cmdMain.PersistentFlags().StringVarP(&courseFlag, "course", "c", "", "course")
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

func devFail() {
	fmt.Fprintln(os.Stderr, "[dev] failing for lack of anything better to do")
	os.Exit(1)
}
