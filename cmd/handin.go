package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/log"
)

var cmdHandin = &cobra.Command{
	// TODO(synful): figure out how to
	// mark the handin as optional in a
	// way that is consistent with Cobra's
	// output.
	Use:   "handin [assignment] [handin]",
	Short: "hand in an assignment",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		common()
		course := requireCourseConfig()
		if len(args) == 0 {
			log.Info.Printf("Usage: %v\n\n", cmd.Use)
			asgns, err := config.ReadAllAssignments(course)
			if err != nil {
				log.Error.Printf("could not read assignments: %v\n", err)
				os.Exit(1)
			}
			log.Info.Println("Available handins:")
			for _, a := range asgns {
				if !a.HasMultipleHandins() {
					log.Info.Printf("  %v\n", a.Code())
				} else {
					// TODO(synful): maybe change the output
					// format? This works for now, but we
					// could think of something better.
					log.Info.Printf("  %v [", a.Code())
					h := a.Handins()
					for _, hh := range h[:len(h)-1] {
						log.Info.Printf("%v | ", hh.Code)
					}
					log.Info.Printf("%v]\n", h[len(h)-1].Code)
				}
			}
			return
		}
		// TODO(synful)
	}
	cmdHandin.Run = f
	addAllGlobalFlagsTo(cmdHandin.Flags())
	cmdMain.AddCommand(cmdHandin)
}
