package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/log"
	"github.com/synful/kudos/lib/user"
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
		switch len(args) {
		case 0:
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
		case 1:
			asgns, err := config.ReadAllAssignments(course)
			if err != nil {
				log.Error.Printf("could not read assignments: %v\n", err)
				os.Exit(1)
			}
			a, ok := config.AssignmentByCode(asgns, args[0])
			if !ok {
				log.Error.Printf("no such assignment: %v\n", args[0])
				os.Exit(1)
			}
			if a.HasMultipleHandins() {
				// TODO(synful): print more useful message,
				// such as available handins?
				log.Error.Printf("assignment has multiple handins; please specify one\n")
				os.Exit(1)
			}
		case 2:
			asgns, err := config.ReadAllAssignments(course)
			if err != nil {
				log.Error.Printf("could not read assignments: %v\n", err)
				os.Exit(1)
			}
			a, ok := config.AssignmentByCode(asgns, args[0])
			if !ok {
				log.Error.Printf("no such assignment: %v\n", args[0])
				os.Exit(1)
			}
			handins := a.Handins()
			h, ok := config.HandinByCode(handins, args[1])
			if !ok {
				log.Error.Printf("no such handin: %v\n", args[1])
				os.Exit(1)
			}
			// TODO(synful): temporary to suppress compiler errors
			_ = h
		default:
			cmd.Help()
		}

		u, err := user.Current()
		if err != nil {
			log.Error.Printf("could not get current user: %v\n", err)
			os.Exit(1)
		}
		// TODO(synful): temporary to suppress compiler errors
		_ = u
	}
	cmdHandin.Run = f
	addAllGlobalFlagsTo(cmdHandin.Flags())
	cmdMain.AddCommand(cmdHandin)
}
