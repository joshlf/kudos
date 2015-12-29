package main

import (
	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration",
	Long: `Perform all of the steps necessary to determine kudos' configuration, 
reporting any errors encountered. This includes reading environment variables,
command line flags, and configuration files. Validate will check not only the 
syntax of these, but will also perform sanity checks to make sure that the 
configuration makes sense.`,
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		ctx := getContext()
		addCourseConfig(ctx)
		ctx.Verbose.Printf("using course %v and course path %v\n", ctx.Course.Code, ctx.CourseRoot())
		_, err := kudos.ParseAllAssignmentFiles(ctx)
		if err != nil {
			ctx.Error.Println("configuration failed to validate")
			dev.Fail()
		}
		ctx.Verbose.Println("configuration validated")
	}
	cmdValidate.Run = f
	addAllGlobalFlagsTo(cmdValidate.Flags())
	cmdMain.AddCommand(cmdValidate)
}
