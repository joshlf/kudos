package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/log"
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
		common()
		requireCourse()
		log.Verbose.Printf("using course %v and course path %v\n", config.Config.Course, config.Config.CoursePath)
		_, err := config.ReadCourseConfig(config.Config.Course, config.Config.CoursePath)
		if err != nil {
			log.Error.Printf("%v\n", err)
			log.Error.Println("configuration failed to validate")
			os.Exit(1)
		}
		log.Verbose.Println("configuration validated")
	}
	cmdValidate.Run = f
	addAllGlobalFlagsTo(cmdValidate.Flags())
	cmdMain.AddCommand(cmdValidate)
}
