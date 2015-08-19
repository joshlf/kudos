package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/log"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize a course's kudos installation",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		common()
		requireCourse()
		err := config.InitCourse(config.Config.Course, config.Config.CoursePath, true)
		if err != nil {
			log.Error.Printf("init: %v\n", err)
			os.Exit(1)
		}
	}
	cmdInit.Run = f
	addAllGlobalFlagsTo(cmdInit.Flags())
	cmdMain.AddCommand(cmdInit)
}
