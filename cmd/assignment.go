package main

import (
	"github.com/joshlf/kudos/lib/kudos"
	"github.com/spf13/cobra"
)

var cmdAssignment = &cobra.Command{
	Use:   "assignment",
	Short: "Manage assignments",
	// TODO(joshlf): long description
}

func init() {
	f := func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		exitUsage()
	}
	cmdAssignment.Run = f
	cmdMain.AddCommand(cmdAssignment)
}

var cmdAssignmentPublish = &cobra.Command{
	Use:   "publish <assignment>",
	Short: "Make an assignment public",
	// TODO(joshlf): long description
}

func init() {
	var forceFlag bool
	f := func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			exitUsage()
		}
		ctx := getContext()
		code := args[0]
		if err := kudos.ValidateCode(code); err != nil {
			ctx.Error.Printf("bad assignment code %q: %v\n", code, err)
			exitUsage()
		}

		addCourseConfig(ctx)
		checkIsTA(ctx)
		openDB(ctx)
		defer cleanupDB(ctx)

		asgn, ok := ctx.DB.Assignments[code]
		if !ok {
			ctx.Error.Println("no such assignment")
			exitLogic()
		}

		openPubDB(ctx)
		defer cleanupPubDB(ctx)

		if _, ok := ctx.PubDB.Assignments[code]; ok {
			if forceFlag {
				ctx.Warn.Println("warning: overwriting previous assignment")
			} else {
				ctx.Error.Println("assignment already in public database; use --force to overwrite")
				exitLogic()
			}
		}

		ctx.PubDB.Assignments[code] = kudos.AssignmentToPub(asgn)
		commitPubDB(ctx)
	}
	cmdAssignmentPublish.Run = f
	addAllGlobalFlagsTo(cmdAssignmentPublish.Flags())
	addAllTAFlagsTo(cmdAssignmentPublish.Flags())
	cmdAssignmentPublish.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite previous version of assignment in public database")
	cmdAssignment.AddCommand(cmdAssignmentPublish)
}
