package main

import "github.com/spf13/cobra"

var cmdMain = &cobra.Command{
	Use:   "kudos",
	Short: "kudos is a simple grading system",
}

func main() {
	cmdMain.Execute()
}
