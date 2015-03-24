package main

import "github.com/spf13/cobra"

var cmdMain = &cobra.Command{
	Use:   "kudos",
	Short: "kudos is a simple grading system",
	Long:  `Made out of love and frustration by m, ezr, and jliebowf`,
}

func main() {
	cmdMain.Execute()
}
