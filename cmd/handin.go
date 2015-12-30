package main

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/joshlf/kudos/lib/dev"
	"github.com/joshlf/kudos/lib/handin"
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/joshlf/kudos/lib/log"
	"github.com/spf13/cobra"
)

var cmdHandin = &cobra.Command{
	// TODO(joshlf): figure out how to
	// mark the handin as optional in a
	// way that is consistent with Cobra's
	// output.
	Use:   "handin [assignment] [handin]",
	Short: "Hand in an assignment",
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		ctx := getContext()
		addCourseConfig(ctx)

		var handinFile string
		switch len(args) {
		case 0:
			ctx.Info.Printf("Usage: %v\n\n", cmd.Use)
			asgns, err := kudos.ParseAllAssignmentFiles(ctx)
			if err != nil {
				ctx.Error.Println("could not read all assignments; aborting")
				dev.Fail()
			}
			ctx.Info.Println("Available handins:")
			for _, a := range asgns {
				if len(a.Handins) == 1 {
					ctx.Info.Printf("  %v\n", a.Code)
				} else {
					// TODO(joshlf): maybe change the output
					// format? This works for now, but we
					// could think of something better.
					ctx.Info.Printf("  %v [", a.Code)
					h := a.Handins
					for _, hh := range h[:len(h)-1] {
						ctx.Info.Printf("%v | ", hh.Code)
					}
					ctx.Info.Printf("%v]\n", h[len(h)-1].Code)
				}
			}
			exitClean()
		case 1:
			asgns, err := kudos.ParseAllAssignmentFiles(ctx)
			if err != nil {
				ctx.Error.Println("could not read all assignments; aborting")
				dev.Fail()
			}
			a, ok := kudos.FindAssignmentByCode(asgns, args[0])
			if !ok {
				ctx.Error.Printf("no such assignment: %v\n", args[0])
				dev.Fail()
			}
			if len(a.Handins) > 1 {
				// TODO(joshlf): print more useful message,
				// such as available handins?
				ctx.Error.Printf("assignment has multiple handins; please specify one\n")
				dev.Fail()
			}
			u, err := user.Current()
			if err != nil {
				ctx.Error.Printf("could not get current user: %v\n", err)
				dev.Fail()
			}
			handinFile = ctx.UserAssignmentHandinFile(args[0], u.Uid)
		case 2:
			asgns, err := kudos.ParseAllAssignmentFiles(ctx)
			if err != nil {
				ctx.Error.Println("could not read all assignments; aborting")
				dev.Fail()
			}
			a, ok := kudos.FindAssignmentByCode(asgns, args[0])
			if !ok {
				ctx.Error.Printf("no such assignment: %v\n", args[0])
				dev.Fail()
			}
			_, ok = a.FindHandinByCode(args[1])
			if !ok {
				ctx.Error.Printf("no such handin: %v\n", args[1])
				dev.Fail()
			}
			u, err := user.Current()
			if err != nil {
				ctx.Error.Printf("could not get current user: %v\n", err)
				dev.Fail()
			}
			handinFile = ctx.UserHandinHandinFile(args[0], args[1], u.Uid)
		default:
			cmd.Help()
			dev.Fail()
		}

		/*
			Perform handin
		*/

		hook := ctx.PreHandinHookFile()
		doHook := true
		_, err := os.Stat(hook)
		if err != nil {
			if os.IsNotExist(err) {
				doHook = false
			} else {
				ctx.Error.Printf("could not stat pre-handin hook: %v\n", err)
				dev.Fail()
			}
		}

		c := exec.Command(hook)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		// TODO(joshlf): Set environment variables

		if doHook {
			err = c.Run()
			if err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					ctx.Warn.Println("pre-handin hook exited with error code; aborting")
					dev.Fail()
				}
				ctx.Error.Printf("could not run pre-handin hook: %v\n", err)
				dev.Fail()
			}
		}

		ctx.Info.Println("Handing in the following files:")
		printFiles := ctx.Logger.GetLevel() <= log.Info
		err = handin.PerformFaclHandin(handinFile, printFiles)
		if err != nil {
			ctx.Error.Printf("could not hand in: %v\n", err)
			dev.Fail()
		}
		ctx.Info.Println("Handin successful.")
	}
	cmdHandin.Run = f
	addAllGlobalFlagsTo(cmdHandin.Flags())
	cmdMain.AddCommand(cmdHandin)
}
