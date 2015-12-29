package main

import (
	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/log"
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize a course's kudos installation",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		ctx := getContext()
		addCourse(ctx)

		err := kudos.InitCourse(ctx)
		if err != nil {
			log.Error.Printf("initialization failed: %v\n", err)
			dev.Fail()
		}
	}
	cmdInit.Run = f
	addAllGlobalFlagsTo(cmdInit.Flags())
	cmdMain.AddCommand(cmdInit)
}
